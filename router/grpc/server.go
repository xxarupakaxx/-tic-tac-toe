package grpc

import (
	"fmt"
	grpc_logrus "github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus"
	logrus2 "github.com/sirupsen/logrus"
	"github.com/xxarupakaxx/tic-tac-toe/gen/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
	"os"
	"os/signal"
)

func main() {
	port := 50051
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatalf("failed to listen :%w", err)
	}
	logusLogger := logrus2.New()
	logrusEnty := logrus2.NewEntry(logusLogger)
	grpc_logrus.ReplaceGrpcLogger(logrusEnty)
	server := grpc.NewServer(
		grpc.UnaryInterceptor(
			grpc_logrus.UnaryServerInterceptor(logrusEnty)))

	proto.RegisterMatchingServiceServer()
	proto.RegisterTicServiceServer()

	reflection.Register(server)

	go func() {
		log.Println("start gRPC server port ", port)

		if err = server.Serve(lis); err != nil {
			log.Println("failed to start server : ", err)
		}

	}()

	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Println("stopping grpc Server")
	server.GracefulStop()
}
