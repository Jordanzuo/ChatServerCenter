package manageCenterBLL

import (
	"encoding/json"
	"errors"
	"fmt"

	"sync"

	"github.com/Jordanzuo/ChatServerCenter/src/bll/configBLL"
	"github.com/Jordanzuo/ManageCenterModel_Go/partner"
	"github.com/Jordanzuo/ManageCenterModel_Go/returnObject"
	"github.com/Jordanzuo/ManageCenterModel_Go/server"
	"github.com/Jordanzuo/goutil/logUtil"
	"github.com/Jordanzuo/goutil/webUtil"
)

var (
	serverMap   = make(map[int]map[int]*server.Server, 128)
	serverMutex sync.RWMutex
)

// 重新加载服务器
func reloadServer() error {
	logUtil.Log("开始刷新服务器列表", logUtil.Debug, true)

	// 获取数据库配置
	configObj := configBLL.GetConfig()

	// 定义请求参数
	postDict := make(map[string]string)
	postDict["GroupType"] = configObj.GetGroupType()

	// 连接服务器，以获取数据
	url := fmt.Sprintf("%s/%s", configObj.GetManageCenterAPI(), configObj.GetServerListAPI())
	returnBytes, err := webUtil.PostWebData(url, postDict, nil)
	if err != nil {
		logUtil.Log(fmt.Sprintf("获取服务器列表出错，错误信息为：%s", err), logUtil.Error, true)
		return err
	}

	// 解析返回值
	returnObj := new(returnObject.ReturnObject)
	if err = json.Unmarshal(returnBytes, &returnObj); err != nil {
		logUtil.Log(fmt.Sprintf("获取服务器列表出错，反序列化返回值出错，错误信息为：%s", err), logUtil.Error, true)
		return err
	}

	// 判断返回状态是否成功
	if returnObj.Code != 0 {
		msg := fmt.Sprintf("获取服务器列表出错，返回状态：%d，信息为：%s", returnObj.Code, returnObj.Message)
		logUtil.Log(msg, logUtil.Error, true)
		return errors.New(msg)
	}

	// 解析Data
	tmpServerList := make([]*server.Server, 0, 1024)
	tmpServerMap := make(map[int]map[int]*server.Server, 128)
	if data, ok := returnObj.Data.(string); !ok {
		msg := "获取服务器列表出错，返回的数据不是string类型"
		logUtil.Log(msg, logUtil.Error, true)
		return errors.New(msg)
	} else {
		if err = json.Unmarshal([]byte(data), &tmpServerList); err != nil {
			logUtil.Log(fmt.Sprintf("获取服务器列表出错，反序列化数据出错，错误信息为：%s", err), logUtil.Error, true)
			return err
		}

		for _, item := range tmpServerList {
			if _, ok := tmpServerMap[item.PartnerId]; !ok {
				tmpServerMap[item.PartnerId] = make(map[int]*server.Server, 512)
			}

			tmpServerMap[item.PartnerId][item.Id] = item
		}
	}

	logUtil.Log(fmt.Sprintf("刷新服务器信息结束，服务器数量:%d", len(tmpServerList)), logUtil.Debug, true)

	// 赋值给最终的serverMap
	serverMutex.Lock()
	defer serverMutex.Unlock()
	serverMap = tmpServerMap

	return nil
}

// 根据合作商对象、服务器Id获取服务器对象
// partnerObj：合作商对象
// serverId：服务器Id
// 返回值：
// 服务器对象
// 是否存在
func getServer(partnerObj *partner.Partner, serverId int) (*server.Server, bool) {
	serverMutex.RLock()
	defer serverMutex.RUnlock()

	if subSserverMap, exists := serverMap[partnerObj.Id]; exists {
		if serverObj, exists := subSserverMap[serverId]; exists {
			return serverObj, true
		}
	}

	return nil, false
}
