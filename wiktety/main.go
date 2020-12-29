package main

import (
	"log"
	"os"

	"github.com/jamespwilliams/etymology/wiktety/wiktlang"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

type textElem struct {
	Data string `xml:",chardata"`
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
	return parseDump(os.Stdin, languages)
}
