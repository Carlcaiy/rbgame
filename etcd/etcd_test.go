package main

import (
	"context"
	"fmt"
	"testing"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/namespace"
)

func Client() *clientv3.Client {
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
	return cli
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}

func TestKV(t *testing.T) {
	cli := Client()
	ctx := context.Background()
	res1, err := cli.Put(ctx, "server/data", "beautiful")
	must(err)
	t.Log(res1.PrevKv.String())

	res, err := cli.Get(ctx, "s", clientv3.WithPrefix())
	must(err)
	for _, kv := range res.Kvs {
		t.Log(kv.String(), string(kv.Key))
	}
}

func TestLease(t *testing.T) {
	cli := Client()
	ctx := context.Background()
	res, err := cli.Grant(ctx, int64(time.Second*3))
	must(err)
	fmt.Println(res.String())
	leaseId := res.ID

	cli.Put(ctx, "sample", "value", clientv3.WithLease(leaseId))
	xx, err := cli.Get(ctx, "sample")
	must(err)
	for _, kv := range xx.Kvs {
		fmt.Println(kv.String())
	}

	xxx, _ := cli.KeepAlive(ctx, leaseId)
	// cli.Revoke(ctx, leaseId)
	// cli.TimeToLive(ctx, leaseId)

	for {
		x, ok := <-xxx
		if ok {
			fmt.Println("for", x)
		} else {
			fmt.Println("failed")
			break
		}
	}
}

func TestGet(t *testing.T) {
	cli := Client()
	ctx := context.Background()
	pres, err := cli.Put(ctx, "李志", "牛逼1", clientv3.WithPrevKV())
	fmt.Println(string(pres.PrevKv.Key), string(pres.PrevKv.Value), err)
	pres, err = cli.Put(ctx, "李志真", "牛逼1", clientv3.WithPrevKV())
	fmt.Println(string(pres.PrevKv.Key), string(pres.PrevKv.Value), err)
	pres, err = cli.Put(ctx, "李志真tm", "牛逼1", clientv3.WithPrevKV())
	fmt.Println(string(pres.PrevKv.Key), string(pres.PrevKv.Value), err)
	pres, err = cli.Put(ctx, "李志真的是", "牛逼1", clientv3.WithPrevKV())
	fmt.Println(string(pres.PrevKv.Key), string(pres.PrevKv.Value), err)
	pres, err = cli.Put(ctx, "李志真真真", "牛逼1", clientv3.WithPrevKV())
	fmt.Println(string(pres.PrevKv.Key), string(pres.PrevKv.Value), err)

	gres, _ := cli.Get(ctx, "李志", clientv3.WithPrefix())
	for _, kv := range gres.Kvs {
		fmt.Println(string(kv.Key), string(kv.Value))
	}
}

func TestDelete(t *testing.T) {
	cli := Client()
	ctx := context.Background()
	res, _ := cli.Delete(ctx, "watching")
	fmt.Println(res)
}

func TestWatch(t *testing.T) {
	cli := Client()
	ctx := context.Background()
	res, _ := cli.Get(ctx, "watching", clientv3.WithPrefix())
	for i := range res.Kvs {
		fmt.Println(res.Kvs[i].String())
	}

	ch := cli.Watch(ctx, "watching", clientv3.WithPrefix())
	for {
		data, ok := <-ch
		if ok {
			for _, events := range data.Events {
				fmt.Println(events.Kv.String())
			}
		} else {
			break
		}
	}
}

func TestTxn(t *testing.T) {
	cli := Client()
	ctx := context.Background()

	res, err := cli.Txn(ctx).If(
		clientv3.Compare(clientv3.Value("watching"), "=", "dasda5"),
	).Then(
		clientv3.OpPut("watching", "heihei"),
	).Else(
		clientv3.OpPut("watching", "heibro"),
	).Commit()
	must(err)
	fmt.Println(res)

	res, err = cli.Txn(ctx).If(
		clientv3.Compare(clientv3.Version("watching"), "=", 129),
	).Then(
		clientv3.OpPut("watching", "hei129"),
	).Else(
		clientv3.OpPut("watching", "hei130"),
	).Commit()

	must(err)
	fmt.Println(res)

	res, err = cli.Txn(ctx).If(
		clientv3.Compare(clientv3.LeaseValue("watching"), "=", 129),
	).Then(
		clientv3.OpPut("watching", "hei129"),
	).Else(
		clientv3.OpPut("watching", "hei130"),
	).Commit()

	must(err)
	fmt.Println(res)
}
