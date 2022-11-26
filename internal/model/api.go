package model

import (
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
	CategoryID   int
	CategoryName string
	MonthPrice   []int // 4月から3月までの数値が配列で返る
	Total        int
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
	} else {
		// YYYYMMDD
		record.Datetime, err = time.ParseInLocation("20060102", req.Datetime, jst)
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
