package rpc

import (
	"time"

	"github.com/Jordanzuo/ChatServerModel/src/transferObject"
	"github.com/Jordanzuo/goutil/logUtil"
)

var (
	// 转发对象的通道
	ForwardObjectChannel = make(chan *transferObject.ForwardObject, 1024*100)
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
			case forwardObj := <-ForwardObjectChannel:
				go push(getClientList(), forwardObj)
			default:
				// 如果channel中没有数据，则休眠5毫秒
				time.Sleep(5 * time.Millisecond)
			}
		}
	}()
}
