package rpc

import (
	"github.com/Jordanzuo/ChatServerModel/src/centerResponseObject"
)

// 请求方法对象
type requestFunc struct {
	// 方法名称
	funcName string

	// 方法定义
	funcDefinition func(*client, []interface{}) *centerResponseObject.ResponseObject

	// 方法参数数量
	funcParamCount int
}

// 检测参数数量
// parameters：参数数组
func (rf *requestFunc) checkParamCount(parameters []interface{}) centerResponseObject.ResultStatus {
	if rf.funcParamCount > 0 && len(parameters) == 0 {
		return centerResponseObject.Con_ParamIsEmpty
	}

	if rf.funcParamCount != len(parameters) {
		return centerResponseObject.Con_ParamNotMatch
	}

	return centerResponseObject.Con_Success
}

// 创建新的请求方法对象
// _funcName：方法名称
// _funcDefinition：方法定义
// _funcParamCount：方法参数数量
func newRequestFunc(_funcName string, _funcDefinition func(*client, []interface{}) *centerResponseObject.ResponseObject, _funcParamCount int) *requestFunc {
	return &requestFunc{
		funcName:       _funcName,
		funcDefinition: _funcDefinition,
		funcParamCount: _funcParamCount,
	}
}
