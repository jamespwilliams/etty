package wiktlang

import (
	"bufio"
	"io"
	"strings"
)

type Languages struct {
	twoFromName   map[string]string
	threeFromName map[string]string
}

func New(twos io.Reader, threes io.Reader) Languages {
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

	// TODO check scanner err

	scanner = bufio.NewScanner(threes)
	for scanner.Scan() {
		line := scanner.Text()
		comps := strings.Split(line, "\t")
		languages.threeFromName[comps[1]] = comps[0]
	}

	// TODO check scanner err

	return languages
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
