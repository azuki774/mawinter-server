package model

import (
	"reflect"
	"testing"
	"time"
)

func TestNewRecordFromReq(t *testing.T) {
	type args struct {
		req RecordRequest
	}
	tests := []struct {
		name       string
		args       args
		wantRecord Recordstruct
		wantErr    bool
	}{
		{
			name: "fix time (YYYYMMDD)",
			args: args{
				req: RecordRequest{
					CategoryID: 100,
					Datetime:   "20220301",
					From:       "fromtest",
					Type:       "typetest",
					Price:      2345,
					Memo:       "memo",
				},
			},
			wantRecord: Recordstruct{
				CategoryID: 100,
				Datetime:   time.Date(2022, 03, 01, 0, 0, 0, 0, jst),
				From:       "fromtest",
				Type:       "typetest",
				Price:      2345,
				Memo:       "memo",
			},
			wantErr: false,
		},
		{
			name: "fix time (2022-11-22T13:54:08+09:00)",
			args: args{
				req: RecordRequest{
					CategoryID: 100,
					Datetime:   "2022-11-22T13:54:08+09:00",
					From:       "fromtest",
					Type:       "typetest",
					Price:      2345,
					Memo:       "memo",
				},
			},
			wantRecord: Recordstruct{
				CategoryID: 100,
				Datetime:   time.Date(2022, 11, 22, 13, 54, 8, 0, jst),
				From:       "fromtest",
				Type:       "typetest",
				Price:      2345,
				Memo:       "memo",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotRecord, err := NewRecordFromReq(tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewRecordFromReq() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotRecord, tt.wantRecord) {
				t.Errorf("NewRecordFromReq() = %v, want %v", gotRecord, tt.wantRecord)
			}
		})
	}
}
