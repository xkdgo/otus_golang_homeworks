package main

import (
	"reflect"
	"testing"
)

// func TestReadDir(t *testing.T) {
// 	// Place your code here
// }

func TestReadDir(t *testing.T) {
	Env := make(map[string]EnvValue, 0)
	Env["HELLO"] = EnvValue{Value: "\"hello\"", NeedRemove: false}
	Env["BAR"] = EnvValue{Value: "bar", NeedRemove: false}
	Env["FOO"] = EnvValue{Value: "   foo\nwith new line", NeedRemove: false}
	Env["EMPTY"] = EnvValue{Value: "", NeedRemove: false}
	Env["UNSET"] = EnvValue{Value: "", NeedRemove: true}

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
			args:    args{dir: "testdata/env"},
			want:    Env,
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
