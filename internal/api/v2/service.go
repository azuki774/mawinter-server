package api

import (
	"context"
	"errors"
	"mawinter-server/internal/model"
	"mawinter-server/internal/openapi"
	"time"

	"github.com/go-sql-driver/mysql"
	"go.uber.org/zap"
)

var jst *time.Location

func init() {
	j, err := time.LoadLocation("Asia/Tokyo")
	if err != nil {
		panic(err)
	}
	jst = j
}

type DBRepository interface {
	CreateTableYYYYMM(yyyymm string) (err error)
	InsertRecord(req openapi.ReqRecord) (rec openapi.Record, err error)
	GetMonthRecords(yyyymm string) (recs []openapi.Record, err error)
	MakeCategoryNameMap() (cnf map[int]string, err error)
}

type APIService struct {
	Logger *zap.Logger
	Repo   DBRepository
}

func (a *APIService) PostRecord(ctx context.Context, req openapi.ReqRecord) (rec openapi.Record, err error) {
	a.Logger.Info("called post record")
	rec, err = a.Repo.InsertRecord(req)
	if err != nil {
		a.Logger.Error("failed to insert", zap.String("msg", err.Error()), zap.Error(err))
		return openapi.Record{}, err
	}

	a.Logger.Info("get category name mapping")
	// categoryNameMap 取得
	cnf, err := a.Repo.MakeCategoryNameMap()
	if err != nil {
		return openapi.Record{}, err
	}
	rec.CategoryName = cnf[rec.CategoryId]

	a.Logger.Info("complete post record")
	return rec, nil
}

func (a *APIService) CreateTableYear(ctx context.Context, year int) (err error) {
	a.Logger.Info("called create year table")
	// yyyyint, err := model.ValidYYYY(year)
	// if err != nil {
	// 	a.Logger.Warn("invalid value detected", zap.Error(err))
	// 	return err
	// }
	yyyymmList := fyInterval(year)
	for _, yyyymm := range yyyymmList {
		err = a.Repo.CreateTableYYYYMM(yyyymm)
		var mysqlError *mysql.MySQLError
		if err != nil && errors.As(err, &mysqlError) {
			if err.(*mysql.MySQLError).Number == 1050 {
				// already exitsts
				a.Logger.Info("already existed Record_YYYYMM table", zap.String("YYYYMM", yyyymm))
				return model.ErrAlreadyRecorded
			}
			// other MySQL error
			a.Logger.Error("failed to create Record_YYYYMM table (MySQL)", zap.String("msg", err.Error()), zap.String("YYYYMM", yyyymm), zap.Error(err))
			return err
		} else if err != nil {
			// internal error
			a.Logger.Error("failed to create Record_YYYYMM table (gorm)", zap.String("msg", err.Error()), zap.String("YYYYMM", yyyymm), zap.Error(err))
			return err
		}
		a.Logger.Info("complete create Record_YYYYMM table", zap.String("YYYYMM", yyyymm))
	}

	return nil
}

// FYyyyy の yyyymm をリストで返す
func fyInterval(yyyy int) (yyyymm []string) {
	t := time.Date(yyyy, 4, 1, 0, 0, 0, 0, jst)
	for i := 0; i < 12; i++ {
		yyyymm = append(yyyymm, t.Format("200601"))
		t = t.AddDate(0, 1, 0)
	}
	return yyyymm
}

// GetYYYYMMRecords は yyyymm 月のレコードを取得する
func (a *APIService) GetYYYYMMRecords(ctx context.Context, yyyymm string) (recs []openapi.Record, err error) {
	a.Logger.Info("called get month records")

	a.Logger.Info("get records from DB")
	recsRaw, err := a.Repo.GetMonthRecords(yyyymm) // category_name なし
	if err != nil {
		return nil, err
	}

	a.Logger.Info("get category name mapping")
	// categoryNameMap 取得
	cnf, err := a.Repo.MakeCategoryNameMap()
	if err != nil {
		return nil, err
	}

	for _, rec := range recsRaw {
		// categoryName を付与
		rec.CategoryName = cnf[rec.CategoryId]
		recs = append(recs, rec)
	}

	a.Logger.Info("complete get month records")
	return recs, nil
}
