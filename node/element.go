package node

import (
	"slices"
)

type Element []any

func (e Element) Type() NodeType {
	return ElementNode
}

func (e Element) RenderTo(w Renderer) (err error) {
	var tag = e[0].(string)

	_, err = w.WriteRune('<')
	if err != nil {
		return
	}

	_, err = w.WriteString(tag)
	if err != nil {
		return
	}

	if len(e) == 1 {
		_, err = w.WriteString(` />`)
		return
	}

	if len(e) == 2 {
		val := e[1].(Node)

		if val.Type() == AttrNode {
			_, err = w.WriteRune(' ')
			if err != nil {
				return
			}

			val.RenderTo(w)

			_, err = w.WriteString(` />`)
			return
		}

		_, err = w.WriteString(`>`)
		if err != nil {
			return
		}

		val.RenderTo(w)

		_, err = w.WriteString(`</`)
		if err != nil {
			return
		}

		_, err = w.WriteString(tag)
		if err != nil {
			return
		}

		_, err = w.WriteRune('>')
		return
	}

	e = e[1:]
	slices.SortStableFunc(e, func(a, b any) int {
		at := a.(Node).Type()
		bt := b.(Node).Type()

		if at == AttrNode || bt == AttrNode {
			return -1
		}

		return 0
	})

	first, last := e[0].(Node), e[len(e)-1].(Node)
	if first.Type() == AttrNode && last.Type() == AttrNode {
		_, err = w.WriteRune(' ')
		if err != nil {
			return
		}

		for _, item := range e {
			err = item.(Node).RenderTo(w)
			if err != nil {
				return
			}

			_, err = w.WriteRune(' ')
			if err != nil {
				return
			}
		}

		_, err = w.WriteString(`/>`)
		return
	}

	if first.Type() != AttrNode && last.Type() != AttrNode {
		_, err = w.WriteRune('>')
		if err != nil {
			return
		}

		for _, item := range e {
			err = item.(Node).RenderTo(w)
			if err != nil {
				return
			}
		}

		_, err = w.WriteString(`</`)
		if err != nil {
			return
		}

		_, err = w.WriteString(tag)
		if err != nil {
			return
		}

		_, err = w.WriteRune('>')
		return
	}

	for idx, item := range e {
		val := item.(Node)

		if idx == 0 || val.Type() == AttrNode {
			_, err = w.WriteRune(' ')
			if err != nil {
				return
			}
		}

		if idx > 0 && e[idx-1].(Node).Type() == AttrNode && val.Type() != AttrNode {
			_, err = w.WriteRune('>')
			if err != nil {
				return
			}
		}

		err = val.RenderTo(w)
		if err != nil {
			return
		}
	}

	_, err = w.WriteString(`</`)
	if err != nil {
		return
	}

	_, err = w.WriteString(tag)
	if err != nil {
		return
	}

	_, err = w.WriteRune('>')
	return
}

type DynamicElementFunc func(w Renderer) error

func (f DynamicElementFunc) Type() NodeType {
	return ElementNode
}

func (f DynamicElementFunc) RenderTo(w Renderer) (err error) {
	return f(w)
}
