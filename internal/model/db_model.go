package model

import "time"

type Types struct {
	Id         int64  `gorm:"primaryKey"`
	TypeString string `gorm:"column:type"`
}

type Categories struct {
	Id         int64 `gorm:"primaryKey"`
	CategoryID int64 `gorm:"column:category_id"`
	Name       string
	TypeString string `gorm:"column:type"`
}

type Records struct {
	Id         int64 `gorm:"primaryKey"`
	CategoryID int64 `gorm:"column:category_id"`
	Date       time.Time
	Price      int64
	Memo       string `gorm:"default:null"`
}

type CategoriesRecords struct {
	Records
	Categories
}
