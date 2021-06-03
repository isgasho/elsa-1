package registry

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

type Application struct {
	segment     string
	serviceName string
	instances   map[string]*Instance
	sync.RWMutex
}

// new a  application
func NewApplication(segment, serviceName string) *Application {

	return &Application{
		segment:     segment,
		serviceName: serviceName,
		instances:   make(map[string]*Instance),
		RWMutex:     sync.RWMutex{},
	}
}

// add a new instance
func (app *Application) addInstance(instance *Instance) (*Instance, bool) {

	app.Lock()
	defer app.Unlock()

	ip := instance.Ip
	port := instance.Port

	in, ok := app.instances[fmt.Sprintf("%s-%d", ip, port)]
	if ok {

		in.UpTimestamp = instance.UpTimestamp

		if in.DirtyTimestamp > instance.DirtyTimestamp {
			instance = in
		}
	}

	app.instances[fmt.Sprintf("%s-%d", ip, port)] = instance

	return instance.Copy(), !ok

}

// cancel the instance from application
func (app *Application) cancel(ip string, port int32) (*Instance, error) {
	app.Lock()
	defer app.Unlock()
	in, ok := app.instances[fmt.Sprintf("%s-%d", ip, port)]
	if !ok {
		return nil, errors.New("instance not found")
	}

	instance := in.Copy()

	delete(app.instances, fmt.Sprintf("%s-%d", ip, port))
	return instance, nil
}

// get instance size
func (app *Application) getInstanceSize() int {

	app.RLock()
	defer app.RUnlock()
	return len(app.instances)
}

// renew the instance
func (app *Application) renew(ip string, port int32) (*Instance, error) {
	app.Lock()
	defer app.Unlock()
	instance, ok := app.instances[fmt.Sprintf("%s-%d", ip, port)]
	if !ok {
		return nil, errors.New("instance not found...")
	}
	instance.RenewTimestamp = time.Now().UnixNano()
	return instance.Copy(), nil

}

// get instances from app
func (app *Application) getInstances() []*Instance {
	app.RLock()
	defer app.RUnlock()

	if len(app.instances) == 0 {
		return make([]*Instance, 0)
	}

	instances := make([]*Instance, 0)
	for _, in := range app.instances {
		instances = append(instances, in.Copy())
	}
	return instances
}
