package main

import (
	"bytes"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRunCmd(t *testing.T) {
	var gotReturnCode int
	var e error
	os.Setenv("HELLO", "SHOULD_REPLACE")
	os.Setenv("FOO", "SHOULD_REPLACE")
	os.Setenv("UNSET", "SHOULD_REMOVE")
	os.Setenv("ADDED", "from original env")
	os.Setenv("EMPTY", "SHOULD_BE_EMPTY")

	type args struct {
		cmd []string
		env Environment
	}
	tests := []struct {
		name           string
		args           args
		wantReturnCode int
		expected       string
	}{
		{
			name: "original variables",
			args: args{
				cmd: []string{"/bin/bash", "testdata/echo.sh", "arg1=1", "arg2=2"},
				env: Environment{},
			},
			wantReturnCode: 0,
			expected: `HELLO is (SHOULD_REPLACE)
BAR is ()
FOO is (SHOULD_REPLACE)
UNSET is (SHOULD_REMOVE)
ADDED is (from original env)
EMPTY is (SHOULD_BE_EMPTY)
arguments are arg1=1 arg2=2`,
		},
		{
			name: "Test echo.sh",
			args: args{
				cmd: []string{"/bin/bash", "testdata/echo.sh", "arg1=1", "arg2=2"},
				env: Environment{
					"HELLO": EnvValue{Value: "\"hello\"", NeedRemove: false},
					"BAR":   EnvValue{Value: "bar", NeedRemove: false},
					"FOO":   EnvValue{Value: "   foo\nwith new line", NeedRemove: false},
					"EMPTY": EnvValue{Value: "", NeedRemove: false},
					"UNSET": EnvValue{Value: "", NeedRemove: true},
				},
			},
			wantReturnCode: 0,
			expected: `HELLO is ("hello")
BAR is (bar)
FOO is (   foo
with new line)
UNSET is ()
ADDED is (from original env)
EMPTY is ()
arguments are arg1=1 arg2=2`,
		},
		{
			name: "cmd not found",
			args: args{
				cmd: []string{},
				env: Environment{
					"HELLO": EnvValue{Value: "\"hello\"", NeedRemove: false},
					"BAR":   EnvValue{Value: "bar", NeedRemove: false},
					"FOO":   EnvValue{Value: "   foo\nwith new line", NeedRemove: false},
					"EMPTY": EnvValue{Value: "", NeedRemove: false},
					"UNSET": EnvValue{Value: "", NeedRemove: true},
				},
			},
			wantReturnCode: 1,
			expected:       "",
		},
		{
			name: "empty env map",
			args: args{
				cmd: []string{},
				env: Environment{},
			},
			wantReturnCode: 1,
			expected:       "",
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			var r *os.File
			func() {
				var w *os.File
				r, w, e = os.Pipe()
				defer w.Close()
				origStdout := os.Stdout
				defer func() { os.Stdout = origStdout }()
				os.Stdout = w
				require.NoError(t, e)
				gotReturnCode = RunCmd(tt.args.cmd, tt.args.env)
			}()
			_, e = buf.ReadFrom(r)
			// fmt.Print("Im from test ", buf.String())
			require.NoError(t, e)
			// gotReturnCode := RunCmd(tt.args.cmd, tt.args.env)
			require.Equalf(t, tt.wantReturnCode, gotReturnCode, "ReturnCode should be %v", tt.wantReturnCode)
			require.Equal(t, tt.expected, strings.TrimRight(buf.String(), "\n"), "must be equal")
		})
	}
}
