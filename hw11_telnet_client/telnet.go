package main

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log"
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
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, timeout)
	return &TelnetConnection{
		Address: address,
		Ctx:     ctx,
		Cancel:  cancel,
		In:      in,
		Out:     out,
	}
}

// Place your code here.
// P.S. Author's solution takes no more than 50 lines.
type TelnetConnection struct {
	Address    string
	Ctx        context.Context
	Cancel     context.CancelFunc
	In         io.ReadCloser
	Out        io.Writer
	Connection net.Conn
}

func (t *TelnetConnection) Connect() error {
	dialer := &net.Dialer{}
	conn, err := dialer.DialContext(t.Ctx, "tcp", t.Address)
	t.Connection = conn
	return err
}

func (t *TelnetConnection) Close() error {
	t.Cancel()
	return t.Connection.Close()
}

func (t *TelnetConnection) Send() error {
	scanner := bufio.NewScanner(t.In)
OUTER:
	for {
		select {
		case <-t.Ctx.Done():
			break OUTER
		default:
			if !scanner.Scan() {
				break OUTER
			}
			str := scanner.Text()
			t.Connection.Write([]byte(fmt.Sprintf("%s\n", str)))
		}
	}
	log.Printf("Finished writeRoutine")
	return nil
}

func (t *TelnetConnection) Receive() error {
	scanner := bufio.NewScanner(t.Connection)
OUTER:
	for {
		select {
		case <-t.Ctx.Done():
			break OUTER
		default:
			if !scanner.Scan() {
				log.Printf("Disconnected from remote server")
				break OUTER
			}
			text := scanner.Text()
			t.Out.Write([]byte(fmt.Sprintf("%s\n", text)))
		}
	}
	log.Printf("Finished readRoutine")
	return nil
}
