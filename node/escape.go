package node

import "strings"

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
