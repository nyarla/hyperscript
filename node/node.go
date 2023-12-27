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
		kv[0] = htmlReplacer.Replace(kv[0])
		return AttrNodeBuilderFunc(func(w *strings.Builder) (int, error) {
			return w.WriteString(kv[0])
		})
	}

	if len(kv)%2 != 0 {
		panic(`node.Attr: odd argument count`)
	}

	for idx := 0; idx < len(kv); idx = idx + 2 {
		kv[idx] = htmlReplacer.Replace(kv[idx]) + `="`

		if idx+1 == len(kv)-1 {
			kv[idx+1] = htmlReplacer.Replace(kv[idx+1]) + `"`
		} else {
			kv[idx+1] = htmlReplacer.Replace(kv[idx+1]) + `" `
		}
	}

	kv[0] = strings.Join(kv, "")
	return AttrNodeBuilderFunc(func(w *strings.Builder) (int, error) {
		return w.WriteString(kv[0])
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
	// element has no contents
	if len(contains) == 0 {
		return ElementNodeBuilderFunc(func(w *strings.Builder) (total int, throw error) {
			var (
				count int
				err   error
			)

			count, err = w.WriteString(`<`)
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

			count, err = w.WriteString(` />`)
			total += count
			if err != nil {
				throw = err
				return
			}

			return
		})
	}

	// element has one attribute
	if len(contains) == 1 && contains[0].Type() == AttrNode {
		return ElementNodeBuilderFunc(func(w *strings.Builder) (total int, throw error) {
			var (
				count int
				err   error
			)

			count, err = w.WriteString(`<`)
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

			count, err = w.WriteString(` `)
			total += count
			if err != nil {
				throw = err
				return
			}

			count, err = contains[0].BuildNode(w)
			total += count
			if err != nil {
				throw = err
				return
			}

			count, err = w.WriteString(` />`)
			total += count
			if err != nil {
				throw = err
				return
			}

			return
		})
	}

	// element has one content
	if len(contains) == 1 && (contains[0].Type() == TextNode || contains[0].Type() == ElementNode) {
		return ElementNodeBuilderFunc(func(w *strings.Builder) (total int, throw error) {
			var (
				count int
				err   error
			)

			count, err = w.WriteString(`<`)
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

			count, err = contains[0].BuildNode(w)
			total += count
			if err != nil {
				throw = err
				return
			}

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

	sort.Sort(ByNodeType(contains))

	// element has attributes without content
	if contains[0].Type() == AttrNode && contains[len(contains)-1].Type() == AttrNode {
		return ElementNodeBuilderFunc(func(w *strings.Builder) (total int, throw error) {
			var (
				count int
				err   error
			)

			count, err = w.WriteString(`<`)
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

			for idx := 0; idx < len(contains); idx++ {
				count, err = w.WriteString(` `)
				total += count
				if err != nil {
					throw = err
					return
				}
				count, err = contains[idx].BuildNode(w)
				total += count
				if err != nil {
					throw = err
					return
				}
			}

			// close tag
			count, err = w.WriteString(` />`)
			total += count
			if err != nil {
				throw = err
				return
			}

			return
		})
	}

	// element has contents without attribute
	if (contains[0].Type() == TextNode || contains[0].Type() == ElementNode) && (contains[len(contains)-1].Type() == TextNode || contains[len(contains)-1].Type() == ElementNode) {
		return ElementNodeBuilderFunc(func(w *strings.Builder) (total int, throw error) {
			var (
				count int
				err   error
			)

			count, err = w.WriteString(`<`)
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

			for idx := 0; idx < len(contains); idx++ {
				count, err = contains[idx].BuildNode(w)
				total += count
				if err != nil {
					throw = err
					return
				}
			}

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

	// element has attributes and contents
	return ElementNodeBuilderFunc(func(w *strings.Builder) (total int, throw error) {
		var (
			count int
			err   error
		)

		// open tag
		count, err = w.WriteString(`<`)
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

		for idx := 0; idx < len(contains); idx++ {
			if idx == 0 || contains[idx].Type() == AttrNode {
				count, err = w.WriteString(` `)
				total += count
				if err != nil {
					throw = err
					return
				}
			}

			if idx > 0 && contains[idx-1].Type() == AttrNode && (contains[idx].Type() == TextNode || contains[idx].Type() == ElementNode) {
				count, err = w.WriteString(`>`)
				total += count
				if err != nil {
					throw = err
					return
				}
			}

			count, err = contains[idx].BuildNode(w)
			total += count
			if err != nil {
				throw = err
				return
			}
		}

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
