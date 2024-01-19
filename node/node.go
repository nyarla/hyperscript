package node

import (
	"io"
	"sort"
	"strings"
)

type NodeType int

const (
	TextNode NodeType = iota
	AttrNode
	ElementNode
)

type Tag string

type NodeBuilder interface {
	Type() NodeType
	WriteString(w io.StringWriter) (int, error)
}

type node struct {
	kind NodeType
	src  strings.Builder
}

func (n *node) Type() NodeType {
	return n.kind
}

func (n *node) WriteString(w io.StringWriter) (int, error) {
	return w.WriteString(n.src.String())
}

var htmlReplacer = strings.NewReplacer(
	`&`, `&amp;`,
	`>`, `&gt;`,
	`<`, `&lt;`,
	`"`, `&#34;`,
	`'`, `&#39;`,
	"`", `&#96;`,
	`{`, `&#123;`,
	`}`, `&#125;`,
)

type ByNodeType []NodeBuilder

func (a ByNodeType) Len() int      { return len(a) }
func (a ByNodeType) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a ByNodeType) Less(i, j int) bool {
	if a[i].Type() == AttrNode {
		return true
	}

	return i < j
}

func Unsafe(src string) NodeBuilder {
	n := new(node)
	n.kind = TextNode
	n.src.WriteString(src)

	return n
}

func Safe(src string) NodeBuilder {
	n := new(node)
	n.kind = TextNode
	n.src.WriteString(htmlReplacer.Replace(src))

	return n
}

func Attr(kv ...string) NodeBuilder {
	n := new(node)
	n.kind = AttrNode

	if len(kv) == 1 {
		n.src.WriteString(htmlReplacer.Replace(kv[0]))
		return n
	}

	if len(kv)%2 != 0 {
		panic(`node.Attr: odd argument count`)
	}

	for idx := 0; idx < len(kv); idx = idx + 2 {
		n.src.WriteString(htmlReplacer.Replace(kv[idx]))
		n.src.WriteString(`="`)
		n.src.WriteString(htmlReplacer.Replace(kv[idx+1]))
		if idx+1 == len(kv)-1 {
			n.src.WriteString(`"`)
		} else {
			n.src.WriteString(`" `)
		}
	}

	return n
}

func Element(el Tag, contains ...NodeBuilder) NodeBuilder {
	n := new(node)
	n.kind = ElementNode
	n.src.WriteString(`<`)
	n.src.WriteString(string(el))

	if len(contains) == 0 {
		n.src.WriteString(` />`)

		return n
	}

	if len(contains) == 1 && contains[0].Type() == AttrNode {
		n.src.WriteString(` `)
		contains[0].WriteString(&n.src)
		n.src.WriteString(` />`)

		return n
	}

	if len(contains) == 1 && contains[0].Type() == TextNode {
		n.src.WriteString(`>`)

		contains[0].WriteString(&n.src)

		n.src.WriteString(`</`)
		n.src.WriteString(string(el))
		n.src.WriteString(`>`)
		return n
	}

	sort.Sort(ByNodeType(contains))
	first, last := contains[0], contains[len(contains)-1]

	if first.Type() == AttrNode && first.Type() == last.Type() {
		n.src.WriteString(` `)
		for idx := 0; idx < len(contains); idx++ {
			contains[idx].WriteString(&n.src)
			n.src.WriteString(` `)
		}

		n.src.WriteString(`/>`)

		return n
	}

	if first.Type() != AttrNode && last.Type() != AttrNode {
		n.src.WriteString(`>`)

		for idx := 0; idx < len(contains); idx++ {
			contains[idx].WriteString(&n.src)
		}

		n.src.WriteString(`</`)
		n.src.WriteString(string(el))
		n.src.WriteString(`>`)

		return n
	}

	for idx := 0; idx < len(contains); idx++ {
		if idx == 0 || contains[idx].Type() == AttrNode {
			n.src.WriteString(` `)
		}

		if idx > 0 && contains[idx-1].Type() == AttrNode && contains[idx].Type() != AttrNode {
			n.src.WriteString(`>`)
		}

		contains[idx].WriteString(&n.src)
	}

	n.src.WriteString(`</`)
	n.src.WriteString(string(el))
	n.src.WriteString(`>`)

	return n
}
