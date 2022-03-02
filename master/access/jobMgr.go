package access

import (
	"crontab/master/config"
	"fmt"
	clientv3 "go.etcd.io/etcd/client/v3"
	"time"
)

type JobMgr struct{
	kv             clientv3.KV
	lease          clientv3.Lease
	client         *clientv3.Client
}

var (
	G_jobMgr JobMgr
)

func InitJobMgr() error{
	var (
		conf         clientv3.Config
		client         *clientv3.Client
		err            error
		kv             clientv3.KV
		lease          clientv3.Lease
	)

	// 设置etcd配置
	conf = clientv3.Config{
		Endpoints:   config.G_config.EtcdEndpoints,
		DialTimeout: time.Duration(config.G_config.EtcdDialTimeOut) * time.Millisecond,
	}

	// etcd客户端连接
	if client, err = clientv3.New(conf); err != nil {
		fmt.Println(err)
		return err
	}

	// 申请一个租约
	lease = clientv3.NewLease(client)

	// 创建kv
	kv = clientv3.NewKV(client)

	G_jobMgr=JobMgr{
		kv:     kv,
		lease:  lease,
		client: client,
	}
	return nil
}



