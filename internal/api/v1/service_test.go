package v1

import (
	"context"
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

func TestAPIService_GetYearCategorySummary(t *testing.T) {
	type fields struct {
		Logger *zap.Logger
		Repo   DBRepository
	}
	type args struct {
		categoryID int
		yyyy       string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantSum *model.CategoryYearSummary
		wantErr bool
	}{
		{
			name: "2022-101",
			fields: fields{
				Logger: l,
				Repo: &mockRepo{
					CategoryMonthSummaryMonthTotal: []int{40, 50, 60, 70, 80, 90, 100, 110, 120, 10, 20, 30},
					CategoryMonthSummaryMonthCount: []int{4, 5, 6, 7, 8, 9, 10, 11, 12, 1, 2, 3},
				},
			},
			args: args{categoryID: 101, yyyy: "2022"},
			wantSum: &model.CategoryYearSummary{
				CategoryID:   101,
				CategoryName: "カテゴリ1",
				MonthPrice:   []int{40, 50, 60, 70, 80, 90, 100, 110, 120, 10, 20, 30},
				Total:        780,
			},
			wantErr: false,
		},
		{
			name: "2022-102",
			fields: fields{
				Logger: l,
				Repo: &mockRepo{
					CategoryMonthSummaryMonthTotal: []int{40, 50, 60, 70, 80, 90, 0, 110, 120, 10, 20, 30}, // 途中 0
					CategoryMonthSummaryMonthCount: []int{4, 5, 6, 7, 8, 9, 0, 11, 12, 1, 2, 3},
				},
			},
			args: args{categoryID: 102, yyyy: "2022"},
			wantSum: &model.CategoryYearSummary{
				CategoryID:   102,
				CategoryName: "カテゴリ1",
				MonthPrice:   []int{40, 50, 60, 70, 80, 90, 0, 110, 120, 10, 20, 30},
				Total:        680,
			},
			wantErr: false,
		},
		{
			name: "2022-103",
			fields: fields{
				Logger: l,
				Repo: &mockRepo{
					CategoryMonthSummaryMonthTotal: []int{40, 50, 60, 70, 80, 90}, // 途中からテーブルなし
					CategoryMonthSummaryMonthCount: []int{4, 5, 6, 7, 8, 9},
				},
			},
			args: args{categoryID: 103, yyyy: "2022"},
			wantSum: &model.CategoryYearSummary{
				CategoryID:   103,
				CategoryName: "カテゴリ1",
				MonthPrice:   []int{40, 50, 60, 70, 80, 90, 0, 0, 0, 0, 0, 0},
				Total:        390,
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
			gotSum, err := a.GetYearCategorySummary(context.Background(), tt.args.categoryID, tt.args.yyyy)
			if (err != nil) != tt.wantErr {
				t.Errorf("APIService.GetYearCategorySummary() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotSum, tt.wantSum) {
				t.Errorf("APIService.GetYearCategorySummary() = %v, want %v", gotSum, tt.wantSum)
			}
		})
	}
}
