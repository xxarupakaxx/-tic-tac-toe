package util

import (
	"github.com/xxarupakaxx/tic-tac-toe/game"
	"github.com/xxarupakaxx/tic-tac-toe/gen/proto"
)

func ConvertPBPlayer(player *game.Player) *proto.Player {
	return &proto.Player{
		PlayerID: player.PlayerID,
		Xo:       ConvertPBXO(player.XO),
	}
}

func ConvertGamePlayer(player *proto.Player) *game.Player {
	return &game.Player{
		PlayerID: player.GetPlayerID(),
		XO:       ConvertGameXO(player.GetXo()),
	}
}

