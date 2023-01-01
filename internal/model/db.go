package model

import "time"

func (Category) TableName() string {
	return "Category"
}

func (MonthlyFixBillingDB) TableName() string {
	return "Monthly_Fix_Billing"
}

func (MonthlyFixDoneDB) TableName() string {
	return "Monthly_Fix_Done"
}

type Category struct {
	ID         int64  `gorm:"id"`
	CategoryID int64  `gorm:"column:category_id"`
	Name       string `gorm:"column:name"`
}

type Record_YYYYMM struct {
	ID         int64     `gorm:"id;primaryKey"`
	CategoryID int64     `gorm:"column:category_id;index;not null"`
	Datetime   time.Time `gorm:"column:datetime;autoCreateTime;index;not null"`
	Price      int64     `gorm:"column:price"`
	From       string    `gorm:"column:from"`
	Type       string    `gorm:"column:type"`
	Memo       string    `gorm:"column:memo"`
	CreatedAt  time.Time `gorm:"column:created_at"`
	UpdatedAt  time.Time `gorm:"column:updated_at"`
}

type MonthlyFixBillingDB struct {
	ID         int64     `gorm:"id;primaryKey"`
	CategoryID int64     `gorm:"column:category_id"`
	Day        int64     `gorm:"column:day"`
	Price      int64     `gorm:"column:price"`
	Type       string    `gorm:"column:type"`
	Memo       string    `gorm:"column:memo"`
	CreatedAt  time.Time `gorm:"column:created_at"`
	UpdatedAt  time.Time `gorm:"column:updated_at"`
}

type MonthlyFixDoneDB struct {
	YYYYMM    string    `gorm:"column:yyyymm;primaryKey"`
	Done      uint8     `gorm:"column:done"`
	CreatedAt time.Time `gorm:"column:created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at"`
}

type SumPriceCategoryID struct {
	CategoryID int64
	Count      int64
	Sum        int64
}
