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
		wantRec model.Record_YYYYMM
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
			wantRec: model.Record_YYYYMM{
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
			wantRec: model.Record_YYYYMM{
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
