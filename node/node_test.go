package node

import (
	"fmt"
	"testing"
)

func ExampleNode() {
	node := Node(`a`, Text(`link`), Attr(`href`, `https://example.com`))

	fmt.Println(node.NodeString())

	// Output:
	// <a href="https://example.com">link</a>
}

func BenchmarkNode(t *testing.B) {
	for idx := 0; idx < t.N; idx++ {
		Node(`a`, Attr(`href`, `http://example.com`), Text(`link`)).NodeString()
	}
}
