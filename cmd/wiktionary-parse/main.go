package main

import (
	"fmt"
	"log"
	"os"

	"github.com/jamespwilliams/etymology/wiktlang"
	"github.com/jamespwilliams/etymology/wiktparse"
	"go.uber.org/zap"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	logger, err := zap.NewDevelopment()
	if err != nil {
		return fmt.Errorf("failed to initialize logger: %w", err)
	}

	twos, err := os.Open("wiktlang/langs_2.txt")
	if err != nil {
		return err
	}

	threes, err := os.Open("wiktlang/langs_3.txt")
	if err != nil {
		return err
	}

	languages, err := wiktlang.New(twos, threes)
	if err != nil {
		return err
	}

	return wiktparse.ParseDump(logger, os.Stdin, os.Stdout, languages)
}
