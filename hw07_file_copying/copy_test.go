package main

import (
	"bytes"
	"errors"
	"os"
	"path"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCopyNoError(t *testing.T) {
	tmpDir, err := os.MkdirTemp("testdata", "tmp_")
	require.NoErrorf(t, err, "Couldn't create tempdir")
	defer os.RemoveAll(tmpDir)

	type args struct {
		fromPath      string
		toPath        string
		offset        int64
		limit         int64
		compareSample string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "offset0_limit0",
			args: args{
				fromPath:      path.Join("testdata", "input.txt"),
				toPath:        path.Join(tmpDir, "result.txt"),
				offset:        0,
				limit:         0,
				compareSample: path.Join("testdata", "out_offset0_limit0.txt"),
			},
			wantErr: false,
		},
		{
			name: "offset0_limit10",
			args: args{
				fromPath:      path.Join("testdata", "input.txt"),
				toPath:        path.Join(tmpDir, "result.txt"),
				offset:        0,
				limit:         10,
				compareSample: path.Join("testdata", "out_offset0_limit10.txt"),
			},
			wantErr: false,
		},
		{
			name: "offset0_limit1000",
			args: args{
				fromPath:      path.Join("testdata", "input.txt"),
				toPath:        path.Join(tmpDir, "result.txt"),
				offset:        0,
				limit:         1000,
				compareSample: path.Join("testdata", "out_offset0_limit1000.txt"),
			},
			wantErr: false,
		},
		{
			name: "offset100_limit1000",
			args: args{
				fromPath:      path.Join("testdata", "input.txt"),
				toPath:        path.Join(tmpDir, "result.txt"),
				offset:        100,
				limit:         1000,
				compareSample: path.Join("testdata", "out_offset100_limit1000.txt"),
			},
			wantErr: false,
		},
		{
			name: "offset6000_limit1000",
			args: args{
				fromPath:      path.Join("testdata", "input.txt"),
				toPath:        path.Join(tmpDir, "result.txt"),
				offset:        6000,
				limit:         1000,
				compareSample: path.Join("testdata", "out_offset6000_limit1000.txt"),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.False(t, tt.wantErr, "This is normal operation tests")
			compareSampleBytes, err := os.ReadFile(tt.args.compareSample)
			require.NoErrorf(t, err, "Couldn't read sample file %v", tt.args.compareSample)
			err = Copy(tt.args.fromPath, tt.args.toPath, tt.args.offset, tt.args.limit)
			require.NoErrorf(t, err, "Should be no error")
			compareResultBytes, err := os.ReadFile(tt.args.toPath)
			require.NoError(t, err, "Couldn`t read copied file %v", tt.args.toPath)
			require.Truef(t,
				bytes.Equal(compareSampleBytes, compareResultBytes),
				"Copied file content and sample content should be equal %v and %v",
				tt.args.compareSample, tt.args.toPath)
		})
	}
}

func TestCopyProgressBar(t *testing.T) {
	tmpDir, err := os.MkdirTemp("testdata", "tmp_")
	require.NoErrorf(t, err, "Couldn't create tempdir")
	defer os.RemoveAll(tmpDir)

	type args struct {
		fromPath          string
		toPath            string
		offset            int64
		limit             int64
		progressBarSample string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "offset100_limit1000",
			args: args{
				fromPath:          path.Join("testdata", "input.txt"),
				toPath:            path.Join(tmpDir, "result.txt"),
				offset:            100,
				limit:             1000,
				progressBarSample: "================] 100.00%",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			var r *os.File
			var e error
			func() {
				var w *os.File
				r, w, e = os.Pipe()
				defer w.Close()
				origStdout := os.Stdout
				defer func() { os.Stdout = origStdout }()
				os.Stdout = w
				require.NoError(t, e)
				err = Copy(tt.args.fromPath, tt.args.toPath, tt.args.offset, tt.args.limit)
			}()
			_, e = buf.ReadFrom(r)
			stdOutResult := buf.String()
			require.Truef(t,
				strings.Contains(stdOutResult, tt.args.progressBarSample),
				"Result:\n %v \nBut should contains:\n%v", stdOutResult, tt.args.progressBarSample)
		})
	}
}

func TestCopyWithErrors(t *testing.T) {
	tmpDir, err := os.MkdirTemp("testdata", "tmp_")
	require.NoErrorf(t, err, "Couldn't create tempdir")
	defer os.RemoveAll(tmpDir)

	type args struct {
		fromPath  string
		toPath    string
		offset    int64
		limit     int64
		errorType error
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Need error about offset",
			args: args{
				fromPath:  path.Join("testdata", "out_offset0_limit10.txt"),
				toPath:    path.Join(tmpDir, "offset101_limit0"),
				offset:    101,
				limit:     0,
				errorType: ErrOffsetExceedsFileSize,
			},
			wantErr: true,
		},
		{
			name: "Need error unsupported file",
			args: args{
				fromPath:  "/dev/urandom",
				toPath:    path.Join(tmpDir, "offset0_limit100"),
				offset:    0,
				limit:     100,
				errorType: ErrUnsupportedFile,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			require.True(t, tt.wantErr)
			err = Copy(tt.args.fromPath, tt.args.toPath, tt.args.offset, tt.args.limit)
			require.Error(t, err)
			require.Truef(t, errors.Is(err, tt.args.errorType), "actual error string %q\nerror type %T", err, err)
		})
	}
}
