package registry

import (
	"encoding/json"
	"github.com/busgo/elsa/pkg/proto/pb"
	"time"
)

type Instance struct {
	Segment         string            `json:"segment"`
	ServiceName     string            `json:"service_name"`
	Ip              string            `json:"ip"`
	Port            int32             `json:"port"`
	Metadata        map[string]string `json:"metadata"`
	RegTimestamp    int64             `json:"reg_timestamp"`
	UpTimestamp     int64             `json:"up_timestamp"`
	RenewTimestamp  int64             `json:"renew_timestamp"`
	DirtyTimestamp  int64             `json:"dirty_timestamp"`
	LatestTimestamp int64             `json:"latest_timestamp"`
}

// copy a new instance
func (in *Instance) Copy() *Instance {

	instance := new(Instance)
	*instance = *in
	return instance
}

// to string
func (in *Instance) String() string {

	content, err := json.Marshal(in)
	if err != nil {
		return ""
	}
	return string(content)
}

func NewInstance(req *pb.RegisterRequest) *Instance {

	now := time.Now().UnixNano()
	return &Instance{
		Segment:         req.Segment,
		ServiceName:     req.ServiceName,
		Ip:              req.Ip,
		Port:            req.Port,
		Metadata:        make(map[string]string),
		RegTimestamp:    now,
		UpTimestamp:     now,
		RenewTimestamp:  now,
		DirtyTimestamp:  now,
		LatestTimestamp: now,
	}
}

// new a service instance
func NewServiceInstance(instance *Instance) *pb.ServiceInstance {

	return &pb.ServiceInstance{
		Segment:         instance.Segment,
		ServiceName:     instance.ServiceName,
		Ip:              instance.Ip,
		Port:            instance.Port,
		Metadata:        instance.Metadata,
		RegTimestamp:    instance.RegTimestamp,
		UpTimestamp:     instance.UpTimestamp,
		RenewTimestamp:  instance.RegTimestamp,
		DirtyTimestamp:  instance.DirtyTimestamp,
		LatestTimestamp: instance.LatestTimestamp,
	}
}

// create renew request
func NewRenewRequest(segment, serviceName, ip string, port int32) *pb.RenewRequest {

	return &pb.RenewRequest{
		Segment:     segment,
		ServiceName: serviceName,
		Ip:          ip,
		Port:        port,
		SyncType:    pb.SyncTypeEnum_None,
	}
}

// create cancel request
func NewCancelRequest(segment, serviceName, ip string, port int32) *pb.CancelRequest {

	return &pb.CancelRequest{
		Segment:     segment,
		ServiceName: serviceName,
		Ip:          ip,
		Port:        port,
		SyncType:    pb.SyncTypeEnum_None,
	}
}
