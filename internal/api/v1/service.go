package v1

import (
	"context"
	"mawinter-server/internal/model"
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
	InsertRecord(req model.Recordstruct) (res model.Recordstruct, err error)
	GetCategoryInfo() (info []model.Category, err error)
	SumPriceForEachCatID(yyyymm string) (sum []model.SumPriceCategoryID, err error) // SELECT category_id, count(*), sum(price) FROM Record_202211 GROUP BY category_id;
}

type APIService struct {
	Logger *zap.Logger
	Repo   DBRepository
}

func (a *APIService) AddRecord(ctx context.Context, req model.RecordRequest) (res model.Recordstruct, err error) {
	err = model.ValidRecordRequest(req)
	if err != nil {
		a.Logger.Warn("invalid value detected", zap.Error(err))
		return model.Recordstruct{}, err
	}

	// デフォルト値挿入
	record, err := model.NewRecordFromReq(req)
	if err != nil {
		a.Logger.Warn("failed to create request", zap.Error(err))
		return model.Recordstruct{}, err
	}

	res, err = a.Repo.InsertRecord(record)
	if err != nil {
		a.Logger.Error("failed to insert record", zap.Error(err))
		return model.Recordstruct{}, err
	}

	return res, nil
}

func (a *APIService) GetYearSummary(ctx context.Context, yyyy string) (sum []*model.CategoryYearSummary, err error) {
	yyyyint, err := model.ValidYYYY(yyyy)
	if err != nil {
		a.Logger.Warn("invalid value detected", zap.Error(err))
		return nil, err
	}

	cats, err := a.Repo.GetCategoryInfo()
	if err != nil {
		return nil, err
	}

	sum = model.NewCategoryYearSummary(cats)

	yyyymmList := fyInterval(yyyyint)
	for _, yyyymm := range yyyymmList {
		monthSums, err := a.Repo.SumPriceForEachCatID(yyyymm)
		if err != nil {
			a.Logger.Error("failed to get info from DB", zap.Error(err))
			return nil, err
		}

		f := loadSumPriceForEachCatID(monthSums)

		for _, s := range sum {
			s.AddMonthPrice(f[s.CategoryID]) // Refresh MonthPrice and Total
		}
	}

	return sum, nil
}

func loadSumPriceForEachCatID(sum []model.SumPriceCategoryID) (f map[int]int) {
	f = make(map[int]int)
	for _, v := range sum {
		f[int(v.CategoryID)] = int(v.Sum)
	}
	return
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
