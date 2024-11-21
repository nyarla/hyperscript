package node

type Raw string

func (t Raw) Type() NodeType {
	return TextNode
}

func (t Raw) RenderTo(w Renderer) (err error) {
	_, err = w.WriteString(string(t))
	return
}

type Text string

func (t Text) Type() NodeType {
	return TextNode
}

func (t Text) RenderTo(w Renderer) (err error) {
	_, err = escape.WriteString(w, string(t))
	return
}

type DynamicTextFunc func(w Renderer) error

func (f DynamicTextFunc) Type() NodeType {
	return TextNode
}

func (f DynamicTextFunc) RenderTo(w Renderer) (err error) {
	return f(w)
}
