package playerBLL

import (
	"github.com/Jordanzuo/ChatServerCenter/src/dal/playerDAL"
	"github.com/Jordanzuo/ChatServerModel/src/player"
	"time"
)

// 根据Id获取玩家对象（先从缓存中取，取不到再从数据库中去取）
// id：玩家Id
// 返回值：
// 玩家对象
// 是否存在该玩家
// 错误对象
func GetPlayer(id string) (playerObj *player.Player, exists bool, err error) {
	if id == "" {
		return
	}

	playerObj, exists, err = playerDAL.GetPlayer(id)

	return
}

// 更新玩家的封号状态
// playerObj：玩家对象
// isForbidden：是否封号
func UpdateForbidStatus(playerObj *player.Player, isForbidden bool) error {
	playerObj.IsForbidden = isForbidden
	if err := playerDAL.UpdateForbiddenStatus(playerObj); err != nil {
		return err
	}

	return nil
}

// 更新玩家的禁言状态
// playerObj：玩家对象
// silentEndTime：禁言结束时间
func UpdateSilentStatus(playerObj *player.Player, silentEndTime time.Time) error {
	playerObj.SilentEndTime = silentEndTime
	return playerDAL.UpdateSilentEndTime(playerObj)
}
