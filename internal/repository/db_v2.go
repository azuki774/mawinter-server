package repository

import (
	"fmt"

	"gorm.io/gorm"
)

type DBRepositoryV2 struct {
	Conn *gorm.DB
}

func (d *DBRepositoryV2) CreateTableYYYYMM(yyyymm string) (err error) {
	baseTableName := "Record_YYYYMM"
	sql := fmt.Sprintf("CREATE TABLE `Record_%s` LIKE %s", yyyymm, baseTableName)
	err = d.Conn.Exec(sql).Error
	return err
}

func (d *DBRepositoryV2) CloseDB() (err error) {
	dbconn, err := d.Conn.DB()
	if err != nil {
		return err
	}
	return dbconn.Close()
}

// func (d *DBRepositoryV2) InsertRecord(req openapi.ReqRecord) (res openapi.Record, err error) {

// }
