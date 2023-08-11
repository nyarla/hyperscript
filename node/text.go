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

// Raw returns non escaped string as NodeStringer.
// This func should not use to rendering of untrusted content.
func Raw(src string) NodeStringer {
	if data, ok := raws[src]; ok {
		return data
	}

	raws[src] = raw(src)

	return raws[src]
}

// Text returns html escaped string as NodeStringer.
// This func is useful for to make text content from unstrusted text. (ex. user input by http query)
func Text(src string) NodeStringer {
	if data, ok := texts[src]; ok {
		return data
	}

	texts[src] = raw(html.EscapeString(src))

	return texts[src]
}
