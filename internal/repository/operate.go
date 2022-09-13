package repository

import (
	"fmt"
	"mawinter-server/internal/azerror"
	"mawinter-server/internal/model"
)

func TransMonthToIndex(month int) (index int) {
	// 4 -> 0,  12 -> 8, 1 -> 9,  3 -> 11
	if month <= 3 {
		return month + 8
	} else {
		return month - 4
	}
}

func (dbR *DBRepository) CreateRecordDB(record model.Records) (retAddRecord model.ShowRecord, err error) {
	res := dbR.Conn.Select("category_id", "price", "date", "memo").Create(&record)
	if res.Error != nil {
		return model.ShowRecord{}, res.Error
	}

	var catrecs model.CategoriesRecords
	res = dbR.Conn.Table("records").Select("records.*, categories.category_id, categories.name").
		Joins("INNER JOIN categories ON records.category_id = categories.category_id").Order("date desc").Where("records.id = ?", record.Id).Take(&catrecs)
	if res.Error != nil {
		return model.ShowRecord{}, res.Error
	}

	retAddRecord.Id = record.Id
	retAddRecord.CategoryID = record.CategoryID
	retAddRecord.CategoryName = catrecs.Name
	retAddRecord.Date = record.Date
	retAddRecord.Price = record.Price
	retAddRecord.Memo = record.Memo

	return retAddRecord, nil
}

// GetYearSummaryDB : category_id, name, year-month, sum をidの昇順で返す
func (dbR *DBRepository) GetYearSummaryDB(year int64) (yearSummaryInters []model.YearSummaryInter, err error) {
	beginDate := fmt.Sprintf("%d-04-01", year)
	endDate := fmt.Sprintf("%d-03-31", year+1)
	sqlText := "SELECT records.category_id , categories.name, DATE_FORMAT(records.date, '%Y%m'), sum(price) FROM records INNER JOIN categories ON categories.category_id = records.category_id WHERE records.date BETWEEN '" + beginDate + "' AND '" + endDate + "' GROUP BY records.category_id , DATE_FORMAT(records.date, '%Y%m') ORDER BY records.category_id;"
	rows, err := dbR.Conn.Raw(sqlText).Rows()
	defer rows.Close()
	if err != nil {
		return nil, azerror.ErrInternal
	}

	for rows.Next() {
		var newSummary model.YearSummaryInter
		if err := rows.Scan(&newSummary.CategoryID, &newSummary.Name, &newSummary.YearMonth, &newSummary.Total); err != nil {
			return nil, azerror.ErrInternal
		}
		yearSummaryInters = append(yearSummaryInters, newSummary)
	}

	return yearSummaryInters, nil
}

func (dbR *DBRepository) DeleteRecordDB(id int64) (err error) {
	res := dbR.Conn.Where("id = ?", id).Delete(&model.Records{})
	if res.Error != nil {
		return res.Error
	}
	return nil
}

// 直近n 個のレコードを取得するSQLを発行する。
func (dbR *DBRepository) GetRecentRecord(n int) (getRecentData []model.ShowRecord, err error) {
	var catrecs []model.CategoriesRecords

	res := dbR.Conn.Table("records").Select("records.*, categories.category_id, categories.name").Joins("INNER JOIN categories ON records.category_id = categories.category_id").Order("date desc").Limit(n).Find(&catrecs)
	if res.Error != nil {
		return nil, azerror.ErrInternal
	}

	for _, v := range catrecs {
		aData := model.ShowRecord{
			Id:           v.Records.Id,
			CategoryID:   v.Records.CategoryID,
			CategoryName: v.Name,
			Date:         v.Records.Date,
			Price:        v.Records.Price,
			Memo:         v.Records.Memo,
		}

		getRecentData = append(getRecentData, aData)
	}

	return getRecentData, nil
}
