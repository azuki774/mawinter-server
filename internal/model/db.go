package model

import "time"

type Category struct {
	ID         int64  `gorm:"id"`
	CategoryID int64  `gorm:"column:category_id"`
	Name       string `gorm:"column:id"`
}

type Record_YYYYMM struct {
	ID         int64     `gorm:"id,primaryKey"`
	CategoryID int64     `gorm:"column:category_id,index,not null"`
	Datetime   time.Time `gorm:"column:datetime,autoCreateTime,index,not null"`
	From       string    `gorm:"column:from"`
	Type       string    `gorm:"column:type"`
	Memo       string    `gorm:"column:memo"`
	CreatedAt  time.Time `gorm:"column:created_at"`
	UpdatedAt  time.Time `gorm:"column:updated_at"`
}

type SumPriceCategoryID struct {
	CategoryID int64 `gorm:"column:category_id"`
	Count      int64 `gorm:"column:count(*)"`
	Sum        int64 `gorm:"column:sum(price)"`
}
