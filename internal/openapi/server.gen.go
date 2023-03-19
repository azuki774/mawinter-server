// Package openapi provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen version v1.12.4 DO NOT EDIT.
package openapi

import (
	"fmt"
	"net/http"

	"github.com/deepmap/oapi-codegen/pkg/runtime"
	"github.com/go-chi/chi/v5"
)

// ServerInterface represents all server handlers.
type ServerInterface interface {
	// health check
	// (GET /)
	Get(w http.ResponseWriter, r *http.Request)
	// post new a record
	// (POST /record/)
	PostRecord(w http.ResponseWriter, r *http.Request)
	// get year records
	// (GET /record/year/{year})
	GetRecordYearYear(w http.ResponseWriter, r *http.Request, year string)
	// create new YYYYMM table
	// (POST /table/{year})
	PostTableYear(w http.ResponseWriter, r *http.Request, year float32)
	// create record
	// (POST /v2/record)
	PostV2Record(w http.ResponseWriter, r *http.Request)
	// create fixmonth record
	// (POST /v2/record/fixmonth)
	PostV2RecordFixmonth(w http.ResponseWriter, r *http.Request, params PostV2RecordFixmonthParams)
	// get year summary
	// (GET /v2/record/summary/{year})
	GetV2RecordYear(w http.ResponseWriter, r *http.Request, year int)
	// get month records
	// (GET /v2/record/{yyyymm})
	GetV2RecordYyyymm(w http.ResponseWriter, r *http.Request, yyyymm string, params GetV2RecordYyyymmParams)
	// create record table
	// (POST /v2/table/{year})
	PostV2TableYear(w http.ResponseWriter, r *http.Request, year int)
	// Your GET endpoint
	// (GET /version)
	GetVersion(w http.ResponseWriter, r *http.Request)
}

// ServerInterfaceWrapper converts contexts to parameters.
type ServerInterfaceWrapper struct {
	Handler            ServerInterface
	HandlerMiddlewares []MiddlewareFunc
	ErrorHandlerFunc   func(w http.ResponseWriter, r *http.Request, err error)
}

type MiddlewareFunc func(http.Handler) http.Handler

// Get operation middleware
func (siw *ServerInterfaceWrapper) Get(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var handler http.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.Get(w, r)
	})

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r.WithContext(ctx))
}

// PostRecord operation middleware
func (siw *ServerInterfaceWrapper) PostRecord(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var handler http.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.PostRecord(w, r)
	})

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r.WithContext(ctx))
}

// GetRecordYearYear operation middleware
func (siw *ServerInterfaceWrapper) GetRecordYearYear(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var err error

	// ------------- Path parameter "year" -------------
	var year string

	err = runtime.BindStyledParameterWithLocation("simple", false, "year", runtime.ParamLocationPath, chi.URLParam(r, "year"), &year)
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "year", Err: err})
		return
	}

	var handler http.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.GetRecordYearYear(w, r, year)
	})

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r.WithContext(ctx))
}

// PostTableYear operation middleware
func (siw *ServerInterfaceWrapper) PostTableYear(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var err error

	// ------------- Path parameter "year" -------------
	var year float32

	err = runtime.BindStyledParameterWithLocation("simple", false, "year", runtime.ParamLocationPath, chi.URLParam(r, "year"), &year)
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "year", Err: err})
		return
	}

	var handler http.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.PostTableYear(w, r, year)
	})

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r.WithContext(ctx))
}

// PostV2Record operation middleware
func (siw *ServerInterfaceWrapper) PostV2Record(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var handler http.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.PostV2Record(w, r)
	})

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r.WithContext(ctx))
}

// PostV2RecordFixmonth operation middleware
func (siw *ServerInterfaceWrapper) PostV2RecordFixmonth(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var err error

	// Parameter object where we will unmarshal all parameters from the context
	var params PostV2RecordFixmonthParams

	// ------------- Optional query parameter "yyyymm" -------------

	err = runtime.BindQueryParameter("form", true, false, "yyyymm", r.URL.Query(), &params.Yyyymm)
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "yyyymm", Err: err})
		return
	}

	var handler http.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.PostV2RecordFixmonth(w, r, params)
	})

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r.WithContext(ctx))
}

// GetV2RecordYear operation middleware
func (siw *ServerInterfaceWrapper) GetV2RecordYear(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var err error

	// ------------- Path parameter "year" -------------
	var year int

	err = runtime.BindStyledParameterWithLocation("simple", false, "year", runtime.ParamLocationPath, chi.URLParam(r, "year"), &year)
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "year", Err: err})
		return
	}

	var handler http.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.GetV2RecordYear(w, r, year)
	})

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r.WithContext(ctx))
}

// GetV2RecordYyyymm operation middleware
func (siw *ServerInterfaceWrapper) GetV2RecordYyyymm(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var err error

	// ------------- Path parameter "yyyymm" -------------
	var yyyymm string

	err = runtime.BindStyledParameterWithLocation("simple", false, "yyyymm", runtime.ParamLocationPath, chi.URLParam(r, "yyyymm"), &yyyymm)
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "yyyymm", Err: err})
		return
	}

	// Parameter object where we will unmarshal all parameters from the context
	var params GetV2RecordYyyymmParams

	// ------------- Optional query parameter "from" -------------

	err = runtime.BindQueryParameter("form", true, false, "from", r.URL.Query(), &params.From)
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "from", Err: err})
		return
	}

	// ------------- Optional query parameter "category_id" -------------

	err = runtime.BindQueryParameter("form", true, false, "category_id", r.URL.Query(), &params.CategoryId)
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "category_id", Err: err})
		return
	}

	var handler http.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.GetV2RecordYyyymm(w, r, yyyymm, params)
	})

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r.WithContext(ctx))
}

