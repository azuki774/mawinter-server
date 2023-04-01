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

func (d *DBRepository) GetMonthRecords(yyyymm string, params openapi.GetV2RecordYyyymmParams) (recs []openapi.Record, err error) {
	var res *gorm.DB

	// TODO: params -> from の実装
	if params.CategoryId != nil {
		// Category ID
		res = d.Conn.Table(fmt.Sprintf("Record_%s", yyyymm)).Where("category_id = ?", *params.CategoryId).Find(&recs)
	} else {
		// No Option
		res = d.Conn.Table(fmt.Sprintf("Record_%s", yyyymm)).Find(&recs)
	}

	if errors.Is(res.Error, gorm.ErrRecordNotFound) { // TODO: 正しくは Error 1146 をハンドリングする
		return []openapi.Record{}, model.ErrNotFound
	} else if res.Error != nil {
		return []openapi.Record{}, res.Error
	}

	return recs, nil
}

func (d *DBRepository) GetMonthRecordsRecent(yyyymm string, num int) (recs []openapi.Record, err error) {
	res := d.Conn.Table(fmt.Sprintf("Record_%s", yyyymm)).Order("id DESC").Limit(num).Find(&recs)
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
	sql := fmt.Sprintf(`SELECT category_id, count(1), sum(price) FROM Record_%s GROUP BY category_id ORDER BY category_id`, yyyymm)

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

// InsertMonthlyFixBilling は Record_YYYYMM に固定費を登録する
func (d *DBRepository) InsertMonthlyFixBilling(yyyymm string) (recs []openapi.Record, err error) {
	var mfb []model.MonthlyFixBillingDB    // DBのモデル
	var fixBills []model.MonthlyFixBilling // 中間構造体
	var records []model.Record_YYYYMM      // レコードDB追加用の構造体

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
		nerr = tx.Table(fmt.Sprintf("Record_%s", yyyymm)).Create(&records).Error // 月固定データ追加
		if nerr != nil {
			return nerr
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

// InsertMonthlyFixBilling は 固定費が登録済かどうかを取得する
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
