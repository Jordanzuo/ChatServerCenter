package rpc

import (
	"fmt"
	"net"
	"sync"
	"time"

	"github.com/Jordanzuo/goutil/logUtil"
)

// 处理客户端收到的数据
// clientObj：客户端对象
func handleReceiveData(clientObj *client) {
	for {
		// 获取有效的消息
		id, message, exists := clientObj.getReceiveData()
		if !exists {
			break
		}

		// 处理数据，如果长度为0则表示心跳包；否则处理请求内容
		if len(message) == 0 {
			continue
		} else {
			handleRequest(clientObj, id, message)
		}
	}
}

// 处理需要客户端发送的数据
// clientObj：客户端对象
func handleSendData(clientObj *client) {
	// 处理内部未处理的异常，以免导致主线程退出，从而导致系统崩溃
	defer func() {
		if r := recover(); r != nil {
			logUtil.LogUnknownError(r)
		}
	}()

	for {
		//连接是否断开
		if clientObj.getConnStatus() == con_Close {
			break
		}

		//是否是最后一条消息
		connStatus := clientObj.getConnStatus()

		// 是否被处理过
		handled := false

		// 优先处理高优先级的数据，如果发送出现错误，表示连接已经断开，则退出方法；如果没有待处理的数据，则退出循环
		for {
			if sendDataItemObj, exists := clientObj.getSendData(); exists {
				handled = true
				if err := clientObj.sendMessage(sendDataItemObj); err != nil {
					return
				}
			} else {
				break
			}
		}

		// 当没有高优先级的数据时才处理低优先级的数据，如果发送出现错误，表示连接已经断开，则退出方法；
		if sendDataItemObj, exists := clientObj.getSendData_LowPriority(); exists {
			handled = true
			if err := clientObj.sendMessage(sendDataItemObj); err != nil {
				return
			}
		}

		// 如果本轮没有被处理过，则休眠5ms
		if !handled {
			time.Sleep(5 * time.Millisecond)
			if connStatus == con_WaitForClose {
				clientObj.setConnStatus(con_Close)
			}
		}
	}
}

// 处理客户端连接
// conn：客户端连接对象
func handleConn(conn net.Conn) {
	// 处理内部未处理的异常，以免导致主线程退出，从而导致系统崩溃
	defer func() {
		if r := recover(); r != nil {
			logUtil.LogUnknownError(r)
		}
	}()

	// 创建客户端对象
	clientObj := newClient(conn)

	// 将客户端对象注册到客户端列表中
	registerClient(clientObj)

	// 启动一个Goroutine来发送数据
	go handleSendData(clientObj)

	// 将客户端对象从注册列表中移除
	defer func() {
		unregisterClient(clientObj)
		clientObj.quit()
	}()

	// 无限循环，不断地读取数据，解析数据，处理数据
	for {
		// 先读取数据，每次读取1024个字节
		readBytes := make([]byte, 1024)

		// Read方法会阻塞，所以不用考虑异步的方式
		n, err := conn.Read(readBytes)
		if err != nil {
			logUtil.Log(fmt.Sprintf("客户端：%s读取数据出错，错误信息为：%s", clientObj.String(), err), logUtil.Error, true)
			break
		}

		// 将读取到的数据追加到已获得的数据的末尾
		clientObj.appendReceiveData(readBytes[:n])

		// 处理数据
		handleReceiveData(clientObj)
	}
}

// 启动服务器
// wg：WaitGroup
func Start(wg *sync.WaitGroup) {
	defer func() {
		wg.Done()
	}()

	logUtil.Log("Socket服务器开始监听...", logUtil.Info, true)

	// 监听指定的端口
	listener, err := net.Listen("tcp", chatServerCenterRpcAddress)
	if err != nil {
		panic(fmt.Errorf("Listen Error: %s", err))
	} else {
		msg := fmt.Sprintf("Got listener for client. (local address: %s)", listener.Addr())

		// 记录和显示日志，并且判断是否需要退出
		logUtil.Log(msg, logUtil.Info, true)
		fmt.Println(msg)
	}

	for {
		// 阻塞直至新连接到来
		conn, err := listener.Accept()
		if err != nil {
			logUtil.Log(fmt.Sprintf("Accept Error: %s", err), logUtil.Error, true)
			continue
		}

		// 启动一个新协程来处理链接(每个客户端对应一个协程)
		go handleConn(conn)
	}
}
