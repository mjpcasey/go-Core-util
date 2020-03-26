package dbmanager

import (
	"database/sql"
	"fmt"
	"gcore/app/share"
	"gcore/gmonitor"

	_ "github.com/go-sql-driver/mysql"
)

// MySQL 客户端封装
type MySQLWrapper struct {
	Name            string
	Username        string
	Password        string
	Addr            string
	Database        string
	ReadTimeoutSec  int
	ConnTimeoutSec  int
	WriteTimeoutSec int
	DB              *sql.DB
}

// 创建 MySQL 客户端封装
func (m *DBManager) createMySQL(app *share.AppShare, confPath string) {
	config := app.GetConfig()

	for i := 0; i < config.GetInt(confPath+"length"); i++ {
		wrapper := &MySQLWrapper{}

		err := config.Scan(fmt.Sprintf("%s%d", confPath, i), &wrapper)
		if err == nil {
			logger.Infof("创建 mysql 数据库封装[%s]", wrapper.Name)

			dsn := fmt.Sprintf(
				"%s:%s@tcp(%s)/%s?timeout=%ds&readTimeout=%ds&writeTimeout=%ds",
				wrapper.Username, wrapper.Password, wrapper.Addr, wrapper.Database,
				wrapper.ConnTimeoutSec,
				wrapper.ReadTimeoutSec,
				wrapper.WriteTimeoutSec,
			)
			logger.Debugf("数据源为[%s]", dsn)

			db, err := sql.Open("mysql", dsn)
			if err == nil {
				err = db.Ping()
			}
			if err == nil {
				wrapper.DB = db
				m.mysqlMap[wrapper.Name] = wrapper

				gmonitor.NewMysqlMonitor(wrapper.DB, wrapper.Name)
			} else {
				logger.Fatalf("mysql 初始化[name=%s]出错: %s", wrapper.Name, err)
			}
		} else {
			logger.Fatalf("解析 mysql 配置%d失败: %s", i, err)
		}
	}

	return
}
