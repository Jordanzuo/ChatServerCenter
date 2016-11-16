package web

import (
	"net/http"

	"github.com/Jordanzuo/ChatServerModel/src/centerResponseObject"
	"github.com/Jordanzuo/ChatServerModel/src/channelType"
)

func init() {
	registerAPI("pushCrossServer", pushCrossServerCallback, "Message")
}

func pushCrossServerCallback(w http.ResponseWriter, r *http.Request) (responseObj *centerResponseObject.ResponseObject) {
	responseObj = centerResponseObject.NewResponseObject()

	// 解析数据
	message := r.Form["Message"][0]

	// 推送数据
	push(channelType.CrossServer, "0", message, "", "")

	return
}
