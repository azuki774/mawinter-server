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

const defaultGetRecordNum = 20

type apigateway struct {
	Logger *zap.Logger
	ap2    APIServiceV2
}

func (a *apigateway) Get(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "It is the root page.\n")
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
	if err != nil && errors.Is(err, model.ErrUnknownCategoryID) {
		// Category ID情報がDBにない場合
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, err.Error())
		return
	} else if err != nil && errors.Is(err, model.ErrAlreadyRecorded) {
		// confirm month 確定済の月だった場合
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "already confirmed month")
		return
	} else if err != nil {
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

// (GET /v2/record/count)
func (a *apigateway) GetV2RecordCount(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	rec, err := a.ap2.GetRecordsCount(ctx)
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
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, string(outputJson))
}

func (a *apigateway) GetV2RecordAvailable(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	res, err := a.ap2.GetRecordsAvailable(ctx)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	outputJson, err := json.Marshal(&res)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
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

// (GET /v2/record)
func (a *apigateway) GetV2Record(w http.ResponseWriter, r *http.Request, params openapi.GetV2RecordParams) {
	ctx := context.Background()
	opts := model.GetRecordOption{ // default value
		Num:    defaultGetRecordNum,
		Offset: 0,
	}

	if params.Num != nil {
		opts.Num = *params.Num
	}

	if params.Offset != nil {
		opts.Offset = *params.Offset
	}

	if params.Yyyymm != nil {
		err := model.ValidYYYYMM(*params.Yyyymm)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprint(w, err.Error())
			return
		}
		opts.YYYYMM = *params.Yyyymm
	}

	if params.CategoryId != nil {
		opts.CategoryID = *params.CategoryId
	}

	recs, err := a.ap2.GetRecords(ctx, opts)
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

// (GET /v2/record/{yyyymm}/confirm)
func (a *apigateway) GetV2RecordYyyymmConfirm(w http.ResponseWriter, r *http.Request, yyyymm string) {
	ctx := context.Background()
	err := model.ValidYYYYMM(yyyymm)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, err.Error())
		return
	}
	yc, err := a.ap2.GetMonthlyConfirm(ctx, yyyymm)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, err.Error())
		return
	}

	outputJson, err := json.Marshal(&yc)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, string(outputJson))
}

// (PUT /v2/record/{yyyymm}/confirm)
func (a *apigateway) PutV2TableYyyymmConfirm(w http.ResponseWriter, r *http.Request, yyyymm string) {
	ctx := context.Background()
	err := model.ValidYYYYMM(yyyymm)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, err.Error())
		return
	}

	var req openapi.ConfirmInfo
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, err.Error())
		return
	}

	yc, err := a.ap2.UpdateMonthlyConfirm(ctx, yyyymm, *req.Status)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, err.Error())
		return
	}

	outputJson, err := json.Marshal(&yc)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, string(outputJson))
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

func (a *apigateway) GetCategories(w http.ResponseWriter, r *http.Request) {
	cats, err := a.ap2.GetCategories(r.Context())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, err.Error())
		return
	}

	outputJson, err := json.Marshal(&cats)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, string(outputJson))
}

// (GET /v2/record/{id})
func (a *apigateway) GetV2RecordId(w http.ResponseWriter, r *http.Request, id int) {
	rec, err := a.ap2.GetRecordByID(r.Context(), id)
	if err != nil {
		if errors.Is(err, model.ErrNotFound) {
			w.WriteHeader(http.StatusNotFound)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, err.Error())
		}
		return
	}

	outputJson, err := json.Marshal(&rec)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, string(outputJson))
}

// (DELETE /v2/record/{id})
func (a *apigateway) DeleteV2RecordId(w http.ResponseWriter, r *http.Request, id int) {
	err := a.ap2.DeleteRecordByID(r.Context(), id)
	if err != nil {
		if errors.Is(err, model.ErrNotFound) {
			w.WriteHeader(http.StatusNotFound)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, err.Error())
		}
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func str2ptr(a string) *string {
	return &a
}
