package main

import (
	"context"
	"github.com/busgo/elsa/example/client/proto/pb"
	"github.com/busgo/elsa/pkg/client"
	"github.com/busgo/elsa/pkg/log"
	"google.golang.org/grpc"
)

type TradeGRPC struct {
	pb.UnimplementedTradeServiceServer
}

func (t *TradeGRPC) Ping(ctx context.Context, request *pb.PingRequest) (*pb.PingResponse, error) {

	log.Infof("receive:%s", request.Ping)
	return &pb.PingResponse{Pong: "ok"}, nil
}

func main() {

	stub, err := client.NewRegistryStub("dev", []string{"127.0.0.1:8005"})
	if err != nil {
		panic(err)
	}
	elsaServer, err := client.NewElsaServer(client.WithName("trade"),
		client.WithServerPort(8001),
		client.WithRegistryStub(stub))
	if err != nil {
		panic(err)
	}

	elsaServer.Init(func(server *grpc.Server) (serverNames []string) {
		pb.RegisterTradeServiceServer(server, new(TradeGRPC))
		return []string{
			pb.TradeService_ServiceDesc.ServiceName,
		}
	})

	if err = elsaServer.Start(); err != nil {
		log.Errorf("start the elsa server fail:%s", err.Error())
		panic(err)
	}

}
