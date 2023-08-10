package node

import (
	"fmt"
	"testing"
)

func ExampleText() {
	text := Text(`<p>`)

	fmt.Println(text)

	// Output:
	// &lt;p&gt;
}

func BenchmarkText(b *testing.B) {
	for idx := 0; idx < b.N; idx++ {
		Text(`<p>`)
	}
}

func ExampleRaw() {
	raw := Raw(`<hr />`)
	fmt.Println(raw)

	// Output:
	// <hr />
}

func BenchmarkRaw(b *testing.B) {
	for idx := 0; idx < b.N; idx++ {
		Raw(`<hr />`)
	}
}
