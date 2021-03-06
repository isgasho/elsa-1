package client

import (
	"context"
	"fmt"
	"github.com/busgo/elsa/pkg/log"
	"google.golang.org/grpc/resolver"
	"sync"
	"time"
)

const (
	DirectScheme = "direct"
	ElsaScheme   = "elsa"
)

func BuildTarget(segment, serviceName string) string {
	return fmt.Sprintf("%s:///%s", segment, serviceName)
}

type DirectResolver struct {
	endpoints []string
}

func NewDirectResolverWithEndpoints(endpoints []string) *DirectResolver {
	return &DirectResolver{endpoints: endpoints}
}

// Build creates a new resolver for the given target.
//
// gRPC dial calls Build synchronously, and fails if the returned error is
// not nil.
func (r *DirectResolver) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOptions) (resolver.Resolver, error) {

	addresses := make([]resolver.Address, 0)

	for _, endpoint := range r.endpoints {
		addresses = append(addresses, resolver.Address{
			Addr: endpoint,
		})
	}
	err := cc.UpdateState(resolver.State{
		Addresses: addresses,
	})
	return r, err
}

// Scheme returns the scheme supported by this resolver.
// Scheme is defined at https://github.com/grpc/grpc/blob/master/doc/naming.md.
func (r *DirectResolver) Scheme() string {
	return DirectScheme
}

// ResolveNow will be called by gRPC to try to resolve the target name
// again. It's just a hint, resolver can ignore this if it's not necessary.
//
// It could be called multiple times concurrently.
func (r *DirectResolver) ResolveNow(opts resolver.ResolveNowOptions) {

	log.Debugf("resolve now...")
}

// Close closes the resolver.
func (r *DirectResolver) Close() {
}

type ElsaResolverBuilder struct {
	resolvers    map[string]*ElsaResolver
	registryStub *RegistryStub
	sync.RWMutex
}

type ElsaResolver struct {
	segment         string
	serviceName     string
	cc              resolver.ClientConn
	registryStub    *RegistryStub
	closedChan      chan bool
	retryChan       chan bool
	latestTimestamp int64
	sync.RWMutex
}

// new a elsa resolver
func NewElsaResolverBuilder(stub *RegistryStub) *ElsaResolverBuilder {
	return &ElsaResolverBuilder{resolvers: make(map[string]*ElsaResolver), RWMutex: sync.RWMutex{}, registryStub: stub}
}

// Build creates a new resolver for the given target.
//
// gRPC dial calls Build synchronously, and fails if the returned error is
// not nil.
func (r *ElsaResolverBuilder) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOptions) (resolver.Resolver, error) {

	r.Lock()
	defer r.Unlock()
	elsaResolver := r.resolvers[target.Endpoint]
	if elsaResolver == nil {
		elsaResolver = NewElsaResolver(target.Endpoint, cc, r.registryStub)
	}
	r.resolvers[target.Endpoint] = elsaResolver
	go elsaResolver.lookup()
	// refresh
	elsaResolver.refresh()
	return elsaResolver, nil
}

// Scheme returns the scheme supported by this resolver.
// Scheme is defined at https://github.com/grpc/grpc/blob/master/doc/naming.md.
func (r *ElsaResolverBuilder) Scheme() string {
	return ElsaScheme
}

// new a elsa resolver
func NewElsaResolver(serviceName string, cli resolver.ClientConn, registryStub *RegistryStub) *ElsaResolver {

	return &ElsaResolver{
		serviceName:  serviceName,
		cc:           cli,
		registryStub: registryStub,
		closedChan:   make(chan bool),
		retryChan:    make(chan bool, 1),
	}
}

// ResolveNow will be called by gRPC to try to resolve the target name
// again. It's just a hint, resolver can ignore this if it's not necessary.
//
// It could be called multiple times concurrently.
func (r *ElsaResolver) ResolveNow(opts resolver.ResolveNowOptions) {
	log.Debugf("resolver now....")
	if !r.checkState() {
		log.Debugf("resolver now not allow refresh....")
		return
	}
	r.refresh()
}

// check refresh state
func (r *ElsaResolver) checkState() bool {

	r.RLock()
	defer r.RUnlock()
	now := time.Now().UnixNano()
	delta := now - r.latestTimestamp
	return delta >= int64(time.Second*3)
}

// Close closes the resolver.
func (r *ElsaResolver) Close() {

	r.closedChan <- true
	log.Infof("the elsa resolver has closed...")
}

func (r *ElsaResolver) lookup() {

	refreshTicker := time.Tick(time.Minute * 5)
	for {
		select {
		case <-refreshTicker:
			r.refresh() // refresh the service instance list
		case <-r.closedChan:
			log.Warn("the elsa resolver has stop...")
			return
		case <-r.retryChan:
			time.Sleep(time.Second * 3)
			r.refresh()
		}
	}
}

// refresh the service instance list
func (r *ElsaResolver) refresh() {

	r.latestTimestamp = time.Now().UnixNano()
	ctx, _ := context.WithTimeout(context.Background(), time.Millisecond*500)
	instances, err := r.registryStub.Fetch(ctx, r.serviceName)
	if err != nil {
		log.Warnf("fetch the service name:%s fail:%s", r.serviceName, err.Error())
		r.retryChan <- true
		return
	}
	r.Lock()
	addresses := make([]resolver.Address, 0)

	for _, instance := range instances {
		addresses = append(addresses, resolver.Address{
			Addr: fmt.Sprintf("%s:%d", instance.Ip, instance.Port),
		})
	}

	err = r.cc.UpdateState(resolver.State{
		Addresses: addresses,
	})
	if err != nil {
		log.Warnf("the elsa resolver segment:%s,serviceName:%s refresh addresses fail:%s", r.segment, r.serviceName, err.Error())
	} else {
		log.Infof("the elsa resolver segment:%s,serviceName:%s refresh addresses success", r.segment, r.serviceName)
	}

	// not found a service instance  must try again
	if len(addresses) == 0 {
		log.Warnf("fetch the service name:%s not found", r.serviceName)
		r.retryChan <- true
	}

	r.Unlock()
}
