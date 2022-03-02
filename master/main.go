package main

import (
	"crontab/master/access"
	"crontab/master/apiserver"
	"crontab/master/config"
	"flag"
	"fmt"
	"runtime"
)

var (
	configPath string
)

func main(){
	var (
		err error
	)

	//初始化命令行参数
	initArgs()

	//初始化最大线程数
	initEnv()

	//初始化配置文件
	if err=config.InitConfig(configPath);err!=nil{
		goto ERR
	}

	//初始化jobMgr(etcd)
	if err=access.InitJobMgr();err!=nil{
		goto ERR
	}

	//初始化httpServer
	if err=apiserver.InitApiServer();err!=nil{
		goto ERR
	}
	return
	ERR:
		fmt.Println(err)
}

func initEnv(){
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func initArgs(){
	flag.StringVar(&configPath,"configFilePath","./src/crontab/master/config/config.yaml","配置参数路径")
}