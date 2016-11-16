package manageCenterBLL

import (
	"encoding/json"
	"errors"
	"fmt"

	"sync"

	"github.com/Jordanzuo/ChatServerCenter/src/bll/configBLL"
	"github.com/Jordanzuo/ManageCenterModel_Go/partner"
	"github.com/Jordanzuo/ManageCenterModel_Go/returnObject"
	"github.com/Jordanzuo/goutil/logUtil"
	"github.com/Jordanzuo/goutil/webUtil"
)

var (
	partnerMap   = make(map[int]*partner.Partner, 128)
	partnerMutex sync.RWMutex
)

// 重新加载合作商
func reloadPartner() error {
	logUtil.Log("开始刷新合作商列表", logUtil.Debug, true)

	// 获取数据库配置
	configObj := configBLL.GetConfig()

	// 连接服务器，以获取数据
	url := fmt.Sprintf("%s/%s", configObj.GetManageCenterAPI(), configObj.GetPartnerListAPI())
	returnBytes, err := webUtil.PostWebData(url, nil, nil)
	if err != nil {
		logUtil.Log(fmt.Sprintf("获取合作商列表出错，错误信息为：%s", err), logUtil.Error, true)
		return err
	}

	// 解析返回值
	returnObj := new(returnObject.ReturnObject)
	if err = json.Unmarshal(returnBytes, &returnObj); err != nil {
		logUtil.Log(fmt.Sprintf("获取合作商列表出错，反序列化返回值出错，错误信息为：%s", err), logUtil.Error, true)
		return err
	}

	// 判断返回状态是否为成功
	if returnObj.Code != 0 {
		msg := fmt.Sprintf("获取合作商列表出错，返回状态：%d，信息为：%s", returnObj.Code, returnObj.Message)
		logUtil.Log(msg, logUtil.Error, true)
		return errors.New(msg)
	}

	// 解析Data
	tmpPartnerList := make([]*partner.Partner, 0, 128)
	tmpPartnerMap := make(map[int]*partner.Partner)
	if data, ok := returnObj.Data.(string); !ok {
		msg := "获取合作商列表出错，返回的数据不是string类型"
		logUtil.Log(msg, logUtil.Error, true)
		return errors.New(msg)
	} else {
		if err = json.Unmarshal([]byte(data), &tmpPartnerList); err != nil {
			logUtil.Log(fmt.Sprintf("获取合作商列表出错，反序列化数据出错，错误信息为：%s", err), logUtil.Error, true)
			return err
		}

		for _, item := range tmpPartnerList {
			tmpPartnerMap[item.Id] = item
		}
	}

	logUtil.Log(fmt.Sprintf("刷新合作商信息结束，合作商的数量:%d", len(tmpPartnerList)), logUtil.Debug, true)

	// 赋值给最终的partnerMap
	partnerMutex.Lock()
	defer partnerMutex.Unlock()
	partnerMap = tmpPartnerMap

	return nil
}

// 根据合作商Id获取合作商对象
// partnerId：合作商Id
// 返回值：
// 合作商对象
// 是否存在
func getPartner(partnerId int) (*partner.Partner, bool) {
	partnerMutex.RLock()
	defer partnerMutex.RUnlock()

	if partnerObj, exists := partnerMap[partnerId]; exists {
		return partnerObj, exists
	}

	return nil, false
}
