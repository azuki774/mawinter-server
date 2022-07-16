package logger

import (
	"go.uber.org/zap"
)

func NewSugarLogger() (Logger *zap.SugaredLogger, err error) {
	l, err := zap.NewDevelopment()
	ls := l.Sugar()
	return ls, err
}
