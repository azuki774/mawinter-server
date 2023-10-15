package server

import (
	"context"
	"mawinter-server/internal/openapi"
	"net/http"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

// Binary Info
var (
	Version  string
	Revision string
	Build    string
)

// V2
type APIServiceV2 interface {
	// V2
	PostRecord(ctx context.Context, req openapi.ReqRecord) (rec openapi.Record, err error)
	PostMonthlyFixRecord(ctx context.Context, yyyymm string) (recs []openapi.Record, err error)
	GetRecords(ctx context.Context, num int) (recs []openapi.Record, err error)
	GetYYYYMMRecords(ctx context.Context, yyyymm string, params openapi.GetV2RecordYyyymmParams) (recs []openapi.Record, err error)
	GetYYYYMMRecordsRecent(ctx context.Context, yyyymm string, num int) (recs []openapi.Record, err error)
	GetV2YearSummary(ctx context.Context, year int) (sums []openapi.CategoryYearSummary, err error)
	GetMonthlyConfirm(ctx context.Context, yyyymm string) (yc openapi.ConfirmInfo, err error)
	UpdateMonthlyConfirm(ctx context.Context, yyyymm string, confirm bool) (yc openapi.ConfirmInfo, err error)
}
type Server struct {
	Logger    *zap.Logger
	Ap2       APIServiceV2
	BasicAuth struct {
		User string
		Pass string
	}
}

func (s *Server) Start(ctx context.Context) error {
	swagger, err := openapi.GetSwagger()
	if err != nil {
		s.Logger.Error("failed to get swagger spec", zap.Error(err))
		return err
	}
	swagger.Servers = nil
	r := chi.NewRouter()
	r.Use(s.middlewareLogging)

	openapi.HandlerFromMux(&apigateway{Logger: s.Logger, ap2: s.Ap2}, r)
	addr := ":8080"
	if err := http.ListenAndServe(addr, r); err != nil {
		s.Logger.Error("failed to listen and serve", zap.Error(err))
		return err
	}

	return nil
}

func (s *Server) middlewareLogging(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			s.Logger.Info("access", zap.String("url", r.URL.Path), zap.String("X-Forwarded-For", r.Header.Get("X-Forwarded-For")))
		}
		h.ServeHTTP(w, r)
	})
}
