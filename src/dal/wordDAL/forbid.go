package wordDAL

import (
	"github.com/Jordanzuo/ChatServerCenter/src/dal"
)

// 初始化屏蔽词列表
func InitForbid() (wordList []string, err error) {
	command := "SELECT Word FROM config_word_forbid;"

	rows, err := dal.GetDB().Query(command)
	if err != nil {
		return
	}

	defer rows.Close()

	for rows.Next() {
		var word string
		if err = rows.Scan(&word); err != nil {
			dal.WriteScanError(command, err)
			return
		}

		wordList = append(wordList, word)
	}

	return
}
