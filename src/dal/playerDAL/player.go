package playerDAL

import (
	"database/sql"
	"time"

	"github.com/Jordanzuo/ChatServerCenter/src/dal"
	"github.com/Jordanzuo/ChatServerModel/src/player"
)

func GetPlayer(id string) (playerObj *player.Player, exists bool, err error) {
	command := "SELECT Name, PartnerId, ServerId, UnionId, ExtraMsg, RegisterTime, LoginTime, IsForbidden, SilentEndTime FROM player WHERE Id = ?;"

	var name string
	var partnerId int
	var serverId int
	var unionId string
	var extraMsg string
	var registerTime time.Time
	var loginTime time.Time
	var isForbidden bool
	var silentEndTime time.Time
	if err = dal.GetDB().QueryRow(command, id).Scan(&name, &partnerId, &serverId, &unionId, &extraMsg, &registerTime, &loginTime, &isForbidden, &silentEndTime); err != nil {
		if err == sql.ErrNoRows {
			// 重置err，使其为nil；因为这代表的是没有查找到数据，而不是真正的错误
			err = nil
			return
		} else {
			dal.WriteScanError(command, err)
			return
		}
	}

	playerObj = player.NewPlayer(id, name, partnerId, serverId, unionId, extraMsg, registerTime, loginTime, isForbidden, silentEndTime)
	exists = true

	return
}

func Insert(player *player.Player) error {
	command := `INSERT INTO 
                player(Id, Name, PartnerId, ServerId, UnionId, ExtraMsg, RegisterTime, LoginTime, IsForbidden, SilentEndTime)
            VALUES
                (?, ?, ?, ?, ?, ?, ?, ?, ?, ?);
    `
	stmt, err := dal.GetDB().Prepare(command)
	if err != nil {
		dal.WritePrepareError(command, err)
		return err
	}

	// 最后关闭
	defer stmt.Close()

	if _, err = stmt.Exec(player.Id, player.Name, player.PartnerId, player.ServerId, player.UnionId, player.ExtraMsg, player.RegisterTime, player.LoginTime, player.IsForbidden, player.SilentEndTime); err != nil {
		dal.WriteExecError(command, err)
		return err
	}

	return nil
}

func UpdateInfo(player *player.Player) error {
	command := "UPDATE player SET Name = ?, UnionId = ?, ExtraMsg = ? WHERE Id = ?"
	stmt, err := dal.GetDB().Prepare(command)
	if err != nil {
		dal.WritePrepareError(command, err)
		return err
	}

	// 最后关闭
	defer stmt.Close()

	if _, err = stmt.Exec(player.Name, player.UnionId, player.ExtraMsg, player.Id); err != nil {
		dal.WriteExecError(command, err)
		return err
	}

	return nil
}

func UpdateLoginTime(player *player.Player) error {
	command := "UPDATE player SET LoginTime = ? WHERE Id = ?"
	stmt, err := dal.GetDB().Prepare(command)
	if err != nil {
		dal.WritePrepareError(command, err)
		return err
	}

	// 最后关闭
	defer stmt.Close()

	if _, err = stmt.Exec(player.LoginTime, player.Id); err != nil {
		dal.WriteExecError(command, err)
		return err
	}

	return nil
}

func UpdateForbiddenStatus(player *player.Player) error {
	command := "UPDATE player SET IsForbidden = ? WHERE Id = ?"
	stmt, err := dal.GetDB().Prepare(command)
	if err != nil {
		dal.WritePrepareError(command, err)
		return err
	}

	// 最后关闭
	defer stmt.Close()

	if _, err = stmt.Exec(player.IsForbidden, player.Id); err != nil {
		dal.WriteExecError(command, err)
		return err
	}

	return nil
}

func UpdateSilentEndTime(player *player.Player) error {
	command := "UPDATE player SET SilentEndTime = ? WHERE Id = ?"
	stmt, err := dal.GetDB().Prepare(command)
	if err != nil {
		dal.WritePrepareError(command, err)
		return err
	}

	// 最后关闭
	defer stmt.Close()

	if _, err = stmt.Exec(player.SilentEndTime, player.Id); err != nil {
		dal.WriteExecError(command, err)
		return err
	}

	return nil
}
