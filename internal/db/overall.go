package db

import (
	"database/sql"
	"fmt"
	"mawinter-expense/internal/logger"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

var DB *gorm.DB

func DBConnect(user string, password string, address string, dbName string) error {
	connectCmd := user + ":" + password + "@(" + address + ")/" + dbName + "?parseTime=true"
	nDB, err := gorm.Open("mysql", connectCmd)
	if err != nil {
		fmt.Println(err)
		logger.FatalPrint("Failed to connect to the database")
		return err
	}
	DB = nDB
	logger.InfoPrint("Connect to the database")
	return nil
}

type DBRepository interface {
	OpenTx() *gorm.DB
	CloseTx(tx *gorm.DB, err error) error
	// category_id 順にその年の月ごとの合計を取得する
	GetYearSummaryDB(tx *gorm.DB, year int64) (yearSummary []GetYearSummaryDBStruct, err error)
	AddRecordDB(tx *gorm.DB, record Records) (retRecords Records, err error)
	DeleteRecordDB(tx *gorm.DB, id int64) (err error)
	GetRecentRecord(tx *gorm.DB, n int64) (getRecentData []RecordsDetails, err error)
}

type dbRepository struct {
	conn *gorm.DB
}

func NewDBRepository(conn *gorm.DB) DBRepository {
	return &dbRepository{conn: conn}
}
func (dbR *dbRepository) OpenTx() *gorm.DB {
	tx := dbR.conn.Begin()
	return tx
}

func (dbR *dbRepository) CloseTx(tx *gorm.DB, err error) error {
	if err != nil {
		logger.ErrorPrint(fmt.Sprintf("CloseTx (Rollback): %s", err.Error()))
		tx.Rollback()
	} else {
		tx.Commit()
	}
	return nil
}

type GetYearSummaryDBStruct struct {
	CategoryId int64  `json:"categoryID"`
	Name       string `json:"name"`
	YearMonth  string `json:"yearmonth"`
	Price      int64  `json:"price"`
}

type Records struct {
	Id         int64          `json:"id"`
	CategoryId int64          `json:"categoryID"`
	Date       time.Time      `json:"date"`
	Price      int64          `json:"price"`
	Memo       sql.NullString `json:"memo"`
}

type Categories struct {
	Id         int64  `json:"id"`
	CategoryId int64  `json:"categoryID"`
	Name       string `json:"name"`
	Type       int64  `json:"type"`
}

type RecordsDetails struct {
	Id         int64          `json:"id"`
	Date       string         `json:"date"`
	CategoryId int64          `json:"categoryID"`
	Name       string         `json:"name"`
	Price      int64          `json:"price"`
	Memo       sql.NullString `json:"memo"`
}
