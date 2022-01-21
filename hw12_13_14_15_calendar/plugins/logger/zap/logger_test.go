package zap

import (
	"bufio"
	"os"
	"path"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/xkdgo/otus_golang_homeworks/hw12_13_14_15_calendar/internal/logger"
)

func TestLogger(t *testing.T) {
	tmpDir, err := os.MkdirTemp(".", "tmp_")
	require.NoErrorf(t, err, "Couldn't create tempdir")
	toPath := path.Join(tmpDir, "loger_result.txt")
	defer os.RemoveAll(tmpDir)

	pluginlogger, err := NewLogger(WithFile(toPath))
	require.NoErrorf(t, err, "Cant initialize zap logger")
	t.Run("Test DEBUG level", func(t *testing.T) {
		logg := logger.New("DEBUG", pluginlogger)

		fd, err := os.Open(toPath)
		require.NoError(t, err)
		scanner := bufio.NewScanner(fd)

		logg.Debug("Debug message")
		scanner.Scan()
		require.True(t, strings.Contains(scanner.Text(), "Debug message"))

		logg.Debugf("Debug message %s", "test")
		scanner.Scan()
		require.True(t, strings.Contains(scanner.Text(), "Debug message test"))

		logg.Info("Info message")
		scanner.Scan()
		require.True(t, strings.Contains(scanner.Text(), "Info message"))

		logg.Infof("Info message %s", "test")
		scanner.Scan()
		require.True(t, strings.Contains(scanner.Text(), "Info message test"))

		logg.Error("Error message")
		scanner.Scan()
		require.True(t, strings.Contains(scanner.Text(), "Error message"))

		logg.Errorf("Error message %s", "test")
		scanner.Scan()
		require.True(t, strings.Contains(scanner.Text(), "Error message test"))
	})

	t.Run("Test ERROR level", func(t *testing.T) {
		logg := logger.New("ERROR", pluginlogger)

		fd, err := os.Open(toPath)
		require.NoError(t, err)
		scanner := bufio.NewScanner(fd)

		logg.Debug("Debug message")
		scanner.Scan()
		require.True(t, strings.Contains(scanner.Text(), ""))

		logg.Debugf("Debug message %s", "test")
		scanner.Scan()
		require.True(t, strings.Contains(scanner.Text(), ""))

		logg.Info("Info message")
		scanner.Scan()
		require.True(t, strings.Contains(scanner.Text(), ""))

		logg.Infof("Info message %s", "test")
		scanner.Scan()
		require.True(t, strings.Contains(scanner.Text(), ""))

		logg.Error("Error message")
		scanner.Scan()
		require.True(t, strings.Contains(scanner.Text(), "Error message"))

		logg.Errorf("Error message %s", "test")
		scanner.Scan()
		require.True(t, strings.Contains(scanner.Text(), "Error message test"))
	})

	t.Run("Test INFO level", func(t *testing.T) {
		logg := logger.New("INFO", pluginlogger)

		fd, err := os.Open(toPath)
		require.NoError(t, err)
		scanner := bufio.NewScanner(fd)

		logg.Debug("Debug message")
		scanner.Scan()
		require.True(t, strings.Contains(scanner.Text(), ""))

		logg.Debugf("Debug message %s", "test")
		scanner.Scan()
		require.True(t, strings.Contains(scanner.Text(), ""))

		logg.Info("Info message")
		scanner.Scan()
		require.True(t, strings.Contains(scanner.Text(), "Info message"))

		logg.Infof("Info message %s", "test")
		scanner.Scan()
		require.True(t, strings.Contains(scanner.Text(), "Info message test"))

		logg.Error("Error message")
		scanner.Scan()
		require.True(t, strings.Contains(scanner.Text(), "Error message"))

		logg.Errorf("Error message %s", "test")
		scanner.Scan()
		require.True(t, strings.Contains(scanner.Text(), "Error message test"))
	})

}
