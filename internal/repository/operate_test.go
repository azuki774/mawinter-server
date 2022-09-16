package repository

import (
	"mawinter-server/internal/model"
	"reflect"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"gorm.io/gorm"
)

func Test_dbRepository_CreateRecordDB(t *testing.T) {
	type fields struct {
		Conn *gorm.DB
	}
	type args struct {
		record model.Records
	}
	tests := []struct {
		name             string
		fields           fields
		args             args
		wantRetAddRecord model.ShowRecord
		wantErr          bool
		mockSetUp        func(mock sqlmock.Sqlmock, args args)
	}{
		{
			name:             "create record 1",
			fields:           fields{},
			args:             args{model.Records{CategoryID: 100, Date: time.Date(2021, 1, 1, 1, 0, 0, 0, time.FixedZone("Asia/Tokyo", 9*60*60)), Price: 1000, Memo: "memo"}},
			wantRetAddRecord: model.ShowRecord{Id: 1, CategoryID: 100, CategoryName: "cat1", Date: time.Date(2021, 1, 1, 1, 0, 0, 0, time.FixedZone("Asia/Tokyo", 9*60*60)), Price: 1000, Memo: "memo"},
			wantErr:          false,
			mockSetUp: func(mock sqlmock.Sqlmock, args args) {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta(
					"INSERT INTO `records` (`category_id`,`date`,`price`,`memo`) VALUES (?,?,?,?)")).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
				mock.ExpectQuery(regexp.QuoteMeta("SELECT records.*, categories.category_id, categories.name FROM `records` INNER JOIN categories ON records.category_id = categories.category_id WHERE records.id = ?")).
					WillReturnRows(sqlmock.NewRows([]string{"Id", "category_id", "date", "price", "memo", "name"}).
						AddRow(1, 100, time.Date(2001, 1, 1, 0, 0, 0, 0, time.FixedZone("Asia/Tokyo", 9*60*60)), 1000, "", "cat1"))
			},
		},
		{
			name:             "create record 2",
			fields:           fields{},
			args:             args{model.Records{CategoryID: 110, Date: time.Date(2021, 2, 1, 1, 0, 0, 0, time.FixedZone("Asia/Tokyo", 9*60*60)), Price: 2000, Memo: "memo"}},
			wantRetAddRecord: model.ShowRecord{Id: 2, CategoryID: 110, CategoryName: "cat2", Date: time.Date(2021, 2, 1, 1, 0, 0, 0, time.FixedZone("Asia/Tokyo", 9*60*60)), Price: 2000, Memo: "memo"},
			wantErr:          false,
			mockSetUp: func(mock sqlmock.Sqlmock, args args) {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta(
					"INSERT INTO `records` (`category_id`,`date`,`price`,`memo`) VALUES (?,?,?,?)")).
					WillReturnResult(sqlmock.NewResult(2, 1))
				mock.ExpectCommit()
				mock.ExpectQuery(regexp.QuoteMeta("SELECT records.*, categories.category_id, categories.name FROM `records` INNER JOIN categories ON records.category_id = categories.category_id WHERE records.id = ?")).
					WillReturnRows(sqlmock.NewRows([]string{"Id", "category_id", "date", "price", "memo", "name"}).
						AddRow(2, 110, time.Date(2001, 1, 1, 0, 0, 0, 0, time.FixedZone("Asia/Tokyo", 9*60*60)), 2000, "memo", "cat2"))
			},
		},
		{
			name:             "create record 3 (ErrInvalidData)",
			fields:           fields{},
			args:             args{model.Records{CategoryID: 120, Date: time.Date(2021, 3, 1, 1, 0, 0, 0, time.FixedZone("Asia/Tokyo", 9*60*60)), Price: 2000, Memo: "abcde"}},
			wantRetAddRecord: model.ShowRecord{},
			wantErr:          true,
			mockSetUp: func(mock sqlmock.Sqlmock, args args) {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta(
					"INSERT INTO `records` (`category_id`,`date`,`price`,`memo`) VALUES (?,?,?,?)")).
					WillReturnError(gorm.ErrInvalidData)
				mock.ExpectRollback()
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gormDB, mock, _ := NewDBMock()
			tt.fields.Conn = gormDB

			dbR := &DBRepository{
				Conn: tt.fields.Conn,
			}

			tt.mockSetUp(mock, tt.args)

			gotRetAddRecord, err := dbR.CreateRecordDB(tt.args.record)
			if (err != nil) != tt.wantErr {
				t.Errorf("dbRepository.CreateRecordDB() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotRetAddRecord, tt.wantRetAddRecord) {
				t.Errorf("dbRepository.CreateRecordDB() = %v, want %v", gotRetAddRecord, tt.wantRetAddRecord)
			}
		})
	}
}

