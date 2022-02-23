package main

import (
	"context"
	"fmt"
	"os/exec"
	"sync"
	"time"
)

func main() {

	// 运行命令
	/*cmd := exec.Command("C:\\Windows\\System32\\bash.exe", "-c", "sleep 5")
	if err := cmd.Run(); err != nil {
		fmt.Println("not normal", err)
	}
	fmt.Println("normal run")*/

	//运行命令并获取返回值
	/*var output []byte
	cmd2 := exec.Command("C:\\Windows\\System32\\bash.exe", "-c", "touch a.txt")
	if output, err := cmd2.CombinedOutput(); err != nil {
		fmt.Println("not normal", string(output),err)
	}
	fmt.Println(string(output))*/


	//运行命令并强制结束
	/*
		1. WithCancel()函数接受一个 Context 并返回其子Context和取消函数cancel
		2. 新创建协程中传入子Context做参数，且需监控子Context的Done通道，若收到消息，则退出
		3. 需要新协程结束时，在外面调用 cancel 函数，即会往子Context的Done通道发送消息
		4. 注意：当 父Context的 Done() 关闭的时候，子 ctx 的 Done() 也会被关闭
	*/
	var wg sync.WaitGroup
	//使用上下文实现强制结束
	ctx,canelFunc:=context.WithCancel(context.TODO())
	wg.Add(1)
	go func(){
		var (
			output []byte
			err error
		)
		cmd2 := exec.CommandContext(ctx,"C:\\Windows\\System32\\bash.exe", "-c", "touch a.txt")
		if output, err = cmd2.CombinedOutput(); err != nil {
			fmt.Println("not normal", string(output),err)
		}
		fmt.Println(string(output))
		wg.Done()
	}()
	time.Sleep(1*time.Second)
	canelFunc()
	wg.Wait()
}
