package web

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/Jordanzuo/goutil/logUtil"
)

// 启动服务器
// wg：WaitGroup对象
func Start(wg *sync.WaitGroup) {
	defer func() {
		wg.Done()
	}()

	logUtil.Log("Web服务器开始监听...", logUtil.Info, true)

	// 启动Web服务器监听
	err := http.ListenAndServe(chatServerWebAddress, new(selfDefineMux))
	if err != nil {
		panic(fmt.Errorf("ListenAndServe失败，错误信息为：%s", err))
	}
}
