package api

import (
	"mawinter-server/internal/model"
	"time"
)

// 正常系
type mockDBRepositry1 struct{}

func (dbR *mockDBRepositry1) CreateRecordDB(record model.Records) (retAddRecord model.ShowRecord, err error) {
	return model.ShowRecord{Id: 1001, CategoryID: 101, CategoryName: "cat1", Date: time.Date(2000, 1, 23, 0, 0, 0, 0, time.FixedZone("Asia/Tokyo", 9*60*60)), Price: 10000}, nil
}

func (dbR *mockDBRepositry1) GetYearSummaryDB(year int64) (yearSummaryInters []model.YearSummaryInter, err error) {
	mock := []model.YearSummaryInter{
		{
			CategoryID: int64(101),
			Name:       "testcatName1",
			YearMonth:  "200004",
			Total:      int64(12345),
		},
		{
			CategoryID: int64(101),
			Name:       "testcatName1",
			YearMonth:  "200103",
			Total:      int64(34567),
		},
		{
			CategoryID: int64(102),
			Name:       "testcatName2",
			YearMonth:  "200008",
			Total:      23456,
		},
	}
	return mock, nil
}

func (dbR *mockDBRepositry1) GetRecentRecord(n int) (getRecentData []model.ShowRecord, err error) {
	getRecentData = []model.ShowRecord{
		{
			Id:           1,
			CategoryID:   100,
			CategoryName: "cat1",
			Date:         time.Date(2000, 1, 2, 0, 0, 0, 0, time.FixedZone("Asia/Tokyo", 9*60*60)),
			Price:        1000,
			Memo:         "",
		},
		{
			Id:           2,
			CategoryID:   200,
			CategoryName: "cat2",
			Date:         time.Date(2000, 1, 1, 0, 0, 0, 0, time.FixedZone("Asia/Tokyo", 9*60*60)),
			Price:        2000,
			Memo:         "memo2",
		},
	}

	return getRecentData, nil
}

func (dbR *mockDBRepositry1) DeleteRecordDB(id int64) (err error) {
	return nil
}

func (dbR *mockDBRepositry1) CloseDB() (err error) {
	return nil
}
