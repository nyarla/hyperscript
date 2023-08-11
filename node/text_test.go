package node

import (
	"fmt"
	"testing"
)

func ExampleText() {
	text := Text(`<p>`).NodeString()

	fmt.Println(text)
	// Output:
	// &lt;p&gt;
}

func BenchmarkText(b *testing.B) {
	for idx := 0; idx < b.N; idx++ {
		Text(`<p>`).NodeString()
	}
}

func ExampleRaw() {
	raw := Raw(`<hr />`).NodeString()
	fmt.Println(raw)
	// Output:
	// <hr />
}

func BenchmarkRaw(b *testing.B) {
	for idx := 0; idx < b.N; idx++ {
		Raw(`<hr />`).NodeString()
	}
}
