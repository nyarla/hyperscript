package node

import (
	"fmt"
	"testing"
)

func ExampleAttr() {
	a1 := Attr(`id`, `bar`)
	a2 := Attr(`id`, `<bar>`)
	a3 := Attr(`<id>`, `bar`)

	fmt.Println(a1.NodeString())
	fmt.Println(a2.NodeString())
	fmt.Println(a3.NodeString())
	// Output:
	//  id="bar"
	//  id="&lt;bar&gt;"
	//  &lt;id&gt;="bar"
}

func BenchmarkAttr(b *testing.B) {
	for idx := 0; idx < b.N; idx++ {
		Attr(`id`, `bar`).NodeString()
	}
}
