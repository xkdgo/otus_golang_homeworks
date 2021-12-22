package main

import (
	"bytes"
	"io/ioutil"
	"net"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestTelnetClient(t *testing.T) {
	t.Run("basic", func(t *testing.T) {
		l, err := net.Listen("tcp", "127.0.0.1:")
		require.NoError(t, err)
		defer func() { require.NoError(t, l.Close()) }()

		var wg sync.WaitGroup
		wg.Add(2)

		go func() {
			defer wg.Done()

			in := &bytes.Buffer{}
			out := &bytes.Buffer{}

			timeout, err := time.ParseDuration("10s")
			require.NoError(t, err)

			client := NewTelnetClient(l.Addr().String(), timeout, ioutil.NopCloser(in), out)
			require.NoError(t, client.Connect())
			defer func() { require.NoError(t, client.Close()) }()

			in.WriteString("hello\n")
			err = client.Send()
			require.NoError(t, err)

			err = client.Receive()
			require.NoError(t, err)
			require.Equal(t, "world\n", out.String())
		}()

		go func() {
			defer wg.Done()

			conn, err := l.Accept()
			require.NoError(t, err)
			require.NotNil(t, conn)
			defer func() { require.NoError(t, conn.Close()) }()

			request := make([]byte, 1024)
			n, err := conn.Read(request)
			require.NoError(t, err)
			require.Equal(t, "hello\n", string(request)[:n])

			n, err = conn.Write([]byte("world\n"))
			require.NoError(t, err)
			require.NotEqual(t, 0, n)
		}()

		wg.Wait()
	})
	t.Run("test unexistent host", func(t *testing.T) {
		netOpErr := &net.OpError{}
		timeout, err := time.ParseDuration("10s")
		require.NoError(t, err)
		in := &bytes.Buffer{}
		out := &bytes.Buffer{}
		client := NewTelnetClient("localhost:4243", timeout, ioutil.NopCloser(in), out)
		err = client.Connect()
		require.ErrorAs(t, err, &netOpErr)
		require.True(t, strings.Contains(netOpErr.Err.Error(), "connect: connection refused"))
	})
	t.Run("test dial timeout", func(t *testing.T) {
		timeout, err := time.ParseDuration("5s")
		netOpErr := &net.OpError{}
		require.NoError(t, err)
		in := &bytes.Buffer{}
		out := &bytes.Buffer{}
		client := NewTelnetClient("1.1.1.1:8080", timeout, ioutil.NopCloser(in), out)
		t1 := time.Now()
		err = client.Connect()
		require.InDelta(t, timeout, time.Since(t1), float64(20*time.Millisecond))
		require.ErrorAs(t, err, &netOpErr)
		require.True(t, strings.Contains(netOpErr.Err.Error(), "i/o timeout"))
	})
}
