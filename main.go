package main

import (
	_ "github.com/Jordanzuo/ChatServerCenter/src/bll/manageCenterBLL"
)

import (
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"sync"
	"syscall"
	"time"

	"github.com/Jordanzuo/ChatServerCenter/src/bll/apiLogBLL"
	"github.com/Jordanzuo/ChatServerCenter/src/bll/configBLL"
	"github.com/Jordanzuo/ChatServerCenter/src/bll/ipBLL"
	"github.com/Jordanzuo/ChatServerCenter/src/bll/messageLogBLL"
	"github.com/Jordanzuo/ChatServerCenter/src/bll/onlineLogBLL"
	"github.com/Jordanzuo/ChatServerCenter/src/bll/reloadBLL"
	"github.com/Jordanzuo/ChatServerCenter/src/bll/wordBLL"
	"github.com/Jordanzuo/ChatServerCenter/src/config"
	"github.com/Jordanzuo/ChatServerCenter/src/rpc"
	"github.com/Jordanzuo/ChatServerCenter/src/web"
	"github.com/Jordanzuo/goutil/debugUtil"
	"github.com/Jordanzuo/goutil/logUtil"
)

var (
	wg sync.WaitGroup
)

func init() {
	// 设置WaitGroup需要等待的数量，只要有一个服务器出现错误都停止服务器
	wg.Add(1)
	runtime.NumGoroutine()
}

// 处理系统信号
func signalProc() {
	// 处理内部未处理的异常，以免导致主线程退出，从而导致系统崩溃
	defer func() {
		if r := recover(); r != nil {
			logUtil.LogUnknownError(r)
		}
	}()

	sigs := make(chan os.Signal)
	signal.Notify(sigs, syscall.SIGTERM, syscall.SIGINT, syscall.SIGHUP)

	for {
		// 准备接收信息
		sig := <-sigs

		// 输出信号
		debugUtil.Println("sig:", sig)

		if sig == syscall.SIGHUP {
			logUtil.Log("收到重启的信号，准备重新加载配置", logUtil.Info, true)

			// 重新加载配置
			reloadBLL.Reload()

			logUtil.Log("收到重启的信号，重新加载配置完成", logUtil.Info, true)
		} else {
			logUtil.Log("收到退出程序的信号，开始退出……", logUtil.Info, true)

			// 做一些收尾的工作

			logUtil.Log("收到退出程序的信号，退出完成……", logUtil.Info, true)

			// 一旦收到信号，则表明管理员希望退出程序，则先保存信息，然后退出
			os.Exit(0)
		}
	}
}

// 记录当前运行的Goroutine数量
func recordGoroutineNum() {
	// 处理内部未处理的异常，以免导致主线程退出，从而导致系统崩溃
	defer func() {
		if r := recover(); r != nil {
			logUtil.LogUnknownError(r)
		}
	}()

	for {
		time.Sleep(5 * time.Minute)

		// 记录当前运行的Goroutine数量
		logUtil.Log(fmt.Sprintf("NumGoroutine:%d", runtime.NumGoroutine()), logUtil.Debug, true)
	}
}

func main() {
	// 处理系统信号
	go signalProc()

	// 记录当前运行的Goroutine数量
	go recordGoroutineNum()

	// 获取数据库配置
	configObj := configBLL.GetConfig()

	// 设置Socket服务器配置，并启动服务器
	rpc.SetConfig(configObj.GetChatServerCenterRpcAddress(),
		wordBLL.HandleSensitiveWords,
		configBLL.HandleMessageLength,
		configBLL.IfClientCountReachMax,
		messageLogBLL.ChatMessageObjectChannel,
		onlineLogBLL.OnlineLogChannel,
		config.DEBUG)
	go rpc.Start(&wg)

	// 设置Web服务器配置，并启动服务器
	web.SetConfig(configObj.GetChatServerCenterWebAddress(),
		ipBLL.IsIPValid,
		apiLogBLL.APILogChannel,
		configBLL.IfRecordAPILog,
		config.DEBUG)
	go web.Start(&wg)

	// 阻塞等待，以免main线程退出
	wg.Wait()
}
