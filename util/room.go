package util

import (
	"github.com/xxarupakaxx/tic-tac-toe/game"
	"github.com/xxarupakaxx/tic-tac-toe/gen/proto"
)

func ConvertPBRoom(room *game.Room) *proto.Room {
	return &proto.Room{
		Id:    room.ID,
		Host:  ConvertPBPlayer(room.Host),
		Guest: ConvertPBPlayer(room.Guest),
	}
}

func ConvertGameRoom(room *proto.Room) *game.Room {
	return &game.Room{
		ID:    room.Id,
		Host:  ConvertGamePlayer(room.GetHost()),
		Guest: ConvertGamePlayer(room.GetGuest()),
	}
}
