package registry

import "encoding/json"

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
