package main

import (
	"log"
	"os"

	"github.com/jamespwilliams/etty/wiktlang"
	"github.com/jamespwilliams/etty/wiktparse"
	"go.uber.org/zap"
)

func main() {
	logger, err := zap.NewDevelopment()
	if err != nil {
		log.Fatal("wiktionary-parse: failed to initialize logger", err)
	}

	if err := run(logger); err != nil {
		logger.Fatal("wiktionary-parse: fatal error", zap.Error(err))
	}
}

func run(logger *zap.Logger) error {
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
