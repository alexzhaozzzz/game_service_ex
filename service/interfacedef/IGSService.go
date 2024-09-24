package interfacedef

import (
	"github.com/alexzhaozzzz/game_service_ex/common/performance"
	"github.com/duanhf2012/origin/v2/rpc"
)

type IGSService interface {
	rpc.IRpcHandler
	GetAnalyzer(analyzerType int, analyzerId int) *performance.Analyzer
	GetGateNodeIdByClientId(clientId string) string
	GetClientPlayer(clientID string) IPlayer
	GetClientIdByPlayerId(playerId string) string
	CloseClient(clientId string)
	DestroyPlayer(playerId string) bool

	GetMsgReceiver() IMsgReceiver
	GetMsgSender() IMsgSender
	GetPlayerTimer() IPlayerTimer
}
