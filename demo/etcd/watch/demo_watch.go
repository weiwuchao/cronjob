package main

import (
	"context"
	"fmt"
	clientv3 "go.etcd.io/etcd/client/v3"
	"time"
)

func main() {

	var (
		config          clientv3.Config
		client          *clientv3.Client
		err             error
		kv              clientv3.KV
		getResp         *clientv3.GetResponse
		watcher         clientv3.Watcher
		watcherChanResp <-chan clientv3.WatchResponse
		watcherResp     clientv3.WatchResponse
		event           *clientv3.Event
	)

	// 设置etcd配置
	config = clientv3.Config{
		Endpoints:   []string{"192.168.92.129:2379"},
		DialTimeout: 5 * time.Second,
	}

	// etcd客户端连接
	if client, err = clientv3.New(config); err != nil {
		fmt.Println(err)
		return
	}

	kv = clientv3.NewKV(client)

	//模拟etcd中的KV变化
	go func() {
		for {
			kv.Put(context.Background(), "demo1", "6666")
			kv.Delete(context.Background(), "demo1")
			time.Sleep(1 * time.Second)
		}
	}()

	//获取到当前值
	if getResp, err = kv.Get(context.Background(), "demo1"); err != nil {
		fmt.Println(err)
		return
	}

	if len(getResp.Kvs) != 0 {
		fmt.Println("当前值是:", string(getResp.Kvs[0].Value))
	}

	revision := getResp.Header.Revision + 1
	fmt.Println("从该版本开始监听：", revision)

	//创建监听器
	watcher = clientv3.NewWatcher(client)
	//设置5秒后取消的context
	ctx, cancelFunc := context.WithCancel(context.TODO())
	time.AfterFunc(5*time.Second, func() {
		cancelFunc()
	})
	watcherChanResp = watcher.Watch(ctx, "demo1", clientv3.WithRev(revision))

	/*for watcherResp=range watcherChanResp{
		for _,event=range watcherResp.Events{
			switch event.Type {
			case clientv3.EventTypePut:
				fmt.Println("修改为:",string(event.Kv.Value),"reversion:",event.Kv.CreateRevision,event.Kv.ModRevision)
			case clientv3.EventTypeDelete:
				fmt.Println("删除了","revision:",event.Kv.ModRevision)
			}
		}
	}*/

	for {
		select {
		case watcherResp= <-watcherChanResp:
			for _, event = range watcherResp.Events {
				switch event.Type {
				case clientv3.EventTypePut:
					fmt.Println("修改为:", string(event.Kv.Value), "reversion:", event.Kv.CreateRevision, event.Kv.ModRevision)
				case clientv3.EventTypeDelete:
					fmt.Println("删除了", "revision:", event.Kv.ModRevision)
				}
			}
		default:
			break
		}
	}
}
