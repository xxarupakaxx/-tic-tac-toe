package handler

import (
	"fmt"
	"github.com/xxarupakaxx/tic-tac-toe/gen/proto"
	"sync"
)

type GameHandler struct {
	sync.RWMutex
}

func (g *GameHandler) Play(server proto.TicService_PlayServer) error {
	for true {
		req, err := server.Recv()
		if err != nil {
			return fmt.Errorf("failed to recv req,%w", err)
		}

		switch req.GetAction().(type){
		case *proto.PlayerRequest_Start:
			return start()
		case *proto.PlayerRequest_Play:
			return play()
		}
	}
}

func start() error {

}

func play() error {

}
