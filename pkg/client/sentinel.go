package client

import (
	"context"
	"github.com/busgo/elsa/pkg/log"
	"github.com/busgo/elsa/pkg/utils"
	"sync"
	"time"
)

const (
	RetryTimeDuration = time.Second * 3        // retry time duration
	RenewTimeDuration = time.Second * 30       // renew time duration
	TimeoutDuration   = time.Millisecond * 500 // renew time duration
)

type ManagedSentinel struct {
	registryStub *RegistryStub
	sentinels    map[string]*Sentinel
	ip           string
	port         int32
	sync.RWMutex
}

type Sentinel struct {
	serviceName    string
	ip             string
	port           int32
	registryStub   *RegistryStub
	registerChan   chan bool
	retryRenewChan chan bool
	closedChan     chan bool
	sync.RWMutex
}

func NewManagedSentinel(serverPort int32, registryStub *RegistryStub) *ManagedSentinel {

	return &ManagedSentinel{
		registryStub: registryStub,
		sentinels:    make(map[string]*Sentinel),
		ip:           utils.GetLocalIp(),
		port:         serverPort,
		RWMutex:      sync.RWMutex{},
	}
}

func (m *ManagedSentinel) PushService(serviceName string) {
	m.Lock()
	defer m.Unlock()
	sentinel := m.sentinels[serviceName]
	if sentinel != nil {
		return
	}

	sentinel = newSentinel(serviceName, m.ip, m.port, m.registryStub)
	m.sentinels[serviceName] = sentinel
	sentinel.register()

	go sentinel.lookup()
	return
}

//  new sentinel
func newSentinel(serviceName, ip string, port int32, registryStub *RegistryStub) *Sentinel {
	return &Sentinel{
		serviceName:    serviceName,
		ip:             ip,
		port:           port,
		registryStub:   registryStub,
		registerChan:   make(chan bool, 10),
		retryRenewChan: make(chan bool, 10),
		closedChan:     make(chan bool),
		RWMutex:        sync.RWMutex{},
	}
}

func (s *Sentinel) lookup() {
	renewTicker := time.Tick(RenewTimeDuration)

	for {
		select {
		case <-renewTicker: //renew
			log.Infof("start renew serviceName:%s,ip:%s,port:%d ...", s.serviceName, s.ip, s.port)
			s.renew()
		case <-s.retryRenewChan:
			time.Sleep(RetryTimeDuration)
			s.renew()
		case <-s.registerChan:
			time.Sleep(RetryTimeDuration)
			s.register()
		case <-s.closedChan:
			log.Warnf("the sentinel has closed serviceName:%s", s.serviceName)
			return

		}
	}

}

//
func (s *Sentinel) register() {
	s.Lock()
	defer s.Unlock()
	ctx, _ := context.WithTimeout(context.Background(), TimeoutDuration)
	state, err := s.registryStub.Register(ctx, s.serviceName, s.ip, s.port)
	if err != nil || !state {
		log.Warnf("register serviceName:%s,ip:%s,port:%d fail,after try again...", s.serviceName, s.ip, s.port)
		s.registerChan <- true
		return
	}
	log.Infof("register serviceName:%s,ip:%s,port:%d success...", s.serviceName, s.ip, s.port)
}

// renew
func (s *Sentinel) renew() {
	s.Lock()
	defer s.Unlock()
	ctx, _ := context.WithTimeout(context.Background(), TimeoutDuration)
	state, err := s.registryStub.Renew(ctx, s.serviceName, s.ip, s.port)
	if err != nil {
		log.Errorf("renew serviceName:%s,ip:%s,port:%d fail,error:%s", s.serviceName, s.ip, s.port, err.Error())
		s.retryRenewChan <- true
		return
	}

	if !state {
		log.Warnf("renew serviceName:%s,ip:%s,port:%d fail", s.serviceName, s.ip, s.port)
		s.registerChan <- true
		return
	}
	log.Infof("renew serviceName:%s,ip:%s,port:%d success", s.serviceName, s.ip, s.port)
}

// cancel
func (s *Sentinel) cancel() {
	ctx, _ := context.WithTimeout(context.Background(), TimeoutDuration)
	state, err := s.registryStub.Cancel(ctx, s.serviceName, s.ip, s.port)
	if err != nil || !state {
		log.Warnf("cancel serviceName:%s,ip:%s,port:%d fail", s.serviceName, s.ip, s.port)
	}
	log.Warnf("cancel serviceName:%s,ip:%s,port:%d success", s.serviceName, s.ip, s.port)
	s.closedChan <- true
}
