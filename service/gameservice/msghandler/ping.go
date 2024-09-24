package msghandler

import (
	"github.com/alexzhaozzzz/game_service_ex/common/proto/msg"
	"github.com/alexzhaozzzz/game_service_ex/service/gameservice/player"
)

func ping(player *player.Player, msg *msg.MsgNil) {
	player.Ping()
}
