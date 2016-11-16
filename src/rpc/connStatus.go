package rpc

// 客户端连接状态
type ConnStatus int

const (
	// 打开状态
	con_Open ConnStatus = 1 + iota

	// 等待关闭
	con_WaitForClose

	// 已经关闭
	con_Close
)
