package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/jamespwilliams/etty"
	. "github.com/logrusorgru/aurora"
)

func main() {
	lang := "eng"
	if len(os.Args) >= 4 {
		lang = os.Args[3]
	}

	if err := run(os.Args[1], os.Args[2], lang); err != nil {
		log.Fatal("etty:", err)
	}
}

func run(wordnetPath, word, lang string) error {
	wordnet, err := os.Open(wordnetPath)
	if err != nil {
		return fmt.Errorf("failed to open wordnet: %v", err)
	}

	tree, err := etty.New(wordnet)
	if err != nil {
		return fmt.Errorf("failed to build etty tree from wordnet: %v", err)
	}

	node := tree.Lookup(etty.Word{Word: word, Language: lang})
	fmt.Println(formatRoot(node))

	return nil
}

func formatRoot(n etty.Node) string {
	return format(n, "", true, true)
}

func format(n etty.Node, indent string, root, lastChild bool) string {
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

	children := n.Etymology

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
