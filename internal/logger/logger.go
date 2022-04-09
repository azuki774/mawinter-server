package logger

import (
	"os"

	"github.com/sirupsen/logrus"
)

func init() {
	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetOutput(os.Stdout)
	logrus.SetLevel(logrus.InfoLevel)
}

func AccessInfoPrint(url string, method string, clientIP string) {
	logrus.WithFields(logrus.Fields{
		"url":      url,
		"method":   method,
		"clientIP": clientIP,
	}).Info("API Access")
}

func DBInfoPrint(funcName string) {
	logrus.WithFields(logrus.Fields{
		"funcName": funcName,
	}).Info("SQL published")
}

func InfoPrint(content string) {
	logrus.Info(content)
}

func WarnPrint(content string) {
	logrus.Warn(content)
}

func ErrorPrint(content string) {
	logrus.Error(content)
}

func FatalPrint(content string) {
	logrus.Fatal(content)
}
