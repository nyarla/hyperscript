package node

import (
	"io"
	"strings"
)

var escape = strings.NewReplacer(
	"&", `&amp;`,
	"<", `&lt;`,
	">", `&gt;`,
	"\"", `&#34;`,
	"'", `&#39;`,
	"`", `&#96;`,
	"{", `&#123;`,
	"}", `&#125;`,
)

func Escape(w io.Writer, s string) error {
	if _, err := escape.WriteString(w, s); err != nil {
		return err
	}

	return nil
}
