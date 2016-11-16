package wordBLL

import (
	"fmt"

	"github.com/Jordanzuo/ChatServerCenter/src/bll/reloadBLL"
	"github.com/Jordanzuo/ChatServerCenter/src/dal/wordDAL"
	"github.com/Jordanzuo/goutil/debugUtil"
	"github.com/Jordanzuo/goutil/dfaUtil"
)

var (
	sensitiveWordList []string = make([]string, 0, 1024)
	sensitiveDFAObj   *dfaUtil.DFAUtil
)

func init() {
	if err := ReloadSensitive(); err != nil {
		panic(fmt.Errorf("初始化敏感词列表失败，错误信息为：%s", err))
	}

	// 注册重新加载的方法
	reloadBLL.RegisterReloadFunc("Sensitive", ReloadSensitive)
}

// 重新加载敏感词列表
func ReloadSensitive() error {
	var err error
	if sensitiveWordList, err = wordDAL.InitSensitive(); err != nil {
		return err
	}

	debugUtil.Printf("SensitiveWordList:%v\n", sensitiveWordList)

	// 构造DFAUtil对象
	sensitiveDFAObj = dfaUtil.NewDFAUtil(sensitiveWordList)

	return nil
}

// 处理屏蔽词汇
// 输入字符串
// 处理屏蔽词汇后的字符串
func HandleSensitiveWords(input string) string {
	if len(sensitiveWordList) == 0 {
		return input
	}

	return sensitiveDFAObj.HandleWord(input, '*')
}
