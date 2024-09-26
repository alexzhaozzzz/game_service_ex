package main

import (
	_ "alexzhaozzzz/game_service_ex/gateservice"

	"time"

	"github.com/duanhf2012/origin/v2/node"
)

func main() {
	node.OpenProfilerReport(time.Second * 10)
	node.Start()
}
