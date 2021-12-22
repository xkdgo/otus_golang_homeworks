package main

import (
	"context"
	"flag"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const defaultTelnetPort = "23"

func main() {
	var timeout time.Duration
	var address string
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT)
	defer cancel()
	flag.DurationVar(&timeout, "timeout", 10*time.Second, "timeout of telnet operation")
	flag.Parse()
	args := flag.Args()
	switch len(args) {
	case 1:
		address = net.JoinHostPort(args[0], defaultTelnetPort)
	case 2:
		address = net.JoinHostPort(args[0], args[1])
	default:
		log.Fatalf("You should enter \"address port\"")
	}
	telnetClient := NewTelnetClient(address, timeout, os.Stdin, os.Stdout)
	err := telnetClient.Connect()
	if err != nil {
		log.Fatalf("Could not connect to server %s\n%q", address, err)
	}
	log.Printf("Connected to %s", address)

	go func() {
		err := telnetClient.Send()
		if err == nil {
			log.Printf("...EOF")
		} else {
			log.Println("send error:", err)
		}
		cancel()
	}()
	go func() {
		err := telnetClient.Receive()
		if err == nil {
			log.Printf("...Connection was closed by peer")
		} else {
			log.Println("recieve error:", err)
		}
		cancel()
	}()

	<-ctx.Done()
}
