package api

import (
	"context"
	"mawinter-server/internal/openapi"
	"reflect"
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

func strPtr(str string) *string {
	return &str
}

func JSTTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	const layout = "2006-01-02T15:04:05+09:00"
	jst := time.FixedZone("Asia/Tokyo", 9*60*60)
	enc.AppendString(t.In(jst).Format(layout))
}

func TestAPIService_PostRecord(t *testing.T) {
	type fields struct {
		Logger *zap.Logger
		Repo   DBRepository
	}
	type args struct {
		ctx context.Context
		req openapi.ReqRecord
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantRec openapi.Record
		wantErr bool
	}{
		{
			name: "ok",
			fields: fields{
				Logger: l,
				Repo:   &mockRepo{},
			},
			args: args{
				ctx: context.Background(),
				req: openapi.ReqRecord{
					CategoryId: 100,
					Datetime:   strPtr("20000123"),
					From:       strPtr("from"),
					Memo:       strPtr("memo"),
					Price:      1234,
					Type:       strPtr("type"),
				},
			},
			wantRec: openapi.Record{
				CategoryId:   100,
				CategoryName: "cat1",
				Datetime:     time.Date(2000, 1, 23, 0, 0, 0, 0, jst),
				From:         "from",
				Id:           1,
				Memo:         "memo",
				Price:        1234,
				Type:         "type",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &APIService{
				Logger: tt.fields.Logger,
				Repo:   tt.fields.Repo,
			}
			gotRec, err := a.PostRecord(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("APIService.PostRecord() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotRec, tt.wantRec) {
				t.Errorf("APIService.PostRecord() = %v, want %v", gotRec, tt.wantRec)
			}
		})
	}
}

func TestAPIService_GetYYYYMMRecords(t *testing.T) {
	type fields struct {
		Logger *zap.Logger
		Repo   DBRepository
	}
	type args struct {
		ctx    context.Context
		yyyymm string
	}
	tests := []struct {
		name     string
		fields   fields
		args     args
		wantRecs []openapi.Record
		wantErr  bool
	}{
		{
			name: "ok",
			fields: fields{
				Logger: l,
				Repo:   &mockRepo{},
			},
			args: args{
				ctx:    context.Background(),
				yyyymm: "200001",
			},
			wantRecs: []openapi.Record{
				{
					CategoryId:   100,
					CategoryName: "cat1",
					Datetime:     time.Date(2000, 1, 23, 0, 0, 0, 0, jst),
					From:         "from",
					Id:           1,
					Memo:         "memo",
					Price:        1234,
					Type:         "type",
				},
				{
					CategoryId:   200,
					CategoryName: "cat2",
					Datetime:     time.Date(2000, 1, 25, 0, 0, 0, 0, jst),
					From:         "",
					Id:           2,
					Memo:         "",
					Price:        2345,
					Type:         "",
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &APIService{
				Logger: tt.fields.Logger,
				Repo:   tt.fields.Repo,
			}
			gotRecs, err := a.GetYYYYMMRecords(tt.args.ctx, tt.args.yyyymm)
			if (err != nil) != tt.wantErr {
				t.Errorf("APIService.GetYYYYMMRecords() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotRecs, tt.wantRecs) {
				t.Errorf("APIService.GetYYYYMMRecords() = %v, want %v", gotRecs, tt.wantRecs)
			}
		})
	}
}

func TestAPIService_GetV2YearSummary(t *testing.T) {
	type fields struct {
		Logger *zap.Logger
		Repo   DBRepository
	}
	type args struct {
		ctx  context.Context
		year int
	}
	tests := []struct {
		name     string
		fields   fields
		args     args
		wantSums []openapi.CategoryYearSummary
		wantErr  bool
	}{
		{
			name: "ok",
			fields: fields{
				Logger: l,
				Repo:   &mockRepo{},
			},
			args: args{
				ctx:  context.Background(),
				year: 2000,
			},
			wantSums: []openapi.CategoryYearSummary{
				{
					CategoryId:   100,
					CategoryName: "cat1",
					Count:        120,
					Price:        []int{1000, 1000, 1000, 1000, 1000, 1000, 1000, 1000, 1000, 1000, 1000, 1000},
					Total:        12000,
				},
				{
					CategoryId:   200,
					CategoryName: "cat2",
					Count:        220, // 20 * 11
					Price:        []int{2000, 2000, 2000, 2000, 0, 2000, 2000, 2000, 2000, 2000, 2000, 2000},
					Total:        22000,
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &APIService{
				Logger: tt.fields.Logger,
				Repo:   tt.fields.Repo,
			}
			gotSums, err := a.GetV2YearSummary(tt.args.ctx, tt.args.year)
			if (err != nil) != tt.wantErr {
				t.Errorf("APIService.GetV2YearSummary() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotSums, tt.wantSums) {
				t.Errorf("APIService.GetV2YearSummary() = %v, want %v", gotSums, tt.wantSums)
			}
		})
	}
}
