package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"mawinter-expense/internal/azerror"
	"mawinter-expense/internal/db"
	"mawinter-expense/internal/logger"
	"net/http"
	"os"
	"strconv"
	"time"

	httpdate "github.com/Songmu/go-httpdate"
	"github.com/gorilla/mux"
)

// HTTPの入力の際に使う構造体
type inWebRecords struct {
	Id         int64  `json:"id"`
	CategoryId int64  `json:"categoryID"`
	Date       string `json:"date"`
	Price      int64  `json:"price"`
	Memo       string `json:"memo"` // sql.NullString のかわり
}

// Date を time.Time に変換、日付指定の場合はTZ分だけ操作 inWebRecords -> inRecords に変換
func (inWebRecords *inWebRecords) inWebRecordToRecord() (record db.Records, err error) {
	record.Id = inWebRecords.Id
	record.CategoryId = inWebRecords.CategoryId
	record.Price = inWebRecords.Price
	if inWebRecords.Memo != "" {
		record.Memo.Valid = true
		record.Memo.String = inWebRecords.Memo
	}

	if inWebRecords.Date == "" {
		record.Date = time.Now()
		record.Date = record.Date.Add(9 * time.Hour)
		err = nil
	} else {
		// UTC -> JST
		record.Date, err = httpdate.Str2Time(inWebRecords.Date, nil)
	}

	if err != nil {
		logger.InfoPrint("inWebRecordToRecord: BadRequest")
		return db.Records{}, azerror.ErrBadRequest
	}

	return record, nil
}

func ServerStart() {
	router := mux.NewRouter()
	router.HandleFunc("/", rootHandler)
	router.Methods("GET").Path("/summary/year/{year}").HandlerFunc(getYearSummaryFunc)
	router.Methods("GET").Path("/record/recent/").HandlerFunc(getRecentRecordFunc)
	router.Methods("POST").Path("/record/").HandlerFunc(addRecordFunc)
	router.Methods("DELETE").Path("/record/{id}").HandlerFunc(deleteRecordFunc)
	logger.InfoPrint("API Start")
	http.ListenAndServe(":80", router)
}

func handleBasicAuth(w http.ResponseWriter, r *http.Request) error {
	// 認証失敗時にerrorを返す
	username := os.Getenv("BASIC_AUTH_USERNAME")
	password := os.Getenv("BASIC_AUTH_PASSWORD")
	inputUsername, inputPassword, ok := r.BasicAuth()
	if !ok {
		w.Header().Set("WWW-Authenticate", `Basic realm="SECRET AREA"`)
		w.WriteHeader(http.StatusUnauthorized) // 401
		fmt.Fprintf(w, "%d Not authorized.", http.StatusUnauthorized)
		logger.ErrorPrint("Not basic authorized")
		return azerror.ErrAuthorized
	}

	if inputUsername != username || inputPassword != password {
		w.Header().Set("WWW-Authenticate", `Basic realm="SECRET AREA"`)
		w.WriteHeader(http.StatusUnauthorized) // 401
		fmt.Fprintf(w, "%d Not authorized.\n", http.StatusUnauthorized)
		logger.ErrorPrint("Basic authorized password unmatched")
		return azerror.ErrAuthorized
	}
	return nil
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	logger.AccessInfoPrint("/mawinter/", "GET", r.Header.Get("X-Forwarded-For"))
	fmt.Fprintf(w, "It is the root page.\n")
}

func getYearSummaryFunc(w http.ResponseWriter, r *http.Request) {
	// /api/summary/year/{year}
	logger.AccessInfoPrint("/mawinter/summary/year/{year}", "GET", r.Header.Get("X-Forwarded-For"))
	err := handleBasicAuth(w, r)
	if err != nil {
		return
	}

	pathParam := mux.Vars(r)

	dbR := db.NewDBRepository(db.DB)
	as := NewAPIService(dbR)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	year, err := strconv.ParseInt(pathParam["year"], 10, 64)
	if err != nil {
		logger.InfoPrint(fmt.Sprintf("getYearSummaryFunc : parameter parse error (year) : %s", err.Error()))
		w.WriteHeader(http.StatusBadRequest)
	}

	yearSummary, err := as.GetYearSummary(year)
	if err != nil {
		logger.ErrorPrint(fmt.Sprintf("getYearSummaryFunc : %s", err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
	}

	if err != nil {
		return
	}

	outputJson, err := json.Marshal(&yearSummary)
	if err != nil {
		logger.ErrorPrint(fmt.Sprintf("getYearSummaryFunc : %s", err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	fmt.Fprint(w, string(outputJson))
}

func addRecordFunc(w http.ResponseWriter, r *http.Request) {
	// /api/record/
	logger.AccessInfoPrint("/mawinter/record/", "POST", r.Header.Get("X-Forwarded-For"))
	err := handleBasicAuth(w, r)
	if err != nil {
		return
	}

	var addRecord db.Records
	var inWebRecords inWebRecords
	json.NewDecoder(r.Body).Decode(&inWebRecords)
	addRecord, err = inWebRecords.inWebRecordToRecord()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	dbR := db.NewDBRepository(db.DB)
	as := NewAPIService(dbR)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	addRecordRes, err := as.AddRecord(addRecord)
	if err != nil {
		// データ追加時エラー
		logger.ErrorPrint(fmt.Sprintf("addRecordFunc : %s", err.Error()))
		w.WriteHeader(http.StatusBadRequest)
	}
	outputJson, err := json.Marshal(&addRecordRes)
	if err != nil {
		logger.ErrorPrint(fmt.Sprintf("addRecordFunc : %s", err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprint(w, string(outputJson))
}

func deleteRecordFunc(w http.ResponseWriter, r *http.Request) {
	logger.AccessInfoPrint("/mawinter/delete/{id}", "GET", r.Header.Get("X-Forwarded-For"))
	err := handleBasicAuth(w, r)
	if err != nil {
		return
	}

	pathParam := mux.Vars(r)
	id, err := strconv.ParseInt(pathParam["id"], 10, 64)
	if err != nil {
		logger.InfoPrint(fmt.Sprintf("deleteRecordFunc : parameter parse error (id) : %s", err.Error()))
		w.WriteHeader(http.StatusBadRequest)
	}

	dbR := db.NewDBRepository(db.DB)
	as := NewAPIService(dbR)

	err = as.DeleteRecord(id)

	if err != nil {
		if errors.Is(err, azerror.ErrRecordNotFound) {
			w.WriteHeader(http.StatusNotFound)
			return
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	w.WriteHeader(http.StatusNoContent)
}

func getRecentRecordFunc(w http.ResponseWriter, r *http.Request) {
	// /api/recent
	logger.AccessInfoPrint("/mawinter/recent", "GET", r.Header.Get("X-Forwarded-For"))
	err := handleBasicAuth(w, r)
	if err != nil {
		return
	}

	// TODO: 表示件数の変更をできるようにする

	dbR := db.NewDBRepository(db.DB)
	as := NewAPIService(dbR)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	getRecentData, err := as.GetRecentRecord(20)
	if err != nil {
		logger.ErrorPrint(fmt.Sprintf("getRecentRecordFunc : %s", err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
	}

	if err != nil {
		return
	}

	outputJson, err := json.Marshal(&getRecentData)
	if err != nil {
		logger.ErrorPrint(fmt.Sprintf("getRecentRecordFunc : %s", err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	fmt.Fprint(w, string(outputJson))
}
