package configDAL

import (
	"github.com/Jordanzuo/ChatServerCenter/src/dal"
	"github.com/Jordanzuo/ChatServerModel/src/config"
)

// 初始化数据库配置
func Init() (configObj *config.Config, err error) {
	command := `SELECT 
					AppId, AppName, AppKey, ManageCenterAPI, PartnerListAPI, ServerListAPI, ServerGroupListAPI, GroupType, 
					ChatServerCenterRpcAddress, ChatServerCenterWebAddress, PlayerInfoAPI, MaxMessageLength, MaxClientCount, IfRecordAPILog
				FROM 
					config;`

	var appId string
	var appName string
	var appKey string
	var manageCenterAPI string
	var partnerListAPI string
	var serverListAPI string
	var serverGroupListAPI string
	var groupType string
	var chatServerCenterRpcAddress string
	var chatServerCenterWebAddress string
	var playerInfoAPI string
	var maxMessageLength int
	var maxClientCount int
	var ifRecordAPILog bool

	err = dal.GetDB().QueryRow(command).Scan(&appId, &appName, &appKey,
		&manageCenterAPI, &partnerListAPI, &serverListAPI, &serverGroupListAPI, &groupType,
		&chatServerCenterRpcAddress, &chatServerCenterWebAddress, &playerInfoAPI, &maxMessageLength,
		&maxClientCount, &ifRecordAPILog)
	if err != nil {
		dal.WriteScanError(command, err)
		return
	}

	// 构造对象
	configObj = config.NewConfig(appId, appName, appKey,
		manageCenterAPI, partnerListAPI, serverListAPI, serverGroupListAPI, groupType,
		chatServerCenterRpcAddress, chatServerCenterWebAddress, playerInfoAPI, maxMessageLength,
		maxClientCount, ifRecordAPILog)

	return
}
