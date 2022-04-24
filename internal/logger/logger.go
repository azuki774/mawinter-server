package logger

import (
	"go.uber.org/zap"
)

var Logger *zap.SugaredLogger

func NewSugarLogger() (err error) {
	l, err := zap.NewDevelopment()
	ls := l.Sugar()
	Logger = ls
	return err
}
