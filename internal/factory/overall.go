package factory

import (
	"fmt"
	"mawinter-server/internal/api"
	"mawinter-server/internal/repository"
	"mawinter-server/internal/server"
	"net"
	"os"
	"time"

	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

const DBConnectRetry = 5
const DBConnectRetryInterval = 10

func NewAPIService(user string, password string, host string, port string, dbName string) (*api.APIService, error) {
	l, err := NewLogger()
	if err != nil {
		return nil, err
	}

	addr := net.JoinHostPort(host, port)
	dsn := user + ":" + password + "@(" + addr + ")/" + dbName + "?parseTime=true&loc=Local"
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
			l.Error("failed to connect (DB)", zap.Error(err))
			return nil, err
		}

		time.Sleep(DBConnectRetryInterval * time.Second)
	}

	dbRepo := &repository.DBRepository{Conn: gormdb}
	return &api.APIService{Logger: l, DBRepo: dbRepo}, nil
}

func NewServer(api *api.APIService) (*server.Server, error) {
	l, err := NewLogger()
	if err != nil {
		return nil, err
	}

	return &server.Server{Logger: l, APIService: api, BasicAuth: struct {
		User string
		Pass string
	}{os.Getenv("BASIC_AUTH_USERNAME"), os.Getenv("BASIC_AUTH_PASSWORD")}}, nil
}

func CloseDB(gormdb *gorm.DB) (err error) {
	sqlDB, err := gormdb.DB()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	return sqlDB.Close()
}
