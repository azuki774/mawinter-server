//go:build integration
// +build integration

package db

import (
	"database/sql"
	httpdate "github.com/Songmu/go-httpdate"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	l "mawinter-expense/internal/logger"
	"os"
	"reflect"
	"testing"
	"time"
)

var dbR DBRepository

func TestMain(m *testing.M) {
	l.NewSugarLogger()

	l.Logger.Info("DB Setup")
	err := DBConnect("root", "password", "localhost", "mawinter-test")
	if err != nil {
		l.Logger.Fatal("DB Setup failed")
		panic(err)
	}

	dbR = NewDBRepository(DB)

	// ローカル実行のため、records初期化
	DB.Raw("TRUNCATE TABLE records")
	l.Logger.Info("DB Setup Completed")

	code := m.Run()

	os.Exit(code)
}

func TestDBConnect(t *testing.T) {
	type args struct {
		user     string
		password string
		address  string
		dbName   string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "TestDB connect",
			args:    args{user: "root", password: "password", address: "localhost", dbName: "mawinter-test"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := DBConnect(tt.args.user, tt.args.password, tt.args.address, tt.args.dbName); (err != nil) != tt.wantErr {
				t.Errorf("DBConnect() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_dbRepository_AddRecordDB(t *testing.T) {
	date3, _ := httpdate.Str2Time("2010-01-01", nil)
	tx := dbR.OpenTx()
	type args struct {
		tx     *gorm.DB
		record Records
	}
	tests := []struct {
		name    string
		args    args
		wantId  int64
		wantErr bool
	}{
		{
			name:    "#1 NowTime + NonMemo",
			args:    args{tx: tx, record: Records{CategoryId: 100, Price: 10000, Date: time.Now()}},
			wantErr: false,
		},
		{
			name:    "#2 NowTime + Memo",
			args:    args{tx: tx, record: Records{CategoryId: 200, Price: 20000, Memo: sql.NullString{String: "Memo", Valid: true}, Date: time.Now()}},
			wantErr: false,
		},
		{
			name:    "#3 OldTime + Memo",
			args:    args{tx: tx, record: Records{CategoryId: 300, Price: 30000, Date: date3, Memo: sql.NullString{String: "Memo", Valid: true}}},
			wantErr: false,
		},
		{
			name:    "#4 Not Exists CategoryId",
			args:    args{tx: tx, record: Records{CategoryId: 777, Price: 10000}},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := dbR.AddRecordDB(tt.args.tx, tt.args.record)
			if (err != nil) != tt.wantErr {
				t.Errorf("dbRepository.AddRecordDB() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
	dbR.CloseTx(tx, nil)
}

func Test_dbRepository_DeleteRecordDB(t *testing.T) {
	tx := dbR.OpenTx()
	record1 := Records{Id: 1201, CategoryId: 200, Price: 1000}
	tx.Select("Id", "CategoryId", "Price").Create(&record1)
	type args struct {
		tx *gorm.DB
		id int64
	}
	tests := []struct {
		name    string
		dbR     *dbRepository
		args    args
		wantErr bool
	}{
		{
			name:    "#1 Normal",
			args:    args{tx: tx, id: 1201},
			wantErr: false,
		},
		{
			name:    "#2 NotFound",
			args:    args{tx: tx, id: 12345},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.dbR.DeleteRecordDB(tt.args.tx, tt.args.id); (err != nil) != tt.wantErr {
				t.Errorf("dbRepository.DeleteRecordDB() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
	dbR.CloseTx(tx, nil)
}

func Test_dbRepository_GetRecentRecord(t *testing.T) {
	tx := dbR.OpenTx()
	type args struct {
		tx *gorm.DB
		n  int64
	}
	tests := []struct {
		name              string
		args              args
		wantGetRecentData []RecordsDetails
		wantErr           bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := dbR.GetRecentRecord(tt.args.tx, tt.args.n)
			if (err != nil) != tt.wantErr {
				t.Errorf("dbRepository.GetRecentRecord() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			// if !reflect.DeepEqual(gotGetRecentData, tt.wantGetRecentData) {
			// 	t.Errorf("dbRepository.GetRecentRecord() = %v, want %v", gotGetRecentData, tt.wantGetRecentData)
			// }
		})
	}
	dbR.CloseTx(tx, nil)
}

func Test_dbRepository_GetYearSummaryDB(t *testing.T) {
	tx := dbR.OpenTx()
	// Make test case
	date1, _ := httpdate.Str2Time("1999-04-01", nil)
	date1 = date1.Add(9 * time.Hour)
	record1 := Records{Id: 1001, CategoryId: 100, Date: date1, Price: 1000}
	tx.Select("CategoryId", "Price", "Date").Create(&record1)

	date2, _ := httpdate.Str2Time("1999-04-11", nil)
	date2 = date2.Add(9 * time.Hour)
	record2 := Records{Id: 1002, CategoryId: 200, Date: date2, Price: 2000}
	tx.Select("CategoryId", "Price", "Date").Create(&record2)

	date3, _ := httpdate.Str2Time("2000-05-01", nil)
	date3 = date3.Add(9 * time.Hour)
	record3 := Records{Id: 1003, CategoryId: 300, Date: date3, Price: 3000}
	tx.Select("CategoryId", "Price", "Date").Create(&record3)

	date4, _ := httpdate.Str2Time("2000-05-11", nil)
	date4 = date4.Add(9 * time.Hour)
	record4 := Records{Id: 1004, CategoryId: 300, Date: date4, Price: 4000}
	tx.Select("CategoryId", "Price", "Date").Create(&record4)

	answer1 := []GetYearSummaryDBStruct{GetYearSummaryDBStruct{CategoryId: 100, Name: "月給", YearMonth: "199904", Price: 1000}, GetYearSummaryDBStruct{CategoryId: 200, Name: "家賃", YearMonth: "199904", Price: 2000}}
	answer2 := []GetYearSummaryDBStruct{GetYearSummaryDBStruct{CategoryId: 300, Name: "保険・税金", YearMonth: "200005", Price: 7000}}

	type args struct {
		tx   *gorm.DB
		year int64
	}
	tests := []struct {
		name            string
		args            args
		wantYearSummary []GetYearSummaryDBStruct
		wantErr         bool
	}{
		{
			name:            "#1 Not Sum Pattern",
			args:            args{tx: tx, year: 1999},
			wantYearSummary: answer1,
			wantErr:         false,
		},
		{
			name:            "#2 Sum Pattern",
			args:            args{tx: tx, year: 2000},
			wantYearSummary: answer2,
			wantErr:         false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotYearSummary, err := dbR.GetYearSummaryDB(tt.args.tx, tt.args.year)
			if (err != nil) != tt.wantErr {
				t.Errorf("dbRepository.GetYearSummaryDB() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotYearSummary, tt.wantYearSummary) {
				t.Errorf("dbRepository.GetYearSummaryDB() = %v, want %v", gotYearSummary, tt.wantYearSummary)
			}
		})
	}
	dbR.CloseTx(tx, nil)
}
