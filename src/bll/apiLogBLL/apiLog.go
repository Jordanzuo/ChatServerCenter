package apiLogBLL

import (
	"github.com/Jordanzuo/ChatServerCenter/src/dal/apiLogDAL"
	"github.com/Jordanzuo/ChatServerModel/src/apiLog"
	"github.com/Jordanzuo/goutil/logUtil"
	"time"
)

var (
	// 存放API日志的通道
	APILogChannel = make(chan *apiLog.ApiLog, 1024*100)
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
			case apiLogObj := <-APILogChannel:
				go apiLogDAL.Insert(apiLogObj)
			default:
				// 如果channel中没有数据，则休眠5毫秒
				time.Sleep(5 * time.Millisecond)
			}
		}
	}()
}
