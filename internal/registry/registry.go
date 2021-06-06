package registry

import (
	"fmt"
	"math/rand"
	"sync"
	"time"

	"github.com/busgo/elsa/internal/registry/census"

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
	c *census.Census
}

func NewRegistry() Registry {
	r := &registry{
		apps:    make(map[string]*Application),
		c:       new(census.Census),
		RWMutex: sync.RWMutex{},
	}

	go r.lookup()

	return r
}

// register a instance
func (r *registry) Register(instance *Instance) (*Instance, error) {
	log.Infof("start register action instance%#v", instance.String())

	segment := instance.Segment
	serviceName := instance.ServiceName
	app, ok := r.getApplication(segment, serviceName)
	if !ok {
		app = NewApplication(segment, serviceName)
	}

	in, created := app.addInstance(instance)
	if created {
		r.c.IncrNeedCount()
	}
	if !ok {
		r.Lock()
		r.apps[fmt.Sprintf("%s-%s", segment, serviceName)] = app
		r.Unlock()
	}
	return in, nil
}

// fetch instances from service
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
		return nil, ApplicationNotFoundError
	}

	in, err := app.cancel(ip, port)

	if err != nil {
		return in, err
	}

	if in != nil {
		r.c.DecrNeedCount()
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
		return nil, ApplicationNotFoundError
	}

	in, err := app.renew(ip, port)
	if err != nil {
		r.c.IncrCount()
	}
	return in, err

}

// get app with segment and service name
func (r *registry) getApplication(segment, serviceName string) (*Application, bool) {
	r.RLock()
	defer r.RUnlock()
	app, ok := r.apps[fmt.Sprintf("%s-%s", segment, serviceName)]
	return app, ok
}

// get all application
func (r *registry) getApplications() []*Application {
	r.RLock()
	defer r.RUnlock()

	if len(r.apps) == 0 {
		return make([]*Application, 0)
	}

	apps := make([]*Application, 0)

	for _, app := range r.apps {

		apps = append(apps, app)

	}
	return apps
}

//------------------------------------evict expired instance task----------------------------------------------------------//
func (r *registry) lookup() {
	evictTicker := time.Tick(census.ScanEvictDuration)
	seekNeedCountTicker := time.Tick(census.ResetRenewNeedCountDuration)
	log.Debug("the registry evict task has start...")
	for {

		select {
		case <-evictTicker:
			r.c.ResetCount()
			r.evict()
		case <-seekNeedCountTicker: //
			r.seekNeedCount()
		}

	}
}

// seek renew need count
func (r *registry) seekNeedCount() {

	var count int64
	apps := r.getApplications()
	for _, app := range apps {
		count += int64(len(app.instances))
	}
	log.Info("start seek need count:%d", count)
	r.c.SeekNeedCount(count)
}

// evict expired instance
func (r *registry) evict() {

	log.Debug("start evict expired task...")
	apps := r.getApplications()
	if len(apps) == 0 {
		log.Warn("the registry apps is nil")
		return
	}

	now := time.Now().UnixNano()
	var instancesSize int64
	expiredInstances := make([]*Instance, 0)
	for _, app := range apps {

		instances := app.getInstances()
		if len(instances) == 0 {
			continue
		}
		instancesSize += int64(len(instances))

		for _, in := range instances {

			delta := now - in.RenewTimestamp
			//  check  expired status
			if (delta > int64(census.InstanceEvictExpiredDuration) && !r.c.ProtectedStatus()) ||
				delta > int64(census.InstanceMaxExpiredDuration) {
				expiredInstances = append(expiredInstances, in)
			}

		}

	}

	// check expire limit
	expiredInstanceLimit := instancesSize - int64(float64(instancesSize)*census.SelfProtectedThreshold)
	expiredInstanceSize := len(expiredInstances)
	if expiredInstanceLimit < int64(expiredInstanceSize) {
		expiredInstanceSize = int(expiredInstanceLimit)
	}

	if expiredInstanceSize <= 0 {
		log.Warn("has no expired instance to evict")
		return
	}

	for i := 0; i < expiredInstanceSize; i++ {
		j := i + rand.Intn(len(expiredInstances)-i)
		expiredInstances[i], expiredInstances[j] = expiredInstances[j], expiredInstances[i]
		expireInstance := expiredInstances[i]
		log.Infof("start evict  the expired instance segment:%s,serviceName:%s,ip:%s,port:%d success", expireInstance.Segment, expireInstance.ServiceName, expireInstance.Ip, expireInstance.Port)
		_, err := r.Cancel(expireInstance.Segment, expireInstance.ServiceName, expireInstance.Ip, expireInstance.Port)
		if err != nil {
			log.Warnf("cancel the expired instance segment:%s,serviceName:%s,ip:%s,port:%d fail:%s", expireInstance.Segment, expireInstance.ServiceName, expireInstance.Ip, expireInstance.Port, err.Error())
			continue
		}

		log.Infof("cancel the expired instance segment:%s,serviceName:%s,ip:%s,port:%d success", expireInstance.Segment, expireInstance.ServiceName, expireInstance.Ip, expireInstance.Port)
	}

}
