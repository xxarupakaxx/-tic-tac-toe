package handler

import (
	"fmt"
	"github.com/xxarupakaxx/tic-tac-toe/game"
	"github.com/xxarupakaxx/tic-tac-toe/gen/proto"
	"github.com/xxarupakaxx/tic-tac-toe/util"
	"sync"
)

type GameHandler struct {
	sync.RWMutex
	games  map[int32]*game.TicTacToe
	client map[int32][]proto.TicService_PlayServer
}

func NewGameHandler() *GameHandler {
	return &GameHandler{
		games:  make(map[int32]*game.TicTacToe),
		client: make(map[int32][]proto.TicService_PlayServer),
	}
}

func (g *GameHandler) Play(stream proto.TicService_PlayServer) error {
	for true {
		req, err := stream.Recv()
		if err != nil {
			return fmt.Errorf("failed to recv req,%w", err)
		}

		roomID := req.GetRoomID()
		player := util.ConvertGamePlayer(req.GetPlayer())

		switch req.GetAction().(type) {
		case *proto.PlayerRequest_Start:
			return g.start(stream, roomID)
		case *proto.PlayerRequest_Play:
			return g.play(roomID, player)
		}
	}
	return nil
}

func (g *GameHandler) start(stream proto.TicService_PlayServer, id int32) error {
	g.Lock()
	defer g.Unlock()

	ga := g.games[id]
	if ga == nil {
		ga = game.NewTicTacToe(game.UNKNOWN)
		g.games[id] = ga
		g.client[id] = make([]proto.TicService_PlayServer, 0, 2)
	}

	g.client[id] = append(g.client[id], stream)

	if len(g.client[id]) == 2 {
		for _, server := range g.client[id] {
			err := server.Send(&proto.PlayerResponse{Event: &proto.PlayerResponse_Ready{
				Ready: &proto.PlayerResponse_ReadyEvent{},
			}})
			if err != nil {
				return err
			}
		}
		fmt.Println("ゲームスタート！")
	} else {
		err := stream.Send(&proto.PlayerResponse{Event: &proto.PlayerResponse_Waiting{
			Waiting: &proto.PlayerResponse_WaitingEvent{},
		}})
		if err != nil {
			return err
		}
	}

	return nil
}

func (g *GameHandler) play(id int32, player *game.Player) error {
	g.Lock()
	defer g.Unlock()

	ga := g.games[id]

	winner := ga.Logic()
	for _, server := range g.client[id] {
		err := server.Send(&proto.PlayerResponse{
			Event: &proto.PlayerResponse_Play{
				Play: &proto.PlayerResponse_PlayEvent{
					Player: util.ConvertPBPlayer(player),
					Board:  util.ConvertGameBoard(g.games[id].Board),
				},
			},
		})
		if err != nil {
			return err
		}
		if winner != game.UNKNOWN {
			err = server.Send(&proto.PlayerResponse{
				Event: &proto.PlayerResponse_Finish{
					Finish: &proto.PlayerResponse_FinishEvent{
						Result: util.ConvertGameResult(game.Winner(winner, player.XO)),
						Board:  util.ConvertGameBoard(g.games[id].Board),
					},
				},
			})
			if err != nil {
				return err
			}

			delete(g.client, id)
		}
	}
	return nil
}
