package api

import "testing"

func Test_isValidYearMonth(t *testing.T) {
	type args struct {
		yyyymm string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "202101",
			args:    args{"202101"},
			wantErr: false,
		},
		{
			name:    "2021012",
			args:    args{"2021012"},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := isValidYearMonth(tt.args.yyyymm); (err != nil) != tt.wantErr {
				t.Errorf("isValidYearMonth() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
