package api

import (
	"context"
	"errors"
	"mawinter-server/internal/model"
	"mawinter-server/internal/openapi"
	"sort"
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
	GetMonthRecords(yyyymm string, params openapi.GetV2RecordYyyymmParams) (recs []openapi.Record, err error)
	GetMonthRecordsRecent(yyyymm string, num int) (recs []openapi.Record, err error)
	MakeCategoryNameMap() (cnf map[int]string, err error)
	GetMonthMidSummary(yyyymm string) (summon []model.CategoryMidMonthSummary, err error) // SELECT category_id, count(*), sum(price) FROM Record_202211 GROUP BY category_id;
	InsertMonthlyFixBilling(yyyymm string) (recs []openapi.Record, err error)
	GetMonthlyFixDone(yyyymm string) (done bool, err error)
	GetMonthlyConfirm(yyyymm string) (yc openapi.ConfirmInfo, err error)
	UpdateMonthlyConfirm(yyyymm string, confirm bool) (yc openapi.ConfirmInfo, err error)
}

type APIService struct {
	Logger *zap.Logger
	Repo   DBRepository
}

func int2ptr(i int) *int {
	return &i
}

func (a *APIService) PostRecord(ctx context.Context, req openapi.ReqRecord) (rec openapi.Record, err error) {
	a.Logger.Info("called post record")
	a.Logger.Info("get monthly confirm")
	// 確定した月でないかを確認する
	yyyymm := (*req.Datetime)[0:5]
	yc, err := a.Repo.GetMonthlyConfirm(yyyymm)
	if err != nil {
		return openapi.Record{}, err
	}
	if *yc.Status {
		a.Logger.Info("already confirm month", zap.String("yyyymm", yyyymm))
		return openapi.Record{}, model.ErrAlreadyRecorded
	}

	a.Logger.Info("get category name mapping")
	// categoryNameMap 取得
	cnf, err := a.Repo.MakeCategoryNameMap()
	if err != nil {
		return openapi.Record{}, err
	}

	_, ok := cnf[req.CategoryId] // DBにCategory IDがあるか確認
	if !ok {
		// Category ID がDBに未登録の場合
		a.Logger.Warn("unknown category ID", zap.Int("category_id", rec.CategoryId))
		return openapi.Record{}, model.ErrUnknownCategoryID
	}

	rec, err = a.Repo.InsertRecord(req)
	if err != nil {
		a.Logger.Error("failed to insert", zap.String("msg", err.Error()), zap.Error(err))
		return openapi.Record{}, err
	}
	rec.CategoryName = cnf[rec.CategoryId]

	a.Logger.Info("complete post record")
	return rec, nil
}

func (a *APIService) PostMonthlyFixRecord(ctx context.Context, yyyymm string) (recs []openapi.Record, err error) {
	a.Logger.Info("called post fixmonth records", zap.String("yyyymm", yyyymm))
	done, err := a.Repo.GetMonthlyFixDone(yyyymm)
	if err != nil {
		a.Logger.Error("failed to get monthly processed data", zap.String("yyyymm", yyyymm), zap.Error(err))
		return []openapi.Record{}, err
	}
	if done {
		// 既に処理済の場合はスキップ
		a.Logger.Info("called post monthly already registed", zap.String("yyyymm", yyyymm))
		return []openapi.Record{}, model.ErrAlreadyRecorded
	}

	recs, err = a.Repo.InsertMonthlyFixBilling(yyyymm)
	if err != nil {
		a.Logger.Error("failed to insert data", zap.String("yyyymm", yyyymm), zap.Error(err))
		return []openapi.Record{}, err
	}

	a.Logger.Info("complete post fixmonth record", zap.String("yyyymm", yyyymm))
	return recs, nil
}

