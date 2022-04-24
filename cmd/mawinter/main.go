package main

import (
	"fmt"
	"mawinter-expense/internal/api"
	"mawinter-expense/internal/db"
	l "mawinter-expense/internal/logger"
	"os"
	"strconv"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	err := l.NewSugarLogger()
	defer l.Logger.Sync()
	if err != nil {
		fmt.Printf("logger failed")
		os.Exit(1)
	}

	l.Logger.Info("Program Start")

	if os.Getenv("BASIC_AUTH_PASSWORD") == "" {
		l.Logger.Warn("No Basic Auth password set")
	} else {
		logMessage := "Basic Authentication username : " + os.Getenv("BASIC_AUTH_USERNAME") + ", password : " + os.Getenv("BASIC_AUTH_PASSWORD")
		l.Logger.Info(logMessage)
	}

	logMessage := "Using Database is " + os.Getenv("MYSQL_DATABASE")
	l.Logger.Info(logMessage)

	var DBSleepTime time.Duration
	if os.Getenv("DB_WAITTIME") == "" {
		DBSleepTime = 0
	} else {
		t, err := strconv.Atoi(os.Getenv("DB_WAITTIME"))
		if err != nil {
			fmt.Printf("DB_WAITTIME is invalid")
			DBSleepTime = 0
		} else {
			fmt.Printf("DB_WAITTIME is %d s", t)
			DBSleepTime = time.Duration(t)
		}
	}

	time.Sleep(time.Second * DBSleepTime)
	db.DBConnect("root", "password", "mawinter-db", os.Getenv("MYSQL_DATABASE"))
	defer db.DB.Close()

	api.ServerStart()
	l.Logger.Info("Program End")
}
