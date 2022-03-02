package apiserver

import (
	"crontab/master/config"
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
		err error
	)

	//配置路由
	serveMux:=http.NewServeMux()

	serveMux.HandleFunc("/job/save",saveJob)

	//启动tcp监听
	if listener,err=net.Listen("tcp",":"+strconv.Itoa(config.G_config.ApiPort));err!=nil{
		return err
	}

	//创建http服务
	httpServer:=&http.Server{
		// 将int转化为time.Duration类型
		ReadTimeout: time.Duration(config.G_config.ApiReadTimeOut)*time.Millisecond,
		WriteTimeout: time.Duration(config.G_config.ApiWriteTimeOut)*time.Millisecond,
		Handler: serveMux,
	}

	//创建单实例apiServer
	G_apiServer=ApiServer{
		httpServer: httpServer,
	}

	//启动服务端
	go httpServer.Serve(listener)

	return nil
}


func saveJob(w http.ResponseWriter, r *http.Request){

}