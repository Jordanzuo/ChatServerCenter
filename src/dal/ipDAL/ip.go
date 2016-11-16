package ipDAL

import (
	"github.com/Jordanzuo/ChatServerCenter/src/dal"
)

// 初始化IP列表
func Init() (ipList []string, err error) {
	command := "SELECT IP FROM config_ip;"

	rows, err := dal.GetDB().Query(command)
	if err != nil {
		return
	}

	defer rows.Close()

	for rows.Next() {
		var ip string
		if err = rows.Scan(&ip); err != nil {
			dal.WriteScanError(command, err)
			return
		}

		ipList = append(ipList, ip)
	}

	return
}
