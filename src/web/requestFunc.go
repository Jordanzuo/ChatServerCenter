package web

import (
	"net/http"

	"github.com/Jordanzuo/ChatServerModel/src/centerResponseObject"
)

// 请求方法对象
type requestFunc struct {
	// API完整路径名称
	apiFullName string

	// 方法定义
	funcDefinition func(http.ResponseWriter, *http.Request) *centerResponseObject.ResponseObject

	// 方法参数名称集合
	funcParamNames []string
}

// 检测参数数量
func (rf *requestFunc) checkParam(r *http.Request) centerResponseObject.ResultStatus {
	if len(rf.funcParamNames) != len(r.Form) {
		return centerResponseObject.Con_APIParamError
	}

	for _, name := range rf.funcParamNames {
		if r.Form[name] == nil || len(r.Form[name]) == 0 {
			return centerResponseObject.Con_APIParamError
		}
	}

	return centerResponseObject.Con_Success
}

// 创建新的请求方法对象
// _apiFullName：API完整路径名称
// _funcDefinition：方法定义
// _funcParamNames：方法参数名称集合
func newRequestFunc(_apiFullName string,
	_funcDefinition func(http.ResponseWriter, *http.Request) *centerResponseObject.ResponseObject,
	_funcParamNames ...string) *requestFunc {
	return &requestFunc{
		apiFullName:    _apiFullName,
		funcDefinition: _funcDefinition,
		funcParamNames: _funcParamNames,
	}
}
