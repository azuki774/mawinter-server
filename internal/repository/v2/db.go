package repository

import (
	"errors"
	"fmt"
	"mawinter-server/internal/model"
	"mawinter-server/internal/openapi"
	"time"

	"gorm.io/gorm"
)

type DBRepository struct {
	Conn *gorm.DB
}

func (d *DBRepository) CloseDB() (err error) {
	dbconn, err := d.Conn.DB()
	if err != nil {
		return err
	}
	return dbconn.Close()
}

func getRecordTable(t time.Time) string {
	YYYYMM := t.Format("200601")
	return fmt.Sprintf("Record_%s", YYYYMM)
}

func (d *DBRepository) CreateTableYYYYMM(yyyymm string) (err error) {
	baseTableName := "Record_YYYYMM"
	sql := fmt.Sprintf("CREATE TABLE `Record_%s` LIKE %s", yyyymm, baseTableName)
	err = d.Conn.Exec(sql).Error
	return err
}

func (d *DBRepository) InsertRecord(req openapi.ReqRecord) (rec openapi.Record, err error) {
	dbRec, err := NewDBModelRecord(req)
	if err != nil {
		return openapi.Record{}, err
	}

	dbRes := d.Conn.Table(getRecordTable(dbRec.Datetime)).Create(&dbRec)
	if dbRes.Error != nil {
		return openapi.Record{}, dbRes.Error
	}

	rec, err = NewRecordFromDB(dbRec)
	if err != nil {
		return openapi.Record{}, err
	}

	return rec, nil
}

func (d *DBRepository) GetMonthRecords(yyyymm string) (recs []openapi.Record, err error) {
	res := d.Conn.Table(fmt.Sprintf("Record_%s", yyyymm)).Find(&recs)
	if errors.Is(res.Error, gorm.ErrRecordNotFound) { // TODO: 正しくは Error 1146 をハンドリングする
		return []openapi.Record{}, model.ErrNotFound
	} else if res.Error != nil {
		return []openapi.Record{}, res.Error
	}

	return recs, nil
}

func (d *DBRepository) MakeCategoryNameMap() (cnf map[int]string, err error) {
	cnf = make(map[int]string)
	var catTable []model.Category
	res := d.Conn.Table("Category").Find(&catTable)
	if res.Error != nil {
		return nil, res.Error
	}

	for _, c := range catTable {
		cnf[int(c.CategoryID)] = c.Name
	}

	return cnf, nil
}
