package main

import (
	_ "github.com/alexzhaozzzz/game_service_ex/service/authservice"
	_ "github.com/alexzhaozzzz/game_service_ex/service/botservice"
	_ "github.com/alexzhaozzzz/game_service_ex/service/centerservice"
	_ "github.com/alexzhaozzzz/game_service_ex/service/dbservice"
	_ "github.com/alexzhaozzzz/game_service_ex/service/gameservice"
	_ "github.com/alexzhaozzzz/game_service_ex/service/gateservice"
	_ "github.com/alexzhaozzzz/game_service_ex/service/hotloadservice"
	_ "github.com/alexzhaozzzz/game_service_ex/service/httpgateservice"
	"github.com/duanhf2012/origin/v2/node"
	"time"
)

func main() {
	node.OpenProfilerReport(time.Second * 10)
	node.Start()
}
