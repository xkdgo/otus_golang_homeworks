package main

import (
	"flag"
	"log"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

const defaultTelnetPort = "23"

func main() {
	var timeout time.Duration
	var address string
	var wg sync.WaitGroup
	stopCh := make(chan struct{}, 1)
	exit := make(chan struct{}, 1)
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT)
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
	wg.Add(2)
	go SendMessage(exit, &wg, telnetClient)
	go ReceiveMessage(exit, &wg, telnetClient)
	go func(wg *sync.WaitGroup) {
		wg.Wait()
		stopCh <- struct{}{}
	}(&wg)
	select {
	case <-stopCh:
	case <-sigCh:
		close(exit)
	}
}

func SendMessage(exit chan struct{}, wg *sync.WaitGroup, client TelnetClient) {
	defer client.Close()
	defer wg.Done()
	select {
	case <-exit:
		return
	default:
		err := client.Send()
		if err != nil {
			return
		}
	}
}

func ReceiveMessage(exit chan struct{}, wg *sync.WaitGroup, client TelnetClient) {
	defer client.Close()
	defer wg.Done()
	select {
	case <-exit:
		return
	default:
		err := client.Receive()
		if err != nil {
			return
		}
	}
}
