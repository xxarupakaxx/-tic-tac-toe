package util

import (
	"github.com/xxarupakaxx/tic-tac-toe/game"
	"github.com/xxarupakaxx/tic-tac-toe/gen/proto"
)

func ConvertPBXO(xo game.XO) proto.XO {
	switch xo {
	case game.X:
		return proto.XO_X
	case game.O:
		return proto.XO_O
	case game.UNKNOWN:
		return proto.XO_UNKNOWN
	}

	panic("あり得ないタイプ")
}

func ConvertGameXO(xo proto.XO) game.XO {
	switch xo {
	case proto.XO_X:
		return game.X
	case proto.XO_O:
		return game.O
	case proto.XO_UNKNOWN:
		return game.UNKNOWN
	}

	panic("あり得ない")
}

func ConvertGameResult(result game.Result) proto.Result {
	switch result {
	case game.WINNER:
		return proto.Result_WIN
	case game.LOSE:
		return proto.Result_LOSE
	case game.DRAW:
		return proto.Result_DRAW
	}

	panic("なにそれ")
}
