package gateservice

import (
	"alexzhaozzzz/game_service_ex/api/proto/msg"
	"alexzhaozzzz/game_service_ex/api/proto/rpc"
	"alexzhaozzzz/game_service_ex/internal/service/gateservice/tcpmodule"
	"alexzhaozzzz/game_service_ex/internal/service/gateservice/wsmodule"
	"alexzhaozzzz/game_service_ex/pkg/util"

	"github.com/duanhf2012/origin/v2/log"
	"github.com/duanhf2012/origin/v2/network/processor"
	"github.com/duanhf2012/origin/v2/node"
	"github.com/duanhf2012/origin/v2/service"
	"google.golang.org/protobuf/proto"

	"errors"
	"fmt"
	"time"
)

func init() {
	node.Setup(&GateService{})
}

type GateService struct {
	service.Service

	netModule      INetModule
	pbRawProcessor processor.PBRawProcessor
	rawPackInfo    processor.PBRawPackInfo

	msgRouter MsgRouter
}

func (s *GateService) OnInit() error {
	s.msgRouter.Init(&s.pbRawProcessor)

	_, _ = s.AddModule(&s.msgRouter)

	iConfig := s.GetService().GetServiceCfg()
	if iConfig == nil {
		return fmt.Errorf("%s config is error", s.GetService().GetName())
	}

	//TcpCfg与WSCfg取其一
	mapTcpCfg := iConfig.(map[string]interface{})
	_, tcpOk := mapTcpCfg["TcpCfg"]
	if tcpOk == true {
		var tcpModule tcpmodule.TcpModule
		tcpModule.SetProcessor(&s.pbRawProcessor)
		_, _ = s.AddModule(&tcpModule)
		s.msgRouter.SetNetModule(&tcpModule)
		s.netModule = &tcpModule
	} else {
		_, wsOk := mapTcpCfg["WSCfg"]
		if wsOk == true {
			var wsModule wsmodule.WSModule
			wsModule.SetProcessor(&s.pbRawProcessor)
			_, _ = s.AddModule(&wsModule)
			s.msgRouter.SetNetModule(&wsModule)
			s.netModule = &wsModule
		} else {
			return errors.New("WSCfg and TcpCfg are not configured")
		}
	}

	s.RegRawRpc(util.RawRpcMsgDispatch, s.RawRpcDispatch)
	s.RegRawRpc(util.RawRpcCloseClient, s.RawCloseClient)
	return nil
}

func (s *GateService) RPC_GSLoginRet(arg *rpc.GsLoginResult, ret *rpc.PlaceHolders) error {
	v, ok := s.msgRouter.mapRouterCache[arg.ClientId]
	if ok == false {
		log.SWarning("Client is close cancel login ")
		return nil
	}

	v.status = Logined
	s.msgRouter.mapRouterCache[arg.ClientId] = v

	var loginRes msg.MsgLoginRes
	loginRes.SessionId = arg.SessionId
	loginRes.UserId = arg.UserId
	loginRes.ServerTime = time.Now().UnixMilli()
	err := s.msgRouter.SendMsg(arg.ClientId, msg.MsgType_LoginRes, &loginRes)
	if err != nil {
		log.SError("GateService.RPC_GSLoginRet ClientId:", arg.ClientId, " SendMsg err:", err.Error())
	}

	return nil
}

func (s *GateService) SendMsg(clientId string, msgType uint16, rawMsg []byte) error {
	s.rawPackInfo.SetPackInfo(msgType, rawMsg)
	bytes, err := s.pbRawProcessor.Marshal(clientId, &s.rawPackInfo)
	if err != nil {
		return err
	}

	err = s.netModule.SendRawMsg(clientId, bytes)
	if err != nil {
		log.Debug("SendMsg fail ", log.ErrorAttr("err", err), log.String("clientId", clientId))
	}

	return err
}

func (s *GateService) RawRpcDispatch(rawInput []byte) {
	var rawInputArgs rpc.RawInputArgs
	err := proto.Unmarshal(rawInput, &rawInputArgs)
	if err != nil {
		log.SError("msg is error:%s", err.Error())
		return
	}

	for _, clientId := range rawInputArgs.ClientIdList {
		err = s.SendMsg(clientId, uint16(rawInputArgs.MsgType), rawInputArgs.RawData)
		if err != nil {
			log.SError("SendRawMsg fail:", err.Error())
		}
	}

	//消息统计
	s.msgRouter.performanceAnalyzer.ChangeDeltaNum(MsgAnalyzer, int(rawInputArgs.MsgType), MsgSendNumColumn, int64(len(rawInputArgs.ClientIdList)))
	s.msgRouter.performanceAnalyzer.ChangeDeltaNum(MsgAnalyzer, int(rawInputArgs.MsgType), MsgSendByteColumn, int64(len(rawInputArgs.ClientIdList)*(len(rawInputArgs.RawData)+2))) //排除掉消息ID长度
}

func (s *GateService) RawCloseClient(rawInput []byte) {
	var rawInputArgs rpc.RawInputArgs
	err := proto.Unmarshal(rawInput, &rawInputArgs)
	if err != nil {
		log.Error("msg is error", log.ErrorAttr("err", err))
		return
	}

	for _, clientId := range rawInputArgs.ClientIdList {
		s.netModule.Close(clientId)
	}
}
