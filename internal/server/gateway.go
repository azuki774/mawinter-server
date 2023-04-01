package server

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"mawinter-server/internal/model"
	"mawinter-server/internal/openapi"
	"mawinter-server/internal/timeutil"
	"net/http"
	"strconv"

	"go.uber.org/zap"
)

type apigateway struct {
	Logger *zap.Logger
	ap1    APIServiceV1
	ap2    APIServiceV2
}

func (a *apigateway) Get(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "It is the root page.\n")
}

// create new YYYYMM table
// (POST /table/{year})
func (a *apigateway) PostTableYear(w http.ResponseWriter, r *http.Request, year float32) {
	err := a.ap1.CreateRecordTableYear(strconv.Itoa(int(year)))
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
// (POST /v1/record/)
func (a *apigateway) PostRecord(w http.ResponseWriter, r *http.Request) {
	var req model.RecordRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, err.Error())
		return
	}

	ctx := context.Background()
	res, err := a.ap1.AddRecord(ctx, req)
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
	yearSummary, err := a.ap1.GetYearSummary(ctx, year)
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

// V2

// (POST /v2/record)
func (a *apigateway) PostV2Record(w http.ResponseWriter, r *http.Request) {
	var req openapi.ReqRecord
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, err.Error())
		return
	}

	ctx := context.Background()
	rec, err := a.ap2.PostRecord(ctx, req)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, err.Error())
		return
	}

	outputJson, err := json.Marshal(&rec)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	fmt.Fprint(w, string(outputJson))
}

// (POST /v2/record/fixmonth)
func (a *apigateway) PostV2RecordFixmonth(w http.ResponseWriter, r *http.Request, params openapi.PostV2RecordFixmonthParams) {
	ctx := context.Background()
	var yms string
	if params.Yyyymm == nil {
		// default value
		yms = timeutil.NowFunc().Format("200601")
	} else {
		yms = strconv.Itoa(*params.Yyyymm)
	}

	recs, err := a.ap2.PostMonthlyFixRecord(ctx, yms)
	if errors.Is(err, model.ErrAlreadyRecorded) {
		w.WriteHeader(http.StatusNoContent)
		return
	} else if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	outputJson, err := json.Marshal(&recs)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	fmt.Fprint(w, string(outputJson))
}

// Your GET endpoint
// (GET /v2/record/summary/{year})
func (a *apigateway) GetV2RecordYear(w http.ResponseWriter, r *http.Request, year int) {
	ctx := context.Background()
	yearSummary, err := a.ap2.GetV2YearSummary(ctx, year)
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

// Your GET endpoint
// (GET /v2/record/{yyyymm})
func (a *apigateway) GetV2RecordYyyymm(w http.ResponseWriter, r *http.Request, yyyymm string, params openapi.GetV2RecordYyyymmParams) {
	ctx := context.Background()

	if params.CategoryId != nil {
		w.WriteHeader(http.StatusNotImplemented)
		return
	}

	err := model.ValidYYYYMM(yyyymm)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, err.Error())
		return
	}

	recs, err := a.ap2.GetYYYYMMRecords(ctx, yyyymm, params)
	if errors.Is(err, model.ErrNotFound) {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprint(w, err.Error())
		return
	} else if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, err.Error())
		return
	}

	outputJson, err := json.Marshal(&recs)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, string(outputJson))
}

// GET v2/record/{yyyymm}/recent?num=x
func (a *apigateway) GetV2RecordYyyymmRecent(w http.ResponseWriter, r *http.Request, yyyymm string, params openapi.GetV2RecordYyyymmRecentParams) {
	const defaultNum = 10
	ctx := context.Background()

	if params.Num == nil {
		params.Num = int2ptr(defaultNum)
	}

	err := model.ValidYYYYMM(yyyymm)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, err.Error())
		return
	}

	recs, err := a.ap2.GetYYYYMMRecordsRecent(ctx, yyyymm, *params.Num)
	if errors.Is(err, model.ErrNotFound) {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprint(w, err.Error())
		return
	} else if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, err.Error())
		return
	}

	outputJson, err := json.Marshal(&recs)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, string(outputJson))
}

// (POST /v2/table/{year})
func (a *apigateway) PostV2TableYear(w http.ResponseWriter, r *http.Request, year int) {
	ctx := context.Background()

	_, err := model.ValidYYYY(strconv.Itoa(year))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, err.Error())
	}

	err = a.ap2.CreateTableYear(ctx, year)
	if errors.Is(err, model.ErrAlreadyRecorded) {
		w.WriteHeader(http.StatusNoContent)
		return
	} else if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, err.Error())
	}

	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "record table created.\n")
}

// (GET /version)
func (a *apigateway) GetVersion(w http.ResponseWriter, r *http.Request) {
	vers := openapi.GetVersionJSONBody{
		Version:  str2ptr(Version),
		Revision: str2ptr(Revision),
		Build:    str2ptr(Build),
	}
	outputJson, err := json.Marshal(&vers)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, string(outputJson))
}

func str2ptr(a string) *string {
	return &a
}

func int2ptr(a int) *int {
	return &a
}
