package main

import (
	"log"
	"os"

	"github.com/jamespwilliams/etymology/wiktlang"
	"github.com/jamespwilliams/etymology/wiktparse"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	twos, err := os.Open("wiktlang/langs_2.txt")
	if err != nil {
		return err
	}

	threes, err := os.Open("wiktlang/langs_3.txt")
	if err != nil {
		return err
	}

	languages := wiktlang.New(twos, threes)
	return wiktparse.ParseDump(os.Stdin, os.Stdout, languages)
}
