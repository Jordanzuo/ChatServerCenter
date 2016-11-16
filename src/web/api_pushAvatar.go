package web

import (
	"net/http"

	"github.com/Jordanzuo/ChatServerModel/src/centerResponseObject"
	"github.com/Jordanzuo/ChatServerModel/src/channelType"
)

func init() {
	registerAPI("pushAvatar", pushAvatarCallback, "ServerGroupIds", "Message", "PlayerIds", "UnionId")
}

func pushAvatarCallback(w http.ResponseWriter, r *http.Request) (responseObj *centerResponseObject.ResponseObject) {
	responseObj = centerResponseObject.NewResponseObject()

	// 解析数据
	serverGroupIds := r.Form["ServerGroupIds"][0]
	message := r.Form["Message"][0]
	playerIds := r.Form["PlayerIds"][0]
	unionId := r.Form["UnionId"][0]

	// 验证服务器组Id
	if serverGroupIds == "" {
		return responseObj.SetResultStatus(centerResponseObject.Con_APIDataError)
	}

	// 判断服务器组是否存在
	if rs := checkServerGroup(serverGroupIds); rs != centerResponseObject.Con_Success {
		return responseObj.SetResultStatus(rs)
	}

	// 推送数据
	push(channelType.Avatar, serverGroupIds, message, playerIds, unionId)

	return
}
