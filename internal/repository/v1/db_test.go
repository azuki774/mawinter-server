package repository

import (
	"mawinter-server/internal/model"
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

// func testDBConnect() *DBRepository {
// 	const DBConnectRetry = 5
// 	const DBConnectRetryInterval = 10
// 	l, _ := zap.NewProduction()
// 	host := "127.0.0.1"
// 	port := "3306"
// 	user := "root"
// 	pass := "password"
// 	name := "mawinter"

// 	addr := net.JoinHostPort(host, port)
// 	dsn := user + ":" + pass + "@(" + addr + ")/" + name + "?parseTime=true&loc=Local"
// 	var gormdb *gorm.DB
// 	var err error
// 	for i := 0; i < DBConnectRetry; i++ {
// 		gormdb, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
// 		if err == nil {
// 			// Success DB connect
// 			l.Info("DB connect")
// 			break
// 		}
// 		l.Warn("DB connection retry")

// 		if i == DBConnectRetry {
// 			l.Error("failed to connect DB", zap.Error(err))
// 			return nil
// 		}

// 		time.Sleep(DBConnectRetryInterval * time.Second)
// 	}

// 	return &DBRepository{Conn: gormdb}
// }

func TestDBRepository_InsertRecord(t *testing.T) {
	type fields struct {
		Conn *gorm.DB
	}
	type args struct {
		req model.Recordstruct
	}
	tests := []struct {
		name      string
		fields    fields
		args      args
		wantRes   model.Recordstruct
		wantErr   bool
		mockSetUp func(mock sqlmock.Sqlmock)
	}{
		{
			name:   "OK",
			fields: fields{},
			args: args{
				req: model.Recordstruct{
					CategoryID: 100,
					Datetime:   time.Date(2000, 10, 1, 0, 0, 0, 0, jst),
					From:       "from",
					Type:       "type",
					Price:      123,
					Memo:       "memo",
				},
			},
			wantRes: model.Recordstruct{
				ID:         1001,
				CategoryID: 100,
				Datetime:   time.Date(2000, 10, 1, 0, 0, 0, 0, jst),
				From:       "from",
				Type:       "type",
				Price:      123,
				Memo:       "memo",
				// CreatedAt
				// UpdatedAt
			},
			wantErr: false,
			mockSetUp: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta(
					"INSERT INTO `Record_200010`")).
					WillReturnResult(sqlmock.NewResult(1001, 1))
				mock.ExpectCommit()
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

			gotRes, err := d.InsertRecord(tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("DBRepository.InsertRecord() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			// CreatedAt
			// UpdatedAt
			tt.wantRes.CreatedAt = gotRes.CreatedAt
			tt.wantRes.UpdatedAt = gotRes.UpdatedAt

			if !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("DBRepository.InsertRecord() = %v, want %v", gotRes, tt.wantRes)
			}
		})
	}
}

