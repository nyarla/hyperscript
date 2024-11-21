package node

type Renderer interface {
	Write(b []byte) (int, error)
	WriteString(s string) (int, error)
	WriteRune(r rune) (int, error)
	String() string
}

type NodeType int

const (
	UnknownNode NodeType = iota
	TextNode
	AttrNode
	ElementNode
)

func (t NodeType) String() string {
	switch t {
	case TextNode:
		return `TextNode`
	case AttrNode:
		return `AttrNode`

	case ElementNode:
		return `ElementNode`
	}

	return `UnknownNode`
}

type Node interface {
	Type() NodeType
	RenderTo(w Renderer) error
}
