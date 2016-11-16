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
	"github.com/Jordanzuo/ManageCenterModel_Go/serverGroup"
	"github.com/Jordanzuo/goutil/logUtil"
	"github.com/Jordanzuo/goutil/webUtil"
)

var (
	serverGroupMap           = make(map[int]*serverGroup.ServerGroup, 512)
	serverGroupChangeFuncMap = make(map[string]func(map[int]*serverGroup.ServerGroup))
	serverGroupMutex         sync.RWMutex
)

// 重新加载服务器组
func reloadServerGroup() error {
	logUtil.Log("开始刷新服务器组信息", logUtil.Debug, true)

	// 获取数据库配置
	configObj := configBLL.GetConfig()

	// 定义请求参数
	postDict := make(map[string]string)
	postDict["GroupType"] = configObj.GetGroupType()

	// 连接服务器，以获取数据
	url := fmt.Sprintf("%s/%s", configObj.GetManageCenterAPI(), configObj.GetServerGroupListAPI())
	returnBytes, err := webUtil.PostWebData(url, postDict, nil)
	if err != nil {
		logUtil.Log(fmt.Sprintf("获取服务器组列表出错，错误信息为：%s", err), logUtil.Error, true)
		return err
	}

	// 解析返回值
	returnObj := new(returnObject.ReturnObject)
	if err = json.Unmarshal(returnBytes, &returnObj); err != nil {
		logUtil.Log(fmt.Sprintf("获取服务器组列表出错，反序列化返回值出错，错误信息为：%s", err), logUtil.Error, true)
		return err
	}

	// 判断返回状态是否为成功
	if returnObj.Code != 0 {
		msg := fmt.Sprintf("获取服务器组列表出错，返回状态：%d，信息为：%s", returnObj.Code, returnObj.Message)
		logUtil.Log(msg, logUtil.Error, true)
		return errors.New(msg)
	}

	// 解析Data
	tmpServerGroupList := make([]*serverGroup.ServerGroup, 0, 512)
	tmpServerGroupMap := make(map[int]*serverGroup.ServerGroup, 512)
	if data, ok := returnObj.Data.(string); !ok {
		msg := "获取服务器组列表出错，返回的数据不是string类型"
		logUtil.Log(msg, logUtil.Error, true)
		return errors.New(msg)
	} else {
		if err = json.Unmarshal([]byte(data), &tmpServerGroupList); err != nil {
			logUtil.Log(fmt.Sprintf("获取服务器组列表出错，反序列化数据出错，错误信息为：%s", err), logUtil.Error, true)
			return err
		}

		for _, item := range tmpServerGroupList {
			tmpServerGroupMap[item.Id] = item
		}
	}

	logUtil.Log(fmt.Sprintf("刷新服务器组信息结束,服务器组的总数量:%d", len(tmpServerGroupList)), logUtil.Debug, true)

	// 判断服务器组是否有变化，如果有变化
	if isServerGroupChanged(tmpServerGroupMap) {
		// 则触发服务器组变化的方法
		triggerServerGroupChangeFunc(tmpServerGroupMap)

		// 赋值给最终的ServerGroupMap
		serverGroupMutex.Lock()
		defer serverGroupMutex.Unlock()
		serverGroupMap = tmpServerGroupMap
	}

	return nil
}

// 判断服务器组是否有变化
// tmpServerGroupMap：临时服务器组集合
func isServerGroupChanged(tmpServerGroupMap map[int]*serverGroup.ServerGroup) bool {
	serverGroupMutex.Lock()
	defer serverGroupMutex.Unlock()

	// 判断数量是否相同
	if len(serverGroupMap) != len(tmpServerGroupMap) {
		return true
	}

	// 判断在serverGroupMap、tmpServerGroupMap中的内容是否完全相同
	for key, _ := range serverGroupMap {
		if _, exists := tmpServerGroupMap[key]; !exists {
			return true
		}
	}

	for key, _ := range tmpServerGroupMap {
		if _, exists := serverGroupMap[key]; !exists {
			return true
		}
	}

	return false
}

// 触发服务器组变化的方法
func triggerServerGroupChangeFunc(tmpServerGroupMap map[int]*serverGroup.ServerGroup) {
	// 如果有注册服务器组变化的方法
	if len(serverGroupChangeFuncMap) > 0 {
		logUtil.Log(fmt.Sprintf("总共有%d个注册方法", len(serverGroupChangeFuncMap)), logUtil.Debug, true)
		for funcName, serverGroupChangeFunc := range serverGroupChangeFuncMap {
			logUtil.Log(fmt.Sprintf("开始触发方法：%s", funcName), logUtil.Debug, true)
			serverGroupChangeFunc(tmpServerGroupMap)
			logUtil.Log(fmt.Sprintf("触发方法：%s执行结束", funcName), logUtil.Debug, true)
		}
	}
}

// 注册服务器组变化方法
// funcName：方法名称
// serverGroupChangeFunc：服务器组变化方法
func RegisterServerGroupChangeFunc(funcName string, serverGroupChangeFunc func(map[int]*serverGroup.ServerGroup)) {
	serverGroupChangeFuncMap[funcName] = serverGroupChangeFunc
}

// 获取服务器组列表
// 返回值：
// 服务器组列表
func GetServerGroupMap() (tmpServerGroupMap map[int]*serverGroup.ServerGroup) {
	serverGroupMutex.RLock()
	defer serverGroupMutex.RUnlock()

	tmpServerGroupMap = make(map[int]*serverGroup.ServerGroup, 128)
	for key, serverGroupItem := range serverGroupMap {
		tmpServerGroupMap[key] = serverGroupItem
	}

	return
}

// 获取服务器组项
// groupId：Id
// 返回值：
// 服务器组对象
// 是否存在
func GetServerGroupItem(groupId int) (*serverGroup.ServerGroup, bool) {
	serverGroupMutex.RLock()
	defer serverGroupMutex.RUnlock()

	if serverGroupObj, exists := serverGroupMap[groupId]; exists {
		return serverGroupObj, exists
	}

	return nil, false
}

// 根据合作商Id、服务器Id获取服务器组对象
// partnerId：合作商Id
// serverId：服务器Id
// 返回值：
// 服务器组对象
// 服务器对象
// 是否存在
func GetServerGroup(partnerId, serverId int) (serverGroupObj *serverGroup.ServerGroup, serverObj *server.Server, exists bool) {
	var partnerObj *partner.Partner

	// 获取合作商对象
	partnerObj, exists = getPartner(partnerId)
	if !exists {
		return
	}

	// 获取服务器对象
	serverObj, exists = getServer(partnerObj, serverId)
	if !exists {
		return
	}

	// 获取服务器组对象
	serverGroupObj, exists = GetServerGroupItem(serverObj.GroupId)

	return
}

// 判断IP是否属于某个服务器组
// ip：指定IP地址
// 返回值：
// 是否属于某个服务器组
func IsIpBelongToServerGroup(ip string) bool {
	// 获取服务器组对象
	serverGroupMutex.RLock()
	defer serverGroupMutex.RUnlock()

	for _, item := range serverGroupMap {
		if item.Ip == ip {
			return true
		}
	}

	return false
}
