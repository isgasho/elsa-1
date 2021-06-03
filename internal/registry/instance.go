package registry

type Instance struct {
	Segment         string
	ServiceName     string
	Ip              string
	Port            int32
	Metadata        map[string]string
	RegTimestamp    int64
	UpTimestamp     int64
	RenewTimestamp  int64
	DirtyTimestamp  int64
	LatestTimestamp int64
}

// copy a new instance
func (in *Instance) Copy() *Instance {

	instance := new(Instance)
	*instance = *in
	return instance
}
