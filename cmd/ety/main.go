package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/jamespwilliams/etymology"
)

func main() {
	if err := run(os.Args[1], os.Args[2]); err != nil {
		log.Fatal("ety:", err)
	}
}

func run(wordnetPath, word string) error {
	wordnet, err := os.Open(wordnetPath)
	if err != nil {
		return fmt.Errorf("failed to open wordnet: %v", err)
	}

	ety, err := etymology.New(wordnet)
	if err != nil {
		return fmt.Errorf("failed to build etymology tree from wordnet: %v", err)
	}

	n := ety.Lookup(etymology.Word{Word: word, Language: "eng"})
	fmt.Println(formatRoot(n))

	return nil
}

func formatRoot(n etymology.Node) string {
	return format(n, "", false, true)
}

func format(n etymology.Node, indent string, root, lastChild bool) string {
	var sb strings.Builder

	if !root {
		sb.WriteString(indent)
		if lastChild {
			sb.WriteString("└── ")
		} else {
			sb.WriteString("├── ")
		}
	}

	sb.WriteString(n.Word.Word)
	sb.WriteString(" (")
	sb.WriteString(n.Word.Language)
	sb.WriteString(")")

	children := append(n.Etymology, n.DerivedFrom...)

	for index, child := range children {
		childIsLast := index == len(children)-1

		newIndent := indent
		if !lastChild {
			newIndent += "│   "
		} else {
			newIndent += "    "
		}

		childStr := format(child, newIndent, false, childIsLast)
		sb.WriteString("\n")
		sb.WriteString(childStr)
	}

	return sb.String()
}
