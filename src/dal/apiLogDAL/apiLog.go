package apiLogDAL

import (
	"github.com/Jordanzuo/ChatServerCenter/src/dal"
	"github.com/Jordanzuo/ChatServerModel/src/apiLog"
)

// 插入数据
func Insert(apiLogObj *apiLog.ApiLog) error {
	command := "INSERT INTO log_api(APIName, Content, Crtime) VALUES(?, ?, ?);"

	stmt, err := dal.GetDB().Prepare(command)
	if err != nil {
		dal.WritePrepareError(command, err)
		return err
	}

	// 最后关闭
	defer stmt.Close()

	if _, err = stmt.Exec(apiLogObj.GetApiName(), apiLogObj.GetContent(), apiLogObj.GetCrtime()); err != nil {
		dal.WriteExecError(command, err)
		return err
	}

	return nil
}
