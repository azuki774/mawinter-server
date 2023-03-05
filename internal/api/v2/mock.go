package api

import (
	"mawinter-server/internal/model"
	"mawinter-server/internal/openapi"
	"time"
)

type mockRepo struct {
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

func (m *mockRepo) GetMonthRecords(yyyymm string) (recs []openapi.Record, err error) {
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
