package rpc

import (
	"github.com/Jordanzuo/ChatServerModel/src/onlineLog"
	"github.com/Jordanzuo/ChatServerModel/src/transferObject"
)

var (
	// 服务器监听地址
	chatServerCenterRpcAddress string

	// 处理敏感词汇
	handleSensitiveWords func(string) string

	// 处理消息长度
	handleMessageLength func(string) string

	// 判断单个ChatServer是否达到最大承载上限
	ifClientCountReachMax func(int) bool

	// 聊天消息通道
	chatMessageObjectChannel chan *transferObject.ChatMessageObject

	// 在线日志通道
	onlineLogChannel chan *onlineLog.OnlineLog

	// 是否DEBUG模式
	debug bool
)

func SetConfig(_chatServerCenterRpcAddress string,
	_handleSensitiveWords func(string) string,
	_handleMessageLength func(string) string,
	_ifClientCountReachMax func(int) bool,
	_chatMessageObjectChannel chan *transferObject.ChatMessageObject,
	_onlineLogChannel chan *onlineLog.OnlineLog,
	_debug bool) {
	chatServerCenterRpcAddress = _chatServerCenterRpcAddress
	handleSensitiveWords = _handleSensitiveWords
	handleMessageLength = _handleMessageLength
	ifClientCountReachMax = _ifClientCountReachMax
	chatMessageObjectChannel = _chatMessageObjectChannel
	onlineLogChannel = _onlineLogChannel
	debug = _debug
}
