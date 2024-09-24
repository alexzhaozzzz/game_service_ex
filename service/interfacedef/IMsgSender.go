package interfacedef

import (
	"github.com/alexzhaozzzz/game_service_ex/common/proto/msg"
	"google.golang.org/protobuf/proto"
)

type IMsgSender interface {
	SendToClient(clientId string, msgType msg.MsgType, message proto.Message) int
	CastToClient(clientIdList []string, msgType msg.MsgType, message proto.Message) int
	SendToPlayer(playerUserId string, msgType msg.MsgType, message proto.Message) int
	CastToPlayer(playerUserId []string, msgType msg.MsgType, message proto.Message) int
}