func (a *APIService) CreateTableYear(ctx context.Context, year int) (err error) {
	a.Logger.Info("called create year table")
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
func (a *APIService) GetYYYYMMRecords(ctx context.Context, yyyymm string, params openapi.GetV2RecordYyyymmParams) (recs []openapi.Record, err error) {
	recs = []openapi.Record{}
	a.Logger.Info("called get month records")

	a.Logger.Info("get records from DB")
	recsRaw, err := a.Repo.GetMonthRecords(yyyymm, params) // category_name なし
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

func (a *APIService) GetYYYYMMRecordsRecent(ctx context.Context, yyyymm string, num int) (recs []openapi.Record, err error) {
	recs = []openapi.Record{}
	a.Logger.Info("called get month recent records")

	a.Logger.Info("get records from DB")
	recsRaw, err := a.Repo.GetMonthRecordsRecent(yyyymm, num)
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

	a.Logger.Info("complete get month recent records")
	return recs, nil
}

func (a *APIService) GetV2YearSummary(ctx context.Context, year int) (sums []openapi.CategoryYearSummary, err error) {
	a.Logger.Info("called get year summary")

	a.Logger.Info("get category name mapping")
	// categoryNameMap 取得
	cnf, err := a.Repo.MakeCategoryNameMap()
	if err != nil {
		return nil, err
	}

	sumsDec := make(map[int]*openapi.CategoryYearSummary) // CatID -> openapi.CategoryYearSummary
	// 初期化
	for catId := range cnf {
		sumsDec[catId] = &openapi.CategoryYearSummary{
			CategoryId:   catId,
			CategoryName: cnf[catId],
			Count:        0,
			Price:        make([]int, 12),
			Total:        0,
		}
	}

	a.Logger.Info("get records from DB")
	// 1月ずつ処理する
	yyyymmList := fyInterval(year)
	for mi, yyyymm := range yyyymmList {
		monthSums, err := a.Repo.GetMonthMidSummary(yyyymm)
		if err != nil {
			a.Logger.Error("failed to get info from DB", zap.Error(err))
			return nil, err
		}

		for _, monthSum := range monthSums {
			catId := monthSum.CategoryId
			count := monthSum.Count
			price := monthSum.Price
			sumsDec[catId].Count += count
			sumsDec[catId].Price[mi] = price
			sumsDec[catId].Total += price
		}
	}

	a.Logger.Info("making month summary")
	// 最終的に出力する構造体に挿入する
	for _, v := range sumsDec {
		newSum := openapi.CategoryYearSummary{
			CategoryId:   v.CategoryId,
			CategoryName: v.CategoryName,
			Count:        v.Count,
			Price:        v.Price,
			Total:        v.Total,
		}
		sums = append(sums, newSum)
	}

	sort.Slice(sums, func(i, j int) bool {
		return sums[i].CategoryId < sums[j].CategoryId
	})

	a.Logger.Info("complete get year summary")

	return sums, nil
}

func (a *APIService) GetMonthlyConfirm(ctx context.Context, yyyymm string) (yc openapi.ConfirmInfo, err error) {
	a.Logger.Info("called GetMonthlyConfirm")
	yc.Yyyymm = &yyyymm
	yc, err = a.Repo.GetMonthlyConfirm(yyyymm)
	if err != nil {
		// Internal error -> error
		a.Logger.Error("failed to get monthly confirm", zap.Error(err))
		return openapi.ConfirmInfo{}, err
	} else {
		// success fetch data or not found (= false)
		a.Logger.Info("fetch monthly confirm successfully", zap.Error(err))
	}

	a.Logger.Info("complete GetMonthlyConfirm")
	return yc, nil
}

func (a *APIService) UpdateMonthlyConfirm(ctx context.Context, yyyymm string, confirm bool) (yc openapi.ConfirmInfo, err error) {
	a.Logger.Info("called UpdateMonthlyConfirm")
	yc.Yyyymm = &yyyymm
	yc, err = a.Repo.UpdateMonthlyConfirm(yyyymm, confirm)
	if err != nil {
		// Internal error -> error
		a.Logger.Error("failed to get monthly confirm", zap.Error(err))
		return openapi.ConfirmInfo{}, err
	}

	a.Logger.Info("update monthly confirm successfully", zap.Error(err))

	a.Logger.Info("complete UpdateMonthlyConfirm")
	return yc, nil
}
