package manageCenterBLL

import (
	"time"

	"github.com/Jordanzuo/ChatServerCenter/src/bll/reloadBLL"
	"github.com/Jordanzuo/goutil/logUtil"
)

func init() {
	// 先初始化一次服务器组
	if err := Reload(); err != nil {
		panic(err)
	}

	// 注册重新加载的方法
	reloadBLL.RegisterReloadFunc("managecenter", Reload)

	go func() {
		// 处理内部未处理的异常，以免导致主线程退出，从而导致系统崩溃
		defer func() {
			if r := recover(); r != nil {
				logUtil.LogUnknownError(r)
			}
		}()

		for {
			// 每小时刷新一次
			time.Sleep(time.Hour)

			// 刷新服务器组
			Reload()
		}
	}()
}

// 刷新服务器组
func Reload() error {
	var err error

	if err = reloadPartner(); err != nil {
		return err
	}

	if err = reloadServer(); err != nil {
		return err
	}

	if err = reloadServerGroup(); err != nil {
		return err
	}

	return nil
}
