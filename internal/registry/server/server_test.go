package server

import (
	"github.com/busgo/elsa/pkg/log"
	"testing"
)

var endpoints = []string{"127.0.0.1:8005"}

func TestNewRegistryServerWithEndpoints(t *testing.T) {

	s, err := NewRegistryServerWithEndpoints(endpoints)
	if err != nil {
		t.Fatal(err)
	}

	log.Infof("new registry server success %#v", s)
}

func TestRegistryServer_Start(t *testing.T) {

	s, err := NewRegistryServerWithEndpoints(endpoints)
	if err != nil {
		t.Fatal(err)
	}

	log.Infof("new registry server success %#v", s)
}
