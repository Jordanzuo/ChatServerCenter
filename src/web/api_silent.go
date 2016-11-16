package web

import (
	"net/http"
	"strconv"
	"time"

	"github.com/Jordanzuo/ChatServerCenter/src/bll/playerBLL"
	"github.com/Jordanzuo/ChatServerCenter/src/rpc"
	"github.com/Jordanzuo/ChatServerModel/src/centerResponseObject"
	"github.com/Jordanzuo/ChatServerModel/src/channelType"
	"github.com/Jordanzuo/ChatServerModel/src/transferObject"
)

func init() {
	registerAPI("silent", silentCallback, "PlayerId", "Type", "Duration")
}

func silentCallback(w http.ResponseWriter, r *http.Request) (responseObj *centerResponseObject.ResponseObject) {
	responseObj = centerResponseObject.NewResponseObject()

	// 解析数据
	playerId := r.Form["PlayerId"][0]
	type_str := r.Form["Type"][0]         // 0:查看；1:禁言； 2:解禁
	duration_str := r.Form["Duration"][0] // 单位：分钟

	// 类型转换
	var type_int int
	var duration_int int
	var err error

	if type_int, err = strconv.Atoi(type_str); err != nil {
		return responseObj.SetResultStatus(centerResponseObject.Con_APIDataError)
	}
	if duration_int, err = strconv.Atoi(duration_str); err != nil {
		return responseObj.SetResultStatus(centerResponseObject.Con_APIDataError)
	}

	// 验证类型是否正确(0:查看禁言状态 1:禁言 2:解禁)
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
		data := make(map[string]interface{}, 2)
		isInSilent, leftMinutes := playerObj.IsInSilent()
		data["Status"] = isInSilent
		if isInSilent {
			data["LeftMinutes"] = leftMinutes
		}
		responseObj.SetData(data)
	} else {
		// 修改禁言状态
		silentEndTime := time.Now()
		if type_int == 1 {
			if duration_int == 0 {
				silentEndTime = silentEndTime.AddDate(10, 0, 0)
			} else {
				silentEndTime = silentEndTime.Add(time.Duration(duration_int) * time.Minute)
			}
		}

		if err := playerBLL.UpdateSilentStatus(playerObj, silentEndTime); err != nil {
			return responseObj.SetResultStatus(centerResponseObject.Con_DataError)
		}

		// 构造聊天消息对象
		chatMessageObj := transferObject.NewChatMessageObject(channelType.System, "", "", nil)
		chatMessageObj.SetSilentInfo(playerId, silentEndTime)

		// 将数据发送到通道中
		rpc.ForwardObjectChannel <- transferObject.NewForwardObject(transferObject.Silent, chatMessageObj)
	}

	return
}
