package game

import "fmt"

type Board struct {
	Line []XO
}

func NewBoard() *Board {
	line := make([]XO, 9)
	for i := 0; i < 9; i++ {
		line[i] = UNKNOWN
	}

	return &Board{Line: line}
}

func (b *Board)DisplayBoard(me XO)  {
	fmt.Println("")
	if me != UNKNOWN {
		fmt.Println("You: ",ConvertToXO(me))
	}

	fmt.Println("┌───┬───┬───┐")
	for i := 0; i < 3; i++ {
		fmt.Print("│ ")
		fmt.Printf("%s ", b.Line[i])
	}
	fmt.Println("│")

	fmt.Println("├───┼───┼───┤")
	for i := 3; i < 6; i++ {
		fmt.Print("│ ")
		fmt.Printf("%s ", b.Line[i])
	}
	fmt.Println("│")

	fmt.Println("├───┼───┼───┤")
	for i := 6; i < 9; i++ {
		fmt.Print("│ ")
		fmt.Printf("%s ", b.Line[i])
	}
	fmt.Println("│")

	fmt.Println("└───┴───┴───┘")
}