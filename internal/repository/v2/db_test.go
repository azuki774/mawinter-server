package repository

import (
	"mawinter-server/internal/model"
	"os"
	"reflect"
	"testing"

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
