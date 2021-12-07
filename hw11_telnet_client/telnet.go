package main

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"time"
)

var (
	errConnectionClosedbyPeer = errors.New("...Connection was closed by peer")
	errEOF                    = errors.New("...EOF")
)

type TelnetClient interface {
	Connect() error
	io.Closer
	Send() error
	Receive() error
}

func NewTelnetClient(address string, timeout time.Duration, in io.ReadCloser, out io.Writer) TelnetClient {
	ctx := context.Background()
	ctx, cancelTimeout := context.WithTimeout(ctx, timeout)
	return &TelnetConnection{
		Address:       address,
		Ctx:           ctx,
		cancelTimeout: cancelTimeout,
		In:            in,
		Out:           out,
	}
}

type TelnetConnection struct {
	Address       string
	Ctx           context.Context
	cancelTimeout context.CancelFunc
	In            io.ReadCloser
	Out           io.Writer
	Connection    net.Conn
}

func (t *TelnetConnection) Connect() error {
	dialer := &net.Dialer{}
	conn, err := dialer.DialContext(t.Ctx, "tcp", t.Address)
	t.Connection = conn
	return err
}

func (t *TelnetConnection) Close() error {
	t.cancelTimeout()
	return t.Connection.Close()
}

func (t *TelnetConnection) Send() error {
	scanner := bufio.NewScanner(t.In)
	for {
		select {
		case <-t.Ctx.Done():
			return nil
		default:
			if !scanner.Scan() {
				fmt.Fprintf(os.Stderr, "%v\n", errEOF)
				return nil
			}
			str := scanner.Text()
			_, err := t.Connection.Write([]byte(fmt.Sprintf("%s\n", str)))
			if err != nil {
				return err
			}
		}
	}
}

func (t *TelnetConnection) Receive() error {
	scanner := bufio.NewScanner(t.Connection)
	for {
		select {
		case <-t.Ctx.Done():
			return nil
		default:
			if !scanner.Scan() {
				fmt.Fprintf(os.Stderr, "%v\n", errConnectionClosedbyPeer)
				return nil
			}
			text := scanner.Text()
			_, err := t.Out.Write([]byte(fmt.Sprintf("%s\n", text)))
			if err != nil {
				return err
			}
		}
	}
}
