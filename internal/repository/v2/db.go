package repository

import (
	"context"
	"errors"
	"fmt"
	"mawinter-server/internal/model"
	"mawinter-server/internal/openapi"
	"mawinter-server/internal/timeutil"

	"gorm.io/gorm"
)

const RecordTableName = "Record"

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

func (d *DBRepository) InsertRecord(req openapi.ReqRecord) (rec openapi.Record, err error) {
	dbRec, err := NewDBModelRecord(req)
	if err != nil {
		return openapi.Record{}, err
	}

	dbRes := d.Conn.Table(RecordTableName).Create(&dbRec)
	if dbRes.Error != nil {
		return openapi.Record{}, dbRes.Error
	}

	rec, err = NewRecordFromDB(dbRec)
	if err != nil {
		return openapi.Record{}, err
	}

	return rec, nil
}

func (d *DBRepository) GetRecords(ctx context.Context, num int) (recs []openapi.Record, err error) {
	res := d.Conn.Table(RecordTableName).Order("id DESC").Limit(num).Find(&recs)
	if res.Error != nil {
		return []openapi.Record{}, res.Error
	}
	return recs, nil
}

func (d *DBRepository) GetMonthRecords(yyyymm string) (recs []openapi.Record, err error) {
	var res *gorm.DB
	startDate, err := yyyymmToInitDayTime(yyyymm)
	if err != nil {
		return []openapi.Record{}, err
	}
	endDate := startDate.AddDate(0, 1, 0)

	res = d.Conn.Debug().Table(RecordTableName).
		Where("datetime >= ? AND datetime < ?", startDate, endDate).
		Find(&recs)

	if errors.Is(res.Error, gorm.ErrRecordNotFound) { // TODO: 正しくは Error 1146 をハンドリングする
		return []openapi.Record{}, model.ErrNotFound
	} else if res.Error != nil {
		return []openapi.Record{}, res.Error
	}

	return recs, nil
}

func (d *DBRepository) GetMonthRecordsRecent(yyyymm string, num int) (recs []openapi.Record, err error) {
	startDate, err := yyyymmToInitDayTime(yyyymm)
	if err != nil {
		return []openapi.Record{}, err
	}
	endDate := startDate.AddDate(0, 1, 0)

	res := d.Conn.Table(RecordTableName).Where("datetime >= ? AND datetime < ?", startDate, endDate).Order("id DESC").Limit(num).Find(&recs)
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

// SumPriceForEachCatID は月間サマリ中間構造体を取得する（category_id 昇順）。
func (d *DBRepository) GetMonthMidSummary(yyyymm string) (summon []model.CategoryMidMonthSummary, err error) {
	startDate, err := yyyymmToInitDayTime(yyyymm)
	if err != nil {
		return []model.CategoryMidMonthSummary{}, err
	}
	endDate := startDate.AddDate(0, 1, 0)
	sqlWhere := fmt.Sprintf("datetime >= \"%s\" AND datetime < \"%s\"", startDate.Format("20060102"), endDate.Format("20060102"))
	sql := fmt.Sprintf(`SELECT category_id, count(1), sum(price) FROM Record WHERE %s GROUP BY category_id ORDER BY category_id`, sqlWhere)

	rows, err := d.Conn.Raw(sql).Rows()
	if err != nil {
		return []model.CategoryMidMonthSummary{}, err
	}
	defer rows.Close()
	for rows.Next() {
		var cm model.CategoryMidMonthSummary
		err = rows.Scan(&cm.CategoryId, &cm.Count, &cm.Price)
		if err != nil {
			return []model.CategoryMidMonthSummary{}, err
		}
		summon = append(summon, cm)
	}

	return summon, nil
}

// InsertMonthlyFixBilling は Record に固定費を登録する
func (d *DBRepository) InsertMonthlyFixBilling(yyyymm string) (recs []openapi.Record, err error) {
	var mfb []model.MonthlyFixBillingDB    // DBのモデル
	var fixBills []model.MonthlyFixBilling // 中間構造体
	var records []model.Record             // レコードDB追加用の構造体

	err = d.Conn.Transaction(func(tx *gorm.DB) error {
		nerr := tx.Table("Monthly_Fix_Billing").Find(&fixBills).Error // 固定費テーブルからデータ取得
		if nerr != nil {
			return nerr
		}
		for _, v := range mfb {
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
		for _, v := range fixBills {
			addrec, nerr := v.ConvAddDBModel(yyyymm)
			if err != nil {
				return nerr
			}
			records = append(records, addrec)
		}

		doneRec := model.MonthlyFixDoneDB{
			YYYYMM: yyyymm,
			Done:   1,
		}

		nerr = tx.Create(&doneRec).Error // 月固定データ追加記録
		if nerr != nil {
			return nerr
		}

		if len(records) > 0 {
			// 挿入すべきデータがある場合: issse #62
			nerr = tx.Table(RecordTableName).Create(&records).Error // 月固定データ追加
			if nerr != nil {
				return nerr
			}
		}

		// commit
		return nil
	})
	if err != nil {
		return []openapi.Record{}, err
	}

	// API返却用に構造体を変換
	cnf, err := d.MakeCategoryNameMap()
	if err != nil {
		return []openapi.Record{}, err
	}

	for _, v := range records {
		rec, err := NewRecordFromDB(v)
		if err != nil {
			return []openapi.Record{}, err
		}
		rec.CategoryName = cnf[rec.CategoryId]
		recs = append(recs, rec)
	}

	return recs, nil
}

// GetMonthlyFixDone は 固定費が登録済かどうかを取得する
// done = false なら未登録
func (d *DBRepository) GetMonthlyFixDone(yyyymm string) (done bool, err error) {
	var mfd model.MonthlyFixDoneDB
	err = d.Conn.Where("yyyymm = ?", yyyymm).Take(&mfd).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		} else {
			return false, err
		}
	}
	// already registed
	return true, nil
}

