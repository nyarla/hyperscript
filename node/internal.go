package node

import (
	"fmt"
	"html"
)

// NodeStringer is a interface of Text, Attribute or Node instances.
// This interface always use inside this library.
type NodeStringer interface {
	NodeString() string
}

type (
	rawStr  string
	rawAttr [1]string
	keyAttr [2]string
	keyNode [3]string
)

var (
	inMemoryRawCache  = make(map[string]NodeStringer)
	inMemoryHTMLCache = make(map[string]NodeStringer)
	inMemoryAttrCache = make(map[keyAttr]NodeStringer)
	inMemoryNodeCache = make(map[keyNode]NodeStringer)
)

func (s rawStr) NodeString() string {
	return string(s)
}

// Raw returns a non-escaped text as it.
// This function should not use to a untrusted html string.
func Raw(src string) NodeStringer {
	if data, ok := inMemoryRawCache[src]; ok {
		return data
	}

	inMemoryRawCache[src] = rawStr(src)

	return inMemoryRawCache[src]
}

// Text returns html-escaped text as NodeStringer.
func Text(src string) NodeStringer {
	if data, ok := inMemoryHTMLCache[src]; ok {
		return data
	}

	inMemoryHTMLCache[src] = rawStr(html.EscapeString(src))

	return inMemoryHTMLCache[src]
}

func (ka keyAttr) string() string {
	return fmt.Sprintf(` %s="%s"`, Text(ka[0]).NodeString(), Text(ka[1]).NodeString())
}

func (ra rawAttr) NodeString() string {
	return ra[0]
}

// Attr returns html tag's attribute pair like as `id="foo"`.
func Attr(k, v string) NodeStringer {
	key := keyAttr{k, v}

	if data, ok := inMemoryAttrCache[key]; ok {
		return data
	}

	inMemoryAttrCache[key] = rawAttr{key.string()}

	return inMemoryAttrCache[key]
}

// Node construct html tags by tag string and NodeStringers (Attr, Text or Raw functions makes it)
func Node(tag string, contains ...NodeStringer) NodeStringer {
	var (
		attrs, contents string
	)

	for _, content := range contains {
		if data, ok := content.(rawAttr); ok {
			attrs += data.NodeString()
			continue
		}

		contents += content.NodeString()
	}

	key := keyNode{tag, attrs, contents}
	if data, ok := inMemoryNodeCache[key]; ok {
		return data
	}

	if len(contents) == 0 {
		inMemoryNodeCache[key] = rawStr(fmt.Sprintf(`<%s%s />`, tag, attrs))
	} else {
		inMemoryNodeCache[key] = rawStr(fmt.Sprintf(`<%s%s>%s</%s>`, tag, attrs, contents, tag))
	}

	return inMemoryNodeCache[key]
}
