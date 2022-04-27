package game

type TicTacToe struct {
	Me    XO
	Board *Board
}

func NewTicTacToe(me XO) *TicTacToe {
	board := NewBoard()
	return &TicTacToe{Me: me, Board: board}
}

func (t *TicTacToe) Logic() XO {
	matrix := [][]int32{
		{0, 1, 2},
		{3, 4, 5},
		{6, 7, 8},
		{0, 3, 6},
		{1, 4, 7},
		{2, 5, 8},
		{0, 4, 8},
		{2, 4, 6},
	}

	for i := 0; i < len(matrix); i++ {
		a := matrix[i][0]
		b := matrix[i][1]
		c := matrix[i][2]
		if t.Board.Line[a] != UNKNOWN && t.Board.Line[a] == t.Board.Line[b] && t.Board.Line[a] == t.Board.Line[c] {
			return t.Board.Line[a]
		}
	}

	return UNKNOWN
}
