package web

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/Jordanzuo/ChatServerCenter/src/bll/manageCenterBLL"
	"github.com/Jordanzuo/ChatServerCenter/src/rpc"
	"github.com/Jordanzuo/ChatServerModel/src/centerResponseObject"
	"github.com/Jordanzuo/ChatServerModel/src/channelType"
	"github.com/Jordanzuo/ChatServerModel/src/transferObject"
)

func init() {
	registerAPI("push", pushCallback, "ServerGroupIds", "Message", "PlayerIds", "UnionId")
}

func pushCallback(w http.ResponseWriter, r *http.Request) (responseObj *centerResponseObject.ResponseObject) {
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
	push(channelType.System, serverGroupIds, message, playerIds, unionId)

	return
}

func checkServerGroup(serverGroupIds string) centerResponseObject.ResultStatus {
	if serverGroupIds != "0" {
		for _, serverGroupId_str := range strings.Split(serverGroupIds, ",") {
			if serverGroupId, err := strconv.Atoi(serverGroupId_str); err != nil {
				return centerResponseObject.Con_APIDataError
			} else {
				if _, exists := manageCenterBLL.GetServerGroupItem(serverGroupId); !exists {
					return centerResponseObject.Con_ServerGroupNotExist
				}
			}
		}
	}

	return centerResponseObject.Con_Success
}

func push(_channelType channelType.ChannelType, serverGroupIds, message, playerIds, unionId string) {
	var toPlayerIds []string = make([]string, 0, 8)
	var toUnionId string = ""

	// 判断是否发给指定玩家；如果不是发给指定玩家，则判断是否发给指定公会
	if playerIds != "" {
		for _, playerId := range strings.Split(playerIds, ",") {
			toPlayerIds = append(toPlayerIds, playerId)
		}
	} else {
		// 判断是否发给指定公会；如果不是发给指定公会，则发给所有玩家
		if unionId != "" {
			toUnionId = unionId
		}
	}

	// 构造聊天消息对象
	chatMessageObj := transferObject.NewChatMessageObject(_channelType, serverGroupIds, message, nil)
	if len(toPlayerIds) > 0 {
		chatMessageObj.SetToPlayerIds(toPlayerIds)
	}
	if toUnionId != "" {
		chatMessageObj.SetToUnionId(toUnionId)
	}

	// 将数据发送到通道中
	rpc.ForwardObjectChannel <- transferObject.NewForwardObject(transferObject.PushMessage, chatMessageObj)
}
