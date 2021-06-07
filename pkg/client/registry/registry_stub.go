package registry

import (
	"context"
	"github.com/busgo/elsa/pkg/client/resolver"
	"github.com/busgo/elsa/pkg/log"
	"github.com/busgo/elsa/pkg/proto/pb"
	"google.golang.org/grpc"
)

type RegistryStub struct {
	endpoints []string
	cli       pb.RegistryServiceClient
}

// new a registry stub
func NewRegistryStub(endpoints []string) (*RegistryStub, error) {

	r := resolver.NewDirectResolverWithEndpoints(endpoints)
	cc, err := grpc.Dial(pb.RegistryService_ServiceDesc.ServiceName, grpc.WithInsecure(), grpc.WithResolvers(r))
	if err != nil {
		return nil, err
	}
	cli := pb.NewRegistryServiceClient(cc)
	return &RegistryStub{
		endpoints: endpoints,
		cli:       cli,
	}, nil

}

// fetch service instance list
func (r *RegistryStub) Fetch(cxt context.Context, segment, serviceName string) ([]*pb.ServiceInstance, error) {
	response, err := r.cli.Fetch(cxt, &pb.FetchRequest{
		Segment:     segment,
		ServiceName: serviceName,
	})
	if err != nil {
		log.Errorf("fetch the segment:%s,serviceName:%s fail:%s", segment, serviceName, err.Error())
		return make([]*pb.ServiceInstance, 0), err
	}

	if response.Code != 0 {
		log.Errorf("fetch the segment:%s,serviceName:%s the the service name not found", segment, serviceName)
		return make([]*pb.ServiceInstance, 0), err
	}
	return response.Instances, nil
}
