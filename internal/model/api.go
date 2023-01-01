package model

import (
	"strconv"
	"time"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
)

var jst *time.Location

func init() {
	j, err := time.LoadLocation("Asia/Tokyo")
	if err != nil {
		panic(err)
	}
	jst = j
}

type RecordRequest struct {
	CategoryID int    `json:"category_id"`
	Datetime   string `json:"datetime"` // YYYYMMDD
	From       string `json:"from"`
	Type       string `json:"type"`
	Price      int    `json:"price"`
	Memo       string `json:"memo"`
}

type Recordstruct struct {
	ID         int       `json:"id"`
	CategoryID int       `json:"category_id"`
	Datetime   time.Time `json:"datetime"`
	From       string    `json:"from"`
	Type       string    `json:"type"`
	Price      int       `json:"price"`
	Memo       string    `json:"memo"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type CategoryMonthSummary struct {
	CategoryID   int
	CategoryName string
	Count        int
	Total        int
}

type CategoryYearSummary struct {
	CategoryID   int    `json:"category_id"`
	CategoryName string `json:"category_name"`
	MonthPrice   []int  `json:"price"` // 4月から3月までの数値が配列で返る
	Total        int    `json:"total"`
}

type MonthlyFixBilling struct {
	CategoryID int
	Day        int
	Price      int
	Type       string
	Memo       string
}

func NewCategoryYearSummary(cats []Category) (caty []*CategoryYearSummary) {
	for _, cat := range cats {
		n := &CategoryYearSummary{
			CategoryID:   int(cat.CategoryID),
			CategoryName: cat.Name,
			MonthPrice:   make([]int, 0),
			Total:        0,
		}
		caty = append(caty, n)
	}
	return caty
}

func NewRecordFromReq(req RecordRequest) (record Recordstruct, err error) {
	record = Recordstruct{
		CategoryID: req.CategoryID,
		// Datetime   time.Time `json:"datetime"`
		From:  req.From,
		Type:  req.Type,
		Price: req.Price,
		Memo:  req.Memo,
	}
	if req.Datetime == "" {
		// now time
		record.Datetime = time.Now()

	} else if validation.Validate(req.Datetime, validation.Length(8, 8), is.Digit) == nil {
		// YYYYMMDD
		record.Datetime, err = time.ParseInLocation("20060102", req.Datetime, jst)
		if err != nil {
			return Recordstruct{}, err
		}
	} else {
		// 2022-11-22T13:54:08+09:00
		record.Datetime, err = time.ParseInLocation(time.RFC3339, req.Datetime, jst)
		if err != nil {
			return Recordstruct{}, err
		}
	}

	return record, nil
}

func (c *CategoryYearSummary) AddMonthPrice(price int) {
	c.MonthPrice = append(c.MonthPrice, price)
	c.Total = c.Total + price
}

func (m *MonthlyFixBilling) ConvAddDBModel(yyyymm string) (Record_YYYYMM, error) {
	yyyynum, err := strconv.Atoi(yyyymm[0:4])
	if err != nil {
		return Record_YYYYMM{}, err
	}

	mmnum, err := strconv.Atoi(yyyymm[5:6])
	if err != nil {
		return Record_YYYYMM{}, err
	}

	return Record_YYYYMM{
		CategoryID: int64(m.CategoryID),
		Datetime:   time.Date(yyyynum, time.Month(mmnum), m.Day, 0, 0, 0, 0, jst),
		From:       "fix-monthly", // 固定値
		Type:       m.Type,
		Memo:       m.Memo,
	}, nil
}
