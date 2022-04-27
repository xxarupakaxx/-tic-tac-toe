package client

import (
	"context"
	"fmt"
	"github.com/xxarupakaxx/tic-tac-toe/game"
	"github.com/xxarupakaxx/tic-tac-toe/gen/proto"
	"google.golang.org/grpc"
	"sync"
)

type TicTacToe struct {
	sync.RWMutex
	started  bool
	finished bool
	me       *game.Player
	Room     *game.Room
	Game     *game.TicTacToe
}

func NewTicTacToe() *TicTacToe {
	return &TicTacToe{}
}

func run() error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	conn, err := grpc.Dial(":50051", grpc.WithInsecure())
	if err != nil {
		return fmt.Errorf("failed to connect grpc server :%w", err)
	}
	defer conn.Close()


}

func matching(ctx context.Context, client proto.MatchingServiceClient) error {
	stream,err := client.JoinRoom(ctx,&proto.JoinRoomRequest{})
}