package api

import (
	"context"
	"fmt"
	"mawinter-server/internal/model"
	"mawinter-server/internal/openapi"
	"mawinter-server/internal/timeutil"
	"time"
)

type mockRepo struct {
	GetMonthlyFixDoneReturn bool
	ReturnConfirm           bool // ここが true ならば各テーブルは confirm していることにする
}

func (m *mockRepo) CreateTableYYYYMM(yyyymm string) (err error) {
	return nil
}

func (m *mockRepo) InsertRecord(req openapi.ReqRecord) (rec openapi.Record, err error) {
	rec = openapi.Record{
		CategoryId: 100,
		// CategoryName string    `json:"category_name"`
		Datetime: time.Date(2000, 1, 23, 0, 0, 0, 0, jst),
		From:     "from",
		Id:       1,
		Memo:     "memo",
		Price:    1234,
		Type:     "type",
	}
	return rec, nil
}

// GetRecords は mock だと num <= 2 && offset = 0 と num = 20 && offset = 1 までしか対応しない
func (m *mockRepo) GetRecords(ctx context.Context, num int, offset int) (recs []openapi.Record, err error) {
	if offset == 1 {
		recs = []openapi.Record{
			{
				CategoryId: 200,
				// CategoryName string    `json:"category_name"`
				Datetime: time.Date(2000, 1, 25, 0, 0, 0, 0, jst),
				From:     "",
				Id:       2,
				Memo:     "",
				Price:    2345,
				Type:     "",
			},
		}
	} else {
		if num >= 2 {
			recs = []openapi.Record{
				{
					CategoryId: 100,
					// CategoryName string    `json:"category_name"`
					Datetime: time.Date(2000, 1, 23, 0, 0, 0, 0, jst),
					From:     "from",
					Id:       1,
					Memo:     "memo",
					Price:    1234,
					Type:     "type",
				},
				{
					CategoryId: 200,
					// CategoryName string    `json:"category_name"`
					Datetime: time.Date(2000, 1, 25, 0, 0, 0, 0, jst),
					From:     "",
					Id:       2,
					Memo:     "",
					Price:    2345,
					Type:     "",
				},
			}
		} else if num == 1 {
			recs = []openapi.Record{
				{
					CategoryId: 100,
					// CategoryName string    `json:"category_name"`
					Datetime: time.Date(2000, 1, 23, 0, 0, 0, 0, jst),
					From:     "from",
					Id:       1,
					Memo:     "memo",
					Price:    1234,
					Type:     "type",
				},
			}
		} else if num == 0 {
			recs = []openapi.Record{}
		} else {
			return []openapi.Record{}, fmt.Errorf("invalid args")
		}
	}

	return recs, nil
}

func (m *mockRepo) GetRecordsCount(ctx context.Context) (num int, err error) {
	return 123, nil // 正常系
}

// empty return
func (m *mockRepo) GetCategories(ctx context.Context) (cats []model.Category, err error) {
	return []model.Category{}, nil
}

func (m *mockRepo) GetMonthRecords(yyyymm string) (recs []openapi.Record, err error) {
	recs = []openapi.Record{
		{
			CategoryId: 100,
			// CategoryName string    `json:"category_name"`
			Datetime: time.Date(2000, 1, 23, 0, 0, 0, 0, jst),
			From:     "ope",
			Id:       1,
			Memo:     "memo",
			Price:    1234,
			Type:     "type",
		},
		{
			CategoryId: 200,
			// CategoryName string    `json:"category_name"`
			Datetime: time.Date(2000, 1, 25, 0, 0, 0, 0, jst),
			From:     "mawinter-web",
			Id:       2,
			Memo:     "",
			Price:    2345,
			Type:     "",
		},
	}

	return recs, nil
}

func (m *mockRepo) GetMonthRecordsRecent(yyyymm string, num int) (recs []openapi.Record, err error) {
	recs = []openapi.Record{
		{
			CategoryId: 100,
			// CategoryName string    `json:"category_name"`
			Datetime: time.Date(2000, 1, 23, 0, 0, 0, 0, jst),
			From:     "from",
			Id:       1,
			Memo:     "memo",
			Price:    1234,
			Type:     "type",
		},
		{
			CategoryId: 200,
			// CategoryName string    `json:"category_name"`
			Datetime: time.Date(2000, 1, 25, 0, 0, 0, 0, jst),
			From:     "",
			Id:       2,
			Memo:     "",
			Price:    2345,
			Type:     "",
		},
	}

	return recs, nil
}

func (m *mockRepo) MakeCategoryNameMap() (cnf map[int]string, err error) {
	cnf = map[int]string{100: "cat1", 200: "cat2"}
	return cnf, nil
}

func (m *mockRepo) GetMonthMidSummary(yyyymm string) (summon []model.CategoryMidMonthSummary, err error) {
	// テストのため、200008 の場合のみ catID: 200 は ノーレコードとする。
	if yyyymm == "200008" {
		return []model.CategoryMidMonthSummary{
			{
				CategoryId: 100,
				Count:      10,
				Price:      1000,
			},
		}, nil
	}

	return []model.CategoryMidMonthSummary{
		{
			CategoryId: 100,
			Count:      10,
			Price:      1000,
		},
		{
			CategoryId: 200,
			Count:      20,
			Price:      2000,
		},
	}, nil
}

func (m *mockRepo) InsertMonthlyFixBilling(yyyymm string) (recs []openapi.Record, err error) {
	return []openapi.Record{
		{
			CategoryId:   100,
			CategoryName: "cat1",
			Datetime:     time.Date(2021, 2, 15, 0, 0, 0, 0, jst),
			From:         "fixmonth",
			Id:           1,
			Memo:         "",
			Price:        1234,
			Type:         "",
		},
		{
			CategoryId:   200,
			CategoryName: "cat2",
			Datetime:     time.Date(2021, 2, 25, 0, 0, 0, 0, jst),
			From:         "fixmonth",
			Id:           2,
			Memo:         "",
			Price:        12345,
			Type:         "",
		},
	}, nil
}

func (m *mockRepo) GetMonthlyFixDone(yyyymm string) (done bool, err error) {
	return m.GetMonthlyFixDoneReturn, nil
}

func (m *mockRepo) GetMonthlyConfirm(yyyymm string) (yc openapi.ConfirmInfo, err error) {
	// 正常系
	t := timeutil.NowFunc() // testconfig
	return openapi.ConfirmInfo{
		ConfirmDatetime: &t,
		Status:          &m.ReturnConfirm,
		Yyyymm:          &yyyymm,
	}, nil
}

func (m *mockRepo) UpdateMonthlyConfirm(yyyymm string, confirm bool) (yc openapi.ConfirmInfo, err error) {
	// 正常系
	t := timeutil.NowFunc() // testconfig
	return openapi.ConfirmInfo{
		ConfirmDatetime: &t,
		Status:          &confirm,
		Yyyymm:          &yyyymm,
	}, nil
}
