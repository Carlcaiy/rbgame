package network

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/namespace"
)

type Etcd struct {
	cli  *clientv3.Client
	conf *ServerConfig
	cb   func(interface{})
}

func NewEtcd(s *ServerConfig, callback func(interface{})) *Etcd {
	ret := new(Etcd)
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:         []string{"127.0.0.1:2379"},
		AutoSyncInterval:  time.Second,
		DialTimeout:       time.Second * 3,
		DialKeepAliveTime: time.Second * 5,
		// Username:          "cyf",
		// Password:          "cyf123",
	})
	if err != nil {
		panic(err)
	}
	cli.KV = namespace.NewKV(cli.KV, "cyf/")
	cli.Watcher = namespace.NewWatcher(cli.Watcher, "cyf/")
	cli.Lease = namespace.NewLease(cli.Lease, "cyf/")
	ret.cli = cli
	ret.conf = s
	ret.cb = callback
	return ret
}

func (p *Etcd) Init() {
	p.Put()

	// 如果有订阅需求的才去获取，正常只需要上报就行
	if len(p.conf.Subs) > 0 {
		p.Get()
		go p.Watch()
	}
}

func (p *Etcd) key() string {
	return fmt.Sprintf("server/%d/%d", p.conf.ServerType, p.conf.ServerId)
}

func (p *Etcd) parseKey(key string) {
	fmt.Println("parseKey", key)
	strs := strings.Split(key, "/")
	if len(strs) != 3 {
		fmt.Println("key struct wrong", key)
		return
	}
	if strs[0] != "server" {
		fmt.Println("without server prefix", strs[0])
		return
	}
	serverType, err := strconv.Atoi(strs[1])
	if err != nil {
		fmt.Println("parse server type error", strs[1])
		return
	}
	serverID, err := strconv.Atoi(strs[2])
	if err != nil {
		fmt.Println("parse server id error", strs[2])
		return
	}

	fmt.Printf("parseKey success serverType=%d serverID=%d\n", serverType, serverID)
}

func (p *Etcd) parseValue(value []byte) {
	fmt.Println("parseValue", string(value))
	if len(value) > 0 {
		conf := new(ServerConfig)
		if err := json.Unmarshal(value, conf); err != nil {
			fmt.Println(err)
			return
		}
		if p.conf.isSub(conf.ServerType) {
			p.cb(conf)
		}
	}
}

func (p *Etcd) Put() {
	value, err := json.Marshal(p.conf)
	if err != nil {
		fmt.Println("server marshal failed")
		return
	}
	fmt.Println("put key", p.key(), "value", string(value))
	p.cli.Put(context.TODO(), p.key(), string(value))
}

func (p *Etcd) Get() {
	res, err := p.cli.Get(context.TODO(), "server/", clientv3.WithPrefix())
	if err != nil {
		return
	}
	for _, kv := range res.Kvs {
		p.parseKey(string(kv.Key))
		p.parseValue(kv.Value)
	}
}

func (p *Etcd) Del() {
	p.cli.Delete(context.TODO(), p.key())
}

func (p *Etcd) Watch() {
	wg.Add(1)
	defer func() {
		wg.Done()
	}()
	watchCh := p.cli.Watch(context.TODO(), "server", clientv3.WithPrefix())
	for {
		select {
		case <-closech:
			return
		case watch := <-watchCh:
			for _, event := range watch.Events {
				switch event.Type {
				case clientv3.EventTypePut:
					fmt.Println("etcd event put")
					p.parseKey(string(event.Kv.Key))
					p.parseValue(event.Kv.Value)
				case clientv3.EventTypeDelete:
					fmt.Println("etcd event delete")
					p.parseKey(string(event.Kv.Key))
					p.parseValue(event.Kv.Value)
				}
			}
		}
	}
}

func (e *Etcd) Close() {
	e.Del()
	e.cli.Close()
}
