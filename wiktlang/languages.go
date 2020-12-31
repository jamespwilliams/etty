package wiktlang

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

type Languages struct {
	twoFromName   map[string]string
	threeFromName map[string]string
}

func New(twos io.Reader, threes io.Reader) (Languages, error) {
	languages := Languages{
		twoFromName:   make(map[string]string),
		threeFromName: make(map[string]string),
	}

	scanner := bufio.NewScanner(twos)
	for scanner.Scan() {
		line := scanner.Text()
		comps := strings.Split(line, "\t")
		languages.twoFromName[comps[1]] = comps[0]
	}

	if err := scanner.Err(); err != nil {
		return Languages{}, fmt.Errorf("failed to scan list of two-letter language codes: %w", err)
	}

	scanner = bufio.NewScanner(threes)
	for scanner.Scan() {
		line := scanner.Text()
		comps := strings.Split(line, "\t")
		languages.threeFromName[comps[1]] = comps[0]
	}

	if err := scanner.Err(); err != nil {
		return Languages{}, fmt.Errorf("failed to scan list of three-letter language codes: %w", err)
	}

	return languages, nil
}

func (l Languages) CodeFromName(name string) string {
	if code, ok := l.twoFromName[name]; ok {
		return code
	}

	if code, ok := l.threeFromName[name]; ok {
		return code
	}

	println("failed to find code for name: " + name)
	return ""
}
