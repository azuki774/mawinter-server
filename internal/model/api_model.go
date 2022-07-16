package model

import "time"

type CreateRecord struct {
	CategoryID int64  `json:"category_id"`
	Date       string `json:"date"` // YYYYMMDD
	Price      int64  `json:"price"`
	Memo       string `json:"memo"`
}

type YearSummary struct {
	CategoryID int64   `json:"category_id"`
	Name       string  `json:"category_name"`
	Price      []int64 `json:"price"`
	Total      int64   `json:"total"`
}

type YearSummaryInter struct {
	CategoryID int64
	Name       string
	YearMonth  string // YYYYMM
	Total      int64
}

type ShowRecord struct {
	Id           int64     `json:"id"`
	CategoryID   int64     `json:"category_id"`
	CategoryName string    `json:"category_name"`
	Date         time.Time `json:"date"`
	Price        int64     `json:"price"`
	Memo         string    `json:"memo"`
}
