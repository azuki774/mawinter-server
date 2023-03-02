// Package server provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen version v1.12.4 DO NOT EDIT.
package server

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"fmt"
	"net/http"
	"net/url"
	"path"
	"strings"

	"github.com/deepmap/oapi-codegen/pkg/runtime"
	"github.com/getkin/kin-openapi/openapi3"
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

	return r
}

// Base64 encoded, gzipped, json marshaled Swagger object
var swaggerSpec = []string{

	"H4sIAAAAAAAC/7RVz28bRRT+V1YDUg+M65m1E5I9gqCKUAWKejEhisbrZ3vK7s54ZjaxG1mizgEEB6Qi",
	"FXGp4FZxiDgAEhKQP2YVtf0v0Jt17LW9gQCpZK92dn689773fd+cklilWmWQOUuiU2LjIaTCv+7DqGEg",
	"VqaHI22UBuMk+LmecOBkCvgOY5HqBEhEQsa2GWchoUQL58BkJCIHrLF7+BahxE00LrLOyGxAppT0jUrx",
	"gI2JFFJVO6GNjKEyIzMHAzA4VX7Z2IMz0vnsylqODIxysG6ZkOo+hBjH44Z1SidyMHR4kuyRiPSzUa7z",
	"0XZ3JHZzf9w+2GthiYWDgTKTI9mrz3KxIBMp1JaIyK6jGvIGw98DxiL/+/hfwbmSTJan3TKX14yy1Sqz",
	"cEOYFT+Odbt9stMVduLPWyA1AWGObJ6mwkxuiviyyH8GfFGtdJCW5Ia+yBNHIkaXbWgxRmswFOO9chsP",
	"KUllVhnNFwtjxMRjp5xIapKsgFdf9M0gbMV6qMf9Y/Fo0t0i0ymeK7N+2eL5+ak4wWaahtCSUHIMxkqF",
	"IuV3GaaoNGQ4FZHWXfzkdTz0oDTxMQBXAmRjI7Ur9w5BJG4YxEOIP/0kE8mJmNjAgMtNFtzZc4G0gRtC",
	"YJRygRYDuEN8ICNw/x6mfg+wqCvKYLgpJYuWrwTwcDVLkvmUtLI1OV2e//ry+Rcvfroozj4vzn4vZhfF",
	"7MnLiz8uv/y+ePxdMfuq+Gy2kcZHyrr9UtmYjbeJd1TPsy5WmYPMRxJaJzL225oPLYZbqNUj9V75HnAc",
	"HIskhw2i8pDRiod6iTPOOLlSMelJO0+klCkhC6aGfMlE4ttcmjZGedNAn0TkjebS1ZtzS29W/LxkRxVw",
	"9O5NGD/8AFkRMn77ELRRT2vqxDFiUHrgdca3gVCJaKsGKh622v8RK1vBiq7B8q4B4aDnubikKTIxyOAk",
	"EIG5YtG44fWWofD7IrFQpS+KvHmKz+m14nq/Uy4IisfnweVvP796+k0x+6U4e1ac/VjMnlx+/fTyz2+v",
	"p/Q9mDO6A8Lgn9S3/X/29mCtuZ6i68198fz81Q/Pls05aNMtuk3fpjt0l3JGOac8pJyGtHW4cMuQt3Z5",
	"OKUbAfgtBWiHnE0PV5mxuAn+jiL1Xj1d9/1N9qCoVogzABfgIXPaWH8rCSNScGCsR1fiRvRiQsm82slV",
	"N0e5NNAjkTM5VKtYv6MPkXpOdBOokO724iyus0NaMWVtIPZamW/bdNwHmFENNVcgir3ivLo6nU7n/v3A",
	"FzKXNJjjqwJyk+CN4ZyOms1ExSIZKuuiHbbDCGZWd3NKMRnLYSsMw/HWIzKd/hUAAP//UOV/5xQLAAA=",
}

// GetSwagger returns the content of the embedded swagger specification file
// or error if failed to decode
func decodeSpec() ([]byte, error) {
	zipped, err := base64.StdEncoding.DecodeString(strings.Join(swaggerSpec, ""))
	if err != nil {
		return nil, fmt.Errorf("error base64 decoding spec: %s", err)
	}
	zr, err := gzip.NewReader(bytes.NewReader(zipped))
	if err != nil {
		return nil, fmt.Errorf("error decompressing spec: %s", err)
	}
	var buf bytes.Buffer
	_, err = buf.ReadFrom(zr)
	if err != nil {
		return nil, fmt.Errorf("error decompressing spec: %s", err)
	}

	return buf.Bytes(), nil
}

var rawSpec = decodeSpecCached()

// a naive cached of a decoded swagger spec
func decodeSpecCached() func() ([]byte, error) {
	data, err := decodeSpec()
	return func() ([]byte, error) {
		return data, err
	}
}

// Constructs a synthetic filesystem for resolving external references when loading openapi specifications.
func PathToRawSpec(pathToFile string) map[string]func() ([]byte, error) {
	var res = make(map[string]func() ([]byte, error))
	if len(pathToFile) > 0 {
		res[pathToFile] = rawSpec
	}

	return res
}

// GetSwagger returns the Swagger specification corresponding to the generated code
// in this file. The external references of Swagger specification are resolved.
// The logic of resolving external references is tightly connected to "import-mapping" feature.
// Externally referenced files must be embedded in the corresponding golang packages.
// Urls can be supported but this task was out of the scope.
func GetSwagger() (swagger *openapi3.T, err error) {
	var resolvePath = PathToRawSpec("")

	loader := openapi3.NewLoader()
	loader.IsExternalRefsAllowed = true
	loader.ReadFromURIFunc = func(loader *openapi3.Loader, url *url.URL) ([]byte, error) {
		var pathToFile = url.String()
		pathToFile = path.Clean(pathToFile)
		getSpec, ok := resolvePath[pathToFile]
		if !ok {
			err1 := fmt.Errorf("path not found: %s", pathToFile)
			return nil, err1
		}
		return getSpec()
	}
	var specData []byte
	specData, err = rawSpec()
	if err != nil {
		return
	}
	swagger, err = loader.LoadFromData(specData)
	if err != nil {
		return
	}
	return
}
