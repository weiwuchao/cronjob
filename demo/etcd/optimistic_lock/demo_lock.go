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
		lease          clientv3.Lease
		leaseGrantResp *clientv3.LeaseGrantResponse
		leaseId        clientv3.LeaseID
		keepResp       *clientv3.LeaseKeepAliveResponse
		keepRespChan   <-chan *clientv3.LeaseKeepAliveResponse
		txn            clientv3.Txn
		txnResp        *clientv3.TxnResponse
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

	// 申请一个租约
	lease = clientv3.NewLease(client)

	// 创建一个10秒的租约
	if leaseGrantResp, err = lease.Grant(context.Background(), 10); err != nil {
		fmt.Println(err)
		return
	}

	// 获取续约ID
	leaseId = leaseGrantResp.ID
	ctx, cancelFunc := context.WithCancel(context.TODO())
	if keepRespChan, err = lease.KeepAlive(ctx, leaseId); err != nil {
		fmt.Println(err)
		return
	}
	go func() {
		for {
			select {
			case keepResp = <-keepRespChan:
				if keepRespChan == nil {
					fmt.Println("租约已经失效")
					goto END
				} else {
					fmt.Println("收到自动续租应答", keepResp.ID)
				}
			}
		}
	END:
	}()

	// 创建kv
	kv = clientv3.NewKV(client)

	//创建事务
	txn = kv.Txn(context.Background())

	//抢锁,先判断是否有lock键，不存在则带租约写入，已经存在则获取对应键当前值
	txn.If(clientv3.Compare(clientv3.CreateRevision("lock"), "=", 0)).
		Then(clientv3.OpPut("lock", "true", clientv3.WithLease(leaseId))).
		Else(clientv3.OpGet("lock")) //抢锁失败

	//提交事务
	if txnResp, err = txn.Commit(); err != nil {
		fmt.Println(err)
		return
	}

	if !txnResp.Succeeded {
		//锁被占用
		fmt.Println("锁被占用")
		fmt.Println(txnResp.Responses[0].GetResponseRange().Kvs[0].Value)
	}

	//释放锁
	defer cancelFunc()
	defer lease.Revoke(context.Background(), leaseId)

	//处理业务
	fmt.Println("处理业务")
	time.Sleep(5 * time.Second)
}
