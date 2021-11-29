package main

import (
	"reflect"
	"testing"
)

func TestReadDir(t *testing.T) {

	type args struct {
		dir string
	}
	tests := []struct {
		name    string
		args    args
		want    Environment
		wantErr bool
	}{
		{name: "test directory env",
			args: args{dir: "testdata/env"},
			want: Environment{
				"HELLO": EnvValue{Value: "\"hello\"", NeedRemove: false},
				"BAR":   EnvValue{Value: "bar", NeedRemove: false},
				"FOO":   EnvValue{Value: "   foo\nwith new line", NeedRemove: false},
				"EMPTY": EnvValue{Value: "", NeedRemove: false},
				"UNSET": EnvValue{Value: "", NeedRemove: true},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, err := ReadDir(tt.args.dir)
			if (err != nil) != tt.wantErr {
				t.Errorf("ReadDir() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ReadDir() = %v, want %v", got, tt.want)
			}
		})
	}
}
