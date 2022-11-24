package v1

import (
	"errors"
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
	InsertRecord(req *model.RecordRequest) (res *model.Record_YYYYMM, err error)
	GetCategoryInfo(categoryID int) (info *model.Category, err error)
	GetCategoryMonthSummary(categoryID int, yyyymm string) (sum *model.CategoryMonthSummary, err error)
}

type APIService struct {
	Logger *zap.Logger
	Repo   DBRepository
}

func (a *APIService) AddRecord(req *model.RecordRequest) (res *model.Record_YYYYMM, err error) {
	err = model.ValidRecordRequest(req)
	if err != nil {
		a.Logger.Warn("invalid value detected", zap.Error(err))
		return nil, err
	}

	res, err = a.Repo.InsertRecord(req)
	if err != nil {
		a.Logger.Error("failed to insert record", zap.Error(err))
		return nil, err
	}

	return res, nil
}

func (a *APIService) GetYearCategorySummary(categoryID int, yyyy string) (sum *model.CategoryYearSummary, err error) {
	yyyyint, err := model.ValidYYYY(yyyy)
	if err != nil {
		a.Logger.Warn("invalid value detected", zap.Error(err))
		return nil, err
	}

	yyyymmList := fyInterval(yyyyint)

	catInfo, err := a.Repo.GetCategoryInfo(categoryID)
	if err != nil {
		a.Logger.Error("failed to category info", zap.Error(err))
		return nil, err
	}

	var monthPrice []int
	var total int

	for _, yyyymm := range yyyymmList { // yyyymm
		monthInfo, err := a.Repo.GetCategoryMonthSummary(categoryID, yyyymm)
		if err != nil && !errors.Is(err, model.ErrTableNotFound) {
			a.Logger.Error("failed to get info from DB", zap.Error(err))
			return nil, err
		}
		if errors.Is(err, model.ErrTableNotFound) {
			// TableNotFound
			monthPrice = append(monthPrice, 0)
		} else {
			// Got Month information
			monthPrice = append(monthPrice, monthInfo.Total)
			total = total + monthInfo.Total
		}
	}

	sum = &model.CategoryYearSummary{
		CategoryID:   categoryID,
		CategoryName: catInfo.Name,
		MonthPrice:   monthPrice,
		Total:        total,
	}
	return sum, nil
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
