package client

import (
	"bufio"
	"context"
	"fmt"
	"github.com/xxarupakaxx/tic-tac-toe/game"
	"github.com/xxarupakaxx/tic-tac-toe/gen/proto"
	"github.com/xxarupakaxx/tic-tac-toe/util"
	"google.golang.org/grpc"
	"os"
	"strconv"
	"sync"
	"time"
)

type TicTacToe struct {
	sync.RWMutex
	started  bool
	finished bool
	me       *game.Player
	room     *game.Room
	game     *game.TicTacToe
}

func NewTicTacToe() *TicTacToe {
	return &TicTacToe{}
}

func (t *TicTacToe) run() error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	conn, err := grpc.Dial(":50051", grpc.WithInsecure())
	if err != nil {
		return fmt.Errorf("faield to connect to grpc server:%w", err)
	}
	defer conn.Close()

	err = t.matching(ctx, proto.NewMatchingServiceClient(conn))
	if err != nil {
		return err
	}
	t.game = game.NewTicTacToe(t.me.XO)

	return t.play(ctx, proto.NewTicServiceClient(conn))

}

func (t *TicTacToe) matching(ctx context.Context, client proto.MatchingServiceClient) error {
	stream, err := client.JoinRoom(ctx)
	if err != nil {
		return err
	}

	defer stream.CloseSend()

	fmt.Println("マッチング相手を探しております...")

	for true {
		resp, err := stream.Recv()
		if err != nil {
			return err
		}

		if resp.GetStatus() == proto.JoinRoomResponse_MATCHED {
			t.room = util.ConvertGameRoom(resp.GetRoom())
			t.me = util.ConvertGamePlayer(resp.GetMe())
			fmt.Printf("Matched roomID=%d\n", resp.GetRoom().GetId())
			return nil
		} else if resp.GetStatus() == proto.JoinRoomResponse_WAITING {
			fmt.Println("waiting matching..")
		}
	}
	return nil
}

func (t *TicTacToe) play(ctx context.Context, client proto.TicServiceClient) error {
	c, cancel := context.WithCancel(ctx)
	defer cancel()
	stream, err := client.Play(c)
	if err != nil {
		return err
	}
	defer stream.CloseSend()

	go func() {
		err = t.send(c, stream)
		if err != nil {
			cancel()
		}
	}()

	err = t.receive(c, stream)
	if err != nil {
		cancel()
		return err
	}

	return nil
}

func (t *TicTacToe) send(ctx context.Context, stream proto.TicService_PlayClient) error {
	for true {
		t.RLock()

		if t.finished {
			t.RUnlock()
			return nil
		} else if !t.started {
			err := stream.Send(&proto.PlayerRequest{
				RoomID: t.room.ID,
				Player: util.ConvertPBPlayer(t.me),
				Action: &proto.PlayerRequest_Start{
					Start: &proto.PlayerRequest_StartAction{},
				},
			})
			t.RUnlock()
			if err != nil {
				return err
			}

			for true {
				t.RLock()
				if t.started {
					t.RUnlock()
					fmt.Printf("対戦見つかったね")
					break
				}
				t.RUnlock()
				fmt.Println("対戦相手が見つかるまで待とうね")
				time.Sleep(time.Second * 1)

			}
		} else {
			t.RUnlock()
			fmt.Println("どの石を動かす？")
			stdin := bufio.NewScanner(os.Stdin)
			stdin.Scan()

			text := stdin.Text()
			number, err := parseInput(text)
			if err != nil {
				fmt.Println(err)
				continue
			}

			t.Lock()
			t.game.Board.DisplayBoard(t.me.XO)
			t.Unlock()

			go func() {
				err = stream.Send(&proto.PlayerRequest{
					RoomID: t.room.ID,
					Player: util.ConvertPBPlayer(t.me),
					Action: &proto.PlayerRequest_Play{
						Play: &proto.PlayerRequest_PlayAction{
							Number: number,
						},
					},
				})
				if err != nil {
					fmt.Println(err)
				}
			}()

			ch := make(chan int)
			go func(ch chan int) {
				fmt.Println("")
				for i := 0; i < 5; i++ {
					fmt.Printf("%d秒間止まります \n", 5-i)
					time.Sleep(1 * time.Second)
				}
				fmt.Println("")
				ch <- 0
			}(ch)
			<-ch

		}

		select {
		case <-ctx.Done():
			return nil
		default:

		}

	}
	return nil
}

func parseInput(text string) (int32, error) {
	number, err := strconv.Atoi(text)
	if err != nil {
		return 0, fmt.Errorf("入力が正しくない :%w", err)
	}

	return int32(number), err
}

func (t *TicTacToe) receive(ctx context.Context, stream proto.TicService_PlayClient) error {
	for true {
		res, err := stream.Recv()
		if err != nil {
			return err
		}

		t.Lock()
		switch res.GetEvent().(type) {
		case *proto.PlayerResponse_Waiting:

		case *proto.PlayerResponse_Ready:
			t.started = true
			t.game.Board.DisplayBoard(t.me.XO)
		case *proto.PlayerResponse_Play:
			xo := util.ConvertGameXO(res.GetPlay().GetPlayer().GetXo())
			if xo != t.me.XO {
				t.game.Board.DisplayBoard(xo)
				fmt.Print("石をどこに動かしますか 例 1")
			}
		case *proto.PlayerResponse_Finish:
			t.finished = true

			winner := util.ConvertGamePlayer(res.GetPlay().GetPlayer())
			fmt.Println("")
			if winner.XO == game.UNKNOWN {
				fmt.Println("Draw!")
			} else if winner.XO == t.me.XO {
				fmt.Println("you win!")
			} else {
				fmt.Println("You Lose!")
			}

			t.Unlock()
			return nil

		}
		t.Unlock()

		select {
		case <-ctx.Done():
			return nil
		default:

		}
	}
	return nil
}
