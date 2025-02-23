package util

import (
	"reflect"
	"testing"
)

func TestInsertStringIfNotExists(t *testing.T) {
	type args struct {
		slice *[]string
		value string
	}
	tests := []struct {
		name string
		args args
		want []string // add
	}{
		{
			name: "Insert new value",
			args: args{
				slice: &[]string{"a", "b", "c"},
				value: "d",
			},
			want: []string{"a", "b", "c", "d"},
		},
		{
			name: "Insert existing value",
			args: args{
				slice: &[]string{"a", "b", "c"},
				value: "b",
			},
			want: []string{"a", "b", "c"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			InsertStringIfNotExists(tt.args.slice, tt.args.value)
			// added below
			if !reflect.DeepEqual(*tt.args.slice, tt.want) {
				t.Errorf("InsertStringIfNotExists() = %v, want %v", *tt.args.slice, tt.want)
			}
		})
	}
}
