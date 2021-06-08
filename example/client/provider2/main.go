package main

import (
	"context"
	"github.com/busgo/elsa/example/client/proto/pb"
	"github.com/busgo/elsa/pkg/client"
	"github.com/busgo/elsa/pkg/log"
	"google.golang.org/grpc"
	"net"
)

type TradeGRPC struct {
	pb.UnimplementedTradeServiceServer
}

func (t *TradeGRPC) Ping(ctx context.Context, request *pb.PingRequest) (*pb.PingResponse, error) {
	log.Infof("receive ping:%s", request.Ping)
	return &pb.PingResponse{Pong: "ok"}, nil
}

func main() {

	stub, err := client.NewRegistryStub("dev", []string{"127.0.0.1:8005"})
	if err != nil {
		panic(err)
	}

	managedSentinel := client.NewManagedSentinel(8002, stub)
	managedSentinel.AddGrpcService(pb.TradeService_ServiceDesc.ServiceName)

	server := grpc.NewServer()
	pb.RegisterTradeServiceServer(server, new(TradeGRPC))

	l, err := net.Listen("tcp", ":8002")
	if err != nil {
		panic(err)
	}

	log.Infof("the trade service has start...")
	if err = server.Serve(l); err != nil {
		panic(err)
	}

}
