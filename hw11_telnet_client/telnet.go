package main

import (
	"io"
	"net"
	"time"
)

type TelnetClient interface {
	Connect() error
	io.Closer
	Send() error
	Receive() error
}

func NewTelnetClient(address string, timeout time.Duration, in io.ReadCloser, out io.Writer) TelnetClient {
	return &TelnetConnection{
		Address: address,
		Timeout: timeout,
		In:      in,
		Out:     out,
	}
}

type TelnetConnection struct {
	Address    string
	Timeout    time.Duration
	In         io.ReadCloser
	Out        io.Writer
	Connection net.Conn
}

func (t *TelnetConnection) Connect() error {
	conn, err := net.DialTimeout("tcp", t.Address, t.Timeout)
	t.Connection = conn
	return err
}

func (t *TelnetConnection) Close() error {
	return t.Connection.Close()
}

func (t *TelnetConnection) Send() error {
	_, err := io.Copy(t.Connection, t.In)
	return err
}

func (t *TelnetConnection) Receive() error {
	_, err := io.Copy(t.Out, t.Connection)
	return err
}
