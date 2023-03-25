package mysql

import (
	"fmt"
	"manage/setting"

	"go.uber.org/zap"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

//_ "github.com/go-sql-driver/mysql" 不要忘记引入驱动

var DB *sqlx.DB

func Init(cfg *setting.MySQLConfig) (err error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.DbName)
	// 也可以使用MustConnect连接不成功就panic
	DB, err = sqlx.Connect("mysql", dsn)
	if err != nil {
		zap.L().Error("connect DB failed", zap.Error(err))
		return
	}
	DB.SetMaxOpenConns(cfg.MaxConn)
	DB.SetMaxIdleConns(cfg.MaxIdleConn)
	return
}

func Close() {
	_ = DB.Close()
}
