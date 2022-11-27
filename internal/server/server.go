package server

import (
	"context"
	"mawinter-server/internal/model"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

type APIService interface {
	CreateRecordTableYear(yyyy string) (err error)
	AddRecord(ctx context.Context, req model.RecordRequest) (res model.Recordstruct, err error)
	GetYearSummary(ctx context.Context, yyyy string) (sum []*model.CategoryYearSummary, err error)
}

type Server struct {
	Logger     *zap.Logger
	APIService APIService
	BasicAuth  struct {
		User string
		Pass string
	}
}

func (s *Server) Start(ctx context.Context) error {
	router := mux.NewRouter()
	s.addRecordFunc(router)

	server := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	ctxIn, stop := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	defer stop()

	var errCh = make(chan error)
	go func() {
		errCh <- server.ListenAndServe()
	}()

	<-ctxIn.Done()
	if nerr := server.Shutdown(ctx); nerr != nil {
		s.Logger.Error("failed to shutdown server", zap.Error(nerr))
		return nerr
	}

	err := <-errCh
	if err != nil && err != http.ErrServerClosed {
		s.Logger.Error("failed to close server", zap.Error(err))
		return err
	}

	s.Logger.Info("http server close gracefully")
	return nil
}

func (s *Server) addRecordFunc(r *mux.Router) {
	r.HandleFunc("/", s.rootHandler)
	r.Use(s.middlewareLogging)
	// Required Basic Auth
	br := r.PathPrefix("/").Subrouter()
	br.Use(s.middlewareBasicAuth)
	br.HandleFunc("/table/{year}", s.createRecordTableHandler).Methods("POST")
	br.HandleFunc("/record/", s.addRecordHandler).Methods("POST")
}
