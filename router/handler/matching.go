package handler

import (
	"context"
	"fmt"
	"github.com/xxarupakaxx/tic-tac-toe/game"
	"github.com/xxarupakaxx/tic-tac-toe/gen/proto"
	"github.com/xxarupakaxx/tic-tac-toe/util"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"sync"
	"time"
)

type MatchingHandler struct {
	sync.RWMutex
	Rooms       map[int32]*game.Room
	maxPlayerID int32
}

func (m *MatchingHandler) JoinRoom(request *proto.JoinRoomRequest, stream proto.MatchingService_JoinRoomServer) error {
	ctx, cancel := context.WithTimeout(stream.Context(), 2*time.Minute)
	defer cancel()

	m.Lock()

	m.maxPlayerID++
	me := &game.Player{PlayerID: m.maxPlayerID}

	for _, room := range m.Rooms {
		if room.Guest == nil {
			me.XO  = game.O
			room.Guest = me
			err := stream.Send(&proto.JoinRoomResponse{
				Room:   util.ConvertPBRoom(room),
				Me:     util.ConvertPBPlayer(room.Guest),
				Status: proto.JoinRoomResponse_MATCHED,
			})
			if err != nil {
				return err
			}

			m.Unlock()
			fmt.Printf("matched roomID = %v\n", room.ID)

			return nil
		}
	}

	me.XO = game.X
	room := &game.Room{
		ID:   int32(len(m.Rooms) + 1),
		Host: me,
	}
	m.Rooms[room.ID] = room
	m.Unlock()

	err := stream.Send(&proto.JoinRoomResponse{
		Room:   util.ConvertPBRoom(room),
		Status: proto.JoinRoomResponse_WAITING,
	})
	if err != nil {
		return err
	}

	ch := make(chan int)
	go func(ch chan<- int) {
		for true {
			m.RLock()
			guest := room.Guest
			m.RUnlock()
			if guest != nil {
				err = stream.Send(&proto.JoinRoomResponse{
					Room:   util.ConvertPBRoom(room),
					Me:     util.ConvertPBPlayer(room.Host),
					Status: proto.JoinRoomResponse_MATCHED,
				})
				if err != nil {
					return
				}
				ch <- 0
				break
			}
			time.Sleep(1 * time.Second)

			select {
			case <-ctx.Done():
				return
			default:

			}
		}
	}(ch)

	select {
	case <-ch:
	case <-ctx.Done():
		return status.Errorf(codes.DeadlineExceeded, "マッチングできませんでした")
	}

	return nil
}

func NewMatchingHandler() *MatchingHandler {
	return &MatchingHandler{
		Rooms: make(map[int32]*game.Room),
	}
}
