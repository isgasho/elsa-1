package main

import (
	"flag"
	"fmt"
	"github.com/busgo/elsa/internal/registry/server"
	"github.com/busgo/elsa/pkg/log"
	"os"
	"strings"
)

const (
	defaultRegistryServerEndpoint = "127.0.0.1:8005"
	defaultVersion                = "1.0"
	defaultLogFile                = "."
)

func main() {

	serverEndpoints := flag.String("registry_server_endpoints", defaultRegistryServerEndpoint, "the registry server endpoints,if multi server endpoint please use '.' split")
	version := flag.String("version", "", "print elsa micro service framework version")
	v := flag.String("v", "", "print elsa micro service framework version")
	flag.String("logfile", defaultLogFile, "set log file path")
	flag.Parse()
	if *version != "" || *v != "" {
		fmt.Printf("elsa micro service framework %s", defaultVersion)
		os.Exit(0)
	}

	endpoints := strings.Split(*serverEndpoints, ",")
	s, err := server.NewRegistryServerWithEndpoints(endpoints)

	if err != nil {
		log.Error("create registry server fail:%#v", err)
		panic(err)
	}

	if err = s.Start(); err != nil {
		log.Errorf("start registry server fail:%#v", err)
		panic(err)
	}

}
