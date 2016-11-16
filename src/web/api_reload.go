package web

import (
	"net/http"

	"github.com/Jordanzuo/ChatServerCenter/src/bll/reloadBLL"
	"github.com/Jordanzuo/ChatServerCenter/src/rpc"
	"github.com/Jordanzuo/ChatServerModel/src/centerResponseObject"
	"github.com/Jordanzuo/ChatServerModel/src/channelType"
	"github.com/Jordanzuo/ChatServerModel/src/transferObject"
)

func init() {
	registerAPI("reload", reloadCallback)
}

func reloadCallback(w http.ResponseWriter, r *http.Request) (responseObj *centerResponseObject.ResponseObject) {
	responseObj = centerResponseObject.NewResponseObject()

	// 重新加载配置
	if errList := reloadBLL.Reload(); len(errList) > 0 {
		responseObj.SetResultStatus(centerResponseObject.Con_ReloadError)
	}

	// 构造聊天消息对象
	chatMessageObj := transferObject.NewChatMessageObject(channelType.System, "", "", nil)

	// 将数据发送到通道中
	rpc.ForwardObjectChannel <- transferObject.NewForwardObject(transferObject.Reload, chatMessageObj)

	return
}