func TestDBRepository_GetCategoryInfo(t *testing.T) {
	type fields struct {
		Conn *gorm.DB
	}
	tests := []struct {
		name      string
		fields    fields
		wantInfo  []model.Category
		wantErr   bool
		mockSetUp func(mock sqlmock.Sqlmock)
	}{
		{
			name:   "get info",
			fields: fields{},
			wantInfo: []model.Category{
				{
					ID:         1,
					CategoryID: 100,
					Name:       "カテゴリ1",
				},
				{
					ID:         2,
					CategoryID: 200,
					Name:       "カテゴリ2",
				},
				{
					ID:         3,
					CategoryID: 300,
					Name:       "カテゴリ3",
				},
			},
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
			name:     "error",
			fields:   fields{},
			wantInfo: []model.Category{},
			wantErr:  true,
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

			gotInfo, err := d.GetCategoryInfo()
			if (err != nil) != tt.wantErr {
				t.Errorf("DBRepository.GetCategoryInfo() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotInfo, tt.wantInfo) {
				t.Errorf("DBRepository.GetCategoryInfo() = %v, want %v", gotInfo, tt.wantInfo)
			}
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestDBRepository_SumPriceForEachCatID(t *testing.T) {
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
		wantSum   []model.SumPriceCategoryID
		wantErr   bool
		mockSetUp func(mock sqlmock.Sqlmock)
	}{
		{
			name:   "OK",
			fields: fields{},
			args:   args{yyyymm: "202210"},
			wantSum: []model.SumPriceCategoryID{
				{
					CategoryID: 100,
					Count:      1,
					Sum:        10000,
				},
				{
					CategoryID: 200,
					Count:      2,
					Sum:        20000,
				},
				{
					CategoryID: 300,
					Count:      3,
					Sum:        30000,
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

			gotSum, err := d.SumPriceForEachCatID(tt.args.yyyymm)
			if (err != nil) != tt.wantErr {
				t.Errorf("DBRepository.SumPriceForEachCatID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotSum, tt.wantSum) {
				t.Errorf("DBRepository.SumPriceForEachCatID() = %v, want %v", gotSum, tt.wantSum)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
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
		wantFlag  bool
		wantErr   bool
		mockSetUp func(mock sqlmock.Sqlmock)
	}{
		{
			name:     "not processed (not found)",
			fields:   fields{},
			args:     args{yyyymm: "202210"},
			wantFlag: false,
			wantErr:  false,
			mockSetUp: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT`).
					WillReturnError(gorm.ErrRecordNotFound)
			},
		},
		{
			name:     "done",
			fields:   fields{},
			args:     args{yyyymm: "202210"},
			wantFlag: true,
			wantErr:  false,
			mockSetUp: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT`).
					WillReturnRows(sqlmock.NewRows([]string{"yyyymm", "done"}).
						AddRow("202210", 1))
			},
		},
		{
			name:     "unknown error",
			fields:   fields{},
			args:     args{yyyymm: "202210"},
			wantFlag: false,
			wantErr:  true,
			mockSetUp: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT`).
					WillReturnError(gorm.ErrInvalidData)
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

			gotFlag, err := d.GetMonthlyFixDone(tt.args.yyyymm)
			if (err != nil) != tt.wantErr {
				t.Errorf("DBRepository.GetMonthlyFixDone() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotFlag != tt.wantFlag {
				t.Errorf("DBRepository.GetMonthlyFixDone() = %v, want %v", gotFlag, tt.wantFlag)
			}
		})
	}
}

func TestDBRepository_InsertMonthlyFixBilling(t *testing.T) {
	type fields struct {
		Conn *gorm.DB
	}
	type args struct {
		yyyymm   string
		fixBills []model.MonthlyFixBilling
	}
	tests := []struct {
		name      string
		fields    fields
		args      args
		wantErr   bool
		mockSetUp func(mock sqlmock.Sqlmock)
	}{
		{
			name:   "ok",
			fields: fields{},
			args: args{
				yyyymm: "202210",
				fixBills: []model.MonthlyFixBilling{
					{
						CategoryID: 100,
						Day:        2,
						Price:      1234,
						Type:       "type1",
						Memo:       "memo1",
					},
					{
						CategoryID: 101,
						Day:        4,
						Price:      2345,
						Type:       "type2",
						Memo:       "memo2",
					},
				},
			},
			wantErr: false,
			mockSetUp: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta(
					"INSERT INTO `Monthly_Fix_Done`")).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectExec(regexp.QuoteMeta(
					"INSERT INTO")).
					WillReturnResult(sqlmock.NewResult(1, 2))
				mock.ExpectCommit()
			},
		},
		{
			name:   "done insert error",
			fields: fields{},
			args: args{
				yyyymm: "202210",
				fixBills: []model.MonthlyFixBilling{
					{
						CategoryID: 100,
						Day:        2,
						Price:      1212,
						Type:       "type1",
						Memo:       "memo1",
					},
					{
						CategoryID: 101,
						Day:        4,
						Price:      2345,
						Type:       "type2",
						Memo:       "memo2",
					},
				},
			},
			wantErr: true,
			mockSetUp: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta(
					"INSERT INTO `Monthly_Fix_Done`")).
					WillReturnError(gorm.ErrRegistered)
				mock.ExpectRollback()
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
			if err := d.InsertMonthlyFixBilling(tt.args.yyyymm, tt.args.fixBills); (err != nil) != tt.wantErr {
				t.Errorf("DBRepository.InsertMonthlyFixBilling() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDBRepository_GetMonthlyFixBilling(t *testing.T) {
	type fields struct {
		Conn *gorm.DB
	}
	tests := []struct {
		name         string
		fields       fields
		wantFixBills []model.MonthlyFixBilling
		wantErr      bool
		mockSetUp    func(mock sqlmock.Sqlmock)
	}{
		{
			name:   "ok",
			fields: fields{},
			wantFixBills: []model.MonthlyFixBilling{
				{
					CategoryID: 100,
					Day:        2,
					Price:      1212,
					Type:       "type1",
					Memo:       "memo1",
				},
				{
					CategoryID: 101,
					Day:        4,
					Price:      2345,
					Type:       "type2",
					Memo:       "memo2",
				},
			},
			wantErr: false,
			mockSetUp: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT`).
					WillReturnRows(sqlmock.NewRows(
						[]string{
							"id",
							"category_id",
							"day",
							"price",
							"type",
							"memo",
							"created_at",
							"updated_at",
						}).
						AddRow(1, 100, 2, 1212, "type1", "memo1", time.Now(), time.Now()).
						AddRow(2, 101, 4, 2345, "type2", "memo2", time.Now(), time.Now()))
			},
		},
		{
			name:         "error",
			fields:       fields{},
			wantFixBills: []model.MonthlyFixBilling{},
			wantErr:      true,
			mockSetUp: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT`).
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
			gotFixBills, err := d.GetMonthlyFixBilling()
			if (err != nil) != tt.wantErr {
				t.Errorf("DBRepository.GetMonthlyFixBilling() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotFixBills, tt.wantFixBills) {
				t.Errorf("DBRepository.GetMonthlyFixBilling() = %v, want %v", gotFixBills, tt.wantFixBills)
			}
		})
	}
}
