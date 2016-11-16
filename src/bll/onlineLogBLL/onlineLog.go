package onlineLogBLL

import (
	"github.com/Jordanzuo/ChatServerCenter/src/dal/onlineLogDAL"
	"github.com/Jordanzuo/ChatServerModel/src/onlineLog"
	"github.com/Jordanzuo/goutil/logUtil"
	"time"
)

var (
	// 存放在线日志的通道
	OnlineLogChannel = make(chan *onlineLog.OnlineLog, 8)
)

func init() {
	go func() {
		// 处理内部未处理的异常，以免导致主线程退出，从而导致系统崩溃
		defer func() {
			if r := recover(); r != nil {
				logUtil.LogUnknownError(r)
			}
		}()

		for {
			select {
			case onlineLogObj := <-OnlineLogChannel:
				onlineLogDAL.Insert(onlineLogObj)
			default:
				// 如果channel中没有数据，则休眠5秒
				time.Sleep(5 * time.Second)
			}
		}
	}()
}
