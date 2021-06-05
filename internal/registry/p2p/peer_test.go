package p2p

import "testing"

var endpoints = []string{"127.0.0.1:8005", "192.168.1.1:8005"}

// test new peer pool
func TestNewPeerPoolWithEndpoints(t *testing.T) {

	pool, err := NewPeerPoolWithEndpoints(endpoints)
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("the peer pool is %#v", pool)
}
