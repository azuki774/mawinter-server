package server

import (
	"fmt"
	"net/http"

	"go.uber.org/zap"
)

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
		if r.URL.Path != "/" {
			s.Logger.Info("access", zap.String("url", r.URL.Path), zap.String("X-Forwarded-For", r.Header.Get("X-Forwarded-For")))
		}
		h.ServeHTTP(w, r)
	})
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "It is the root page.\n")
}
