package ipBLL

import (
	"fmt"

	"github.com/Jordanzuo/ChatServerCenter/src/bll/manageCenterBLL"
	"github.com/Jordanzuo/ChatServerCenter/src/bll/reloadBLL"
	"github.com/Jordanzuo/ChatServerCenter/src/dal/ipDAL"
	"github.com/Jordanzuo/goutil/debugUtil"
)

var (
	ipList []string = make([]string, 0, 8)
)

func init() {
	if err := Reload(); err != nil {
		panic(fmt.Errorf("初始化IP列表失败，错误信息为：%s", err))
	}

	// 注册重新加载的方法
	reloadBLL.RegisterReloadFunc("ip", Reload)
}

// 重新加载IP列表
func Reload() error {
	var err error
	if ipList, err = ipDAL.Init(); err != nil {
		return err
	}

	debugUtil.Printf("IPList:%v\n", ipList)

	return nil
}

// 判断IP是否有效
// ip：IP地址
// 返回值：
// 是否有效
func IsIPValid(ip string) bool {
	// 判断IP是否在配置的可允许的列表中
	for _, value := range ipList {
		if value == ip {
			return true
		}
	}

	// 判断IP是否是某个服务器组的IP
	if manageCenterBLL.IsIpBelongToServerGroup(ip) {
		return true
	}

	return false
}
