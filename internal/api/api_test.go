package api

import (
	"database/sql"
	"mawinter-expense/internal/db"
	l "mawinter-expense/internal/logger"
	"os"
	"reflect"
	"testing"
	"time"

	httpdate "github.com/Songmu/go-httpdate"
	"github.com/jinzhu/gorm"
)

// -------------------------------------
func TestMain(m *testing.M) {
	l.NewSugarLogger()
	code := m.Run()
	os.Exit(code)
}

// 正常系
type mockDBRepositry1 struct{}

var collect1_GetYearSummaryStructs = []GetYearSummaryStruct{
	{
		CategoryId: 101,
		Name:       "testcatName1",
		Price:      []int64{12345, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 34567},
		Total:      int64(46912),
	},
	{
		CategoryId: 102,
		Name:       "testcatName2",
		Price:      []int64{0, 0, 0, 0, 23456, 0, 0, 0, 0, 0, 0, 0},
		Total:      int64(23456),
	},
}

var collect1_GetRecentRecord = []db.RecordsDetails{
	{
		Id:         1001,
		CategoryId: 101,
		Name:       "testcatName1",
		Date:       "2000-01-01",
		Price:      10001,
		Memo:       sql.NullString{Valid: false},
	},
	{
		Id:         1002,
		CategoryId: 102,
		Name:       "testcatName2",
		Date:       "2000-01-02",
		Price:      10002,
		Memo:       sql.NullString{String: "memo2", Valid: true},
	},
	{
		Id:         1005,
		CategoryId: 103,
		Name:       "testcatName3",
		Date:       "2000-01-05",
		Price:      10003,
		Memo:       sql.NullString{Valid: false},
	},
}

func (dbR *mockDBRepositry1) OpenTx() *gorm.DB {
	return nil
}

func (dbR *mockDBRepositry1) CloseTx(tx *gorm.DB, err error) error {
	return nil
}

func (dbR *mockDBRepositry1) GetYearSummaryDB(tx *gorm.DB, year int64) (yearSummary []db.GetYearSummaryDBStruct, err error) {
	mock := []db.GetYearSummaryDBStruct{
		{
			CategoryId: int64(101),
			Name:       "testcatName1",
			YearMonth:  "200004",
			Price:      int64(12345),
		},
		{
			CategoryId: int64(101),
			Name:       "testcatName1",
			YearMonth:  "200103",
			Price:      int64(34567),
		},
		{
			CategoryId: int64(102),
			Name:       "testcatName2",
			YearMonth:  "200008",
			Price:      23456,
		},
	}
	return mock, nil
}

func (dbR *mockDBRepositry1) AddRecordDB(tx *gorm.DB, record db.Records) (retRecords db.Records, err error) {
	return db.Records{Id: 1001, CategoryId: 101, Date: time.Now(), Price: 10000}, nil
}

func (dbR *mockDBRepositry1) DeleteRecordDB(tx *gorm.DB, id int64) (err error) {
	return nil
}

func (dbR *mockDBRepositry1) GetRecentRecord(tx *gorm.DB, n int64) (getRecentData []db.RecordsDetails, err error) {
	// n = 3 でテスト
	mock := []db.RecordsDetails{
		{
			Id:         1001,
			CategoryId: 101,
			Name:       "testcatName1",
			Date:       "2000-01-01",
			Price:      10001,
		},
		{
			Id:         1002,
			CategoryId: 102,
			Name:       "testcatName2",
			Date:       "2000-01-02",
			Price:      10002,
			Memo:       sql.NullString{String: "memo2", Valid: true},
		},
		{
			Id:         1005,
			CategoryId: 103,
			Name:       "testcatName3",
			Date:       "2000-01-05",
			Price:      10003,
		},
	}
	return mock, nil
}

//-------------------------------------

