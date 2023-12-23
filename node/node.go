package node

import (
	"fmt"
	"html"
	"strings"
)

type Raw string

func Text(src string) Raw {
	return Raw(html.EscapeString(src))
}

type Attr string

func Pair(name, value string) Attr {
	return Attr(fmt.Sprintf(`%s="%s"`, html.EscapeString(name), html.EscapeString(value)))
}

func Flag(name string) Attr {
	return Attr(html.EscapeString(name))
}

func Node(el string, contains ...interface{}) Raw {
	if len(contains) == 0 {
		return Raw(fmt.Sprintf(`<%s />`, html.EscapeString(el)))
	}

	attrs := make([]string, 0, len(contains))
	contents := make([]string, 0, len(contains))

	for _, contain := range contains {
		if data, ok := contain.(Attr); ok {
			attrs = append(attrs, string(data))
			continue
		}

		if data, ok := contain.(Raw); ok {
			contents = append(contents, string(data))
			continue
		}

		panic(fmt.Sprintf(`Unsupported %+v`, contain))
	}

	if len(contents) == 0 {
		return Raw(fmt.Sprintf(`<%s %s />`, html.EscapeString(el), strings.Join(attrs, " ")))
	}

	return Raw(fmt.Sprintf(`<%s %s>%s</%s>`, html.EscapeString(el), strings.Join(attrs, " "), strings.Join(contents, ""), html.EscapeString(el)))
}
