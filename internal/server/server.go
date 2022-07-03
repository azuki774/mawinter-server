package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"mawinter-server/internal/api"
	"mawinter-server/internal/azerror"
	"mawinter-server/internal/model"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

var logger *zap.SugaredLogger
var apis api.APIService

func Start(as api.APIService, l *zap.SugaredLogger) {
	logger = l
	apis = as

	router := mux.NewRouter()
	router.HandleFunc("/", rootHandler)
	router.Methods("GET").Path("/summary/year/{year}").HandlerFunc(getYearSummaryFunc)
	router.Methods("GET").Path("/record/recent/").HandlerFunc(getRecentRecordFunc)
	router.Methods("POST").Path("/record/").HandlerFunc(CreateRecordFunc)
	router.Methods("DELETE").Path("/record/{id}").HandlerFunc(deleteRecordFunc)
	logger.Info("WebServer Start")
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
		logger.Warnw("not authorized access (no basic auth access)", "Web", "Basic authorized password unmatched access")
		return azerror.ErrAuthorized
	}

	if inputUsername != username || inputPassword != password {
		w.Header().Set("WWW-Authenticate", `Basic realm="SECRET AREA"`)
		w.WriteHeader(http.StatusUnauthorized) // 401
		fmt.Fprintf(w, "%d Not authorized.\n", http.StatusUnauthorized)
		logger.Warnw("not authorized access", "Web", "Basic authorized password unmatched access")
		return azerror.ErrAuthorized
	}
	return nil
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	logger.Infow("access", "url", "/mawinter/", "GET", r.Header.Get("X-Forwarded-For"))
	fmt.Fprintf(w, "It is the root page.\n")
}

func CreateRecordFunc(w http.ResponseWriter, r *http.Request) {
	// /api/record/
	logger.Infow("Access", "url", "/mawinter/record/", "method", "POST", "ip", r.Header.Get("X-Forwarded-For"))
	err := handleBasicAuth(w, r)
	if err != nil {
		return
	}

	var addRecord model.CreateRecord

	err = json.NewDecoder(r.Body).Decode(&addRecord)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		logger.Warnw(err.Error())
		return
	}

	addRecordRes, err := apis.CreateRecord(addRecord)
	if err != nil {
		// データ追加時エラー
		w.WriteHeader(http.StatusBadRequest)
		logger.Warnw(err.Error())
		return
	}
	outputJson, err := json.Marshal(&addRecordRes)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		logger.Warnw(err.Error())
		return
	}
	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprint(w, string(outputJson))
}

func deleteRecordFunc(w http.ResponseWriter, r *http.Request) {
	logger.Infow("access", "url", "/mawinter/delete/{id}", "method", "DELETE", "ip", r.Header.Get("X-Forwarded-For"))
	err := handleBasicAuth(w, r)
	if err != nil {
		return
	}

	pathParam := mux.Vars(r)
	id, err := strconv.ParseInt(pathParam["id"], 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}

	err = apis.DeleteRecord(id)

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

func getYearSummaryFunc(w http.ResponseWriter, r *http.Request) {
	// /api/summary/year/{year}
	logger.Infow("access", "url", "/mawinter/summary/year/{year}", "method", "GET", "ip", r.Header.Get("X-Forwarded-For"))
	err := handleBasicAuth(w, r)
	if err != nil {
		return
	}

	pathParam := mux.Vars(r)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	year, err := strconv.ParseInt(pathParam["year"], 10, 64)
	if err != nil {
		logger.Infow("access", "Web", "getYearSummaryFunc", "parameter parse error (year)", zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
	}

	yearSummary, err := apis.GetYearSummary(year)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	if err != nil {
		return
	}

	outputJson, err := json.Marshal(&yearSummary)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	fmt.Fprint(w, string(outputJson))
}

func getRecentRecordFunc(w http.ResponseWriter, r *http.Request) {
	// /api/recent
	logger.Infow("Access", "url", "/mawinter/recent", "method", "DELETE", "ip", r.Header.Get("X-Forwarded-For"))
	err := handleBasicAuth(w, r)
	if err != nil {
		return
	}

	// TODO: 表示件数の変更をできるようにする

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	getRecentData, err := apis.GetRecentRecord(20)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	if err != nil {
		return
	}

	outputJson, err := json.Marshal(&getRecentData)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	fmt.Fprint(w, string(outputJson))
}
