package internalhttp

import (
	"reflect"
	"testing"
)

func TestNewRootHandler(t *testing.T) {
	type args struct {
		app    Application
		logger Logger
	}
	tests := []struct {
		name string
		args args
		want *RootHandler
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewRootHandler(tt.args.app, tt.args.logger); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewRootHandler() = %v, want %v", got, tt.want)
			}
		})
	}
}
