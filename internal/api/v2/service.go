package v2

import (
	"go.uber.org/zap"
)

// var jst *time.Location

// func init() {
// 	j, err := time.LoadLocation("Asia/Tokyo")
// 	if err != nil {
// 		panic(err)
// 	}
// 	jst = j
// }

type DBRepository interface {
}

type APIService struct {
	Logger *zap.Logger
	Repo   DBRepository
}
