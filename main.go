package main

import (
	"encoding/hex"
	"flag"
	"log"
	"net"
	"strings"
)

var (
	defaultServers = []string{"149.154.175.50:443", "149.154.167.51:443", "149.154.175.100:443", "149.154.167.91:443", "149.154.171.5:443"}
	localAddr      = flag.String("l", ":9999", "local address")
	servers        = flag.String("s", "149.154.175.50:443,149.154.167.51:443,149.154.175.100:443,149.154.167.91:443,149.154.171.5:443", "servers")
	secret         = flag.String("secret", "", "secret key")
	verbose        = flag.Bool("v", false, "display server actions")
)

func main() {
	flag.Parse()

	if !*verbose {
		log.SetFlags(0)
	}

	if *servers != "" {
		defaultServers = strings.Split(*servers, ",")
	}

	secret_hex, err := hex.DecodeString(*secret)
	if err != nil {
		log.Fatal(err)
	}

	laddr, err := net.ResolveTCPAddr("tcp", *localAddr)
	if err != nil {
		log.Fatal(err)
	}

	listener, err := net.ListenTCP("tcp", laddr)
	if err != nil {
		log.Fatal(err)
	}

	network := NewNetwork(defaultServers)

	log.Printf("Server running on %v", *localAddr)

	for {
		conn, err := listener.AcceptTCP()
		if err != nil {
			log.Printf("Failed to accept connection '%s'", err)
			continue
		}

		client := NewClient(conn, network, secret_hex)
		go client.Do()
	}
}
