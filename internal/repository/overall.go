package repository

import (
	"gorm.io/gorm"
)

type DBRepository struct {
	Conn *gorm.DB
}

func (dbR *DBRepository) CloseDB() (err error) {
	sqlDB, err := dbR.Conn.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}
