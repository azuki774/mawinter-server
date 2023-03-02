package server

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"mawinter-server/internal/model"
	"net/http"
	"strconv"

	"go.uber.org/zap"
)

type apigateway struct {
	Logger     *zap.Logger
	APIService APIService
}

func (a *apigateway) Get(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "It is the root page.\n")
}

// create new YYYYMM table
// (POST /table/{year})
func (a *apigateway) PostTableYear(w http.ResponseWriter, r *http.Request, year float32) {
	err := a.APIService.CreateRecordTableYear(strconv.Itoa(int(year)))
	if err != nil {
		if errors.Is(err, model.ErrInvalidValue) {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "Error: %s\n", err.Error())
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error: %s\n", err.Error())
		return
	}

	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "record table created.\n")
}

// post new a record
// (POST /v1/record)
func (a *apigateway) PostRecord(w http.ResponseWriter, r *http.Request) {
	var req model.RecordRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, err.Error())
		return
	}

	ctx := context.Background()
	res, err := a.APIService.AddRecord(ctx, req)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, err.Error())
		return
	}

	outputJson, err := json.Marshal(&res)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	fmt.Fprint(w, string(outputJson))
}

// get year records
// (GET /v1/record/year/{year})
func (a *apigateway) GetRecordYearYear(w http.ResponseWriter, r *http.Request, year string) {
	ctx := context.Background()
	yearSummary, err := a.APIService.GetYearSummary(ctx, year)
	if err != nil {
		if errors.Is(err, model.ErrInvalidValue) {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprint(w, err.Error())
			return
		}

		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, err.Error())
		return
	}

	outputJson, err := json.Marshal(&yearSummary)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	fmt.Fprint(w, string(outputJson))
}
