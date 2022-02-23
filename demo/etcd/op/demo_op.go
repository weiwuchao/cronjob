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
		putOp clientv3.Op
		getOp clientv3.Op
		opResp clientv3.OpResponse
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

	//创建写入op
	putOp=clientv3.OpPut("demo2","11")

	//执行写入op
	if opResp,err=kv.Do(context.Background(),putOp);err!=nil{
		fmt.Println(err)
		return
	}

	fmt.Println("写入Reversion:",opResp.Put().Header.Revision)

	//创建获取op
	getOp=clientv3.OpGet("demo2")

	//执行获取op
	if opResp,err=kv.Do(context.Background(),getOp);err!=nil{
		fmt.Println(err)
		return
	}
	fmt.Println("获取Reversion:",opResp.Get().Kvs[0].ModRevision)
	fmt.Println("获取数据:",string(opResp.Get().Kvs[0].Value))

}
