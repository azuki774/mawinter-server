package factory

import (
	"net"
	"time"

	v2db "mawinter-server/internal/repository/v2"

	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

const DBConnectRetry = 5
const DBConnectRetryInterval = 10

func NewDBRepositoryV1(host, port, user, pass, name string) (dbR *v2db.DBRepository, err error) {
	l, err := NewLogger()
	if err != nil {
		return nil, err
	}

	addr := net.JoinHostPort(host, port)
	dsn := user + ":" + pass + "@(" + addr + ")/" + name + "?parseTime=true&loc=Local"
	var gormdb *gorm.DB
	for i := 0; i < DBConnectRetry; i++ {
		gormdb, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
		if err == nil {
			// Success DB connect
			l.Info("DB connect")
			break
		}
		l.Warn("DB connection retry")

		if i == DBConnectRetry {
			l.Error("failed to connect DB", zap.Error(err))
			return nil, err
		}

		time.Sleep(DBConnectRetryInterval * time.Second)
	}

	return &v2db.DBRepository{Conn: gormdb}, nil
}

func NewDBRepositoryV2(host, port, user, pass, name string) (dbR *v2db.DBRepository, err error) {
	l, err := NewLogger()
	if err != nil {
		return nil, err
	}

	addr := net.JoinHostPort(host, port)
	dsn := user + ":" + pass + "@(" + addr + ")/" + name + "?parseTime=true&loc=Local"
	var gormdb *gorm.DB
	for i := 0; i < DBConnectRetry; i++ {
		gormdb, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
		if err == nil {
			// Success DB connect
			l.Info("DB connect")
			break
		}
		l.Warn("DB connection retry")

		if i == DBConnectRetry {
			l.Error("failed to connect DB", zap.Error(err))
			return nil, err
		}

		time.Sleep(DBConnectRetryInterval * time.Second)
	}

	return &v2db.DBRepository{Conn: gormdb}, nil
}
