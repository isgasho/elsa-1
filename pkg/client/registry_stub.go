package client

import (
	"context"
	"github.com/busgo/elsa/pkg/log"
	"github.com/busgo/elsa/pkg/proto/pb"
	"google.golang.org/grpc"
	"time"
)

type RegistryStub struct {
	segment   string
	endpoints []string
	cli       pb.RegistryServiceClient
}

// new a registry stub
func NewRegistryStub(segment string, endpoints []string) (*RegistryStub, error) {

	r := NewDirectResolverWithEndpoints(endpoints)
	cc, err := grpc.Dial(pb.RegistryService_ServiceDesc.ServiceName, grpc.WithInsecure(), grpc.WithResolvers(r))
	if err != nil {
		return nil, err
	}
	cli := pb.NewRegistryServiceClient(cc)
	return &RegistryStub{
		segment:   segment,
		endpoints: endpoints,
		cli:       cli,
	}, nil

}

func (r *RegistryStub) GetSegment() string {
	return r.segment
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

// register a service instance
func (r *RegistryStub) Register(ctx context.Context, serviceName, ip string, port int32) (bool, error) {

	response, err := r.cli.Register(ctx, &pb.RegisterRequest{
		Segment:         r.segment,
		ServiceName:     serviceName,
		Ip:              ip,
		Port:            port,
		Metadata:        make(map[string]string),
		RegTimestamp:    time.Now().UnixNano(),
		UpTimestamp:     time.Now().UnixNano(),
		RenewTimestamp:  time.Now().UnixNano(),
		DirtyTimestamp:  time.Now().UnixNano(),
		LatestTimestamp: time.Now().UnixNano(),
		SyncType:        pb.SyncTypeEnum_Yes,
	})

	if err != nil {
		log.Errorf("register the segment:%s,serviceName:%s,ip:%s,port:%32 fail:%s", r.segment, serviceName, ip, port, err.Error())
		return false, err
	}

	if response.Code != 0 {
		log.Warnf("register the segment:%s,serviceName:%s,ip:%s,port:%32 fail code:%d", r.segment, serviceName, ip, port, response.Code)
		return false, nil
	}

	return true, nil
}

// renew a service instance
func (r *RegistryStub) Renew(ctx context.Context, serviceName, ip string, port int32) (bool, error) {

	response, err := r.cli.Renew(ctx, &pb.RenewRequest{
		Segment:     r.segment,
		ServiceName: serviceName,
		Ip:          ip,
		Port:        port,
		SyncType:    pb.SyncTypeEnum_Yes,
	})

	if err != nil {
		log.Errorf("renew the segment:%s,serviceName:%s,ip:%s,port:%d fail:%s", r.segment, serviceName, ip, port, err.Error())
		return false, err
	}

	if response.Code != 0 {
		log.Warnf("renew the segment:%s,serviceName:%s,ip:%s,port:%d code:%d", r.segment, serviceName, ip, port, response.Code)
		return false, nil
	}
	return true, nil
}

// cancel a service instance
func (r *RegistryStub) Cancel(ctx context.Context, serviceName, ip string, port int32) (bool, error) {
	response, err := r.cli.Cancel(ctx, &pb.CancelRequest{
		Segment:     r.segment,
		ServiceName: serviceName,
		Ip:          ip,
		Port:        port,
		SyncType:    pb.SyncTypeEnum_Yes,
	})

	if err != nil {
		log.Errorf("cancel the [segment:%s,serviceName:%s,ip:%s,port:%s] fail:%s", r.segment, serviceName, ip, port, err.Error())
		return false, err
	}

	if response.Code != 0 {
		log.Warnf("cancel the [segment:%s,serviceName:%s,ip:%s,port:%s code:%d] fail", r.segment, serviceName, ip, port, response.Code)
		return false, nil
	}
	return true, nil
}
