package factory

import (
	"fmt"
	v1 "mawinter-server/internal/api/v1"
	"mawinter-server/internal/repository"
	"mawinter-server/internal/server"
	"os"

	"go.uber.org/zap"
)

func NewLogger() (*zap.Logger, error) {
	config := zap.NewProductionConfig()
	// config.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
	l, err := config.Build()

	l.WithOptions(zap.AddStacktrace(zap.ErrorLevel))
	if err != nil {
		fmt.Printf("failed to create logger: %v\n", err)
	}
	return l, err
}

func NewService(l *zap.Logger, db *repository.DBRepository) (ap *v1.APIService) {
	return &v1.APIService{Logger: l, Repo: db}
}

func NewServer(l *zap.Logger, ap *v1.APIService) *server.Server {
	return &server.Server{Logger: l, APIService: ap, BasicAuth: struct {
		User string
		Pass string
	}{os.Getenv("BASIC_AUTH_USERNAME"), os.Getenv("BASIC_AUTH_PASSWORD")}}
}
