package node

import (
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
	BuildNode(w *strings.Builder) (int, error)
}

type TextNodeBuilderFunc func(w *strings.Builder) (int, error)

func (f TextNodeBuilderFunc) Type() NodeType { return TextNode }
func (f TextNodeBuilderFunc) BuildNode(w *strings.Builder) (int, error) {
	return f(w)
}

type AttrNodeBuilderFunc func(w *strings.Builder) (int, error)

func (f AttrNodeBuilderFunc) Type() NodeType { return AttrNode }
func (f AttrNodeBuilderFunc) BuildNode(w *strings.Builder) (int, error) {
	return f(w)
}

type ElementNodeBuilderFunc func(w *strings.Builder) (int, error)

func (f ElementNodeBuilderFunc) Type() NodeType { return ElementNode }
func (f ElementNodeBuilderFunc) BuildNode(w *strings.Builder) (int, error) {
	return f(w)
}

func Unsafe(src string) NodeBuilder {
	return TextNodeBuilderFunc(func(w *strings.Builder) (int, error) {
		return w.WriteString(src)
	})
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

func Safe(src string) NodeBuilder {
	src = htmlReplacer.Replace(src)
	return TextNodeBuilderFunc(func(w *strings.Builder) (int, error) {
		return w.WriteString(src)
	})
}

func Attr(kv ...string) NodeBuilder {
	if len(kv) == 1 {
		return AttrNodeBuilderFunc(func(w *strings.Builder) (int, error) {
			return htmlReplacer.WriteString(w, kv[0])
		})
	}

	if len(kv)%2 != 0 {
		panic(`node.Attr: odd argument count`)
	}

	return AttrNodeBuilderFunc(func(w *strings.Builder) (total int, throw error) {
		for idx := 0; idx < len(kv); idx = idx + 2 {
			var (
				count int
				err   error
			)

			count, err = htmlReplacer.WriteString(w, kv[idx])
			total += count
			if err != nil {
				throw = err
				return
			}

			count, err = w.WriteString(`=`)
			total += count
			if err != nil {
				throw = err
				return
			}
			count, err = w.WriteString(`"`)
			total += count
			if err != nil {
				throw = err
				return
			}

			count, err = htmlReplacer.WriteString(w, kv[idx+1])
			total += count
			if err != nil {
				throw = err
				return
			}

			count, err = w.WriteString(`"`)
			total += count
			if err != nil {
				throw = err
				return
			}
		}

		return
	})
}

type ByNodeType []NodeBuilder

func (a ByNodeType) Len() int      { return len(a) }
func (a ByNodeType) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a ByNodeType) Less(i, j int) bool {
	if a[i].Type() == AttrNode {
		return true
	}

	return i < j
}

func Element(el Tag, contains ...NodeBuilder) NodeBuilder {
	sort.Sort(ByNodeType(contains))
	return ElementNodeBuilderFunc(func(w *strings.Builder) (total int, throw error) {
		var (
			count int
			err   error
		)

		// begin tag start
		count, err = w.WriteString(`<`)
		total += count
		if err != nil {
			throw = err
			return
		}

		// write tag name
		count, err = w.WriteString(string(el))
		total += count
		if err != nil {
			throw = err
			return
		}

		// write contains
		for idx, contain := range contains {
			if (idx == 0 && contain.Type() == TextNode || contain.Type() == ElementNode) || (idx > 0 && contains[idx-1].Type() == AttrNode && (contain.Type() == TextNode || contain.Type() == ElementNode)) {
				count, err = w.WriteString(`>`)
				total += count
				if err != nil {
					throw = err
					return
				}
			}

			if contain.Type() == AttrNode {
				count, err = w.WriteString(` `)
				total += count
				if err != nil {
					throw = err
					return
				}
			}

			count, err = contain.BuildNode(w)
			total += count
			if err != nil {
				throw = err
				return
			}
		}

		// close tag if content is empty
		if len(contains) == 0 || contains[len(contains)-1].Type() == AttrNode {
			count, err = w.WriteString(` />`)
			total += count
			if err != nil {
				throw = err
				return
			}

			return
		}

		// close tag
		count, err = w.WriteString(`</`)
		total += count
		if err != nil {
			throw = err
			return
		}

		count, err = w.WriteString(string(el))
		total += count
		if err != nil {
			throw = err
			return
		}

		count, err = w.WriteString(`>`)
		total += count
		if err != nil {
			throw = err
			return
		}

		return
	})
}
