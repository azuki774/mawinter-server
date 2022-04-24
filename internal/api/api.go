package api

import (
	"mawinter-expense/internal/db"
	l "mawinter-expense/internal/logger"
	"strconv"

	"go.uber.org/zap"
)

type GetYearSummaryStruct struct {
	CategoryId int64   `json:"categoryID"`
	Name       string  `json:"name"`
	Price      []int64 `json:"price"`
	Total      int64   `json:"total"`
}

type APIService interface {
	AddRecord(addRecord db.Records) (retAddRecord db.Records, err error)
	DeleteRecord(id int64) (err error)
	GetRecentRecord(dataNum int64) (getRecentData []db.RecordsDetails, err error)
	GetYearSummary(year int64) (yearSummary []GetYearSummaryStruct, err error)
}

type apiService struct {
	dbR db.DBRepository
}

func NewAPIService(dbR_in db.DBRepository) APIService {
	return &apiService{dbR: dbR_in}
}

func (apis *apiService) GetYearSummary(year int64) (yearSummary []GetYearSummaryStruct, err error) {
	l.Logger.Info("API", "GetYearSummary called")
	tx := apis.dbR.OpenTx()
	defer apis.dbR.CloseTx(tx, err)

	fetchDBData, err := apis.dbR.GetYearSummaryDB(tx, year) // fetchDBData is sorted by category_id
	if err != nil {
		l.Logger.Errorw("API error", "name", "GetYearSummary", "error", zap.Error(err))
		return nil, err
	}

	fetchIndexes := getYearSummaryMakeIndex(fetchDBData)
	if len(fetchIndexes) == 0 {
		return []GetYearSummaryStruct{}, nil
	}

	yearSummary = make([]GetYearSummaryStruct, fetchIndexes[len(fetchIndexes)-1]+1)

	// Price fieldを初期化
	for i, _ := range yearSummary {
		yearSummary[i].Price = make([]int64, 12)
	}

	for i, v := range fetchDBData {
		ind := fetchIndexes[i]
		yearSummary[ind].CategoryId = v.CategoryId
		yearSummary[ind].Name = v.Name
		// 20xxyy -> yy = datenumに変換
		datenum, err := strconv.Atoi(v.YearMonth[4:6])
		if err != nil {
			l.Logger.Errorw("API error", "name", "GetYearSummary", "error", zap.Error(err))
			return nil, err
		}

		yearSummary[ind].Price[db.TransMonthToIndex(datenum)] = v.Price
		yearSummary[ind].Total += v.Price
	}

	return yearSummary, nil
}

func getYearSummaryMakeIndex(yearSummary []db.GetYearSummaryDBStruct) (indexes []int) {
	arr := make([]int64, len(yearSummary))
	for i, v := range yearSummary {
		arr[i] = v.CategoryId
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

	for i, _ := range arr {
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

func (apis *apiService) AddRecord(addRecord db.Records) (retAddRecord db.Records, err error) {
	l.Logger.Infow("API called", "name", "AddRecord")
	tx := apis.dbR.OpenTx()
	defer apis.dbR.CloseTx(tx, err)
	retAddRecord, err = apis.dbR.AddRecordDB(tx, addRecord)
	if err != nil {
		l.Logger.Errorw("API error", "name", "AddRecord", "error", zap.Error(err))
		return db.Records{}, err
	}

	return retAddRecord, nil
}

func (apis *apiService) DeleteRecord(id int64) (err error) {
	l.Logger.Infow("API called", "name", "DeleteRecord")
	tx := apis.dbR.OpenTx()
	defer apis.dbR.CloseTx(tx, err)
	err = apis.dbR.DeleteRecordDB(tx, id)
	if err != nil {
		l.Logger.Errorw("API error", "name", "DeleteRecord", "error", zap.Error(err))
	}

	return nil
}

// 直近の Record データをdataNum件分取得。
func (apis *apiService) GetRecentRecord(dataNum int64) (getRecentData []db.RecordsDetails, err error) {
	l.Logger.Infow("API called", "name", "GetRecentRecord")
	tx := apis.dbR.OpenTx()
	defer apis.dbR.CloseTx(tx, err)

	if dataNum == 0 {
		dataNum = 20
	}

	getRecentDataDB, err := apis.dbR.GetRecentRecord(tx, dataNum)

	if err != nil {
		l.Logger.Errorw("API error", "name", "GetRecentRecord", "error", zap.Error(err))
		return nil, err
	}

	// Memo が NULL なら空文字にする
	for i := range getRecentDataDB {
		var agetRecentData db.RecordsDetails
		agetRecentData.Id = getRecentDataDB[i].Id
		agetRecentData.CategoryId = getRecentDataDB[i].CategoryId
		agetRecentData.Name = getRecentDataDB[i].Name
		agetRecentData.Date = getRecentDataDB[i].Date
		agetRecentData.Price = getRecentDataDB[i].Price
		agetRecentData.Memo = getRecentDataDB[i].Memo
		getRecentData = append(getRecentData, agetRecentData)
	}

	return getRecentData, nil
}
