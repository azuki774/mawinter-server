package server

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"mawinter-server/internal/azerror"
	"mawinter-server/internal/model"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

type APIService interface {
	CreateRecord(addRecord model.CreateRecord) (retAddRecord model.ShowRecord, err error)
	GetYearSummary(year int64) (yearSummary []model.YearSummary, err error)
	GetRecentRecord(n int) (getRecentData []model.ShowRecord, err error)
	DeleteRecord(id int64) (err error)
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
		Addr:    ":80",
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
	r.HandleFunc("/", rootHandler)
	r.Use(s.middlewareLogging)
	// Required Basic Auth
	br := r.PathPrefix("/").Subrouter()
	br.Methods("GET").Path("/summary/year/{year}").HandlerFunc(s.getYearSummaryFunc)
	br.Methods("GET").Path("/record/recent/").HandlerFunc(s.getRecentRecordFunc)
	br.Methods("POST").Path("/record/").HandlerFunc(s.CreateRecordFunc)
	br.Methods("DELETE").Path("/record/{id}").HandlerFunc(s.deleteRecordFunc)
	br.Use(s.middlewareBasicAuth)
}

func (s *Server) middlewareBasicAuth(h http.Handler) http.Handler {
	// 認証失敗時にerrorを返す
	username := s.BasicAuth.User
	password := s.BasicAuth.Pass
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		inputUsername, inputPassword, ok := r.BasicAuth()
		if !ok {
			w.Header().Set("WWW-Authenticate", `Basic realm="SECRET AREA"`)
			w.WriteHeader(http.StatusUnauthorized) // 401
			fmt.Fprintf(w, "%d Not authorized.", http.StatusUnauthorized)
			return
		}

		if inputUsername != username || inputPassword != password {
			w.Header().Set("WWW-Authenticate", `Basic realm="SECRET AREA"`)
			w.WriteHeader(http.StatusUnauthorized) // 401
			fmt.Fprintf(w, "%d Not authorized.\n", http.StatusUnauthorized)
			return
		}
		// Basic Auth OK
		h.ServeHTTP(w, r)
	})
}

func (s *Server) middlewareLogging(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		s.Logger.Info("access", zap.String("url", r.URL.Path), zap.String("X-Forwarded-For", r.Header.Get("X-Forwarded-For")))
		h.ServeHTTP(w, r)
	})
}

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
