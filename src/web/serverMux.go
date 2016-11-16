package web

import (
	"encoding/json"
	"fmt"
	"net/http"

	"time"

	"github.com/Jordanzuo/ChatServerModel/src/apiLog"
	"github.com/Jordanzuo/ChatServerModel/src/centerResponseObject"
	"github.com/Jordanzuo/goutil/logUtil"
	"github.com/Jordanzuo/goutil/webUtil"
)

// 定义自定义的Mux对象
type selfDefineMux struct {
}

func (mux *selfDefineMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	responseObj := centerResponseObject.NewResponseObject()
	startTime := time.Now().Unix()
	r.ParseForm()

	// 获取输入参数的字符串形式
	parameter := ""
	if len(r.Form) > 0 {
		parameter_byte, _ := json.Marshal(r.Form)
		parameter = string(parameter_byte)
	}

	// 在输出结果给客户端之后再来处理日志的记录，以便于可以尽快地返回给客户端
	defer func() {
		// 记录DEBUG日志
		if debug {
			result, _ := json.Marshal(responseObj)

			msg := fmt.Sprintf("%s-->", r.RequestURI)
			msg += fmt.Sprintf("请求数据：%v;", parameter)
			msg += fmt.Sprintf("返回数据：%s;", string(result))
			logUtil.Log(msg, logUtil.Debug, true)
		}
	}()

	// 判断是否是POST方法
	if r.Method != "POST" {
		responseResult(w, responseObj.SetResultStatus(centerResponseObject.Con_OnlySupportPOST))
		return
	}

	// 验证IP是否正确
	if isIPValidFunc(webUtil.GetRequestIP(r)) == false {
		logUtil.Log(fmt.Sprintf("请求的IP：%s无效", webUtil.GetRequestIP(r)), logUtil.Error, true)
		responseResult(w, responseObj.SetResultStatus(centerResponseObject.Con_InvalidIP))
		return
	}

	// 根据路径选择不同的处理方法
	var requestFunc *requestFunc
	var exists bool
	if requestFunc, exists = funcMap[r.RequestURI]; !exists {
		responseResult(w, responseObj.SetResultStatus(centerResponseObject.Con_APINotDefined))
		return
	}

	// 判断参数是否正确
	if rs := requestFunc.checkParam(r); rs != centerResponseObject.Con_Success {
		responseResult(w, responseObj.SetResultStatus(rs))
		return
	}

	// 调用方法
	responseObj = requestFunc.funcDefinition(w, r)
	endTime := time.Now().Unix()

	// 在输出结果给客户端之后再来处理日志的记录，以便于可以尽快地返回给客户端
	defer func() {
		// 向API日志通道中写入API日志
		if ifRecordAPILogFunc() {
			parameter = fmt.Sprintf("%s:start:%d,end:%d,used:%d", parameter, startTime, endTime, endTime-startTime)
			apiLogChannel <- apiLog.NewApiLog(r.RequestURI, parameter)
		}
	}()

	// 输出结果
	responseResult(w, responseObj)
}
