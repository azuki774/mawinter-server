package api

import (
	"fmt"
	"mawinter-server/internal/model"
	"os"
	"reflect"
	"testing"
	"time"

	"go.uber.org/zap"
)

var testLogger *zap.Logger

func newLogger() (*zap.Logger, error) {
	config := zap.NewProductionConfig()
	// config.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
	l, err := config.Build()

	l.WithOptions(zap.AddStacktrace(zap.ErrorLevel))
	if err != nil {
		fmt.Printf("failed to create logger: %v\n", err)
	}
	return l, err
}

func init() {
	l, err := newLogger()
	defer l.Sync()
	if err != nil {
		fmt.Printf("logger failed")
		os.Exit(1)
	}
	testLogger = l
}

func Test_apiService_CreateRecord(t *testing.T) {
	type fields struct {
		Logger *zap.Logger
		DBRepo DBRepository
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
			fields:           fields{DBRepo: &mockDBRepositry1{}, Logger: testLogger},
			args:             args{model.CreateRecord{CategoryID: 101, Date: "20000123", Price: 10000}},
			wantRetAddRecord: model.ShowRecord{Id: 1001, CategoryID: 101, CategoryName: "cat1", Date: time.Date(2000, 1, 23, 0, 0, 0, 0, time.FixedZone("Asia/Tokyo", 9*60*60)), Price: 10000},
			wantErr:          false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ap := &APIService{
				DBRepo: tt.fields.DBRepo,
				Logger: testLogger,
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
		apis            *APIService
		args            args
		wantYearSummary []model.YearSummary
		wantErr         bool
	}{
		{
			name:            "test #1",
			apis:            &APIService{DBRepo: &mockDBRepositry1{}, Logger: testLogger},
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
		apis    *APIService
		args    args
		wantErr bool
	}{
		{
			name:    "success",
			apis:    &APIService{DBRepo: &mockDBRepositry1{}, Logger: testLogger},
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
