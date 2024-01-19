package attrs

import (
	"fmt"
	"strings"

	"github.com/nyarla/hyperscript/node"
)

var Text = node.Safe

func ExampleId() {
	var out strings.Builder
	var tag = node.Element(`p`, Id(`msg`), Text(`hi,`))

	tag.WriteString(&out)

	fmt.Println(out.String())
	// Ouput:
	// <p id="msg">hi,</p>
}
