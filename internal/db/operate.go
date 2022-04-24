package db

import (
	"errors"
	"fmt"
	"mawinter-expense/internal/azerror"
	l "mawinter-expense/internal/logger"

	"github.com/jinzhu/gorm"
	"go.uber.org/zap"
)

func TransMonthToIndex(month int) (index int) {
	// 4 -> 0,  12 -> 8, 1 -> 9,  3 -> 11
	if month <= 3 {
		return month + 8
	} else {
		return month - 4
	}
}

// GetYearSummaryDB : category_id, name, year-month, sum をidの昇順で返す
func (dbR *dbRepository) GetYearSummaryDB(tx *gorm.DB, year int64) (yearSummary []GetYearSummaryDBStruct, err error) {
	beginDate := fmt.Sprintf("%d-04-01", year)
	endDate := fmt.Sprintf("%d-03-31", year+1)
	sqlText := "SELECT records.category_id , categories.name, DATE_FORMAT(records.date, '%Y%m'), sum(price) FROM records INNER JOIN categories ON categories.category_id = records.category_id WHERE records.date BETWEEN '" + beginDate + "' AND '" + endDate + "' GROUP BY records.category_id , DATE_FORMAT(records.date, '%Y%m') ORDER BY records.category_id;"

	l.Logger.Info("SQL", "GetYearSummaryDB")

	rows, err := tx.Raw(sqlText).Rows()
	defer rows.Close()
	if err != nil {
		l.Logger.Error("SQL", "GetYearSummaryDB", zap.Error(err))
		return nil, azerror.ErrInternal
	}

	for rows.Next() {
		var newSummary GetYearSummaryDBStruct
		if err := rows.Scan(&newSummary.CategoryId, &newSummary.Name, &newSummary.YearMonth, &newSummary.Price); err != nil {
			l.Logger.Error("SQL", "GetYearSummaryDB", zap.Error(err))
			return nil, azerror.ErrInternal
		}
		yearSummary = append(yearSummary, newSummary)
	}

	return yearSummary, nil
}

func (dbR *dbRepository) AddRecordDB(tx *gorm.DB, record Records) (retRecords Records, err error) {
	l.Logger.Info("SQL", "AddRecordDB")
	// Web API側でDateを埋めておく
	res := tx.Select("CategoryId", "Price", "Date", "Memo").Create(&record)

	if res.Error != nil {
		l.Logger.Error("SQL", "AddRecordDB", zap.Error(res.Error))
		return Records{}, azerror.ErrInternal
	}

	return record, nil
}

func (dbR *dbRepository) DeleteRecordDB(tx *gorm.DB, id int64) (err error) {
	l.Logger.Info("SQL", "DeleteRecordDB")
	record := Records{Id: id}
	res := tx.First(&record)
	if res.Error != nil {
		l.Logger.Error("SQL", "DeleteRecordDB", zap.Error(res.Error))
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return azerror.ErrRecordNotFound
		} else {
			return azerror.ErrInternal
		}
	}

	res = tx.Delete(&record)
	if res.Error != nil {
		l.Logger.Error("SQL", "DeleteRecordDB", zap.Error(res.Error))
		return azerror.ErrInternal
	}

	return nil
}

// 直近n 個のレコードを取得するSQLを発行する。
func (dbR *dbRepository) GetRecentRecord(tx *gorm.DB, n int64) (getRecentData []RecordsDetails, err error) {
	l.Logger.Info("SQL", "GetRecentRecord")
	sqlText := "SELECT records.id, records.date, records.category_id, categories.name, records.price, records.memo FROM records LEFT OUTER JOIN categories ON records.category_id = categories.category_id ORDER BY records.id DESC LIMIT ?"
	rows, err := tx.Raw(sqlText, n).Rows()
	if err != nil {
		l.Logger.Error("SQL", "GetRecentRecord", zap.Error(err))
		return nil, azerror.ErrInternal
	}
	defer rows.Close()

	for rows.Next() {
		var aRecentData RecordsDetails
		err := rows.Scan(&aRecentData.Id, &aRecentData.Date, &aRecentData.CategoryId, &aRecentData.Name, &aRecentData.Price, &aRecentData.Memo)
		if err != nil {
			l.Logger.Error("SQL", "GetRecentRecord", zap.Error(err))
			return nil, azerror.ErrInternal
		}
		getRecentData = append(getRecentData, aRecentData)
	}

	return getRecentData, nil
}
