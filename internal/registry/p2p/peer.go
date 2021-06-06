package p2p

import (
	"context"
	"github.com/busgo/elsa/pkg/log"
	"github.com/busgo/elsa/pkg/proto/pb"
	"google.golang.org/grpc"
	"time"
)

type SyncMsgType int32

const (
	SyncMsgRegType = iota
	SyncMsgRenewType
	SyncMsgCancelType
)

// peer pool
type PeerPool struct {
	endpoints   []string
	peers       []*Peer
	state       bool
	syncMsgChan chan *SyncMsg
}

// sync message
type SyncMsg struct {
	Type    SyncMsgType
	Content interface{}
}

// peer
type Peer struct {
	endpoint string
	cli      pb.RegistryServiceClient
}

// new a peer pool with endpoints
func NewPeerPoolWithEndpoints(endpoints []string) (*PeerPool, error) {

	peers := make([]*Peer, 0)

	for _, endpoint := range endpoints {
		p, err := NewPeerEndpoint(endpoint)
		if err != nil {
			return nil, err
		}
		log.Debugf("the peer pool add peer:%s success", endpoint)
		peers = append(peers, p)
	}

	return &PeerPool{
		endpoints:   endpoints,
		peers:       peers,
		state:       false,
		syncMsgChan: make(chan *SyncMsg, 128),
	}, nil
}

// push a sync message
func (pool *PeerPool) PushMsg(msg *SyncMsg) {
	pool.syncMsgChan <- msg
}

// start the peer pool
func (pool *PeerPool) Start() {
	if pool.state {
		log.Warnf("the peer pool %s has start", pool.endpoints)
		return
	}
	pool.state = true
	go pool.lookup()
}

// lookup the sync message
func (pool *PeerPool) lookup() {

	log.Debugf("the peer pool:%s has start...", pool.endpoints)
	for {
		select {
		case msg, ok := <-pool.syncMsgChan:
			if !ok {
				log.Warnf("the peer pool [%s] has closed...", pool.endpoints)
				return
			}
			// handle the sync message
			pool.handleSyncMsg(msg)
		}
	}
}

// handle the sync message
func (pool *PeerPool) handleSyncMsg(msg *SyncMsg) {
	if msg == nil || msg.Content == nil {
		log.Warn("the sync message is nil")
		return
	}
	if len(pool.peers) == 0 {
		log.Warn("the peer pool has no peer")
		return
	}

	switch msg.Type {
	case SyncMsgRegType: // reg
		pool.handleRegMsg(msg)
	case SyncMsgRenewType: // renew
		pool.handleRenewMsg(msg)
	case SyncMsgCancelType: // cancel
		pool.handleCancelMsg(msg)
	}

}

// handle the reg message
func (pool *PeerPool) handleRegMsg(msg *SyncMsg) {

	log.Debugf("handle the reg message %#v", msg)
	req := msg.Content.(*pb.RegisterRequest)
	for _, peer := range pool.peers {

		ctx, _ := context.WithTimeout(context.Background(), time.Millisecond*500)
		_, err := peer.cli.Register(ctx, req)
		if err != nil {
			log.Warnf("the peer:%s register sync instance fail:%s", peer.endpoint, err.Error())
			continue
		}

		log.Debugf("the peer:%s register sync instance success", peer.endpoint)
	}

}

// handle the renew msg
func (pool *PeerPool) handleRenewMsg(msg *SyncMsg) {

	log.Debugf("handle the renew message %#v", msg)
	req := msg.Content.(*pb.RenewRequest)

	for _, peer := range pool.peers {

		ctx, _ := context.WithTimeout(context.Background(), time.Millisecond*500)
		_, err := peer.cli.Renew(ctx, req)
		if err != nil {
			log.Warnf("the peer:%s renew sync instance fail:%s", peer.endpoint, err.Error())
			continue
		}

		log.Debugf("the peer:%s renew sync instance success", peer.endpoint)
	}
}

// handle the cancel msg
func (pool *PeerPool) handleCancelMsg(msg *SyncMsg) {

	log.Debugf("handle the cancel message %#v", msg)
	req := msg.Content.(*pb.CancelRequest)

	for _, peer := range pool.peers {

		ctx, _ := context.WithTimeout(context.Background(), time.Millisecond*500)
		_, err := peer.cli.Cancel(ctx, req)
		if err != nil {
			log.Warnf("the peer:%s cancel sync instance fail:%s", peer.endpoint, err.Error())
			continue
		}
		log.Debugf("the peer:%s cancel sync instance success", peer.endpoint)
	}
}

// new a peer with endpoint
func NewPeerEndpoint(endpoint string) (*Peer, error) {

	cc, err := grpc.Dial(endpoint, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	return &Peer{
		endpoint: endpoint,
		cli:      pb.NewRegistryServiceClient(cc),
	}, nil
}
