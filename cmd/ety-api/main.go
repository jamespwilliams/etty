package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os"

	"github.com/jamespwilliams/etymology"
	"github.com/jamespwilliams/etymology/api"
)

func main() {
	if err := run(os.Args[1], os.Args[2], os.Args[3]); err != nil {
		log.Fatal("ety:", err)
	}
}

func run(wordnetPath, network, address string) error {
	wordnet, err := os.Open(wordnetPath)
	if err != nil {
		return fmt.Errorf("failed to open wordnet: %v", err)
	}

	ety, err := etymology.New(wordnet)
	if err != nil {
		return fmt.Errorf("failed to build etymology tree from wordnet: %v", err)
	}

	s := api.NewServer(ety)

	l, err := net.Listen(network, address)
	if err != nil {
		return fmt.Errorf("failed to create listener: %v", err)
	}

	return http.Serve(l, s)
}
