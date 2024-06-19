package node

type UnsafeFlag string

func (a UnsafeFlag) Type() NodeType {
	return AttrNode
}

func (a UnsafeFlag) RenderTo(w Renderer) (err error) {
	_, err = w.WriteString(string(a))
	return
}

type Flag string

func (a Flag) Type() NodeType {
	return AttrNode
}

func (a Flag) RenderTo(w Renderer) (err error) {
	_, err = escape.WriteString(w, string(a))
	return
}

type UnsafeAttr [2]string

func (a UnsafeAttr) Type() NodeType {
	return AttrNode
}

func (a UnsafeAttr) RenderTo(w Renderer) (err error) {
	_, err = w.WriteString(a[0])
	if err != nil {
		return err
	}

	_, err = w.WriteString(`="`)
	if err != nil {
		return err
	}

	_, err = w.WriteString(a[1])
	if err != nil {
		return err
	}

	_, err = w.WriteRune('"')
	if err != nil {
		return err
	}

	return
}

type Attr [2]string

func (a Attr) Type() NodeType {
	return AttrNode
}

func (a Attr) RenderTo(w Renderer) (err error) {
	_, err = escape.WriteString(w, a[0])
	if err != nil {
		return err
	}

	_, err = w.WriteString(`="`)
	if err != nil {
		return err
	}

	_, err = escape.WriteString(w, a[1])
	if err != nil {
		return err
	}

	_, err = w.WriteRune('"')
	if err != nil {
		return err
	}

	return
}

type DynamicAttrFunc func(w Renderer) error

func (f DynamicAttrFunc) Type() NodeType {
	return AttrNode
}

func (f DynamicAttrFunc) RenderTo(w Renderer) (err error) {
	return f(w)
}
