package reloadBLL

import (
	"fmt"

	"github.com/Jordanzuo/goutil/logUtil"
)

var (
	reloadFuncMap = make(map[string]func() error)
)

func RegisterReloadFunc(funcName string, reloadFunc func() error) {
	reloadFuncMap[funcName] = reloadFunc
}

// 重新加载配置
func Reload() (errList []error) {
	for funcName, reloadFunc := range reloadFuncMap {
		if err := reloadFunc(); err != nil {
			logUtil.Log(fmt.Sprintf("%s Reload Fail, Error:%s", funcName, err), logUtil.Error, true)
			errList = append(errList, err)
		}
	}

	return
}
