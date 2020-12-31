package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os"

	"github.com/jamespwilliams/etymology"
	"github.com/jamespwilliams/etymology/api"
	"go.uber.org/zap"
)

func main() {
	logger, err := zap.NewDevelopment()
	if err != nil {
		log.Fatal("ety-api: failed to initialise logger", err)
	}

	if err := run(logger, os.Args[1], os.Args[2], os.Args[3]); err != nil {
		logger.Fatal("ety-api: fatal error", zap.Error(err))
	}
}

func run(logger *zap.Logger, wordnetPath, network, address string) error {
	wordnet, err := os.Open(wordnetPath)
	if err != nil {
		return fmt.Errorf("failed to open wordnet: %v", err)
	}

	ety, err := etymology.New(wordnet)
	if err != nil {
		return fmt.Errorf("failed to build etymology tree from wordnet: %v", err)
	}

	s := api.NewServer(logger, ety)

	l, err := net.Listen(network, address)
	if err != nil {
		return fmt.Errorf("failed to create listener: %v", err)
	}

	return http.Serve(l, s)
}
