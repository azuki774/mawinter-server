package server

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"mawinter-server/internal/model"
	"net/http"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

func (s *Server) middlewareBasicAuth(h http.Handler) http.Handler {
	// 認証失敗時にerrorを返す
	username := s.BasicAuth.User
	password := s.BasicAuth.Pass

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if username == "" {
			// Basic Auth Skip
			h.ServeHTTP(w, r)
			return
		}
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
		if r.URL.Path != "/" {
			s.Logger.Info("access", zap.String("url", r.URL.Path), zap.String("X-Forwarded-For", r.Header.Get("X-Forwarded-For")))
		}
		h.ServeHTTP(w, r)
	})
}

func (s *Server) rootHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "It is the root page.\n")
}

func (s *Server) createRecordTableHandler(w http.ResponseWriter, r *http.Request) {
	pathParam := mux.Vars(r)
	yyyy := pathParam["year"]
	err := s.APIService.CreateRecordTableYear(yyyy)
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

func (s *Server) addRecordHandler(w http.ResponseWriter, r *http.Request) {
	var req model.RecordRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, err.Error())
		return
	}

	ctx := context.Background()
	res, err := s.APIService.AddRecord(ctx, req)
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
