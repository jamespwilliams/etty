package etty

type Node struct {
	Word      Word
	Etymology []Node
}

func contains(nodes []Node, node Node) bool {
	for _, n := range nodes {
		if n.Word == node.Word {
			return true
		}
	}

	return false
}

type sortedNodes []Node

func (n sortedNodes) Len() int {
	return len(n)
}

func (n sortedNodes) Less(i, j int) bool {
	ni := n[i].Word.Word
	nj := n[j].Word.Word

	niSuffix := ni[0] == '-'
	njSuffix := nj[0] == '-'

	niPrefix := ni[len(ni)-1] == '-'
	njPrefix := nj[len(nj)-1] == '-'

	if niSuffix && njSuffix || niPrefix && njPrefix {
		return ni < nj
	}

	if niSuffix || njSuffix {
		return njSuffix
	}

	if niPrefix || njPrefix {
		return niPrefix
	}

	return ni < nj
}

func (n sortedNodes) Swap(i, j int) {
	n[i], n[j] = n[j], n[i]
}
