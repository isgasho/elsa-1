package main

import "github.com/busgo/elsa/internal/registry/server"

func main() {

	endpoints := []string{"127.0.0.1:8005"}
	s, err := server.NewRegistryServerWithEndpoints(endpoints)

	if err != nil {
		panic(err)
	}

	if err = s.Start(); err != nil {
		panic(err)
	}

}
