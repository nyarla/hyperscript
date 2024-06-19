package node

import (
	"strconv"
	"strings"
	"testing"
)

func TestRaw(t *testing.T) {
	var out strings.Builder
	var r = Raw(`hello, world`)
	var expected = `hello, world`

	if r.Type() != TextNode {
		t.Errorf(`unexpected node type: got %s, expected %s`, r.Type().String(), TextNode.String())
	}

	r.RenderTo(&out)

	if out.String() != expected {
		t.Errorf(`unexpected output: got %s, expected %s`, out.String(), expected)
	}
}

func BenchmarkRaw(b *testing.B) {
	var out strings.Builder

	b.ReportAllocs()
	b.ResetTimer()

	for idx := 0; idx < b.N; idx++ {
		var r = Raw(`hello, world`)
		r.RenderTo(&out)
	}
}

func TestText(t *testing.T) {
	var out strings.Builder
	var r = Text(`1 & 0`)
	var expected = `1 &amp; 0`

	if r.Type() != TextNode {
		t.Errorf(`unexpected node type: got %s, expected %s`, r.Type().String(), TextNode.String())
	}

	r.RenderTo(&out)

	if out.String() != expected {
		t.Errorf(`unexpected output: got %s, expected %s`, out.String(), expected)
	}
}

func BenchmarkText(b *testing.B) {
	var out strings.Builder

	b.ReportAllocs()
	b.ResetTimer()

	for idx := 0; idx < b.N; idx++ {
		var r = Text(`1 & 0`)
		r.RenderTo(&out)
	}
}

func TestDynamicTextFunc(t *testing.T) {
	var out strings.Builder
	var count = 0
	var r = DynamicTextFunc(func(w Renderer) error {
		count++

		_, err := w.WriteString(strconv.Itoa(count))
		return err
	})

	if r.Type() != TextNode {
		t.Errorf(`unexpected node type: got %s, expected %s`, r.Type().String(), TextNode.String())
	}

	for idx := 0; idx < 5; idx++ {
		r.RenderTo(&out)

		if out.String() != strconv.Itoa(idx+1) {
			t.Errorf(`unexpected output: got %s, expected %s`, out.String(), strconv.Itoa(idx+1))
		}

		out.Reset()
	}
}

func BenchmarkDynamicTextFunc(b *testing.B) {
	var out strings.Builder
	var count = 0

	b.ReportAllocs()
	b.ResetTimer()

	for idx := 0; idx < b.N; idx++ {
		var r = DynamicTextFunc(func(w Renderer) error {
			count++

			_, err := w.WriteString(strconv.Itoa(count))
			return err
		})
		r.RenderTo(&out)
	}
}
