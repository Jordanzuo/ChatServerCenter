package messageLogDAL

import (
	"github.com/Jordanzuo/ChatServerCenter/src/dal"
	"github.com/Jordanzuo/ChatServerModel/src/transferObject"
)

// 插入数据
func Insert(chatMessageObj *transferObject.ChatMessageObject) error {
	command := "INSERT INTO log_message(PlayerId, Name, PartnerId, ServerId, ServerGroupId, Message, ChannelType, ToPlayerId, Crtime) VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?);"

	stmt, err := dal.GetDB().Prepare(command)
	if err != nil {
		dal.WritePrepareError(command, err)
		return err
	}

	// 最后关闭
	defer stmt.Close()

	if _, err = stmt.Exec(chatMessageObj.Player.Id, chatMessageObj.Player.Name, chatMessageObj.Player.PartnerId, chatMessageObj.Player.ServerId, chatMessageObj.Player.ServerGroupId,
		chatMessageObj.Message, int(chatMessageObj.ChannelType), chatMessageObj.ToPlayerId, chatMessageObj.Crtime); err != nil {
		dal.WriteExecError(command, err)
		return err
	}

	return nil
}
