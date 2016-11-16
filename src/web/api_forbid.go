package web

import (
	"net/http"
	"strconv"

	"github.com/Jordanzuo/ChatServerCenter/src/bll/playerBLL"
	"github.com/Jordanzuo/ChatServerCenter/src/rpc"
	"github.com/Jordanzuo/ChatServerModel/src/centerResponseObject"
	"github.com/Jordanzuo/ChatServerModel/src/channelType"
	"github.com/Jordanzuo/ChatServerModel/src/transferObject"
)

func init() {
	registerAPI("forbid", forbidCallback, "PlayerId", "Type")
}

func forbidCallback(w http.ResponseWriter, r *http.Request) (responseObj *centerResponseObject.ResponseObject) {
	responseObj = centerResponseObject.NewResponseObject()

	// 解析数据
	playerId := r.Form["PlayerId"][0]
	type_str := r.Form["Type"][0] // 0:查看；1:封号； 2:解封

	// 类型转换
	var type_int int
	var err error
	if type_int, err = strconv.Atoi(type_str); err != nil {
		return responseObj.SetResultStatus(centerResponseObject.Con_APIDataError)
	}

	// 验证类型是否正确(0:查看封号状态 1:封号 2:解封)
	if type_int != 0 && type_int != 1 && type_int != 2 {
		return responseObj.SetResultStatus(centerResponseObject.Con_APIDataError)
	}

	// 判断玩家是否存在
	playerObj, exists, err := playerBLL.GetPlayer(playerId)
	if err != nil {
		return responseObj.SetResultStatus(centerResponseObject.Con_DataError)
	}
	if !exists {
		return responseObj.SetResultStatus(centerResponseObject.Con_PlayerNotExist)
	}

	// 判断是否为查询状态
	if type_int == 0 {
		responseObj.SetData(playerObj.IsForbidden)
	} else {
		// 修改封号状态
		if err := playerBLL.UpdateForbidStatus(playerObj, type_int == 1); err != nil {
			return responseObj.SetResultStatus(centerResponseObject.Con_DataError)
		}
	}

	// 构造聊天消息对象
	chatMessageObj := transferObject.NewChatMessageObject(channelType.System, "", "", nil)
	chatMessageObj.SetForbidPlayerId(playerId)

	// 将数据发送到通道中
	rpc.ForwardObjectChannel <- transferObject.NewForwardObject(transferObject.Forbid, chatMessageObj)

	return
}
