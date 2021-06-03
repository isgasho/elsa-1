package registry

import (
	"sync"
	"testing"
	"time"
)

var instance1 = &Instance{

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
}

var instance2 = &Instance{

	Segment:         "dev",
	ServiceName:     "com.busgo.trade.proto.TradeService",
	Ip:              "192.168.1.2",
	Port:            8001,
	Metadata:        make(map[string]string),
	RegTimestamp:    time.Now().UnixNano(),
	UpTimestamp:     time.Now().UnixNano(),
	RenewTimestamp:  time.Now().UnixNano(),
	DirtyTimestamp:  time.Now().UnixNano(),
	LatestTimestamp: time.Now().UnixNano(),
}

const (
	segment     = "dev"
	serviceName = "com.busgo.trade.proto.TradeService"
)

func initRegistry() Registry {
	return NewRegistry()
}

// register a instance
func TestRegistry_Register(t *testing.T) {

	r := initRegistry()
	//t.Logf("registry:%#v", r)

	in, err := r.Register(instance1)
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("register the instance success in:%#v", in)

}

// test fetch instances
func TestRegistry_Fetch(t *testing.T) {

	r := initRegistry()

	in, err := r.Register(instance1)
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("register the instance success in:%#v", in)

	in, err = r.Register(instance2)
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("register the instance success in:%#v", in)

	ins, err := r.Fetch(segment, serviceName)

	if err != nil {
		t.Fatal(err)
	}

	for _, instance := range ins {

		t.Logf("fetch instance:%#v", instance)
	}

}

func TestRegistry_Cancel(t *testing.T) {

	r := initRegistry()

	in, err := r.Register(instance1)
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("register the instance success in:%#v", in)

	in, err = r.Register(instance2)
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("register the instance success in:%#v", in)

	ins, err := r.Fetch(segment, serviceName)

	if err != nil {
		t.Fatal(err)
	}

	for _, instance := range ins {

		t.Logf("fetch instance:%#v", instance)
	}

	in, err = r.Cancel(segment, serviceName, "192.168.1.1", 8001)

	if err != nil {

		t.Fatal(err)
	}

	t.Logf("cancel the instance:%#v", in)

	ins, err = r.Fetch(segment, serviceName)

	if err != nil {
		t.Fatal(err)
	}

	for _, instance := range ins {

		t.Logf("fetch instance:%#v", instance)
	}

}

func TestRegistry_Renew(t *testing.T) {

	r := initRegistry()

	in, err := r.Register(instance1)
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("register the instance success in:%#v", in)

	in, err = r.Register(instance2)
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("register the instance success in:%#v", in)

	ins, err := r.Fetch(segment, serviceName)

	if err != nil {
		t.Fatal(err)
	}

	for _, instance := range ins {

		t.Logf("fetch instance:%#v", instance)
	}

	group := sync.WaitGroup{}
	group.Add(1)

	go func() {

		time.Sleep(time.Second * 2)
		in, err := r.Renew(segment, serviceName, "192.168.1.1", 8001)
		if err != nil {
			t.Error(err)
		}

		t.Logf("renew the instance:%#v,success", in)
		time.Sleep(time.Second)
		group.Done()
	}()

	group.Wait()
	ins, err = r.Fetch(segment, serviceName)

	if err != nil {
		t.Fatal(err)
	}

	for _, instance := range ins {

		t.Logf("fetch instance:%#v", instance)
	}
}
