package sqlx_client

import (
	"fmt"
	"iot_server8/config"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

var db *sqlx.DB

func SetUp() {
	db = sqlx.MustConnect("mysql", fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?loc=Local",
		config.Conf.MysqlUserName, config.Conf.MysqlPassword, config.Conf.MysqlAddr, config.Conf.MysqlPort, config.Conf.MysqlDatabase))
}
