package v1

import (
	"context"
	"fmt"
	"mawinter-server/internal/model"
	"reflect"
	"testing"

	"go.uber.org/zap"
)

var l *zap.Logger

func init() {
	l, _ = zap.NewProduction()
}
func Test_fyInterval(t *testing.T) {
	type args struct {
		yyyy int
	}
	tests := []struct {
		name       string
		args       args
		wantYyyymm []string
	}{
		{
			name: "2022",
			args: args{yyyy: 2022},
			wantYyyymm: []string{
				"202204",
				"202205",
				"202206",
				"202207",
				"202208",
				"202209",
				"202210",
				"202211",
				"202212",
				"202301",
				"202302",
				"202303",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotYyyymm := fyInterval(tt.args.yyyy); !reflect.DeepEqual(gotYyyymm, tt.wantYyyymm) {
				t.Errorf("fyInterval() = %v, want %v", gotYyyymm, tt.wantYyyymm)
			}
		})
	}
}

func TestAPIService_GetYearSummary(t *testing.T) {
	type fields struct {
		Logger *zap.Logger
		Repo   DBRepository
	}
	type args struct {
		ctx  context.Context
		yyyy string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantSum []*model.CategoryYearSummary
		wantErr bool
	}{
		{
			name: "full",
			fields: fields{
				Logger: l,
				Repo: &mockRepo{
					RecordYYYYMMNum: 12,
				},
			},
			args: args{
				ctx:  context.Background(),
				yyyy: "2021",
			},
			wantSum: []*model.CategoryYearSummary{
				{
					CategoryID:   100,
					CategoryName: "カテゴリ1",
					MonthPrice:   []int{10, 10, 10, 10, 10, 10, 10, 10, 10, 10, 10, 10},
					Total:        120,
				},
				{
					CategoryID:   200,
					CategoryName: "カテゴリ2",
					MonthPrice:   []int{100, 100, 100, 100, 100, 100, 100, 100, 100, 100, 100, 100},
					Total:        1200,
				},
				{
					CategoryID:   300,
					CategoryName: "カテゴリ3",
					MonthPrice:   []int{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
					Total:        0,
				},
			},
			wantErr: false,
		},
		{
			name: "half",
			fields: fields{
				Logger: l,
				Repo: &mockRepo{
					RecordYYYYMMNum: 6,
				},
			},
			args: args{
				ctx:  context.Background(),
				yyyy: "2022",
			},
			wantSum: []*model.CategoryYearSummary{
				{
					CategoryID:   100,
					CategoryName: "カテゴリ1",
					MonthPrice:   []int{10, 10, 10, 10, 10, 10, 0, 0, 0, 0, 0, 0},
					Total:        60,
				},
				{
					CategoryID:   200,
					CategoryName: "カテゴリ2",
					MonthPrice:   []int{100, 100, 100, 100, 100, 100, 0, 0, 0, 0, 0, 0},
					Total:        600,
				},
				{
					CategoryID:   300,
					CategoryName: "カテゴリ3",
					MonthPrice:   []int{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
					Total:        0,
				},
			},
			wantErr: false,
		},
		{
			name: "get category error",
			fields: fields{
				Logger: l,
				Repo: &mockRepo{
					errGetCategoryInfo: fmt.Errorf("error"),
				},
			},
			args: args{
				ctx:  context.Background(),
				yyyy: "2021",
			},
			wantSum: nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &APIService{
				Logger: tt.fields.Logger,
				Repo:   tt.fields.Repo,
			}
			gotSum, err := a.GetYearSummary(tt.args.ctx, tt.args.yyyy)
			if (err != nil) != tt.wantErr {
				t.Errorf("APIService.GetYearSummary() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotSum, tt.wantSum) {
				t.Errorf("APIService.GetYearSummary() = %v, want %v", gotSum, tt.wantSum)
			}
		})
	}
}