// PostV2TableYear operation middleware
func (siw *ServerInterfaceWrapper) PostV2TableYear(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var err error

	// ------------- Path parameter "year" -------------
	var year int

	err = runtime.BindStyledParameterWithLocation("simple", false, "year", runtime.ParamLocationPath, chi.URLParam(r, "year"), &year)
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "year", Err: err})
		return
	}

	var handler http.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.PostV2TableYear(w, r, year)
	})

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r.WithContext(ctx))
}

// GetVersion operation middleware
func (siw *ServerInterfaceWrapper) GetVersion(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var handler http.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.GetVersion(w, r)
	})

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r.WithContext(ctx))
}

type UnescapedCookieParamError struct {
	ParamName string
	Err       error
}

func (e *UnescapedCookieParamError) Error() string {
	return fmt.Sprintf("error unescaping cookie parameter '%s'", e.ParamName)
}

func (e *UnescapedCookieParamError) Unwrap() error {
	return e.Err
}

type UnmarshallingParamError struct {
	ParamName string
	Err       error
}

func (e *UnmarshallingParamError) Error() string {
	return fmt.Sprintf("Error unmarshalling parameter %s as JSON: %s", e.ParamName, e.Err.Error())
}

func (e *UnmarshallingParamError) Unwrap() error {
	return e.Err
}

type RequiredParamError struct {
	ParamName string
}

func (e *RequiredParamError) Error() string {
	return fmt.Sprintf("Query argument %s is required, but not found", e.ParamName)
}

type RequiredHeaderError struct {
	ParamName string
	Err       error
}

func (e *RequiredHeaderError) Error() string {
	return fmt.Sprintf("Header parameter %s is required, but not found", e.ParamName)
}

func (e *RequiredHeaderError) Unwrap() error {
	return e.Err
}

type InvalidParamFormatError struct {
	ParamName string
	Err       error
}

func (e *InvalidParamFormatError) Error() string {
	return fmt.Sprintf("Invalid format for parameter %s: %s", e.ParamName, e.Err.Error())
}

func (e *InvalidParamFormatError) Unwrap() error {
	return e.Err
}

type TooManyValuesForParamError struct {
	ParamName string
	Count     int
}

func (e *TooManyValuesForParamError) Error() string {
	return fmt.Sprintf("Expected one value for %s, got %d", e.ParamName, e.Count)
}

// Handler creates http.Handler with routing matching OpenAPI spec.
func Handler(si ServerInterface) http.Handler {
	return HandlerWithOptions(si, ChiServerOptions{})
}

type ChiServerOptions struct {
	BaseURL          string
	BaseRouter       chi.Router
	Middlewares      []MiddlewareFunc
	ErrorHandlerFunc func(w http.ResponseWriter, r *http.Request, err error)
}

// HandlerFromMux creates http.Handler with routing matching OpenAPI spec based on the provided mux.
func HandlerFromMux(si ServerInterface, r chi.Router) http.Handler {
	return HandlerWithOptions(si, ChiServerOptions{
		BaseRouter: r,
	})
}

func HandlerFromMuxWithBaseURL(si ServerInterface, r chi.Router, baseURL string) http.Handler {
	return HandlerWithOptions(si, ChiServerOptions{
		BaseURL:    baseURL,
		BaseRouter: r,
	})
}

// HandlerWithOptions creates http.Handler with additional options
func HandlerWithOptions(si ServerInterface, options ChiServerOptions) http.Handler {
	r := options.BaseRouter

	if r == nil {
		r = chi.NewRouter()
	}
	if options.ErrorHandlerFunc == nil {
		options.ErrorHandlerFunc = func(w http.ResponseWriter, r *http.Request, err error) {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
	}
	wrapper := ServerInterfaceWrapper{
		Handler:            si,
		HandlerMiddlewares: options.Middlewares,
		ErrorHandlerFunc:   options.ErrorHandlerFunc,
	}

	r.Group(func(r chi.Router) {
		r.Get(options.BaseURL+"/", wrapper.Get)
	})
	r.Group(func(r chi.Router) {
		r.Post(options.BaseURL+"/record/", wrapper.PostRecord)
	})
	r.Group(func(r chi.Router) {
		r.Get(options.BaseURL+"/record/year/{year}", wrapper.GetRecordYearYear)
	})
	r.Group(func(r chi.Router) {
		r.Post(options.BaseURL+"/table/{year}", wrapper.PostTableYear)
	})
	r.Group(func(r chi.Router) {
		r.Post(options.BaseURL+"/v2/record", wrapper.PostV2Record)
	})
	r.Group(func(r chi.Router) {
		r.Post(options.BaseURL+"/v2/record/fixmonth", wrapper.PostV2RecordFixmonth)
	})
	r.Group(func(r chi.Router) {
		r.Get(options.BaseURL+"/v2/record/summary/{year}", wrapper.GetV2RecordYear)
	})
	r.Group(func(r chi.Router) {
		r.Get(options.BaseURL+"/v2/record/{yyyymm}", wrapper.GetV2RecordYyyymm)
	})
	r.Group(func(r chi.Router) {
		r.Post(options.BaseURL+"/v2/table/{year}", wrapper.PostV2TableYear)
	})
	r.Group(func(r chi.Router) {
		r.Get(options.BaseURL+"/version", wrapper.GetVersion)
	})

	return r
}
