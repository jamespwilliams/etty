package etymology

type Node struct {
	Word        Word
	DerivedFrom []Node
	Etymology   []Node
}

func contains(nodes []Node, node Node) bool {
	for _, n := range nodes {
		if n.Word == node.Word {
			return true
		}
	}

	return false
}

type SortedNodes []Node

func (n SortedNodes) Len() int {
	return len(n)
}

func (n SortedNodes) Less(i, j int) bool {
	ni := n[i].Word.Word
	nj := n[j].Word.Word

	if ni[0] == '-' {
		return false
	}

	if ni[len(ni)-1] == '-' {
		return true
	}

	return ni < nj
}

func (n SortedNodes) Swap(i, j int) {
	n[i], n[j] = n[j], n[i]
}
