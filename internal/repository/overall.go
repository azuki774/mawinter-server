package repository

import (
	"gorm.io/gorm"
)

type DBRepository struct {
	Conn *gorm.DB
}
