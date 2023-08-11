package node

import (
	"fmt"
	"strings"
)

// NodeStringer is interface for generates a parts of html string.
type NodeStringer interface {
	NodeString() string
}

type attrNode [128]string

func (a attrNode) String() string {
	return strings.Join(a[:], ``)
}

type textNode [256]string

func (t textNode) String() string {
	return strings.Join(t[:], ``)
}

type node struct {
	tag   string
	attrs attrNode
	texts textNode
	aIdx  int
	tIdx  int
}

func (n node) NodeString() string {
	if n.tIdx == 0 && n.aIdx == 0 {
		return fmt.Sprintf(`<%s />`, n.tag)
	}

	if n.tIdx == 0 {
		return fmt.Sprintf(`<%s %s/>`, n.tag, n.attrs.String())
	}

	if n.aIdx == 0 {
		return fmt.Sprintf(`<%s>%s</%s>`, n.tag, n.texts.String(), n.tag)
	}

	return fmt.Sprintf(`<%s%s>%s</%s>`, n.tag, n.attrs.String(), n.texts.String(), n.tag)
}

var nodes = make(map[node]NodeStringer)

// Node returns html string as NodeStringer
func Node(tag string, contains ...NodeStringer) NodeStringer {
	n := node{
		tag:   tag,
		attrs: [128]string{},
		texts: [256]string{},
		aIdx:  0,
		tIdx:  0,
	}

	for _, val := range contains {
		if k, ok := val.(attr); ok {
			n.attrs[n.aIdx] = k.NodeString()
			n.aIdx++

			continue
		}

		n.texts[n.tIdx] = val.NodeString()
		n.tIdx++
	}

	if data, ok := nodes[n]; ok {
		return data
	}

	nodes[n] = raw(n.NodeString())

	return nodes[n]
}