func Test_dbRepository_GetYearSummaryDB(t *testing.T) {
	type fields struct {
		Conn *gorm.DB
	}
	type args struct {
		year int64
	}
	tests := []struct {
		name                  string
		fields                fields
		args                  args
		wantYearSummaryInters []model.YearSummaryInter
		wantErr               bool
		mockSetUp             func(mock sqlmock.Sqlmock, args args)
	}{
		{
			name:                  "GetYearSummary 1",
			fields:                fields{},
			args:                  args{year: int64(2022)},
			wantYearSummaryInters: []model.YearSummaryInter{{CategoryID: 100, Name: "catname", YearMonth: "202101", Total: 12345}},
			wantErr:               false,
			mockSetUp: func(mock sqlmock.Sqlmock, args args) {
				mock.ExpectQuery(regexp.QuoteMeta("SELECT records.category_id , categories.name, DATE_FORMAT(records.date, '%Y%m'), sum(price) FROM records INNER JOIN categories ON categories.category_id = records.category_id WHERE records.date")).
					WillReturnRows(sqlmock.NewRows([]string{"category_id", "name", "date", "sum(price)"}).AddRow(100, "catname", "202101", 12345))
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gormDB, mock, _ := NewDBMock()
			tt.fields.Conn = gormDB

			dbR := &DBRepository{
				Conn: tt.fields.Conn,
			}

			tt.mockSetUp(mock, tt.args)

			gotYearSummaryInters, err := dbR.GetYearSummaryDB(tt.args.year)
			if (err != nil) != tt.wantErr {
				t.Errorf("dbRepository.GetYearSummaryDB() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotYearSummaryInters, tt.wantYearSummaryInters) {
				t.Errorf("dbRepository.GetYearSummaryDB() = %v, want %v", gotYearSummaryInters, tt.wantYearSummaryInters)
			}
		})
	}
}

func Test_dbRepository_GetRecentRecord(t *testing.T) {
	type fields struct {
		Conn *gorm.DB
	}
	type args struct {
		n int
	}
	tests := []struct {
		name              string
		fields            fields
		args              args
		wantGetRecentData []model.ShowRecord
		wantErr           bool
		mockSetUp         func(mock sqlmock.Sqlmock, args args)
	}{
		{
			name:   "success",
			fields: fields{},
			args:   args{n: 2},
			wantGetRecentData: []model.ShowRecord{
				{
					Id:           1,
					CategoryID:   100,
					CategoryName: "cat1",
					Date:         time.Date(2000, 1, 2, 0, 0, 0, 0, time.FixedZone("Asia/Tokyo", 9*60*60)),
					Price:        1000,
					Memo:         "",
				},
				{
					Id:           2,
					CategoryID:   200,
					CategoryName: "cat2",
					Date:         time.Date(2000, 1, 1, 0, 0, 0, 0, time.FixedZone("Asia/Tokyo", 9*60*60)),
					Price:        2000,
					Memo:         "memo2",
				},
			},
			wantErr: false,
			mockSetUp: func(mock sqlmock.Sqlmock, args args) {
				mock.ExpectQuery(regexp.QuoteMeta("SELECT records.*, categories.category_id, categories.name FROM `records` INNER JOIN categories ON records.category_id = categories.category_id")).
					WillReturnRows(sqlmock.NewRows([]string{"Id", "category_id", "date", "price", "memo", "name", "type"}).
						AddRow(1, 100, time.Date(2000, 1, 2, 0, 0, 0, 0, time.FixedZone("Asia/Tokyo", 9*60*60)), 1000, "", "cat1", "type").
						AddRow(2, 200, time.Date(2000, 1, 1, 0, 0, 0, 0, time.FixedZone("Asia/Tokyo", 9*60*60)), 2000, "memo2", "cat2", "type"))
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gormDB, mock, _ := NewDBMock()
			tt.fields.Conn = gormDB

			dbR := &DBRepository{
				Conn: tt.fields.Conn,
			}

			tt.mockSetUp(mock, tt.args)

			gotGetRecentData, err := dbR.GetRecentRecord(tt.args.n)
			if (err != nil) != tt.wantErr {
				t.Errorf("dbRepository.GetRecentRecord() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotGetRecentData, tt.wantGetRecentData) {
				t.Errorf("dbRepository.GetRecentRecord() = %v, want %v", gotGetRecentData, tt.wantGetRecentData)
			}
		})
	}
}

func Test_dbRepository_DeleteRecordDB(t *testing.T) {
	type fields struct {
		Conn *gorm.DB
	}
	type args struct {
		id int64
	}
	tests := []struct {
		name      string
		fields    fields
		args      args
		wantErr   bool
		mockSetUp func(mock sqlmock.Sqlmock, args args)
	}{
		{
			name:    "sucess",
			fields:  fields{},
			args:    args{id: int64(1)},
			wantErr: false,
			mockSetUp: func(mock sqlmock.Sqlmock, args args) {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta(
					"DELETE FROM `records` WHERE id = ?")).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
		},
		{
			name:    "not found",
			fields:  fields{},
			args:    args{id: int64(2)},
			wantErr: true,
			mockSetUp: func(mock sqlmock.Sqlmock, args args) {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta(
					"DELETE FROM `records` WHERE id = ?")).
					WillReturnError(gorm.ErrRecordNotFound)
				mock.ExpectRollback()
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gormDB, mock, _ := NewDBMock()
			tt.fields.Conn = gormDB

			dbR := &DBRepository{
				Conn: tt.fields.Conn,
			}

			tt.mockSetUp(mock, tt.args)

			if err := dbR.DeleteRecordDB(tt.args.id); (err != nil) != tt.wantErr {
				t.Errorf("dbRepository.DeleteRecordDB() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
