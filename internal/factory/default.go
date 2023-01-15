package factory

import (
	"fmt"
	v1 "mawinter-server/internal/api/v1"
	"mawinter-server/internal/client"
	"mawinter-server/internal/register"
	"mawinter-server/internal/repository"
	"mawinter-server/internal/server"
	"os"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func NewLogger() (*zap.Logger, error) {
	config := zap.NewProductionConfig()
	// config.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
	config.EncoderConfig.EncodeTime = JSTTimeEncoder
	l, err := config.Build()

	l.WithOptions(zap.AddStacktrace(zap.ErrorLevel))
	if err != nil {
		fmt.Printf("failed to create logger: %v\n", err)
	}
	return l, err
}

func JSTTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	const layout = "2006-01-02T15:04:05+09:00"
	jst := time.FixedZone("Asia/Tokyo", 9*60*60)
	enc.AppendString(t.In(jst).Format(layout))
}

func NewService(l *zap.Logger, db *repository.DBRepository) (ap *v1.APIService) {
	return &v1.APIService{Logger: l, Repo: db}
}

func NewFetcherBill(host string, port string) *client.BillFetcher {
	return &client.BillFetcher{Host: host, Port: port}
}

func NewRegisterService(l *zap.Logger, db *repository.DBRepository, fet *client.BillFetcher) (ap *register.RegisterService) {
	return &register.RegisterService{Logger: l, DB: db, BillFetcher: fet}
}

func NewServer(l *zap.Logger, ap *v1.APIService) *server.Server {
	return &server.Server{Logger: l, APIService: ap, BasicAuth: struct {
		User string
		Pass string
	}{os.Getenv("BASIC_AUTH_USERNAME"), os.Getenv("BASIC_AUTH_PASSWORD")}}
}
