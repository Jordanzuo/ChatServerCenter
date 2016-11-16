package rpc

import (
	"encoding/json"
	"fmt"

	"github.com/Jordanzuo/ChatServerModel/src/centerRequestObject"
	"github.com/Jordanzuo/ChatServerModel/src/centerResponseObject"
	"github.com/Jordanzuo/ChatServerModel/src/transferObject"
	"github.com/Jordanzuo/goutil/logUtil"
)

var (
	// 所有对外提供的方法列表
	funcMap = make(map[transferObject.TransferType]*requestFunc)
)

func init() {
	funcMap[transferObject.Login] = newRequestFunc("Login", login, 1)
	funcMap[transferObject.Forward] = newRequestFunc("Forward", forward, 1)
	funcMap[transferObject.UpdateClientAndPlayerCount] = newRequestFunc("UpdateClientAndPlayerCount", updateClientAndPlayerCount, 2)
}

// 处理请求
// clientObj：对应的客户端对象
// id：客户端请求唯一标识
// request：请求内容字节数组(json格式)
// 返回值：无
func handleRequest(clientObj *client, id int32, request []byte) {
	responseObj := centerResponseObject.NewResponseObject()

	// 提取请求内容
	requestObj := new(centerRequestObject.RequestObject)
	if err := json.Unmarshal(request, requestObj); err != nil {
		logUtil.Log(fmt.Sprintf("反序列化%s出错，错误信息为：%s", string(request), err), logUtil.Error, true)
		return
	}

	// 对requestObj的属性Id赋值
	requestObj.Id = id

	// 查找方法
	if requestFuncObj, exists := funcMap[transferObject.TransferType(requestObj.MethodName)]; !exists {
		responseResult(clientObj, requestObj, responseObj.SetResultStatus(centerResponseObject.Con_MethodNotDefined), Con_HighPriority)
		return
	} else {
		// 检测参数数量
		if rs := requestFuncObj.checkParamCount(requestObj.Parameters); rs != centerResponseObject.Con_Success {
			responseResult(clientObj, requestObj, responseObj.SetResultStatus(rs), Con_HighPriority)
			return
		}

		// 调用方法
		responseObj = requestFuncObj.funcDefinition(clientObj, requestObj.Parameters)

		// 输出结果
		responseResult(clientObj, requestObj, responseObj, Con_HighPriority)
	}
}

// SocketServer登录
// clientObj：客户端对象
// parameters：参数
// 返回值：响应对象
func login(clientObj *client, parameters []interface{}) *centerResponseObject.ResponseObject {
	responseObj := centerResponseObject.NewResponseObject()

	// 获取ServerId
	if serverId, ok := parameters[0].(string); !ok {
		return responseObj.SetResultStatus(centerResponseObject.Con_ParamTypeError)
	} else {
		clientObj.login(serverId)
	}

	return responseObj
}

// 转发消息
// clientObj：客户端对象
// parameters：参数
// 返回值：响应对象
func forward(clientObj *client, parameters []interface{}) *centerResponseObject.ResponseObject {
	responseObj := centerResponseObject.NewResponseObject()

	if msg, ok := parameters[0].(string); !ok {
		return responseObj.SetResultStatus(centerResponseObject.Con_ParamTypeError)
	} else {
		// 解析数据
		chatMessageObj := new(transferObject.ChatMessageObject)
		if err := json.Unmarshal([]byte(msg), chatMessageObj); err != nil {
			logUtil.Log(fmt.Sprintf("反序列化%s为ChatMessageObject出错，错误信息为：%s", msg, err), logUtil.Error, true)
			return responseObj.SetResultStatus(centerResponseObject.Con_DataError)
		}

		// 将聊天消息添加到日志通道中(保存原始数据)
		chatMessageObjectChannel <- chatMessageObj

		// 处理敏感词汇
		chatMessageObj.Message = handleSensitiveWords(chatMessageObj.Message)

		// 处理消息长度
		chatMessageObj.Message = handleMessageLength(chatMessageObj.Message)

		// 将数据发送到通道中
		ForwardObjectChannel <- transferObject.NewForwardObject(transferObject.ChatMessage, chatMessageObj)
	}

	return responseObj
}

// 更新客户端和玩家数量
// clientObj：客户端对象
// parameters：参数
// 返回值：响应对象
func updateClientAndPlayerCount(clientObj *client, parameters []interface{}) *centerResponseObject.ResponseObject {
	responseObj := centerResponseObject.NewResponseObject()
	clientCount := 0
	playerCount := 0

	if clientCount_float64, ok := parameters[0].(float64); !ok {
		return responseObj.SetResultStatus(centerResponseObject.Con_ParamTypeError)
	} else {
		clientCount = int(clientCount_float64)
	}

	if playerCount_float64, ok := parameters[1].(float64); !ok {
		return responseObj.SetResultStatus(centerResponseObject.Con_ParamTypeError)
	} else {
		playerCount = int(playerCount_float64)
	}

	clientObj.updateClientAndPlayerCount(clientCount, playerCount)

	return responseObj
}
