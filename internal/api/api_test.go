package api

import (
	"fmt"
	"mawinter-server/internal/logger"
	"mawinter-server/internal/model"
	"mawinter-server/internal/repository"
	"os"
	"reflect"
	"testing"
	"time"

	"go.uber.org/zap"
)

var testLogger *zap.SugaredLogger

func init() {
	l, err := logger.NewSugarLogger()
	defer l.Sync()
	if err != nil {
		fmt.Printf("logger failed")
		os.Exit(1)
	}
	testLogger = l
}

// 正常系
type mockDBRepositry1 struct{}

func (dbR *mockDBRepositry1) CreateRecordDB(record model.Records) (retAddRecord model.ShowRecord, err error) {
	return model.ShowRecord{Id: 1001, CategoryID: 101, CategoryName: "cat1", Date: time.Date(2000, 1, 23, 0, 0, 0, 0, time.FixedZone("Asia/Tokyo", 9*60*60)), Price: 10000}, nil
}

func (dbR *mockDBRepositry1) GetYearSummaryDB(year int64) (yearSummaryInters []model.YearSummaryInter, err error) {
	mock := []model.YearSummaryInter{
		{
			CategoryID: int64(101),
			Name:       "testcatName1",
			YearMonth:  "200004",
			Total:      int64(12345),
		},
		{
			CategoryID: int64(101),
			Name:       "testcatName1",
			YearMonth:  "200103",
			Total:      int64(34567),
		},
		{
			CategoryID: int64(102),
			Name:       "testcatName2",
			YearMonth:  "200008",
			Total:      23456,
		},
	}
	return mock, nil
}

func (dbR *mockDBRepositry1) GetRecentRecord(n int) (getRecentData []model.ShowRecord, err error) {
	getRecentData = []model.ShowRecord{
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
	}

	return getRecentData, nil
}

func (dbR *mockDBRepositry1) DeleteRecordDB(id int64) (err error) {
	return nil
}

func Test_apiService_CreateRecord(t *testing.T) {
	type fields struct {
		dbR repository.DBRepository
	}
	type args struct {
		addRecord model.CreateRecord
	}
	tests := []struct {
		name             string
		fields           fields
		args             args
		wantRetAddRecord model.ShowRecord
		wantErr          bool
	}{
		{
			name:             "success",
			fields:           fields{&mockDBRepositry1{}},
			args:             args{model.CreateRecord{CategoryID: 101, Date: "20000123", Price: 10000}},
			wantRetAddRecord: model.ShowRecord{Id: 1001, CategoryID: 101, CategoryName: "cat1", Date: time.Date(2000, 1, 23, 0, 0, 0, 0, time.FixedZone("Asia/Tokyo", 9*60*60)), Price: 10000},
			wantErr:          false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ap := &apiService{
				dbR:    tt.fields.dbR,
				logger: testLogger,
			}
			gotRetAddRecord, err := ap.CreateRecord(tt.args.addRecord)
			if (err != nil) != tt.wantErr {
				t.Errorf("apiService.CreateRecord() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotRetAddRecord, tt.wantRetAddRecord) {
				t.Errorf("apiService.CreateRecord() = %v, want %v", gotRetAddRecord, tt.wantRetAddRecord)
			}
		})
	}
}

func Test_apiService_GetYearSummary(t *testing.T) {
	var collect1_GetYearSummary = []model.YearSummary{
		{
			CategoryID: 101,
			Name:       "testcatName1",
			Price:      []int64{12345, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 34567},
			Total:      int64(46912),
		},
		{
			CategoryID: 102,
			Name:       "testcatName2",
			Price:      []int64{0, 0, 0, 0, 23456, 0, 0, 0, 0, 0, 0, 0},
			Total:      int64(23456),
		},
	}

	type args struct {
		year int64
	}
	tests := []struct {
		name            string
		apis            *apiService
		args            args
		wantYearSummary []model.YearSummary
		wantErr         bool
	}{
		{
			name:            "test #1",
			apis:            &apiService{dbR: &mockDBRepositry1{}, logger: testLogger},
			args:            args{year: 2000},
			wantYearSummary: collect1_GetYearSummary,
			wantErr:         false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotYearSummary, err := tt.apis.GetYearSummary(tt.args.year)
			if (err != nil) != tt.wantErr {
				t.Errorf("apiService.GetYearSummary() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotYearSummary, tt.wantYearSummary) {
				t.Errorf("apiService.GetYearSummary() = %v, want %v", gotYearSummary, tt.wantYearSummary)
			}
		})
	}
}

func Test_apiService_DeleteRecord(t *testing.T) {
	type args struct {
		id int64
	}
	tests := []struct {
		name    string
		apis    *apiService
		args    args
		wantErr bool
	}{
		{
			name:    "success",
			apis:    &apiService{dbR: &mockDBRepositry1{}, logger: testLogger},
			args:    args{id: 1},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if err := tt.apis.DeleteRecord(tt.args.id); (err != nil) != tt.wantErr {
				t.Errorf("apiService.DeleteRecord() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
