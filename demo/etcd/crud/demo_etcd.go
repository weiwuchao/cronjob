package main

import (
	"context"
	"fmt"
	clientv3 "go.etcd.io/etcd/client/v3"
	"time"
)

func main()  {

	var (
		config clientv3.Config
		client *clientv3.Client
		err error
		kv clientv3.KV
		//putResp *clientv3.PutResponse
		getResp *clientv3.GetResponse
		delResp *clientv3.DeleteResponse
	)

	//设置etcd配置
	config=clientv3.Config{
		Endpoints: []string{"192.168.92.128:2379"},
		DialTimeout: 5*time.Second,
	}

	//etcd客户端连接
	if client,err=clientv3.New(config);err!=nil{
		fmt.Println(err)
		return
	}

	//用KV操作etcd键值
	kv=clientv3.NewKV(client)

	// put操作
	/*if putResp,err=kv.Put(context.Background(),"test2","129",clientv3.WithPrevKV());err!=nil{
		fmt.Println(err)
	}else{
		fmt.Println("Revision",putResp.Header.Revision)
		//获取历史值
		if putResp.PrevKv!=nil{
			fmt.Println("PrevKv:",string(putResp.PrevKv.Value))
		}
	}*/

	//get操作
	if getResp,err=kv.Get(context.Background(),"test2",clientv3.WithCountOnly());err!=nil{
		fmt.Println(err)
	}else{
		fmt.Println(getResp.Kvs,getResp.Count)
	}

	//delete操作
	if delResp,err=kv.Delete(context.Background(),"test2",clientv3.WithPrevKV());err!=nil{
		fmt.Println(err)
	}
	//查看删除的值
	if len(delResp.PrevKvs)>0{
		for _,kv:=range delResp.PrevKvs{
			fmt.Println("删除了:",string(kv.Key),string(kv.Value))
		}
	}



	defer client.Close()
}