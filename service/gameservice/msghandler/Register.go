package msghandler

import (
	"github.com/alexzhaozzzz/game_service_ex/common/proto/msg"
	"github.com/alexzhaozzzz/game_service_ex/service/gameservice/msgrouter"
)

func init() {
	msgrouter.RegMsgHandler(msg.MsgType_Ping, ping)
}
