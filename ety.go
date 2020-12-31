package etymology

import (
	"bufio"
	"fmt"
	"io"
	"sort"
	"strings"
)

var relationWhitelist = map[string]struct{}{
	"rel:etymology":       {},
	"rel:is_derived_from": {},
	"rel:inherited":       {},
	"rel:borrowed":        {},
	"rel:derived":         {},
	"rel:prefix":          {},
	"rel:suffix":          {},
	"rel:component":       {},
}

type Word struct {
	Language string
	Word     string
}

type Etymology struct {
	edges map[Word][]Word
}

func New(data io.Reader) (Etymology, error) {
	edges := make(map[Word][]Word)

	scanner := bufio.NewScanner(data)
	for scanner.Scan() {
		components := strings.Split(scanner.Text(), "\t")

		if _, ok := relationWhitelist[components[1]]; !ok {
			continue
		}

		from := strings.Split(components[0], ":")
		to := strings.Split(components[2], ":")

		fromWord := Word{
			Language: from[0],
			Word:     strings.TrimSpace(from[1]),
		}

		edges[fromWord] = append(edges[fromWord], Word{
			Language: to[0],
			Word:     strings.TrimSpace(to[1]),
		})
	}
	fmt.Println(edges[Word{
		Language: "eng",
		Word:     "aerodynamically",
	}])

	if err := scanner.Err(); err != nil {
		return Etymology{}, fmt.Errorf("failed to scan file: %w", err)
	}

	return Etymology{edges: edges}, nil
}

type Relation int

const (
	RelationEtymology = iota
)

func (e Etymology) Lookup(word Word) Node {
	node := Node{
		Word: word,
	}

	for _, etym := range e.edges[word] {
		n := e.Lookup(etym)
		node.Etymology = append(node.Etymology, n)
	}

	sort.Sort(sortedNodes(node.Etymology))

	return node
}
