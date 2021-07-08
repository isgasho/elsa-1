package etcd

import (
	"context"
	"github.com/busgo/elsa/pkg/log"
	"go.etcd.io/etcd/api/v3/mvccpb"
	"go.etcd.io/etcd/client/v3"
)

type EventType int32

const (
	CreateKeyChangeEvent EventType = 1
	UpdateKeyChangeEvent EventType = 2
	DeleteKeyChangeEvent EventType = 3
)

type Cli struct {
	kv    clientv3.KV
	lease clientv3.Lease
	c     *clientv3.Client
}

type KeyChangeEvent struct {
	Type  EventType
	Key   []byte
	Value []byte
}

type WatchKeyResponse struct {
	Watcher clientv3.Watcher
	Id      clientv3.LeaseID
	Event   <-chan *KeyChangeEvent
}

// new a etcd client
func NewEtcdClient(endpoints []string, userName, password string) (*Cli, error) {

	c, err := clientv3.New(clientv3.Config{
		Endpoints: endpoints,
		Username:  userName,
		Password:  password,
	})
	if err != nil {
		log.Errorf("create etcd client fail:%s", err.Error())
		return nil, err
	}
	cli := &Cli{
		kv:    clientv3.NewKV(c),
		lease: clientv3.NewLease(c),
		c:     c,
	}
	log.Debug("the etcd endpoints:%+v ,username:%s,password:%s init success", endpoints, userName, password)
	return cli, err
}

// put a key with value
func (cli *Cli) Put(ctx context.Context, key, value string) error {

	_, err := cli.kv.Put(ctx, key, value)
	return err
}

// put not exist key
func (cli *Cli) PutWithNotExist(ctx context.Context, key, value string) (bool, error) {

	tx := cli.c.Txn(ctx)

	txnResponse, err := tx.If(clientv3.Compare(clientv3.Version(key), "=", 0)).
		Then(clientv3.OpPut(key, value)).
		Else(clientv3.OpGet(key)).
		Commit()
	if err != nil {
		return false, err
	}
	return txnResponse.Succeeded, err
}

// delete a key
func (cli *Cli) Delete(ctx context.Context, key string) error {
	_, err := cli.kv.Delete(ctx, key)
	return err
}

// delete with prefix
func (cli *Cli) DeleteWithPrefix(ctx context.Context, key string) error {
	_, err := cli.kv.Delete(ctx, key, clientv3.WithPrefix())
	return err
}

// keep alive once with key
func (cli *Cli) KeepAliveOnce(ctx context.Context, key, value string, ttl int64) error {

	leaseResponse, err := cli.lease.Grant(ctx, ttl)
	if err != nil {
		log.Warnf("the key:%s,value:%s,lease grant fail:%s", key, value, err.Error())
		return err
	}
	leaseId := leaseResponse.ID
	_, err = cli.kv.Put(ctx, key, value, clientv3.WithLease(leaseId))
	if err != nil {
		log.Warnf("the key:%s,value:%s,put fail:%s", key, value, err.Error())
		return err
	}
	_, err = cli.lease.KeepAliveOnce(ctx, leaseId)
	if err != nil {
		log.Warnf("the key:%s,value:%s,keepalive once fail:%s", key, value, err.Error())
	}
	return err
}

// keep alive once with key
func (cli *Cli) KeepAlive(ctx context.Context, key, value string, ttl int64) (clientv3.LeaseID, error) {

	leaseResponse, err := cli.lease.Grant(ctx, ttl)
	if err != nil {
		log.Warnf("the key:%s,value:%s,lease grant fail:%s", key, value, err.Error())
		return 0, err
	}
	leaseId := leaseResponse.ID
	_, err = cli.kv.Put(ctx, key, value, clientv3.WithLease(leaseId))
	if err != nil {
		log.Warnf("the key:%s,value:%s,put fail:%s", key, value, err.Error())
		return leaseId, err
	}
	leaseKeepAliveResponseChan, err := cli.lease.KeepAlive(ctx, leaseId)
	if err != nil {
		log.Warnf("the key:%s,value:%s,keepalive fail:%s", key, value, err.Error())
	}

	go handleLeaseKeepaliveChan(key, leaseKeepAliveResponseChan)
	return leaseId, err
}

// handle the lease keepalive chan
func handleLeaseKeepaliveChan(key string, ch <-chan *clientv3.LeaseKeepAliveResponse) {

	go func() {
		for {
			select {
			case c, ok := <-ch:
				if !ok {
					log.Warnf("keepalive chan key:%s has closed", key)
					return
				}
				log.Debug("keepalive chan  key:%s,leaseId:%d", key, c.ID)
			}
		}
	}()

}

func (cli *Cli) Watch(ctx context.Context, key string) *WatchKeyResponse {

	watchCh := cli.c.Watch(ctx, key)
	changeEventCh := make(chan *KeyChangeEvent, 32)

	go func() {
		for {
			ch := <-watchCh
			if ch.Canceled {
				log.Warnf("the watcher key:%s has canceled...", key)
				break
			}
			for _, event := range ch.Events {
				handleWatchEvent(event, changeEventCh)
			}
		}
	}()

	w := clientv3.NewWatcher(cli.c)
	return &WatchKeyResponse{
		Watcher: w,
		Event:   changeEventCh,
	}
}

// watch with prefix
func (cli *Cli) WatchWithPrefix(ctx context.Context, prefix string) *WatchKeyResponse {

	watchCh := cli.c.Watch(ctx, prefix, clientv3.WithPrefix())
	changeEventCh := make(chan *KeyChangeEvent, 32)

	go func() {
		for {
			ch := <-watchCh
			if ch.Canceled {
				log.Warnf("the watcher prefix:%s has canceled...", prefix)
				break
			}
			for _, event := range ch.Events {
				handleWatchEvent(event, changeEventCh)
			}
		}
	}()

	w := clientv3.NewWatcher(cli.c)
	return &WatchKeyResponse{
		Watcher: w,
		Event:   changeEventCh,
	}
}

// revoke  lease
func (cli *Cli) Revoke(ctx context.Context, id clientv3.LeaseID) error {

	_, err := cli.lease.Revoke(ctx, id)
	return err
}

// handle the watch key change  event
func handleWatchEvent(event *clientv3.Event, ch chan *KeyChangeEvent) {

	e := &KeyChangeEvent{}
	switch event.Type {
	case mvccpb.PUT:
		e.Type = CreateKeyChangeEvent
		if event.IsModify() {
			e.Type = UpdateKeyChangeEvent
		}
		e.Key = event.Kv.Key
		e.Value = event.Kv.Value
	case mvccpb.DELETE:
		e.Type = DeleteKeyChangeEvent
		e.Key = event.PrevKv.Key
		e.Value = event.PrevKv.Value
	}
	ch <- e
}
