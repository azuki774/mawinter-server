package service

import (
	"context"
	"mawinter-server/internal/openapi"
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

func Test_judgeDuplicateRecords(t *testing.T) {
	type args struct {
		d1 openapi.Record
		d2 openapi.Record
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "duplicate",
			args: args{
				d1: openapi.Record{
					CategoryId: 100,
					Datetime:   time.Date(2000, 1, 23, 12, 0, 0, 0, jst),
					From:       "from1",
					Price:      1234,
				},
				d2: openapi.Record{
					CategoryId: 100,
					Datetime:   time.Date(2000, 1, 23, 0, 0, 0, 0, jst),
					From:       "from2",
					Price:      1234,
				},
			},
			want: true,
		},
		{
			name: "no duplicate",
			args: args{
				d1: openapi.Record{
					CategoryId: 100,
					Datetime:   time.Date(2000, 1, 24, 12, 0, 0, 0, jst),
					From:       "from1",
					Price:      1234,
				},
				d2: openapi.Record{
					CategoryId: 100,
					Datetime:   time.Date(2000, 1, 23, 0, 0, 0, 0, jst),
					From:       "from1",
					Price:      1234,
				},
			},
			want: false,
		},
		{
			name: "same",
			args: args{
				d1: openapi.Record{
					CategoryId: 100,
					Datetime:   time.Date(2000, 1, 23, 12, 0, 0, 0, jst),
					From:       "from1",
					Price:      1234,
				},
				d2: openapi.Record{
					CategoryId: 100,
					Datetime:   time.Date(2000, 1, 23, 12, 0, 0, 0, jst),
					From:       "from1",
					Price:      1234,
				},
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := judgeDuplicateRecords(tt.args.d1, tt.args.d2); got != tt.want {
				t.Errorf("judgeDuplicateRecords() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDuplicateCheckService_DuplicateCheck(t *testing.T) {
	type fields struct {
		Logger *zap.Logger
		Ap     APIServiceDup
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
			name: "duplicate_exist",
			fields: fields{
				Logger: l,
				Ap:     &mockAp{},
			},
			args:    args{},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &DuplicateCheckService{
				Logger: tt.fields.Logger,
				Ap:     tt.fields.Ap,
			}
			if err := d.DuplicateCheck(tt.args.ctx, tt.args.yyyymm); (err != nil) != tt.wantErr {
				t.Errorf("DuplicateCheckService.DuplicateCheck() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
