package rpc

import (
	"encoding/binary"
	"fmt"
	"net"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/Jordanzuo/goutil/fileUtil"
	"github.com/Jordanzuo/goutil/intAndBytesUtil"
	"github.com/Jordanzuo/goutil/logUtil"
	"github.com/Jordanzuo/goutil/timeUtil"
)

const (
	// 包头的长度
	con_HEADER_LENGTH = 4

	// 定义请求、响应数据的前缀的长度
	con_ID_LENGTH = 4

	// 客户端失效的秒数
	con_CLIENT_EXPIRE_SECONDS = 300 * time.Second
)

var (
	// 字节的大小端顺序
	byterOrder = binary.LittleEndian

	// 全局客户端的id，从1开始进行自增
	globalClientId int32 = 0
)

// 定义客户端对象，以实现对客户端连接的封装
type client struct {
	// 唯一标识
	id int32

	// 客户端连接对象
	conn net.Conn

	// 客户端连接状态
	connStatus ConnStatus

	// 接收到的消息内容
	receiveData []byte

	// 待发送的数据
	sendData []*sendDataItem

	// 低优先级的待发送的数据
	sendData_LowPriority []*sendDataItem

	// 上次活跃时间
	activeTime time.Time

	// 为客户端提供服务的ChatServer对象
	*chatServer

	// 锁对象（用于控制对sendDatap、sendData_LowPriority的并发访问；receiveData不需要，因为是同步访问）
	mutex sync.Mutex
}

// 获取连接状态
func (clientObj *client) getConnStatus() ConnStatus {
	return clientObj.connStatus
}

// 设置连接状态
func (clientObj *client) setConnStatus(status ConnStatus) {
	clientObj.connStatus = status
}

// 获取待发送的数据
// 返回值：
// 待发送数据项
// 是否含有有效数据
func (clientObj *client) getSendData() (sendDataItemObj *sendDataItem, exists bool) {
	clientObj.mutex.Lock()
	defer clientObj.mutex.Unlock()

	// 如果没有数据则直接返回
	if len(clientObj.sendData) == 0 {
		return
	}

	// 取出第一条数据,并为返回值赋值
	sendDataItemObj = clientObj.sendData[0]
	exists = true

	// 删除已经取出的数据
	clientObj.sendData = clientObj.sendData[1:]

	return
}

// 获取低优先级待发送的数据
// 返回值：
// 待发送数据项
// 是否含有有效数据
func (clientObj *client) getSendData_LowPriority() (sendDataItemObj *sendDataItem, exists bool) {
	clientObj.mutex.Lock()
	defer clientObj.mutex.Unlock()

	// 如果没有数据则直接返回
	if len(clientObj.sendData_LowPriority) == 0 {
		return
	}

	// 取出第一条数据,并为返回值赋值
	sendDataItemObj = clientObj.sendData_LowPriority[0]
	exists = true

	// 删除已经取出的数据
	clientObj.sendData_LowPriority = clientObj.sendData_LowPriority[1:]

	return
}

// 追加发送的数据
// sendDataItemObj:待发送数据项
// priority:优先级
// 返回值：无
func (clientObj *client) appendSendData(sendDataItemObj *sendDataItem, priority Priority) {
	clientObj.mutex.Lock()
	defer clientObj.mutex.Unlock()

	if priority == Con_LowPriority {
		clientObj.sendData_LowPriority = append(clientObj.sendData_LowPriority, sendDataItemObj)
	} else {
		clientObj.sendData = append(clientObj.sendData, sendDataItemObj)
	}
}

// 获取远程地址（IP_Port）
func (clientObj *client) GetRemoteAddr() string {
	items := strings.Split(clientObj.conn.RemoteAddr().String(), ":")

	return fmt.Sprintf("%s_%s", items[0], items[1])
}

// 获取远程地址（IP）
func (clientObj *client) GetRemoteShortAddr() string {
	items := strings.Split(clientObj.conn.RemoteAddr().String(), ":")

	return items[0]
}

// 获取client的Id
// 返回值：
// Id
func (clientObj *client) GetId() int32 {
	return clientObj.id
}

// 获取接收到的数据
// 返回值：
// 消息对应客户端的唯一标识
// 消息内容
// 是否含有有效数据
func (clientObj *client) getReceiveData() (id int32, message []byte, exists bool) {
	// 判断是否包含头部信息
	if len(clientObj.receiveData) < con_HEADER_LENGTH {
		return
	}

	// 获取头部信息
	header := clientObj.receiveData[:con_HEADER_LENGTH]

	// 将头部数据转换为内容的长度
	contentLength := intAndBytesUtil.BytesToInt32(header, byterOrder)

	// 判断长度是否满足
	if len(clientObj.receiveData) < con_HEADER_LENGTH+int(contentLength) {
		return
	}

	// 运行到此处，标识有数据
	exists = true

	// 提取消息内容
	content := clientObj.receiveData[con_HEADER_LENGTH : con_HEADER_LENGTH+contentLength]

	// 将对应的数据截断，以得到新的内容
	clientObj.receiveData = clientObj.receiveData[con_HEADER_LENGTH+contentLength:]

	// 判断是否为心跳包，如果是心跳包，则不解析，直接返回
	if contentLength == 0 || len(content) == 0 {
		return
	}

	// 判断内容的长度是否足够
	if len(content) < con_ID_LENGTH {
		logUtil.Log(fmt.Sprintf("内容数据不正确；con_ID_LENGTH=%d,len(content)=%d", con_ID_LENGTH, len(content)), logUtil.Warn, true)
	}

	// 将内容分隔为2部分
	idBytes, content := content[:con_ID_LENGTH], content[con_ID_LENGTH:]

	// 提取id、message
	id = intAndBytesUtil.BytesToInt32(idBytes, byterOrder)
	message = content

	return
}

