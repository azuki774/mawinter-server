package repository

import (
	"fmt"

	"gorm.io/gorm"
)

type DBRepository struct {
	Conn *gorm.DB
}

func (d *DBRepository) CreateTableYYYYMM(yyyymm string) (err error) {
	baseTableName := "Record_YYYYMM"
	sql := fmt.Sprintf("CREATE TABLE `Record_%s` LIKE %s", yyyymm, baseTableName)
	err = d.Conn.Exec(sql).Error
	return err
}

func (d *DBRepository) CloseDB() (err error) {
	dbconn, err := d.Conn.DB()
	if err != nil {
		return err
	}
	return dbconn.Close()
}

// func (d *DBRepository) InsertRecord(req openapi.ReqRecord) (rec openapi.Record, err error) {

// }
