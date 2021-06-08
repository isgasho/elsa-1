package main

import (
	"context"
	"github.com/busgo/elsa/example/client/proto/pb"
	"github.com/busgo/elsa/pkg/client"
	"github.com/busgo/elsa/pkg/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/resolver"
	"time"
)

func main() {

	stub, err := client.NewRegistryStub("dev", []string{"127.0.0.1:8005"})
	if err != nil {
		panic(err)
	}
	builder := client.NewElsaResolverBuilder(stub)
	resolver.Register(builder)
	cc, err := grpc.Dial(client.BuildTarget(builder.Scheme(), pb.TradeService_ServiceDesc.ServiceName), grpc.WithInsecure(), grpc.WithBalancerName("round_robin"))
	if err != nil {
		panic(err)
	}

	cli := pb.NewTradeServiceClient(cc)

	for {
		response, err := cli.Ping(context.Background(), &pb.PingRequest{Ping: "ok"})
		if err != nil {
			log.Errorf("response error:%s", err.Error())
			continue
		}

		log.Infof("call trade service ping method  success pong:%s", response.Pong)
		time.Sleep(time.Millisecond * 200)
	}

}
