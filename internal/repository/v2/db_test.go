package repository

import (
	"fmt"
	"mawinter-server/internal/model"
	"mawinter-server/internal/openapi"
	"os"
	"reflect"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func NewDBMock() (*gorm.DB, sqlmock.Sqlmock, error) {
	sqlDB, mock, err := sqlmock.New()
	if err != nil {
		os.Exit(1)
	}
	mockDB, err := gorm.Open(mysql.New(mysql.Config{
		Conn:                      sqlDB,
		DriverName:                "mysql",
		SkipInitializeWithVersion: true,
	}), &gorm.Config{})
	return mockDB, mock, err
}

func TestDBRepository_GetMonthMidSummary(t *testing.T) {
	type fields struct {
		Conn *gorm.DB
	}
	type args struct {
		yyyymm string
	}
	tests := []struct {
		name       string
		fields     fields
		args       args
		wantSummon []model.CategoryMidMonthSummary
		wantErr    bool
		mockSetUp  func(mock sqlmock.Sqlmock)
	}{
		{
			name:   "OK",
			fields: fields{},
			args:   args{yyyymm: "199903"},
			wantSummon: []model.CategoryMidMonthSummary{
				{
					CategoryId: 100,
					Count:      1,
					Price:      10000,
				},
				{
					CategoryId: 200,
					Count:      2,
					Price:      20000,
				},
				{
					CategoryId: 300,
					Count:      3,
					Price:      30000,
				},
			},
			wantErr: false,
			mockSetUp: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT`).
					WillReturnRows(sqlmock.NewRows([]string{"category_id", "count(1)", "sum(price)"}).
						AddRow("100", "1", "10000").
						AddRow("200", "2", "20000").
						AddRow("300", "3", "30000"))
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gormDB, mock, _ := NewDBMock()
			tt.fields.Conn = gormDB

			d := &DBRepository{
				Conn: tt.fields.Conn,
			}

			tt.mockSetUp(mock)

			gotSummon, err := d.GetMonthMidSummary(tt.args.yyyymm)
			if (err != nil) != tt.wantErr {
				t.Errorf("DBRepository.GetMonthMidSummary() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotSummon, tt.wantSummon) {
				t.Errorf("DBRepository.GetMonthMidSummary() = %v, want %v", gotSummon, tt.wantSummon)
			}
		})
	}
}

func TestDBRepository_MakeCategoryNameMap(t *testing.T) {
	type fields struct {
		Conn *gorm.DB
	}
	tests := []struct {
		name      string
		fields    fields
		wantCnf   map[int]string
		wantErr   bool
		mockSetUp func(mock sqlmock.Sqlmock)
	}{
		{
			name:    "get info",
			fields:  fields{},
			wantCnf: map[int]string{100: "カテゴリ1", 200: "カテゴリ2", 300: "カテゴリ3"},
			wantErr: false,
			mockSetUp: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT").
					WillReturnRows(sqlmock.NewRows([]string{"id", "category_id", "name"}).
						AddRow("1", "100", "カテゴリ1").
						AddRow("2", "200", "カテゴリ2").
						AddRow("3", "300", "カテゴリ3"))
			},
		},
		{
			name:    "error",
			fields:  fields{},
			wantCnf: nil,
			wantErr: true,
			mockSetUp: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT").
					WillReturnError(gorm.ErrInvalidDB)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gormDB, mock, _ := NewDBMock()
			tt.fields.Conn = gormDB

			d := &DBRepository{
				Conn: tt.fields.Conn,
			}

			tt.mockSetUp(mock)

			gotCnf, err := d.MakeCategoryNameMap()
			if (err != nil) != tt.wantErr {
				t.Errorf("DBRepository.MakeCategoryNameMap() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotCnf, tt.wantCnf) {
				t.Errorf("DBRepository.MakeCategoryNameMap() = %v, want %v", gotCnf, tt.wantCnf)
			}
		})
	}
}

func TestDBRepository_GetMonthlyFixDone(t *testing.T) {
	type fields struct {
		Conn *gorm.DB
	}
	type args struct {
		yyyymm string
	}
	tests := []struct {
		name      string
		fields    fields
		args      args
		wantDone  bool
		wantErr   bool
		mockSetUp func(mock sqlmock.Sqlmock)
	}{
		{
			name:   "none",
			fields: fields{},
			args: args{
				yyyymm: "200001",
			},
			wantDone: false,
			wantErr:  false,
			mockSetUp: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT").WillReturnError(gorm.ErrRecordNotFound)
			},
		},
		{
			name:   "already registed",
			fields: fields{},
			args: args{
				yyyymm: "200002",
			},
			wantDone: true,
			wantErr:  false,
			mockSetUp: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT").
					WillReturnRows(sqlmock.NewRows([]string{"yyyymm", "done"}).
						AddRow("200002", "1"))
			},
		},
		{
			name:   "already registed",
			fields: fields{},
			args: args{
				yyyymm: "200003",
			},
			wantDone: false,
			wantErr:  true,
			mockSetUp: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT").
					WillReturnError(fmt.Errorf("error"))
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gormDB, mock, _ := NewDBMock()
			tt.fields.Conn = gormDB

			d := &DBRepository{
				Conn: tt.fields.Conn,
			}

			tt.mockSetUp(mock)

			gotDone, err := d.GetMonthlyFixDone(tt.args.yyyymm)
			if (err != nil) != tt.wantErr {
				t.Errorf("DBRepository.GetMonthlyFixDone() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotDone != tt.wantDone {
				t.Errorf("DBRepository.GetMonthlyFixDone() = %v, want %v", gotDone, tt.wantDone)
			}
		})
	}
}

func TestDBRepository_InsertMonthlyFixBilling(t *testing.T) {
	type fields struct {
		Conn *gorm.DB
	}
	type args struct {
		yyyymm string
	}
	tests := []struct {
		name      string
		fields    fields
		args      args
		wantRecs  []openapi.Record
		wantErr   bool
		mockSetUp func(mock sqlmock.Sqlmock)
	}{
		{
			name:   "ok",
			fields: fields{},
			args: args{
				yyyymm: "202102",
			},
			wantErr: false,
			wantRecs: []openapi.Record{
				{
					CategoryId:   100,
					CategoryName: "cat1",
					Datetime:     time.Date(2021, 2, 15, 0, 0, 0, 0, jst),
					From:         "fixmonth",
					Id:           1,
					Memo:         "",
					Price:        1234,
					Type:         "",
				},
				{
					CategoryId:   200,
					CategoryName: "cat2",
					Datetime:     time.Date(2021, 2, 25, 0, 0, 0, 0, jst),
					From:         "fixmonth",
					Id:           2,
					Memo:         "",
					Price:        12345,
					Type:         "",
				},
			},
			mockSetUp: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectQuery("SELECT").
					WillReturnRows(sqlmock.NewRows([]string{"category_id", "day", "price"}).
						AddRow("100", "15", "1234").
						AddRow("200", "25", "12345"))
				mock.ExpectExec(regexp.QuoteMeta(
					"INSERT INTO `Monthly_Fix_Done`")).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectExec(regexp.QuoteMeta(
					"INSERT INTO ")).
					WillReturnResult(sqlmock.NewResult(1, 2))
				mock.ExpectCommit()
				mock.ExpectQuery("SELECT").
					WillReturnRows(sqlmock.NewRows([]string{"id", "category_id", "name"}).
						AddRow("1", "100", "cat1").
						AddRow("2", "200", "cat2"))
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gormDB, mock, _ := NewDBMock()
			tt.fields.Conn = gormDB

			d := &DBRepository{
				Conn: tt.fields.Conn,
			}

			tt.mockSetUp(mock)

			gotRecs, err := d.InsertMonthlyFixBilling(tt.args.yyyymm)
			if (err != nil) != tt.wantErr {
				t.Errorf("DBRepository.InsertMonthlyFixBilling() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotRecs, tt.wantRecs) {
				t.Errorf("DBRepository.InsertMonthlyFixBilling() = %v, want %v", gotRecs, tt.wantRecs)
			}
		})
	}
}
