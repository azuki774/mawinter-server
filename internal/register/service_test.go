package register

import (
	"context"
	"fmt"
	"testing"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var l *zap.Logger

func init() {
	config := zap.NewProductionConfig()
	// config.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
	config.EncoderConfig.EncodeTime = JSTTimeEncoder
	l, _ = config.Build()

	l.WithOptions(zap.AddStacktrace(zap.ErrorLevel))
}

func JSTTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	const layout = "2006-01-02T15:04:05+09:00"
	jst := time.FixedZone("Asia/Tokyo", 9*60*60)
	enc.AppendString(t.In(jst).Format(layout))
}

func TestAPIService_InsertMonthlyFixBilling(t *testing.T) {
	type fields struct {
		Logger     *zap.Logger
		DB         DBRepository
		MailClient MailClient
	}
	type args struct {
		ctx    context.Context
		yyyymm string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:    "ok",
			fields:  fields{Logger: l, DB: &mockRepo{}, MailClient: &mockMailClient{}},
			args:    args{ctx: context.Background(), yyyymm: "202201"},
			wantErr: false,
		},
		{
			name:    "error",
			fields:  fields{Logger: l, DB: &mockRepo{errGetMonthly: fmt.Errorf("error")}, MailClient: &mockMailClient{}},
			args:    args{ctx: context.Background(), yyyymm: "202201"},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &RegisterService{
				Logger: tt.fields.Logger,
				DB:     tt.fields.DB,
			}
			if err := a.InsertMonthlyFixBilling(tt.args.ctx, tt.args.yyyymm); (err != nil) != tt.wantErr {
				t.Errorf("APIService.InsertMonthlyFixBilling() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
