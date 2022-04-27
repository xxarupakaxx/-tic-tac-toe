package game

type XO int32
type Result int32
const (
	X XO =iota
	O
)

const (
	WINNER = iota
	LOSE
	DRAW
)

func ConvertToXO(a XO) string {
	switch a {
	case X:
		return "X"
	case O:
		return "O"
	}

	return ""
}

func ConvertWinner(a Result) string {
	switch a {
	case WINNER:
		return "勝ち"
	case LOSE:
		return "負け"
	default:
		return "引き分け"
	}
}


