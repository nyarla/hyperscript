package node

import "html"

var (
	raws  = make(map[string]NodeStringer)
	texts = make(map[string]NodeStringer)
)

type raw string

func (r raw) NodeString() string {
	return string(r)
}

func Raw(src string) NodeStringer {
	if data, ok := raws[src]; ok {
		return data
	}

	raws[src] = raw(src)

	return raws[src]
}

func Text(src string) NodeStringer {
	if data, ok := texts[src]; ok {
		return data
	}

	texts[src] = raw(html.EscapeString(src))

	return texts[src]
}
