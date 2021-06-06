package main

import (
	"context"
	"github.com/busgo/elsa/pkg/log"
	"github.com/busgo/elsa/pkg/proto/pb"
	"google.golang.org/grpc"
	"time"
)

func main() {

	cc, err := grpc.Dial("127.0.0.1:8005", grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	cli := pb.NewRegistryServiceClient(cc)

	response, err := cli.Register(context.Background(), &pb.RegisterRequest{
		Segment:         "dev",
		ServiceName:     "com.busgo.trade.proto.TradeService",
		Ip:              "192.168.1.1",
		Port:            8001,
		Metadata:        make(map[string]string),
		RegTimestamp:    time.Now().UnixNano(),
		UpTimestamp:     time.Now().UnixNano(),
		RenewTimestamp:  time.Now().UnixNano(),
		DirtyTimestamp:  time.Now().UnixNano(),
		LatestTimestamp: time.Now().UnixNano(),
		SyncType:        pb.SyncTypeEnum_None,
	})

	if err != nil {
		panic(err)
	}

	log.Infof("response code:%d,instance:%v", response.Code, response.Instance)

}
