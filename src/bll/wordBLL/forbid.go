package wordBLL

import (
	"fmt"

	"github.com/Jordanzuo/ChatServerCenter/src/bll/reloadBLL"
	"github.com/Jordanzuo/ChatServerCenter/src/dal/wordDAL"
	"github.com/Jordanzuo/goutil/debugUtil"
	"github.com/Jordanzuo/goutil/dfaUtil"
)

var (
	forbidWordList []string = make([]string, 0, 1024)
	forbidDFAObj   *dfaUtil.DFAUtil
)

func init() {
	if err := ReloadForbid(); err != nil {
		panic(fmt.Errorf("初始化屏蔽词列表失败，错误信息为：%s", err))
	}

	// 注册重新加载的方法
	reloadBLL.RegisterReloadFunc("Forbid", ReloadForbid)
}

// 重新加载屏蔽词列表
func ReloadForbid() error {
	var err error
	if forbidWordList, err = wordDAL.InitForbid(); err != nil {
		return err
	}

	debugUtil.Printf("ForbidWordList:%v\n", forbidWordList)

	// 构造DFAUtil对象
	forbidDFAObj = dfaUtil.NewDFAUtil(forbidWordList)

	return nil
}

// 是否包含屏蔽词
func IfContainsForbidWords(input string) bool {
	if len(forbidWordList) == 0 {
		return false
	}

	return forbidDFAObj.IsMatch(input)
}
