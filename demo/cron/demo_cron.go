package main

import (
	"fmt"
	"github.com/gorhill/cronexpr"
	"time"
)

type CronJob struct{
	expr *cronexpr.Expression
	nextTime time.Time
}
func main(){

	var (
		cronJob *CronJob
		expr *cronexpr.Expression
		nowTime time.Time
		scheduleTable map[string]*CronJob
	)

	scheduleTable=make(map[string]*CronJob)
	nowTime=time.Now()

	//定义两个cronJob
	expr=cronexpr.MustParse("*/5 * * * * * *")
	cronJob=&CronJob{
		expr:expr,
		nextTime: expr.Next(nowTime),
	}
	scheduleTable["job1"]=cronJob
	expr=cronexpr.MustParse("*/5 * * * * * *")
	cronJob=&CronJob{
		expr:expr,
		nextTime: expr.Next(nowTime),
	}
	scheduleTable["job2"]=cronJob

	//启动一个调度协程
	go func(){
		var (
			jobName string
			cronJob *CronJob
			nowTime time.Time
		)
		for{
			nowTime=time.Now()
			//遍历job
			for jobName,cronJob=range scheduleTable{
				//执行时间>=当前时间
				if cronJob.nextTime.Before(nowTime)|| cronJob.nextTime.Equal(nowTime){
					//启动一个协程，执行这个任务
					go func(jobName string){
						fmt.Println("执行这个任务",jobName)
						//计算下次调度时间
						cronJob.nextTime=cronJob.expr.Next(nowTime)
						fmt.Println("下次执行时间",cronJob.nextTime)
					}(jobName)
				}
			}

			//休眠100毫秒
			select {
			case <-time.NewTimer(1000*time.Millisecond).C:
			}
		}
	}()

	time.Sleep(10000*time.Millisecond)
}
