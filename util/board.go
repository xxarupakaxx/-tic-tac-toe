package util

import (
	"github.com/xxarupakaxx/tic-tac-toe/game"
	"github.com/xxarupakaxx/tic-tac-toe/gen/proto"
)

func ConvertGameBoard(board *game.Board) *proto.Board {
	xos := make([]proto.XO, 0, 10)

	for _, xo := range board.Line {
		xos = append(xos, ConvertPBXO(xo))
	}

	return &proto.Board{Xos: xos}
}
