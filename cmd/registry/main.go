package main

import (
	"github.com/busgo/elsa/internal/registry"
	"github.com/busgo/elsa/pkg/log"
	"time"
)

func main() {

	log.Info("ddd")
	log.Debug("debug...")

	r := registry.NewRegistry()
	log.Infof("registry:%#v", r)
	r.Register(&registry.Instance{

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
	})
	for {

		time.Sleep(time.Second)
	}
}