// func TestNewAPIService(t *testing.T) {
// 	type args struct {
// 		dbR_in db.DBRepository
// 	}
// 	tests := []struct {
// 		name string
// 		args args
// 		want APIService
// 	}{
// 		// TODO: Add test cases.
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			if got := NewAPIService(tt.args.dbR_in); !reflect.DeepEqual(got, tt.want) {
// 				t.Errorf("NewAPIService() = %v, want %v", got, tt.want)
// 			}
// 		})
// 	}
// }

func Test_apiService_GetYearSummary(t *testing.T) {
	type args struct {
		year int64
	}
	tests := []struct {
		name            string
		apis            *apiService
		args            args
		wantYearSummary []GetYearSummaryStruct
		wantErr         bool
	}{
		{
			name:            "test #1",
			apis:            &apiService{dbR: &mockDBRepositry1{}},
			args:            args{year: 2000},
			wantYearSummary: collect1_GetYearSummaryStructs,
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

func Test_getYearSummaryMakeIndex(t *testing.T) {
	type args struct {
		yearSummary []db.GetYearSummaryDBStruct
	}
	tests := []struct {
		name        string
		args        args
		wantIndexes []int
	}{
		{
			name: "test #1",
			args: args{
				yearSummary: []db.GetYearSummaryDBStruct{
					{
						CategoryId: int64(100),
					},
					{
						CategoryId: int64(100),
					},
					{
						CategoryId: int64(100),
					},
					{
						CategoryId: int64(200),
					},
					{
						CategoryId: int64(300),
					},
				},
			},
			wantIndexes: []int{0, 0, 0, 1, 2},
		},
		{
			name: "test #2",
			args: args{
				yearSummary: []db.GetYearSummaryDBStruct{},
			},
			wantIndexes: []int{},
		},
		{
			name: "test #3",
			args: args{
				yearSummary: []db.GetYearSummaryDBStruct{
					{
						CategoryId: int64(100),
					},
					{
						CategoryId: int64(100),
					},
					{
						CategoryId: int64(200),
					},
					{
						CategoryId: int64(300),
					},
					{
						CategoryId: int64(300),
					},
				},
			},
			wantIndexes: []int{0, 0, 1, 2, 2},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotIndexes := getYearSummaryMakeIndex(tt.args.yearSummary); !reflect.DeepEqual(gotIndexes, tt.wantIndexes) {
				t.Errorf("getYearSummaryMakeIndex() = %v, want %v", gotIndexes, tt.wantIndexes)
			}
		})
	}
}

func Test_posCompSlideint64(t *testing.T) {
	type args struct {
		arr []int64
	}
	tests := []struct {
		name        string
		args        args
		wantIndexes []int
	}{
		{
			name:        "test #1",
			args:        args{arr: []int64{100, 200, 300, 400, 500}},
			wantIndexes: []int{0, 1, 2, 3, 4},
		},
		{
			name:        "test #2",
			args:        args{arr: []int64{}},
			wantIndexes: []int{},
		},
		{
			name:        "test #3",
			args:        args{arr: []int64{100, 100, 100, 200, 300}},
			wantIndexes: []int{0, 0, 0, 1, 2},
		},
		{
			name:        "test #4",
			args:        args{arr: []int64{100, 200, 300, 300, 300}},
			wantIndexes: []int{0, 1, 2, 2, 2},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotIndexes := posCompSlideint64(tt.args.arr); !reflect.DeepEqual(gotIndexes, tt.wantIndexes) {
				t.Errorf("posCompSlideint64() = %v, want %v", gotIndexes, tt.wantIndexes)
			}
		})
	}
}

func Test_apiService_AddRecord(t *testing.T) {
	date1, _ := httpdate.Str2Time("2010-01-01", nil)

	type args struct {
		addRecord db.Records
	}
	tests := []struct {
		name             string
		apis             *apiService
		args             args
		wantRetAddRecord db.Records
		wantErr          bool
	}{
		{
			name:             "test #1",
			apis:             &apiService{dbR: &mockDBRepositry1{}},
			args:             args{addRecord: db.Records{CategoryId: 100, Price: 10000, Date: date1}},
			wantRetAddRecord: db.Records{},
			wantErr:          false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := tt.apis.AddRecord(tt.args.addRecord)
			if (err != nil) != tt.wantErr {
				t.Errorf("apiService.AddRecord() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			// if !reflect.DeepEqual(gotRetAddRecord, tt.wantRetAddRecord) {
			// 	t.Errorf("apiService.AddRecord() = %v, want %v", gotRetAddRecord, tt.wantRetAddRecord)
			// }
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
			name:    "test #1",
			apis:    &apiService{dbR: &mockDBRepositry1{}},
			args:    args{id: 100},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.apis.DeleteRecord(tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("apiService.DeleteRecord() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func Test_apiService_GetRecentRecord(t *testing.T) {
	type args struct {
		dataNum int64
	}
	tests := []struct {
		name              string
		apis              *apiService
		args              args
		wantGetRecentData []db.RecordsDetails
		wantErr           bool
	}{
		{
			name:              "test #1",
			apis:              &apiService{dbR: &mockDBRepositry1{}},
			args:              args{dataNum: 3},
			wantGetRecentData: collect1_GetRecentRecord,
			wantErr:           false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotGetRecentData, err := tt.apis.GetRecentRecord(tt.args.dataNum)
			if (err != nil) != tt.wantErr {
				t.Errorf("apiService.GetRecentRecord() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotGetRecentData, tt.wantGetRecentData) {
				t.Errorf("apiService.GetRecentRecord() = %v, want %v", gotGetRecentData, tt.wantGetRecentData)
			}
		})
	}
}
