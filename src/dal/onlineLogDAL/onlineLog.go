package onlineLogDAL

import (
	"github.com/Jordanzuo/ChatServerCenter/src/dal"
	"github.com/Jordanzuo/ChatServerModel/src/onlineLog"
)

// 插入数据
func Insert(onlineLogObj *onlineLog.OnlineLog) error {
	command := "INSERT INTO log_online(OnlineTime, Sid, ServerAddress, ClientCount, PlayerCount, TotalCount) VALUES(?, ?, ?, ?, ?, ?);"

	stmt, err := dal.GetDB().Prepare(command)
	if err != nil {
		dal.WritePrepareError(command, err)
		return err
	}

	// 最后关闭
	defer stmt.Close()

	if _, err = stmt.Exec(onlineLogObj.GetOnlineTime(), onlineLogObj.GetSid(), onlineLogObj.GetServerAddress(), onlineLogObj.GetClientCount(), onlineLogObj.GetPlayerCount(), onlineLogObj.GetTotalCount()); err != nil {
		dal.WriteExecError(command, err)
		return err
	}

	return nil
}
