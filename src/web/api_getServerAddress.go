package web

import (
	"net/http"

	"github.com/Jordanzuo/ChatServerCenter/src/rpc"
	"github.com/Jordanzuo/ChatServerModel/src/centerResponseObject"
	"github.com/Jordanzuo/goutil/logUtil"
)

func init() {
	registerAPI("getServerAddress", getServerAddressCallback)
}

func getServerAddressCallback(w http.ResponseWriter, r *http.Request) (responseObj *centerResponseObject.ResponseObject) {
	responseObj = centerResponseObject.NewResponseObject()

	data := make(map[string]interface{})
	if address, exists := rpc.GetAvailableServer(); exists {
		data["ChatServerUrl"] = address
	} else {
		data["ChatServerUrl"] = ""
		logUtil.Log("没有找到可用的聊天服务器", logUtil.Warn, true)
	}

	// 组装返回值
	responseObj.SetData(data)

	return
}
