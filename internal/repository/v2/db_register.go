package repository

import (
	"errors"
	"mawinter-server/internal/model"
	"mawinter-server/internal/openapi"
	"mawinter-server/internal/register"

	"gorm.io/gorm"
)

// InsertUniqueCatIDRecord は 同一のカテゴリIDがない場合ときに挿入、既にあればエラーを返す
func (d *DBRepository) InsertUniqueCatIDRecord(req openapi.Record) (res openapi.Record, err error) {
	yyyymm := req.Datetime.Format("200601")
	startDate, err := yyyymmToInitDayTime(yyyymm)
	if err != nil {
		return openapi.Record{}, err
	}
	endDate := startDate.AddDate(0, 1, 0)

	err = d.Conn.Table(RecordTableName).
		Where("category_id = ?", req.CategoryId).
		Where("datetime >= ? AND datetime < ?", startDate, endDate).
		Take(&model.Record{}).Error

	if err == nil {
		// already recorded
		return openapi.Record{}, register.ErrAlreadyRegisted
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		// unknown error
		return openapi.Record{}, err
	}
	dbres := d.Conn.Table(RecordTableName).Create(&req)
	if dbres.Error != nil {
		return openapi.Record{}, dbres.Error
	}

	res = req
	return res, nil
}

func (d *DBRepository) GetMonthlyFixBilling() (fixBills []model.MonthlyFixBilling, err error) {
	var recs []model.MonthlyFixBillingDB
	err = d.Conn.Find(&recs).Error
	if err != nil {
		return []model.MonthlyFixBilling{}, err
	}
	for _, v := range recs {
		fixBills = append(fixBills,
			model.MonthlyFixBilling{
				CategoryID: int(v.CategoryID),
				Day:        int(v.Day),
				Price:      int(v.Price),
				Type:       v.Type,
				Memo:       v.Memo,
			},
		)
	}
	return fixBills, nil
}
