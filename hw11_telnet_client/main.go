package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"sync"
	"time"
)

func main() {
	// Place your code here,
	// P.S. Do not rush to throw context down, think think if it is useful with blocking operation?
	var userTimeout string
	var address string
	var wg sync.WaitGroup
	flag.StringVar(&userTimeout, "timeout", "10s", "timeout of telnet operation")
	flag.Parse()
	args := flag.Args()
	switch len(args) {
	case 1:
		address = fmt.Sprintf("%s:", args[0])
	case 2:
		address = fmt.Sprintf("%s:%s", args[0], args[1])
	default:
		log.Fatalf("You should enter \"address port\"")
	}
	timeout, err := time.ParseDuration(userTimeout)
	if err != nil {
		log.Fatalf("You should enter timeout like \"10s\"")
	}

	telnetClient := NewTelnetClient(address, timeout, os.Stdin, os.Stdout)
	err = telnetClient.Connect()
	if err != nil {
		log.Fatalf("Could not connect to server %s", address)
	}
	log.Printf("Connected to server %s", address)
	wg.Add(2)
	go func(wg *sync.WaitGroup, client TelnetClient) {
		defer client.Close()
		defer wg.Done()
		client.Send()
	}(&wg, telnetClient)
	go func(wg *sync.WaitGroup, client TelnetClient) {
		defer client.Close()
		defer wg.Done()
		client.Receive()
	}(&wg, telnetClient)
	wg.Wait()
}
