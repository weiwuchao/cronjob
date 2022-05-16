package access

import (
	"context"
	"crontab/master/common"
	"crontab/master/config"
	"encoding/json"
	"fmt"
	clientv3 "go.etcd.io/etcd/client/v3"
	"time"
)

type JobMgr struct {
	kv     clientv3.KV
	lease  clientv3.Lease
	client *clientv3.Client
}

var (
	G_jobMgr JobMgr
)

func InitJobMgr() error {
	var (
		conf   clientv3.Config
		client *clientv3.Client
		err    error
		kv     clientv3.KV
		lease  clientv3.Lease
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

	G_jobMgr = JobMgr{
		kv:     kv,
		lease:  lease,
		client: client,
	}
	return nil
}

//保存job
func (jobMgr JobMgr) SaveJob(job *common.Job) (*common.Job, error) {
	var (
		err     error
		jobByte []byte
		putResp *clientv3.PutResponse
		oldJob  *common.Job
	)
	jobKey := "/cron/jobs/" + job.Name
	//反序列化job
	if jobByte, err = json.Marshal(job); err != nil {
		return nil, err
	}
	if putResp, err = jobMgr.kv.Put(context.Background(), jobKey, string(jobByte), clientv3.WithPrevKV()); err != nil {
		return nil, err
	}

	//如果是更新，返回旧值
	if putResp.PrevKv != nil {
		json.Unmarshal(putResp.PrevKv.Value, oldJob)
	}
	return oldJob, err
}

//删除job
func (jobMgr JobMgr) DeleteJob(name string) (*common.Job, error) {
	var (
		err     error
		delResp *clientv3.DeleteResponse
		oldJob  *common.Job
	)
	jobKey := "/cron/jobs/" + name

	if delResp, err = jobMgr.kv.Delete(context.Background(), jobKey, clientv3.WithPrevKV()); err != nil {
		return nil, err
	}

	//如果是更新，返回旧值
	if delResp.PrevKvs != nil && len(delResp.PrevKvs) > 0 {
		json.Unmarshal(delResp.PrevKvs[0].Value, oldJob)
	}
	return oldJob, err
}

//查询job列表
func (jobMgr JobMgr) ListJob() ([]*common.Job, error) {
	var (
		err     error
		getResp *clientv3.GetResponse
		job     *common.Job
	)
	jobDir := "/cron/jobs/"

	if getResp, err = jobMgr.kv.Get(context.Background(), jobDir, clientv3.WithPrefix()); err != nil {
		return nil, err
	}

	jobs := make([]*common.Job, 0)
	//遍历目录
	for _, kv := range getResp.Kvs {
		job = &common.Job{}
		if err = json.Unmarshal(kv.Value, job); err != nil {
			err = nil
			continue
		}
		jobs = append(jobs, job)
	}
	return jobs, err
}

//杀死任务
func (jobMgr JobMgr) KillJob(name string) error {
	var (
		err     error
		leaseResp *clientv3.LeaseGrantResponse
	)
	jobKey := "/cron/killer/" + name

	//让worker监听到一次操作即可，为了不占etcd存储，设置租约，监听到后即失效
	if leaseResp,err=jobMgr.lease.Grant(context.Background(),1);err!=nil{
		return  err
	}

	//租约id
	leaseId:=leaseResp.ID

	_, err = jobMgr.kv.Put(context.Background(), jobKey, "",clientv3.WithLease(leaseId))
	return err

}
