package configBLL

import (
	"fmt"

	"github.com/Jordanzuo/ChatServerCenter/src/bll/reloadBLL"
	"github.com/Jordanzuo/ChatServerCenter/src/dal/configDAL"
	"github.com/Jordanzuo/ChatServerModel/src/config"
	"github.com/Jordanzuo/goutil/debugUtil"
	"github.com/Jordanzuo/goutil/stringUtil"
)

func init() {
	if err := Reload(); err != nil {
		panic(fmt.Errorf("初始化数据库配置失败，错误信息为：%s", err))
	}

	// 注册重新加载的方法
	reloadBLL.RegisterReloadFunc("config", Reload)
}

var (
	configObj *config.Config
)

// 获取数据库配置
func GetConfig() *config.Config {
	return configObj
}

// 初始化数据库连接
func Reload() error {
	var err error
	if configObj, err = configDAL.Init(); err != nil {
		return err
	}

	debugUtil.Printf("Config:%v\n", configObj)

	return nil
}

// 处理消息长度
// message：消息
// 返回值：
// 处理后的消息
func HandleMessageLength(message string) string {
	if len(message) > configObj.GetMaxMessageLength() {
		return stringUtil.Substring(message, 0, configObj.GetMaxMessageLength())
	}

	return message
}

// 判断单个ChatServer是否达到最大承载上限
// clientCount：客户端数量
// 返回值：
// 单个ChatServer是否达到最大承载上限
func IfClientCountReachMax(clientCount int) bool {
	return clientCount >= configObj.GetMaxClientCount()
}

// 判断是否记录API日志
// 返回值：
// 是否记录API日志
func IfRecordAPILog() bool {
	return configObj.GetIfRecordAPILog()
}
