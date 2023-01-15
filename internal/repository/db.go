package repository

import (
	"errors"
	"fmt"
	"mawinter-server/internal/model"
	"time"

	"gorm.io/gorm"
)

var jst *time.Location

func init() {
	j, err := time.LoadLocation("Asia/Tokyo")
	if err != nil {
		panic(err)
	}
	jst = j
}

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

func (d *DBRepository) CreateRecordTable(yyyymm string) (err error) {
	sql := fmt.Sprintf("CREATE TABLE `Record_%s` (", yyyymm)
	sql = sql + "`id` int NOT NULL AUTO_INCREMENT,"
	sql = sql + "`category_id`  int NOT NULL,"
	sql = sql + "`datetime` datetime NOT NULL default current_timestamp,"
	sql = sql + "`from` varchar(64) NOT NULL,"
	sql = sql + "`type` varchar(64) NOT NULL,"
	sql = sql + "`price` int NOT NULL,"
	sql = sql + "`memo` varchar(255) NOT NULL,"
	sql = sql + "`created_at` datetime default current_timestamp,"
	sql = sql + "`updated_at` timestamp default current_timestamp on update current_timestamp,"
	sql = sql + "PRIMARY KEY (`id`),"
	sql = sql + "index idx_cat (`category_id`),"
	sql = sql + "index idx_date (`datetime`) )"

	err = d.Conn.Exec(sql).Error
	return err
}

func getRecordTable(t time.Time) string {
	YYYYMM := t.Format("200601")
	return fmt.Sprintf("Record_%s", YYYYMM)
}

func (d *DBRepository) InsertRecord(req model.Recordstruct) (res model.Recordstruct, err error) {
	dbres := d.Conn.Table(getRecordTable(req.Datetime)).Create(&req)
	if dbres.Error != nil {
		return model.Recordstruct{}, dbres.Error
	}

	res = req
	return res, err
}



func (d *DBRepository) GetCategoryInfo() (info []model.Category, err error) {
	info = []model.Category{}
	err = d.Conn.Find(&info).Error
	return info, err
}

func (d *DBRepository) SumPriceForEachCatID(yyyymm string) (sum []model.SumPriceCategoryID, err error) {
	sql := fmt.Sprintf(`SELECT category_id, count(1), sum(price) FROM Record_%s GROUP BY category_id`, yyyymm)

	rows, err := d.Conn.Raw(sql).Rows()
	if err != nil {
		return []model.SumPriceCategoryID{}, err
	}
	defer rows.Close()
	for rows.Next() {
		var s model.SumPriceCategoryID
		err = rows.Scan(&s.CategoryID, &s.Count, &s.Sum)
		if err != nil {
			return []model.SumPriceCategoryID{}, err
		}
		sum = append(sum, s)
	}

	return sum, nil
}

// GetMonthlyFixDone は yyyymm 月のレコードがあれば true を返す
func (d *DBRepository) GetMonthlyFixDone(yyyymm string) (flag bool, err error) {
	var doneRec model.MonthlyFixDoneDB
	err = d.Conn.Where("yyyymm = ?", yyyymm).Take(&doneRec).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return false, err
	} else if errors.Is(err, gorm.ErrRecordNotFound) {
		// 未処理
		return false, nil
	}

	// 処理済
	return true, nil
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

func (d *DBRepository) InsertMonthlyFixBilling(yyyymm string, fixBills []model.MonthlyFixBilling) (err error) {
	doneRec := model.MonthlyFixDoneDB{
		YYYYMM: yyyymm,
		Done:   1,
	}

	var insRecs []model.Record_YYYYMM

	for _, v := range fixBills {
		addrec, err := v.ConvAddDBModel(yyyymm)
		if err != nil {
			return err
		}
		insRecs = append(insRecs, addrec)
	}

	err = d.Conn.Transaction(func(tx *gorm.DB) error {
		nerr := tx.Create(&doneRec).Error
		if nerr != nil {
			return nerr
		}

		nerr = tx.Table(fmt.Sprintf("Record_%s", yyyymm)).Create(&insRecs).Error
		if nerr != nil {
			return nerr
		}

		// commit
		return nil
	})

	return err
}
