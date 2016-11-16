package dal

import (
	"database/sql"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/Jordanzuo/ChatServerCenter/src/config"
	"github.com/Jordanzuo/goutil/logUtil"
	_ "github.com/go-sql-driver/mysql"
)

var (
	// 数据库对象
	db *sql.DB
)

func init() {
	db = openMysqlConnection(config.DBConnection)
}

// 获取数据库对象
// 返回值：
// 数据库对象
func GetDB() *sql.DB {
	return db
}

// 打开数据库连接
// connectionString：数据库连接字符串
// 返回值：
// 数据库对象
func openMysqlConnection(connectionString string) *sql.DB {
	connectionSlice := strings.Split(connectionString, "||")
	if len(connectionSlice) != 3 {
		panic(fmt.Errorf("数据库连接配置不完整，当前的为：%s", connectionString))
	}

	// 建立数据库连接
	db, err := sql.Open("mysql", connectionSlice[0])
	if err != nil {
		panic(fmt.Errorf("打开游戏数据库失败,连接字符串为：%s", connectionString))
	}

	// 设置连接池相关
	maxOpenConns_string := strings.Replace(connectionSlice[1], "MaxOpenConns=", "", 1)
	maxOpenCons, err := strconv.Atoi(maxOpenConns_string)
	if err != nil {
		panic(fmt.Errorf("MaxOpenConns必须为int型,连接字符串为：%s", connectionString))
	}

	maxIdleConns_string := strings.Replace(connectionSlice[2], "MaxIdleConns=", "", 1)
	maxIdleConns, err := strconv.Atoi(maxIdleConns_string)
	if err != nil {
		panic(fmt.Errorf("MaxIdleConns必须为int型,连接字符串为：%s", connectionString))
	}

	if maxOpenCons > 0 && maxIdleConns > 0 {
		db.SetMaxOpenConns(maxOpenCons)
		db.SetMaxIdleConns(maxIdleConns)

		go func() {
			// 处理内部未处理的异常，以免导致主线程退出，从而导致系统崩溃
			defer func() {
				if r := recover(); r != nil {
					logUtil.LogUnknownError(r)
				}
			}()

			for {
				// 每分钟ping一次数据库
				time.Sleep(time.Minute)

				if err := db.Ping(); err != nil {
					logUtil.Log(fmt.Sprintf("Ping数据库失败,连接字符串为：%s,错误信息为：%s", connectionString, err), logUtil.Error, true)
				}
			}
		}()
	}

	if err := db.Ping(); err != nil {
		panic(fmt.Errorf("Ping数据库失败,连接字符串为：%s,错误信息为：%s", connectionString, err))
	}

	return db
}

// 记录Prepare错误
// command：执行的SQL语句
// err：错误对象
func WritePrepareError(command string, err error) {
	logUtil.Log(fmt.Sprintf("Prepare失败，错误信息：%s，command:%s", err, command), logUtil.Error, true)
}

// 记录Exec错误
// command：执行的SQL语句
// err：错误对象
func WriteExecError(command string, err error) {
	logUtil.Log(fmt.Sprintf("Exec失败，错误信息：%s，command:%s", err, command), logUtil.Error, true)
}

// 记录Scan错误
// command：执行的SQL语句
// err：错误对象
func WriteScanError(command string, err error) {
	logUtil.Log(fmt.Sprintf("Scan失败，错误信息：%s，command:%s", err, command), logUtil.Error, true)
}
