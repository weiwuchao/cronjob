package apiserver

import (
	"crontab/master/access"
	"crontab/master/common"
	"crontab/master/config"
	"encoding/json"
	"net"
	"net/http"
	"strconv"
	"time"
)

type ApiServer struct {
	httpServer *http.Server
}

var (
	G_apiServer ApiServer
)

func InitApiServer() error {

	var (
		listener net.Listener
		err      error
	)

	//配置路由
	serveMux := http.NewServeMux()

	serveMux.HandleFunc("/job/save", saveJob)
	serveMux.HandleFunc("/job/delete", deleteJob)
	serveMux.HandleFunc("/job/list", listJob)
	serveMux.HandleFunc("/job/list", killJob)

	//启动tcp监听
	if listener, err = net.Listen("tcp", ":"+strconv.Itoa(config.G_config.ApiPort)); err != nil {
		return err
	}

	//创建http服务
	httpServer := &http.Server{
		// 将int转化为time.Duration类型
		ReadTimeout:  time.Duration(config.G_config.ApiReadTimeOut) * time.Millisecond,
		WriteTimeout: time.Duration(config.G_config.ApiWriteTimeOut) * time.Millisecond,
		Handler:      serveMux,
	}

	// 创建单实例apiServer
	G_apiServer = ApiServer{
		httpServer: httpServer,
	}

	//启动服务端
	go httpServer.Serve(listener)

	return nil
}

func saveJob(resp http.ResponseWriter, req *http.Request) {

	var (
		err     error
		postjob string
		job     common.Job
		oldJob  *common.Job
		bytes   []byte
	)
	//解析表单
	if err = req.ParseForm(); err != nil {
		goto ERR
	}

	//获取表单中job字段
	postjob = req.PostForm.Get("job")
	//反序列化job
	if err = json.Unmarshal([]byte(postjob), &job); err != nil {
		goto ERR
	}
	//保存job到etcd
	if oldJob, err = access.G_jobMgr.SaveJob(&job); err != nil {
		goto ERR
	}
	if bytes, err = common.BuildResp("200", "success", oldJob); err == nil {
		resp.Write(bytes)
	}
	return
ERR:
	if bytes, err = common.BuildResp("500", err.Error(), oldJob); err == nil {
		resp.Write(bytes)
	}
}

func deleteJob(resp http.ResponseWriter, req *http.Request) {

	var (
		err    error
		oldJob *common.Job
		bytes  []byte
		name   string
	)
	//解析表单
	if err = req.ParseForm(); err != nil {
		goto ERR
	}

	//获取表单中job字段
	name = req.PostForm.Get("name")

	//删除job从etcd
	if oldJob, err = access.G_jobMgr.DeleteJob(name); err != nil {
		goto ERR
	}
	if bytes, err = common.BuildResp("200", "success", oldJob); err == nil {
		resp.Write(bytes)
	}
	return
ERR:
	if bytes, err = common.BuildResp("500", err.Error(), oldJob); err == nil {
		resp.Write(bytes)
	}
}

// 杀死任务
func killJob(resp http.ResponseWriter, req *http.Request) {

	var (
		err    error
		bytes  []byte
		name   string
	)
	//解析表单
	if err = req.ParseForm(); err != nil {
		goto ERR
	}

	//获取表单中job字段
	name = req.PostForm.Get("name")

	//通知worker杀死job
	if  err = access.G_jobMgr.KillJob(name); err != nil {
		goto ERR
	}
	if bytes, err = common.BuildResp("200", "success",""); err == nil {
		resp.Write(bytes)
	}
	return
ERR:
	if bytes, err = common.BuildResp("500", err.Error(), ""); err == nil {
		resp.Write(bytes)
	}
}

// 查询任务列表
func listJob(resp http.ResponseWriter, req *http.Request) {

	var (
		err    error
		jobList []*common.Job
		bytes  []byte
	)

	//删除job从etcd
	if jobList, err = access.G_jobMgr.ListJob(); err != nil {
		goto ERR
	}
	if bytes, err = common.BuildResp("200", "success", jobList); err == nil {
		resp.Write(bytes)
	}
	return
ERR:
	if bytes, err = common.BuildResp("500", err.Error(), jobList); err == nil {
		resp.Write(bytes)
	}
}