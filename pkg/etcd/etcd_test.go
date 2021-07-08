package etcd

import (
	"context"
	"github.com/busgo/elsa/pkg/log"
	"testing"
)

func TestNewEtcdClient(t *testing.T) {

	cli, err := NewEtcdClient([]string{"127.0.0.1:2379"}, "", "")
	if err != nil {
		t.Fatal(err)
	}

	t.Log("cli:", cli)
	err = cli.Put(context.Background(), "/one", "echo")
	if err != nil {
		t.Fatal(err)
	}

	response := cli.WatchWithPrefix(context.Background(), "/one")

	for {

		select {
		case event := <-response.Event:
			log.Infof("event:%+v", event)

		}
	}

}
