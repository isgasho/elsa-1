package p2p

import (
	"github.com/busgo/elsa/pkg/proto/pb"
	"testing"
	"time"
)

var endpoints = []string{"127.0.0.1:8005", "192.168.1.1:8005"}

// test new peer pool
func TestNewPeerPoolWithEndpoints(t *testing.T) {

	pool, err := NewPeerPoolWithEndpoints(endpoints)
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("the peer pool is %#v", pool)
}

func TestPeerPool_Start(t *testing.T) {

	pool, err := NewPeerPoolWithEndpoints(endpoints)
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("the peer pool is %#v", pool)

	go pool.Start()
	t.Logf("the peer pool has start...")
}

// test push sync message
func TestPeerPool_PushMsg(t *testing.T) {

	pool, err := NewPeerPoolWithEndpoints(endpoints)
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("the peer pool is %#v", pool)

	go pool.Start()
	t.Logf("the peer pool has start...")

	pool.PushMsg(&SyncMsg{
		Type: SyncMsgRegType,
		Content: &pb.RegisterRequest{
			Segment:         "dev",
			ServiceName:     pb.RegistryService_ServiceDesc.ServiceName,
			Ip:              "192.168.1.1",
			Port:            8001,
			Metadata:        make(map[string]string),
			RegTimestamp:    time.Now().UnixNano(),
			UpTimestamp:     time.Now().UnixNano(),
			RenewTimestamp:  time.Now().UnixNano(),
			DirtyTimestamp:  time.Now().UnixNano(),
			LatestTimestamp: time.Now().UnixNano(),
			SyncType:        pb.SyncTypeEnum_None,
		},
	})

	pool.PushMsg(&SyncMsg{
		Type: SyncMsgRenewType,
		Content: &pb.RenewRequest{
			Segment:     "dev",
			ServiceName: pb.RegistryService_ServiceDesc.ServiceName,
			Ip:          "192.168.1.1",
			Port:        8001,
			SyncType:    pb.SyncTypeEnum_None,
		},
	})

	pool.PushMsg(&SyncMsg{
		Type: SyncMsgCancelType,
		Content: &pb.CancelRequest{
			Segment:     "dev",
			ServiceName: pb.RegistryService_ServiceDesc.ServiceName,
			Ip:          "192.168.1.1",
			Port:        8001,
			SyncType:    pb.SyncTypeEnum_None,
		},
	})
	time.Sleep(time.Second * 2)
}
