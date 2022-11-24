package model

import "time"

type Category struct {
	ID         int64  `gorm:"id"`
	CategoryID int64  `gorm:"column:category_id"`
	Name       string `gorm:"column:id"`
}

type Record_YYYYMM struct {
	ID         int64     `gorm:"id"`
	CategoryID int64     `gorm:"column:category_id"`
	From       string    `gorm:"column:from"`
	Type       string    `gorm:"column:type"`
	CreatedAt  time.Time `gorm:"column:created_at"`
	UpdatedAt  time.Time `gorm:"column:updated_at"`
}
