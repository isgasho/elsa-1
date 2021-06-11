package client

import (
	"errors"
	"fmt"
	"github.com/busgo/elsa/pkg/log"
	"google.golang.org/grpc"
	"net"
)

type ElsaServer struct {
	managedSentinel *ManagedSentinel
	opts            ServerOptions
	server          *grpc.Server
	state           bool
}

type InitMethod func(server *grpc.Server) (serverNames []string)

type ServerOptions struct {
	name         string
	segment      string
	serverPort   int32
	registryStub *RegistryStub
}

type ServerOption func(options *ServerOptions)

func WithServerPort(serverPort int32) ServerOption {
	return func(options *ServerOptions) {
		options.serverPort = serverPort
	}
}

func WithRegistryStub(stub *RegistryStub) ServerOption {
	return func(options *ServerOptions) {
		options.registryStub = stub
	}
}

func WithName(name string) ServerOption {
	return func(options *ServerOptions) {
		options.name = name
	}
}

// new elsa server
func NewElsaServer(options ...ServerOption) (*ElsaServer, error) {

	opts := ServerOptions{
		name:         DefaultServerName,
		segment:      DefaultSegment,
		serverPort:   DefaultServerPort,
		registryStub: nil,
	}
	for _, opt := range options {
		opt(&opts)
	}

	if opts.registryStub == nil {
		stub, err := NewRegistryStub(opts.segment, []string{DefaultRegistryEndpoint})
		if err != nil {
			return nil, err
		}
		opts.registryStub = stub
	}
	return &ElsaServer{
		managedSentinel: NewManagedSentinel(opts.serverPort, opts.registryStub),
		server:          grpc.NewServer(),
		opts:            opts,
		state:           false,
	}, nil
}

func (s *ElsaServer) Init(m InitMethod) {
	serviceNames := m(s.server)
	s.state = true
	for _, serviceName := range serviceNames {
		s.managedSentinel.PushService(serviceName)
	}

	log.Infof("the %s server initialize success", s.opts.name)
}

// start the elsa server
func (s *ElsaServer) Start() error {
	if !s.state {
		log.Error(fmt.Sprintf("the %s server  has not initialize", s.opts.name))
		return errors.New(fmt.Sprintf("the %s server  has not initialize", s.opts.name))
	}

	l, err := net.Listen("tcp", fmt.Sprintf(":%d", s.opts.serverPort))
	if err != nil {
		return err
	}
	log.Infof("the %s server has start...", s.opts.name)
	if err = s.server.Serve(l); err != nil {
		return err
	}
	return nil
}
