package node

import (
	"strings"
	"testing"
)

func TestEscape(t *testing.T) {
	test := "this is a test &<>\"'`{}"
	expected := `this is a test &amp;&lt;&gt;&#34;&#39;&#96;&#123;&#125;`

	var out strings.Builder
	escape.WriteString(&out, test)

	if out.String() != expected {
		t.Errorf(`unexpected string: %s, want %s`, out.String(), expected)
	}
}

func BenchmarkEscape(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()

	var out strings.Builder
	for idx := 0; idx < b.N; idx++ {
		escape.WriteString(&out, `aaa & bbb <> ccc "' ddd {}`)
	}
}
