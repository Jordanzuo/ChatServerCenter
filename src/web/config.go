package web

import (
	"github.com/Jordanzuo/ChatServerModel/src/apiLog"
	"github.com/Jordanzuo/goutil/debugUtil"
)

var (
	// 服务器监听地址
	chatServerWebAddress string

	// 验证IP是否有效的方法
	isIPValidFunc func(string) bool

	// 存放API日志的通道
	apiLogChannel chan *apiLog.ApiLog

	// 是否记录API日志
	ifRecordAPILogFunc func() bool

	// 是否DEBUG模式
	debug bool
)

// 设置Config信息
// _chatServerWebAddress:服务器监听地址
// _isIPValid:验证IP是否有效的方法
// _apiLogChannel:存放API日志的通道
// _ifRecordAPILog:是否记录API日志
// _debug:是否DEBUG模式
func SetConfig(_chatServerWebAddress string, _isIPValidFunc func(string) bool, _apiLogChannel chan *apiLog.ApiLog, _ifRecordAPILogFunc func() bool, _debug bool) {
	chatServerWebAddress = _chatServerWebAddress
	isIPValidFunc = _isIPValidFunc
	apiLogChannel = _apiLogChannel
	ifRecordAPILogFunc = _ifRecordAPILogFunc
	debug = _debug

	debugUtil.Println("chatServerWebAddress:", chatServerWebAddress)
	debugUtil.Println("debug:", debug)
}
