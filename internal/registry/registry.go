package registry

import (
	"errors"
	"fmt"
	"sync"

	"github.com/busgo/elsa/pkg/log"
)

type Registry interface {

	// register a instance
	Register(instance *Instance) (*Instance, error)

	// fetch with segment and service name
	Fetch(segment, serviceName string) ([]*Instance, error)

	// cancel the instance
	Cancel(segment, serviceName, ip string, port int32) (*Instance, error)

	// renew the instance
	Renew(segment, serviceName, ip string, port int32) (*Instance, error)
}

type registry struct {
	apps map[string]*Application
	sync.RWMutex
}

func NewRegistry() Registry {
	return &registry{
		apps:    make(map[string]*Application),
		RWMutex: sync.RWMutex{},
	}
}

// register a instance
func (r *registry) Register(instance *Instance) (*Instance, error) {
	log.Infof("start register action instance%#v", instance)

	segment := instance.Segment
	serviceName := instance.ServiceName
	app, ok := r.getApplication(segment, serviceName)
	if !ok {
		app = NewApplication(segment, serviceName)
	}

	in, _ := app.addInstance(instance)
	if !ok {
		r.Lock()
		r.apps[fmt.Sprintf("%s-%s", segment, serviceName)] = app
		r.Unlock()
	}
	return in, nil
}

func (r *registry) Fetch(segment, serviceName string) ([]*Instance, error) {
	app, ok := r.getApplication(segment, serviceName)
	if !ok {
		log.Warnf("the application not found segment:%s,serviceName:%s", segment, serviceName)
		return make([]*Instance, 0), nil
	}
	return app.getInstances(), nil
}

// cancel the instance
func (r *registry) Cancel(segment, serviceName, ip string, port int32) (*Instance, error) {
	app, ok := r.getApplication(segment, serviceName)
	if !ok {
		log.Warnf("the application not found segment:%s,serviceName:%s", segment, serviceName)
		return nil, errors.New("app not found")
	}

	in, err := app.cancel(ip, port)

	if err != nil {
		return in, err
	}

	if app.getInstanceSize() == 0 { // must delete the application
		r.Lock()
		delete(r.apps, fmt.Sprintf("%s-%s", segment, serviceName))
		r.Unlock()
	}
	return in, nil
}

// renew the instance
func (r *registry) Renew(segment, serviceName, ip string, port int32) (*Instance, error) {
	app, ok := r.getApplication(segment, serviceName)
	if !ok {
		log.Warnf("the application not found segment:%s,serviceName:%s", segment, serviceName)
		return nil, errors.New("app not found")
	}

	in, err := app.renew(ip, port)
	return in, err

}

// get app with segment and service name
func (r *registry) getApplication(segment, serviceName string) (*Application, bool) {
	r.RLock()
	defer r.RUnlock()
	app, ok := r.apps[fmt.Sprintf("%s-%s", segment, serviceName)]
	return app, ok
}
