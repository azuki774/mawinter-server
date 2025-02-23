package timeutil

import "testing"

func TestConvertToFiscalYear(t *testing.T) {
	type args struct {
		yyyymm string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "202504",
			args: args{
				yyyymm: "202504",
			},
			want: "2025",
		},
		{
			name: "202503",
			args: args{
				yyyymm: "202503",
			},
			want: "2024",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ConvertToFiscalYear(tt.args.yyyymm); got != tt.want {
				t.Errorf("ConvertToFiscalYear() = %v, want %v", got, tt.want)
			}
		})
	}
}
