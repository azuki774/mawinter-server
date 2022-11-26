package repository

import (
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
