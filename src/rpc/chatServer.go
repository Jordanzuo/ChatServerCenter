package rpc

import (
	"fmt"
	"sort"
)

// 为客户端提供服务的chatServer对象
type chatServer struct {
	// ChatServer唯一标识（监听地址）
	id string

	// 连上ChatServer的客户端数量
	clientCount int

	// 连上ChatServer的玩家数量
	playerCount int
}

// 更新客户端数量
// clientCount：客户端数量
// playerCount：玩家数量
// 返回值：无
func (server *chatServer) updateClientAndPlayerCount(clientCount int, playerCount int) {
	server.clientCount = clientCount
	server.playerCount = playerCount
}

// 格式化字符串
// 返回值：
// 格式化后的字符串
func (server *chatServer) String() string {
	return fmt.Sprintf("Id:%s, ClientCount:%d, PlayerCount:%d", server.id, server.clientCount, server.playerCount)
}

// 创建新的chatServer对象
// _id：唯一标识
// 返回值：chatServer指针
func newChatServer(_id string) *chatServer {
	return &chatServer{
		id:          _id,
		clientCount: 0,
		playerCount: 0,
	}
}

// 获取适合的服务器的监听地址
// 返回值：
// 服务器监听地址
// 是否存在对应的服务器
func GetAvailableServer() (string, bool) {
	clientList := make([]*client, 0, 4)

	// 获得Client列表
	tmpClientList := getClientList()
	for _, clientObj := range tmpClientList {
		if clientObj.chatServer != nil {
			clientList = append(clientList, clientObj)
		}
	}

	if len(clientList) == 0 {
		return "", false
	}

	// 升序排序
	sortOfClientList := sortOfClientList(clientList)
	sort.Sort(sortOfClientList)

	// 判断是否最小的已经达到了最大值
	if ifClientCountReachMax(sortOfClientList[0].chatServer.clientCount) {
		return "", false
	}

	return sortOfClientList[0].chatServer.id, true
}
