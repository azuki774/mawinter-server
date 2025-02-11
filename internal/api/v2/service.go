package api

import (
	"context"
	"mawinter-server/internal/model"
	"mawinter-server/internal/openapi"
	"mawinter-server/internal/timeutil"
	"sort"
	"time"

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
	InsertRecord(req openapi.ReqRecord) (rec openapi.Record, err error)
	GetRecords(ctx context.Context, GetRecordOpt model.GetRecordOption) (recs []openapi.Record, err error)
	GetRecordByID(ctx context.Context, id int) (rec openapi.Record, err error)
	DeleteRecordByID(ctx context.Context, id int) (err error)
	GetRecordsCount(ctx context.Context) (num int, err error)
	GetCategories(ctx context.Context) (cats []model.Category, err error)
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

	// FIX: req.Datetime の変換タイミングが悪いので暫定対応
	var yyyymm string
	if req.Datetime == nil {
		// Datetime が未設定なら現時刻がDBに挿入されるはずなので、今の時点でのYYYYMMをセットする
		yyyymm = timeutil.NowFunc().Format("200601")
	} else {
		yyyymm = (*req.Datetime)[0:6]
	}

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

// GetRecords は num の数だけ ID 降順に Record を取得する
func (a *APIService) GetRecords(ctx context.Context, GetRecordOpt model.GetRecordOption) (recs []openapi.Record, err error) {
	a.Logger.Info("called GetRecordsRecent", zap.Int("num", GetRecordOpt.Num))
	recsRaw, err := a.Repo.GetRecords(ctx, GetRecordOpt)
	if err != nil {
		a.Logger.Error("failed to get records")
		return []openapi.Record{}, err
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

	a.Logger.Info("complete GetRecordsRecent", zap.Int("num", GetRecordOpt.Num))
	return recs, nil
}

func (a *APIService) GetRecordByID(ctx context.Context, id int) (rec openapi.Record, err error) {
	a.Logger.Info("called GetRecordByID")
	rec, err = a.Repo.GetRecordByID(ctx, id)
	if err != nil {
		a.Logger.Error("failed to get record by ID", zap.Int("id", id), zap.Error(err))
		return openapi.Record{}, err
	}

	return rec, nil
}

func (a *APIService) DeleteRecordByID(ctx context.Context, id int) (err error) {
	a.Logger.Info("called DeleteRecordByID")
	err = a.Repo.DeleteRecordByID(ctx, id)
	if err != nil {
		a.Logger.Error("failed to delete record by ID", zap.Int("id", id), zap.Error(err))
		return err
	}

	return nil
}

// GetRecordsCount は レコード件数の総数を返す
func (a *APIService) GetRecordsCount(ctx context.Context) (rec openapi.RecordCount, err error) {
	a.Logger.Info("called GetRecordsCount")
	num, err := a.Repo.GetRecordsCount(ctx)
	if err != nil {
		a.Logger.Error("failed to get the number of records", zap.Error(err))
		return openapi.RecordCount{}, err
	}

	rec.Num = int2ptr(num)
	return rec, nil
}

// GetCategories は管理しているカテゴリ情報を返却する
func (a *APIService) GetCategories(ctx context.Context) (cats []openapi.Category, err error) {
	a.Logger.Info("called get categories")
	cs, err := a.Repo.GetCategories(ctx)
	if err != nil {
		a.Logger.Error("failed to get categories from DB", zap.Error(err))
		return []openapi.Category{}, err
	}

	for _, c := range cs {
		var tmpC openapi.Category
		tmpC.CategoryId = int(c.CategoryID)
		tmpC.CategoryName = c.Name
		cats = append(cats, tmpC)
	}

	return cats, nil
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

// FYyyyy の yyyymm をリストで返す
func fyInterval(yyyy int) (yyyymm []string) {
	t := time.Date(yyyy, 4, 1, 0, 0, 0, 0, jst)
	for i := 0; i < 12; i++ {
		yyyymm = append(yyyymm, t.Format("200601"))
		t = t.AddDate(0, 1, 0)
	}
	return yyyymm
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
