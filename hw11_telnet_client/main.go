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
	flag.DurationVar(&timeout, "timeout", 10*time.Second, "timeout of telnet operation")
	flag.Parse()
	args := flag.Args()
	switch len(args) {
	case 1:
		address = net.JoinHostPort(args[0], defaultTelnetPort)
	case 2:
		address = net.JoinHostPort(args[0], args[1])
	default:
		log.Printf("You should enter \"address port\"")
		return
	}
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT)
	defer cancel()
	telnetClient := NewTelnetClient(address, timeout, os.Stdin, os.Stdout)
	if err := telnetClient.Connect(); err != nil {
		log.Printf("Could not connect to server %s\n%q", address, err)
		return
	}
	log.Printf("Connected to %s", address)

	go func() {
		if err := telnetClient.Send(); err == nil {
			log.Printf("...EOF")
		} else {
			log.Println("send error:", err)
		}
		cancel()
	}()
	go func() {
		if err := telnetClient.Receive(); err == nil {
			log.Printf("...Connection was closed by peer")
		} else {
			log.Println("receive error:", err)
		}
		cancel()
	}()

	<-ctx.Done()
}
