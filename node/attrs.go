package node

import (
	"fmt"
	"html"
)

type kv [2]string
type attr [1]string

var attrs = make(map[kv]NodeStringer)

func (a kv) String() string {
	return fmt.Sprintf(` %s="%s"`, html.EscapeString(a[0]), html.EscapeString(a[1]))
}

func (a attr) NodeString() string {
	return a[0]
}

// Attr returns string of html attribute like as `id="foo"`
func Attr(k, v string) NodeStringer {
	kv := kv{k, v}

	if data, ok := attrs[kv]; ok {
		return data
	}

	attrs[kv] = attr{kv.String()}

	return attrs[kv]
}
