package factory

import (
	"fmt"
	v1 "mawinter-server/internal/api/v1"
	v2 "mawinter-server/internal/api/v2"
	"mawinter-server/internal/client"
	"mawinter-server/internal/register"
	v1db "mawinter-server/internal/repository/v1"
	v2db "mawinter-server/internal/repository/v2"
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

func NewServiceV1(l *zap.Logger, db *v1db.DBRepository) (ap *v1.APIService) {
	return &v1.APIService{Logger: l, Repo: db}
}

func NewServiceV2(l *zap.Logger, db *v2db.DBRepository) (ap *v2.APIService) {
	return &v2.APIService{Logger: l, Repo: db}
}

func NewFetcherBill(billEndpoint string) *client.BillFetcher {
	return &client.BillFetcher{BillEndpoint: billEndpoint}
}

func NewRegisterService(l *zap.Logger, db *v1db.DBRepository, fet *client.BillFetcher, mc *client.MailClient) (ap *register.RegisterService) {
	return &register.RegisterService{Logger: l, DB: db, BillFetcher: fet, MailClient: mc}
}

func NewServer(l *zap.Logger, ap1 *v1.APIService, ap2 *v2.APIService) *server.Server {
	return &server.Server{Logger: l, Ap1: ap1, Ap2: ap2, BasicAuth: struct {
		User string
		Pass string
	}{os.Getenv("BASIC_AUTH_USERNAME"), os.Getenv("BASIC_AUTH_PASSWORD")}}
}

func NewMailClient() *client.MailClient {
	host := os.Getenv("MAIL_HOST")
	port := os.Getenv("MAIL_PORT")
	user := os.Getenv("MAIL_USER")
	pass := os.Getenv("MAIL_PASS")
	return &client.MailClient{
		SMTPHost: host,
		SMTPPort: port,
		SMTPUser: user,
		SMTPPass: pass,
	}
}
