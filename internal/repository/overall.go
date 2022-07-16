package repository

import (
	"mawinter-server/internal/model"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type DBRepository interface {
	// category_id 順にその年の月ごとの合計を取得する
	CreateRecordDB(record model.Records) (retAddRecord model.ShowRecord, err error)
	GetYearSummaryDB(year int64) (yearSummaryInters []model.YearSummaryInter, err error)
	GetRecentRecord(n int) (getRecentData []model.ShowRecord, err error)
	DeleteRecordDB(id int64) (err error)
}

type dbRepository struct {
	conn gorm.DB
}

func NewDBRepository(conn *gorm.DB) DBRepository {
	return &dbRepository{conn: *conn}
}

func DBConnect(user string, password string, address string, dbName string) (gormDB *gorm.DB, err error) {
	dsn := user + ":" + password + "@(" + address + ")/" + dbName + "?parseTime=true&loc=Local"
	DB, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return DB, nil
}
