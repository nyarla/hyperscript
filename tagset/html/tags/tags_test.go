package tags

import (
	"fmt"
	"strings"

	"github.com/nyarla/hyperscript/node"
)

var Text = node.Safe

func ExampleP() {
	var out strings.Builder
	var p = P(Text("hello, world!"))

	p.WriteString(&out)

	fmt.Println(out.String())
	// Output:
	// <p>hello, world!</p>
}
