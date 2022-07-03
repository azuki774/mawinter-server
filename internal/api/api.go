package api

import (
	"errors"
	"fmt"
	"mawinter-server/internal/azerror"
	"mawinter-server/internal/model"
	"mawinter-server/internal/repository"
	"strconv"
	"time"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

func init() {
	jst, err := time.LoadLocation("Asia/Tokyo")
	if err != nil {
		panic(err)
	}
	time.Local = jst
}

type APIService interface {
	CreateRecord(addRecord model.CreateRecord) (retAddRecord model.ShowRecord, err error)
	GetYearSummary(year int64) (yearSummary []model.YearSummary, err error)
	GetRecentRecord(n int) (getRecentData []model.ShowRecord, err error)
	DeleteRecord(id int64) (err error)
}

type apiService struct {
	dbR    repository.DBRepository
	logger *zap.SugaredLogger
}

func NewAPIService(dbR_in repository.DBRepository, l *zap.SugaredLogger) APIService {
	return &apiService{dbR: dbR_in, logger: l}
}

func (ap *apiService) CreateRecord(addRecord model.CreateRecord) (retAddRecord model.ShowRecord, err error) {
	ap.logger.Infow("API called", "func", "CreateRecord")
	var record model.Records = model.Records{
		CategoryID: addRecord.CategoryID,
		Price:      addRecord.Price,
		Memo:       addRecord.Memo,
	}

	if addRecord.Date == "" {
		record.Date = time.Now()
	} else {
		YYYYs := addRecord.Date[:4]
		YYYY, err := strconv.Atoi(YYYYs)
		if err != nil {
			return model.ShowRecord{}, fmt.Errorf("failed to parse Year: %w", err)
		}

		MMs := addRecord.Date[4:6]
		MM, err := strconv.Atoi(MMs)
		if err != nil {
			return model.ShowRecord{}, fmt.Errorf("failed to parse Month: %w", err)
		}
		DDs := addRecord.Date[6:8]
		DD, err := strconv.Atoi(DDs)
		if err != nil {
			return model.ShowRecord{}, fmt.Errorf("failed to parse Day: %w", err)
		}

		record.Date = time.Date(YYYY, time.Month(MM), DD, 0, 0, 0, 0, time.Local)
	}
	
	retAddRecord, err = ap.dbR.CreateRecordDB(record)
	if err != nil {
		ap.logger.Infow("API error", "func", "CreateRecord", "error", zap.Error(err))
		return model.ShowRecord{}, err
	}

	return retAddRecord, nil
}

func (ap *apiService) GetYearSummary(year int64) (yearSummary []model.YearSummary, err error) {
	ap.logger.Infow("API called", "func", "GetYearSummary")
	var fetchDBData []model.YearSummaryInter
	fetchDBData, err = ap.dbR.GetYearSummaryDB(year) // fetchDBData is sorted by category_id
	if err != nil {
		ap.logger.Infow("API error", "func", "GetYearSummary", "error", zap.Error(err))
		return nil, err
	}

	fetchIndexes := getYearSummaryMakeIndex(fetchDBData)
	if len(fetchIndexes) == 0 {
		return []model.YearSummary{}, nil
	}

	yearSummary = make([]model.YearSummary, fetchIndexes[len(fetchIndexes)-1]+1)

	// Price fieldを初期化
	for i := range yearSummary {
		yearSummary[i].Price = make([]int64, 12)
	}

	for i, v := range fetchDBData {
		ind := fetchIndexes[i]
		yearSummary[ind].CategoryID = v.CategoryID
		yearSummary[ind].Name = v.Name
		// 20xxyy -> yy = datenumに変換
		datenum, err := strconv.Atoi(v.YearMonth[4:6])
		if err != nil {
			return nil, err
		}

		yearSummary[ind].Price[repository.TransMonthToIndex(datenum)] = v.Total
		yearSummary[ind].Total += v.Total
	}

	return yearSummary, nil
}

func getYearSummaryMakeIndex(yearSummary []model.YearSummaryInter) (indexes []int) {
	arr := make([]int64, len(yearSummary))
	for i, v := range yearSummary {
		arr[i] = v.CategoryID
	}

	indexes = posCompSlideint64(arr)
	return indexes
}

// posCompSlideint64 : arr is sorted
func posCompSlideint64(arr []int64) (indexes []int) {
	if len(arr) == 0 {
		return []int{}
	} else if len(arr) == 1 {
		return []int{0}
	}
	var ind int = 0

	for i := range arr {
		if i == 0 {
			indexes = append(indexes, ind)
			continue
		}

		if arr[i-1] == arr[i] {
			indexes = append(indexes, ind)
		} else {
			ind++
			indexes = append(indexes, ind)
		}
	}
	return indexes
}

// 直近の Record データをdataNum 件分取得。
func (ap *apiService) GetRecentRecord(n int) (getRecentData []model.ShowRecord, err error) {
	ap.logger.Infow("API called", "func", "GetRecentRecord")
	if n == 0 {
		n = 20 // default
	}

	getRecentData, err = ap.dbR.GetRecentRecord(n)

	if err != nil {
		ap.logger.Errorw("API error", "func", "GetRecentRecord", "error", zap.Error(err))
		return nil, err
	}

	return getRecentData, nil
}

func (ap *apiService) DeleteRecord(id int64) (err error) {
	ap.logger.Infow("API called", "func", "DeleteRecord")
	err = ap.dbR.DeleteRecordDB(id)
	if err != nil {
		ap.logger.Errorw("API error", "func", "DeleteRecord", "error", zap.Error(err))
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return azerror.ErrRecordNotFound
		}
	}

	return nil
}
