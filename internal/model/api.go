package model

import "time"

type RecordRequest struct {
	CategoryID int    `json:"category_id"`
	Datetime   string `json:"datetime"` // YYYYMMDD
	From       string `json:"from"`
	Type       string `json:"type"`
	Memo       string `json:"memo"`
}

type Recordstruct struct {
	ID         int       `json:"id"`
	CategoryID int       `json:"category_id"`
	Datetime   time.Time `json:"datetime"`
	From       string    `json:"from"`
	Type       string    `json:"type"`
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
