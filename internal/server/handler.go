package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"mawinter-server/internal/model"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func rootHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "It is the root page.\n")
}

func (s *Server) CreateRecordFunc(w http.ResponseWriter, r *http.Request) {
	// /api/record/
	var addRecord model.CreateRecord

	err := json.NewDecoder(r.Body).Decode(&addRecord)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	addRecordRes, err := s.APIService.CreateRecord(addRecord)
	if err != nil {
		// データ追加時エラー
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	outputJson, err := json.Marshal(&addRecordRes)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprint(w, string(outputJson))
}

func (s *Server) deleteRecordFunc(w http.ResponseWriter, r *http.Request) {
	pathParam := mux.Vars(r)
	id, err := strconv.ParseInt(pathParam["id"], 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}

	err = s.APIService.DeleteRecord(id)

	if err != nil {
		if errors.Is(err, model.ErrRecordNotFound) {
			w.WriteHeader(http.StatusNotFound)
			return
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	w.WriteHeader(http.StatusNoContent)
}

func (s *Server) getYearSummaryFunc(w http.ResponseWriter, r *http.Request) {
	pathParam := mux.Vars(r)

	year, err := strconv.ParseInt(pathParam["year"], 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}

	yearSummary, err := s.APIService.GetYearSummary(year)
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

func (s *Server) getRecentRecordFunc(w http.ResponseWriter, r *http.Request) {
	// TODO: 表示件数の変更をできるようにする
	getRecentData, err := s.APIService.GetRecentRecord(20)
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
