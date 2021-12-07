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

func main() {
	// Place your code here,
	// P.S. Do not rush to throw context down, think think if it is useful with blocking operation?
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
		address = net.JoinHostPort(args[0], "23")
	case 2:
		address = net.JoinHostPort(args[0], args[1])
	default:
		log.Fatalf("You should enter \"address port\"")
	}
	telnetClient := NewTelnetClient(address, timeout, os.Stdin, os.Stdout)
	if telnetClient.Connect() != nil {
		log.Fatalf("Could not connect to server %s", address)
	}
	log.Printf("Connected to %s", address)
	wg.Add(2)
	go func(wg *sync.WaitGroup, client TelnetClient) {
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
	}(&wg, telnetClient)
	go func(wg *sync.WaitGroup, client TelnetClient) {
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
	}(&wg, telnetClient)
	go func(wg *sync.WaitGroup) {
		wg.Wait()
		stopCh <- struct{}{}
	}(&wg)
	select {
	case <-stopCh:
	case <-sigCh:
		exit <- struct{}{}
	}
}
