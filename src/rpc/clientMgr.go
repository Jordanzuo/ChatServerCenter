package rpc

import (
	"fmt"
	"sync"
	"time"

	"github.com/Jordanzuo/ChatServerModel/src/onlineLog"
	"github.com/Jordanzuo/goutil/logUtil"
)

var (
	// 客户端连接集合
	clientMap = make(map[int32]*client)

	// 锁对象
	mutex sync.RWMutex
)

func init() {
	// 处理client过期
	go func() {
		// 处理内部未处理的异常，以免导致主线程退出，从而导致系统崩溃
		defer func() {
			if r := recover(); r != nil {
				logUtil.LogUnknownError(r)
			}
		}()

		for {
			// 因为刚开始时不存在过期，所以先暂停5分钟
			time.Sleep(5 * time.Minute)

			// 获取客户端连接列表
			clientList := getClientList()

			// 记录日志
			logUtil.Log(fmt.Sprintf("当前客户端数量为：%d，准备清理过期的客户端", len(clientList)), logUtil.Debug, true)

			onlineLogList := make([]*onlineLog.OnlineLog, 0, 8)
			totalCount := 0
			sid := 1
			for _, clientObj := range clientList {
				// 如果有ChatServer连接，则记录在线日志
				if clientObj.chatServer != nil {
					onlineLogObj := onlineLog.NewOnlineLog(sid, clientObj.chatServer.id, clientObj.chatServer.clientCount, clientObj.chatServer.playerCount)
					onlineLogList = append(onlineLogList, onlineLogObj)
					totalCount += onlineLogObj.GetClientCount()
					sid = sid + 1
				}

				if clientObj.expired() {
					logUtil.Log(fmt.Sprintf("客户端超时被断开，对应的信息为：%s", clientObj.String()), logUtil.Debug, true)
					clientObj.quit()
				} else {
					logUtil.Log(fmt.Sprintf("客户端的信息为：%s", clientObj.String()), logUtil.Debug, true)
				}
			}

			// 重新计算所有的在线人数
			for i := 0; i < len(onlineLogList); i++ {
				onlineLogList[i].SetTotalCount(totalCount)
			}

			// 将onlineLogList写入通道
			for _, onlineLogObj := range onlineLogList {
				onlineLogChannel <- onlineLogObj
			}
		}
	}()
}

// 获取客户端列表
// 返回值：
// 客户端列表
func getClientList() (clientList []*client) {
	mutex.RLock()
	defer mutex.RUnlock()

	for _, value := range clientMap {
		clientList = append(clientList, value)
	}

	return
}

// 根据玩家对象获取对应的客户端对象
// id：客户端Id
// 返回值：
// 客户端对象
// 是否存在客户端对象
func getClientById(id int32) (*client, bool) {
	mutex.RLock()
	defer mutex.RUnlock()

	if clientObj, exists := clientMap[id]; exists {
		return clientObj, true
	}

	return nil, false
}

// 根据服务器Id获取客户端对象
// serverId：服务器Id
// 返回值：
// 客户端对象
// 是否存在
func getClientByServerId(serverId string) (*client, bool) {
	mutex.RLock()
	defer mutex.RUnlock()

	for _, value := range clientMap {
		if value.chatServer.id == serverId {
			return value, true
		}
	}

	return nil, false
}

// 根据服务器Ip来判断服务器是否存在
// socketServerAddress：服务器的地址
func IfServerExists(socketServerAddress string) bool {
	_, exists := getClientByServerId(socketServerAddress)

	return exists
}

// 注册客户端对象
// clientObj：客户端对象
// 返回值：
// 无
func registerClient(clientObj *client) {
	mutex.Lock()
	defer mutex.Unlock()

	clientMap[clientObj.id] = clientObj
	logUtil.Log(fmt.Sprintf("收到新的客户端连接，详细信息为：%s", clientObj.String()), logUtil.Info, true)
}

// 取消客户端注册
// clientObj：客户端对象
// 返回值：
// 无
func unregisterClient(clientObj *client) {
	mutex.Lock()
	defer mutex.Unlock()

	delete(clientMap, clientObj.id)
	logUtil.Log(fmt.Sprintf("收到客户端断开连接，详细信息为：%s", clientObj.String()), logUtil.Info, true)
}
