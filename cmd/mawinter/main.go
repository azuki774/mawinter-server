package main

import (
	"fmt"
	"mawinter-server/internal/api"
	"mawinter-server/internal/logger"
	"mawinter-server/internal/repository"
	"mawinter-server/internal/server"
	"os"
	"strconv"

	"time"
)

func main() {
	l, err := logger.NewSugarLogger()
	defer l.Sync()
	if err != nil {
		fmt.Printf("logger failed")
		os.Exit(1)
	}

	l.Infow("Program Start")

	if os.Getenv("BASIC_AUTH_PASSWORD") == "" {
		l.Warnw("No Basic Auth password set")
	} else {
		l.Infow("Basic Authentication info", "username", os.Getenv("BASIC_AUTH_USERNAME"), "password", os.Getenv("BASIC_AUTH_PASSWORD"))
	}

	l.Infow("database info", "name", os.Getenv("MYSQL_DATABASE"))

	var DBSleepTime time.Duration
	if os.Getenv("DB_WAITTIME") == "" {
		DBSleepTime = 0
	} else {
		t, err := strconv.Atoi(os.Getenv("DB_WAITTIME"))
		if err != nil {
			l.Errorw("DB_WAITTIME is invalid")
			DBSleepTime = 0
		} else {
			l.Infof("DB_WAITTIME is %d s", t)
			DBSleepTime = time.Duration(t)
		}
	}

	time.Sleep(time.Second * DBSleepTime)
	gormdb, err := repository.DBConnect("root", "password", "mawinter-db", os.Getenv("MYSQL_DATABASE"))
	if err != nil {
		l.Errorw(err.Error())
		os.Exit(1)
	}
	sqlDB, err := gormdb.DB()
	if err != nil {
		l.Error(err.Error())
		os.Exit(1)
	}
	defer sqlDB.Close()

	dbR := repository.NewDBRepository(gormdb)
	as := api.NewAPIService(dbR, l)
	server.Start(as, l)
	l.Info("Program End")
}
