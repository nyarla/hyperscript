package node

import "fmt"

// NodeStringer is interface for generates a parts of html string.
type NodeStringer interface {
	NodeString() string
}

type node [3]string

var nodes = make(map[node]NodeStringer)

func Node(tag string, contains ...NodeStringer) NodeStringer {
	var (
		a, c string
	)

	for _, content := range contains {
		if data, ok := content.(attr); ok {
			a += data.NodeString()
			continue
		}

		c += content.NodeString()
	}

	var node = [3]string{tag, a, c}

	if data, ok := nodes[node]; ok {
		return data
	}

	if len(c) == 0 {
		nodes[node] = raw(fmt.Sprintf(`<%s%s/>`, tag, a))
	}

	nodes[node] = raw(fmt.Sprintf(`<%s%s>%s</%s>`, tag, a, c, tag))

	return nodes[node]
}