// 追加接收到的数据
// receiveData：接收到的数据
// 返回值：无
func (clientObj *client) appendReceiveData(receiveData []byte) {
	clientObj.receiveData = append(clientObj.receiveData, receiveData...)
	clientObj.activeTime = time.Now()
}

// 判断客户端是否超时
// 返回值：
// 是否超时
func (clientObj *client) expired() bool {
	return time.Now().Unix() > clientObj.activeTime.Add(con_CLIENT_EXPIRE_SECONDS).Unix()
}

// 发送消息
// id：需要添加到内容前的数据
// sendDataItemObj：待发送数据项
// 返回值：
// 错误对象
func (clientObj *client) sendMessage(sendDataItemObj *sendDataItem) error {
	clientObj.WriteLog(fmt.Sprintf("发送消息,%d:%s", sendDataItemObj.id, string(sendDataItemObj.data)))

	idBytes := intAndBytesUtil.Int32ToBytes(sendDataItemObj.id, byterOrder)

	// 将idByte和内容合并
	content := append(idBytes, sendDataItemObj.data...)

	// 获得数组的长度
	contentLength := len(content)

	// 将长度转化为字节数组
	header := intAndBytesUtil.Int32ToBytes(int32(contentLength), byterOrder)

	// 将头部与内容组合在一起
	message := append(header, content...)

	// 发送消息
	_, err := clientObj.conn.Write(message)
	if err != nil {
		clientObj.WriteLog(fmt.Sprintf("发送消息,%d:%s,出现错误：%s", sendDataItemObj.id, string(sendDataItemObj.data), err))
	}

	return err
}

// 登录
// serverId：SocketServer的Id
// 返回值：无
func (clientObj *client) login(serverId string) {
	clientObj.chatServer = newChatServer(serverId)
}

// 退出
// 返回值：无
func (clientObj *client) quit() {
	clientObj.conn.Close()
	clientObj.setConnStatus(con_Close)
}

// 格式化客户端对象
// 返回值：
// 格式化的字符串
func (clientObj *client) String() string {
	if clientObj.chatServer == nil {
		return fmt.Sprintf("{Id:%d, RemoteAddr:%d, ActiveTime:%s, ChatServer尚未登录}", clientObj.id, clientObj.GetRemoteAddr(), timeUtil.Format(clientObj.activeTime, "yyyy-MM-dd HH:mm:ss"))
	} else {
		return fmt.Sprintf("{Id:%d, RemoteAddr:%d, ActiveTime:%s, ChatServer:%s}", clientObj.id, clientObj.GetRemoteAddr(), timeUtil.Format(clientObj.activeTime, "yyyy-MM-dd HH:mm:ss"), clientObj.chatServer.String())
	}
}

// 记录日志
// log：日志内容
func (clientObj *client) WriteLog(log string) {
	if debug {
		fileUtil.WriteFile("Log", clientObj.GetRemoteAddr(), true,
			timeUtil.Format(time.Now(), "yyyy-MM-dd HH:mm:ss"),
			" ",
			fmt.Sprintf("client:%s", clientObj.String()),
			" ",
			log,
			"\r\n",
			"\r\n",
		)
	}
}

// 新建客户端对象
// _conn：连接对象
// 返回值：客户端对象的指针
func newClient(_conn net.Conn) *client {
	// 生成Id的方法
	generateId := func() int32 {
		atomic.AddInt32(&globalClientId, 1)
		return globalClientId
	}

	return &client{
		id:                   generateId(),
		conn:                 _conn,
		connStatus:           con_Open,
		receiveData:          make([]byte, 0, 1024),
		sendData:             make([]*sendDataItem, 0, 16),
		sendData_LowPriority: make([]*sendDataItem, 0, 16),
		activeTime:           time.Now(),
		// 其他保持默认值
	}
}

/*-----------------------------------------------------------------------------------------------------------------------------------------------------------------*/

// 定义用于排序的新类型
type sortOfClientList []*client

// 获取数据长度
func (list sortOfClientList) Len() int {
	return len(list)
}

// 按照玩家人数来进行升序排序
func (list sortOfClientList) Less(i, j int) bool {
	if list[i].playerCount < list[j].playerCount {
		return true
	}

	return false
}

// 交换数据
func (list sortOfClientList) Swap(i, j int) {
	list[i], list[j] = list[j], list[i]
}

/*-----------------------------------------------------------------------------------------------------------------------------------------------------------------*/
