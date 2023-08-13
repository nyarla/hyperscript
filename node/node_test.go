package node

import (
	"fmt"
	"testing"
)

func ExampleNode() {
	link := Node(`a`, Text(`link`), Attr(`href`, `https://example.com`)).NodeString()
	hr := Node(`hr`).NodeString()
	img := Node(`img`, Attr(`width`, `100`), Attr(`height`, `100`)).NodeString()

	fmt.Println(link)
	fmt.Println(hr)
	fmt.Println(img)
	// Output:
	// <a href="https://example.com">link</a>
	// <hr />
	// <img width="100" height="100" />
}

func BenchmarkNode(t *testing.B) {
	for idx := 0; idx < t.N; idx++ {
		Node(`a`, Attr(`href`, `http://example.com`), Text(`link`)).NodeString()
	}
}
