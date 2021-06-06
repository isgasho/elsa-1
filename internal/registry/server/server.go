package server

import (
	"context"
	"github.com/busgo/elsa/internal/registry"
	"github.com/busgo/elsa/internal/registry/p2p"
	"github.com/busgo/elsa/pkg/log"
	"github.com/busgo/elsa/pkg/proto/pb"
	"github.com/busgo/elsa/pkg/utils"
	"google.golang.org/grpc"
	"net"
	"strings"
)

type RegistryServer struct {
	endpoint string
	r        registry.Registry
	pool     *p2p.PeerPool
	server   *grpc.Server
	pb.UnimplementedRegistryServiceServer
}

// new  registry server
func NewRegistryServerWithEndpoints(endpoints []string) (*RegistryServer, error) {
	pool, err := p2p.NewPeerPoolWithEndpoints(endpoints)
	if err != nil {
		return nil, err
	}
	return &RegistryServer{
		endpoint: getLocalEndpoint(endpoints),
		r:        registry.NewRegistry(),
		pool:     pool,
		server:   grpc.NewServer(),
	}, nil
}

// start registry server
func (s *RegistryServer) Start() error {

	l, err := net.Listen("tcp", s.endpoint)
	if err != nil {
		return err
	}
	pb.RegisterRegistryServiceServer(s.server, s)
	s.pool.Start()
	log.Infof("start the registry server endpoint:%s success", s.endpoint)
	if err = s.server.Serve(l); err != nil {
		return err
	}
	return nil
}

// get local endpoint
func getLocalEndpoint(endpoints []string) string {

	if len(endpoints) == 0 {
		return p2p.DefaultEndpoint
	}
	ip := utils.GetLocalIp()
	for _, endpoint := range endpoints {
		if strings.HasPrefix(endpoint, ip) {
			return endpoint
		}
	}
	return p2p.DefaultEndpoint
}

// register a service instance
func (s *RegistryServer) Register(ctx context.Context, request *pb.RegisterRequest) (*pb.RegisterResponse, error) {

	instance := registry.NewInstance(request)
	in, _ := s.r.Register(instance)

	return &pb.RegisterResponse{
		Code:     0,
		Message:  "",
		Instance: registry.NewServiceInstance(in),
	}, nil
}

// renew a service instance
func (s *RegistryServer) Renew(ctx context.Context, request *pb.RenewRequest) (*pb.RenewResponse, error) {

	in, err := s.r.Renew(request.Segment, request.ServiceName, request.Ip, request.Port)
	if err != nil {
		e := err.(registry.RegistryError)
		return &pb.RenewResponse{
			Code:     e.Code,
			Message:  e.Message,
			Instance: nil,
		}, nil
	}

	// sync other peer
	if request.SyncType == pb.SyncTypeEnum_Yes {
		s.pool.PushMsg(&p2p.SyncMsg{
			Type:    p2p.SyncMsgRenewType,
			Content: registry.NewRenewRequest(in.Segment, in.ServiceName, in.Ip, in.Port),
		})
	}

	return &pb.RenewResponse{
		Code:     0,
		Message:  "",
		Instance: registry.NewServiceInstance(in),
	}, nil
}

// cancel a service instance
func (s *RegistryServer) Cancel(ctx context.Context, request *pb.CancelRequest) (*pb.CancelResponse, error) {
	in, err := s.r.Renew(request.Segment, request.ServiceName, request.Ip, request.Port)
	if err != nil {
		e := err.(registry.RegistryError)
		return &pb.CancelResponse{
			Code:     e.Code,
			Message:  e.Message,
			Instance: nil,
		}, nil
	}

	// sync other peer
	if request.SyncType == pb.SyncTypeEnum_Yes {
		s.pool.PushMsg(&p2p.SyncMsg{
			Type:    p2p.SyncMsgCancelType,
			Content: registry.NewCancelRequest(in.Segment, in.ServiceName, in.Ip, in.Port),
		})
	}

	return &pb.CancelResponse{
		Code:     0,
		Message:  "",
		Instance: registry.NewServiceInstance(in),
	}, nil
}
