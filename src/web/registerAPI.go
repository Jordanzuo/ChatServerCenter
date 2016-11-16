package web

import (
	"fmt"
	"net/http"

	"github.com/Jordanzuo/ChatServerModel/src/centerResponseObject"
)

var (
	// 处理的方法映射表
	funcMap = make(map[string]*requestFunc)
)

// 注册API
// apiName：API名称
// callback：回调方法
// paramNames：参数名称集合
func registerAPI(apiName string,
	callback func(http.ResponseWriter, *http.Request) *centerResponseObject.ResponseObject,
	paramNames ...string) {
	apiFullName := fmt.Sprintf("/API/%s", apiName)
	funcMap[apiFullName] = newRequestFunc(apiFullName, callback, paramNames...)
}