func dbModelToConfirmInfo(mc model.MonthlyConfirm) (yc openapi.ConfirmInfo) {
	var statusBool bool
	if mc.Confirm == uint8(1) {
		statusBool = true
	} else {
		statusBool = false
	}
	yc = openapi.ConfirmInfo{
		Status: &statusBool,
		Yyyymm: &mc.YYYYMM,
	}
	// status = false の場合は、ConfirmDatetime はDBにあっても無視する
	if statusBool {
		yc.ConfirmDatetime = &mc.ConfirmDatetime
	}

	return yc
}

func (d *DBRepository) GetMonthlyConfirm(yyyymm string) (yc openapi.ConfirmInfo, err error) {
	var mc model.MonthlyConfirm
	boolFalse := false
	err = d.Conn.Debug().Table("Monthly_Confirm").Where("yyyymm = ?", yyyymm).Take(&mc).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return openapi.ConfirmInfo{
				ConfirmDatetime: nil,
				Status:          &boolFalse,
				Yyyymm:          &yyyymm,
			}, nil
		} else {
			return openapi.ConfirmInfo{}, err
		}
	}

	yc = dbModelToConfirmInfo(mc)
	return yc, nil
}

func (d *DBRepository) UpdateMonthlyConfirm(yyyymm string, confirm bool) (yc openapi.ConfirmInfo, err error) {
	var mc model.MonthlyConfirm
	var confirmNum uint8
	t := timeutil.NowFunc()
	// confirm
	if confirm {
		confirmNum = uint8(1)
	} else {
		confirmNum = uint8(0)
	}

	err = d.Conn.Transaction(func(tx *gorm.DB) error {
		// GET
		nerr := tx.Debug().Table("Monthly_Confirm").Where("yyyymm = ?", yyyymm).Take(&mc).Error
		if nerr != nil && !errors.Is(nerr, gorm.ErrRecordNotFound) {
			return nerr
		}

		// UPSERT
		mc = model.MonthlyConfirm{
			YYYYMM:          yyyymm,
			Confirm:         confirmNum,
			ConfirmDatetime: t,
		}

		nerr = tx.Debug().Table("Monthly_Confirm").Save(&mc).Error
		if nerr != nil {
			return nerr
		}

		// commit
		return nil
	})
	if err != nil {
		return openapi.ConfirmInfo{}, err
	}

	yc = dbModelToConfirmInfo(mc)
	return yc, nil
}
