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

func (d *DBRepository) CloseDB() (err error) {
	dbconn, err := d.Conn.DB()
	if err != nil {
		return err
	}
	return dbconn.Close()
}

func (d *DBRepository) CreateRecordTable(yyyymm string) (err error) {
	sql := fmt.Sprintf(`CREATE TABLE Record_%s (
		id int NOT NULL AUTO_INCREMENT,
		category_id int NOT NULL,
		datetime datetime NOT NULL default current_timestamp,
		from varchar(64) NOT NULL,
		type varchar(64) NOT NULL,
		price int NOT NULL,
		created_at datetime default current_timestamp,
		updated_at timestamp default current_timestamp on update current_timestamp,
		PRIMARY KEY (id),
		index idx_cat (category_id),
		index idx_date (datetime)
	  )`, yyyymm)

	err = d.Conn.Debug().Raw(sql).Error
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
