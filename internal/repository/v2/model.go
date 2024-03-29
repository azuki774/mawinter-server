package repository

import (
	"mawinter-server/internal/model"
	"mawinter-server/internal/openapi"
	"mawinter-server/internal/timeutil"
	"time"
)

var jst *time.Location

func init() {
	j, err := time.LoadLocation("Asia/Tokyo")
	if err != nil {
		panic(err)
	}
	jst = j
}

func yyyymmToInitDayTime(yyyymm string) (t time.Time, err error) {
	// yyyymm -> time.Time
	// yyyymm から その月の1日目の time.Time を返す
	t, err = time.ParseInLocation("200601", yyyymm, jst)
	if err != nil {
		return time.Time{}, err
	}
	return t, nil
}

func NewDBModelRecord(req openapi.ReqRecord) (rec model.Record, err error) {
	// ID is not set
	rec.CategoryID = int64(req.CategoryId)

	if req.Datetime != nil {
		// YYYYMMDD
		rec.Datetime, err = time.ParseInLocation("20060102", *req.Datetime, jst)
		if err != nil {
			return model.Record{}, nil
		}
	} else {
		// default
		rec.Datetime = timeutil.NowFunc()
	}

	rec.Price = int64(req.Price)
	if req.From != nil {
		rec.From = *req.From
	} else {
		rec.From = ""
	}

	if req.Type != nil {
		rec.Type = *req.Type
	} else {
		rec.Type = ""
	}

	if req.Memo != nil {
		rec.Memo = *req.Memo
	} else {
		rec.Memo = ""
	}

	// CreatedAt  time.Time `gorm:"column:created_at"`
	// UpdatedAt  time.Time `gorm:"column:updated_at"`

	return rec, nil
}

// NewRecordFromDB では Record テーブルをもとに、API Structを出力する。
func NewRecordFromDB(req model.Record) (rec openapi.Record, err error) {
	rec = openapi.Record{
		CategoryId: int(req.CategoryID),
		// CategoryName: req.CategoryName : ここでは取得しない
		Datetime: req.Datetime,
		From:     req.From,
		Id:       int(req.ID),
		Memo:     req.Memo,
		Price:    int(req.Price),
		Type:     req.Type,
	}
	return rec, nil
}
