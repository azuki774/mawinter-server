package repository

import (
	"mawinter-server/internal/model"
	"mawinter-server/internal/openapi"
	"mawinter-server/internal/timeutil"
	"reflect"
	"testing"
	"time"
)

func strPtr(str string) *string {
	return &str
}

func TestNewDBModelRecord(t *testing.T) {
	type args struct {
		req openapi.ReqRecord
	}
	tests := []struct {
		name    string
		args    args
		wantRec model.Record
		wantErr bool
		nowTime time.Time
	}{
		{
			name: "minimum",
			args: args{
				req: openapi.ReqRecord{
					CategoryId: 100,
					Price:      10000,
				},
			},
			wantRec: model.Record{
				CategoryID: 100,
				Datetime:   time.Date(2010, 12, 1, 0, 0, 0, 0, jst),
				Price:      10000,
			},
			wantErr: false,
			nowTime: time.Date(2010, 12, 1, 0, 0, 0, 0, jst),
		},
		{
			name: "full",
			args: args{
				req: openapi.ReqRecord{
					CategoryId: 200,
					Datetime:   strPtr("20111201"),
					From:       strPtr("from"),
					Memo:       strPtr("memo"),
					Price:      20000,
					Type:       strPtr("type"),
				},
			},
			wantRec: model.Record{
				CategoryID: 200,
				Datetime:   time.Date(2011, 12, 1, 0, 0, 0, 0, jst),
				From:       "from",
				Memo:       "memo",
				Price:      20000,
				Type:       "type",
			},
			wantErr: false,
			nowTime: time.Date(2011, 12, 1, 0, 0, 0, 0, jst),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// set test time
			timeutil.NowFunc = func() time.Time {
				return tt.nowTime
			}

			gotRec, err := NewDBModelRecord(tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewDBModelRecord() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotRec, tt.wantRec) {
				t.Errorf("NewDBModelRecord() = %v, want %v", gotRec, tt.wantRec)
			}
		})
	}
}

func Test_yyyymmToInitDayTime(t *testing.T) {
	type args struct {
		yyyymm string
	}
	tests := []struct {
		name    string
		args    args
		wantT   time.Time
		wantErr bool
	}{
		{
			name: "200102",
			args: args{
				yyyymm: "200102",
			},
			wantT:   time.Date(2001, 2, 1, 0, 0, 0, 0, jst),
			wantErr: false,
		},
		{
			name: "202511",
			args: args{
				yyyymm: "202511",
			},
			wantT:   time.Date(2025, 11, 1, 0, 0, 0, 0, jst),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotT, err := yyyymmToInitDayTime(tt.args.yyyymm)
			if (err != nil) != tt.wantErr {
				t.Errorf("yyyymmToInitDayTime() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotT, tt.wantT) {
				t.Errorf("yyyymmToInitDayTime() = %v, want %v", gotT, tt.wantT)
			}
		})
	}
}
