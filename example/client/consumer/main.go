package main

import (
	"context"
	"github.com/busgo/elsa/example/client/proto/pb"
	"github.com/busgo/elsa/pkg/client"
	"github.com/busgo/elsa/pkg/log"
	"google.golang.org/grpc"
	"time"
)

func main() {

	stub, err := client.NewRegistryStub("dev", []string{"127.0.0.1:8005"})
	if err != nil {
		panic(err)
	}

	elsaServer, err := client.NewElsaServer(client.WithServerPort(8002), client.WithRegistryStub(stub), client.WithName("consumer"))
	if err != nil {
		panic(err)
	}

	elsaServer.Init(func(server *grpc.Server) (serverNames []string) {
		return make([]string, 0)
	})

	// build the stub
	cli := elsaServer.BuildStub(pb.TradeService_ServiceDesc.ServiceName, func(cc *grpc.ClientConn) interface{} {
		return pb.NewTradeServiceClient(cc)
	}).(pb.TradeServiceClient)

	go func() {
		for {
			response, err := cli.Ping(context.Background(), &pb.PingRequest{Ping: "ok"})
			if err != nil {
				log.Errorf("response error:%s", err.Error())
				time.Sleep(time.Millisecond * 500)
				continue
			}

			log.Infof("call trade service ping method  success pong:%s", response.Pong)
			time.Sleep(time.Millisecond * 500)
		}
	}()

	if err = elsaServer.Start(); err != nil {
		panic(err)
	}

}
