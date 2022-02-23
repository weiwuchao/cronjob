package main

import (
	"context"
	"fmt"
	clientv3 "go.etcd.io/etcd/client/v3"
	"time"
)

func main() {

	var (
		config         clientv3.Config
		client         *clientv3.Client
		err            error
		kv             clientv3.KV
		putResp        *clientv3.PutResponse
		getResp        *clientv3.GetResponse
		lease          clientv3.Lease
		leaseGrantResp *clientv3.LeaseGrantResponse
		leaseId        clientv3.LeaseID
		keepResp	*clientv3.LeaseKeepAliveResponse
		keepRespChan  <-chan *clientv3.LeaseKeepAliveResponse
	)

	// 设置etcd配置
	config = clientv3.Config{
		Endpoints:   []string{"192.168.92.128:2379"},
		DialTimeout: 5 * time.Second,
	}

	// etcd客户端连接
	if client, err = clientv3.New(config); err != nil {
		fmt.Println(err)
		return
	}

	// 申请一个租约
	lease = clientv3.NewLease(client)

	// 创建一个10秒的租约
	if leaseGrantResp, err = lease.Grant(context.Background(), 10); err != nil {
		fmt.Println(err)
		return
	}

	// 获取续约ID
	leaseId = leaseGrantResp.ID

	if keepRespChan,err=lease.KeepAlive(context.Background(),leaseId);err!=nil{
		fmt.Println(err)
		return
	}

	go func() {
		for{
			select {
			case keepResp=<-keepRespChan:
				if keepRespChan==nil{
					fmt.Println("租约已经失效")
					goto END
				}else{
					fmt.Println("收到自动续租应答",keepResp.ID)
				}
			}
		}
		END:
	}()

	// 创建kv
	kv = clientv3.NewKV(client)

	// 写入kv
	if putResp, err = kv.Put(context.Background(), "test3", "9999", clientv3.WithLease(leaseId)); err != nil {
		fmt.Println(err)
		return
	} else {
		fmt.Println("写入成功:", putResp.Header.Revision)
	}

	// 定时读取key是否过期
	for {
		if getResp, err = kv.Get(context.Background(), "test3"); err != nil {
			fmt.Println(err)
			return
		}
		if getResp.Count == 0 {
			fmt.Println("过期了..")
			break
		}
		fmt.Println("还没过期", getResp.Kvs)
		time.Sleep(2 * time.Second)
	}

}
