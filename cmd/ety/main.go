package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/jamespwilliams/etymology"
	. "github.com/logrusorgru/aurora"
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
	return format(n, "", true, true)
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

	sb.WriteString(formatWord(n.Word.Word))
	sb.WriteString(" (")
	sb.WriteString(Magenta(n.Word.Language).String())
	sb.WriteString(")")

	children := append(n.Etymology)

	for index, child := range children {
		childIsLast := index == len(children)-1

		newIndent := indent
		if !root && !lastChild {
			newIndent += "│   "
		} else if !root {
			newIndent += "    "
		}

		childStr := format(child, newIndent, false, childIsLast)
		sb.WriteString("\n")
		sb.WriteString(childStr)
	}

	return sb.String()
}

func formatWord(word string) string {
	if word[len(word)-1] == '-' {
		return Yellow(word[:len(word)-1]).String() + "-"
	}

	if word[0] == '-' {
		return "-" + Yellow(word[1:]).String()
	}

	return Yellow(word).String()
}
