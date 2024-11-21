package node

import (
	"strconv"
	"strings"
	"testing"
)

func TestUnsafeFlag(t *testing.T) {
	var out strings.Builder
	var r = UnsafeFlag(`async`)
	var expected = `async`

	if r.Type() != AttrNode {
		t.Errorf(`unexpected node type: got %s, expected %s`, r.Type().String(), AttrNode.String())
	}

	r.RenderTo(&out)

	if out.String() != expected {
		t.Errorf(`unexpected output: got %s, expected %s`, out.String(), expected)
	}
}

func BenchmarkUnsafeFlag(b *testing.B) {
	var out strings.Builder

	b.ReportAllocs()
	b.ResetTimer()

	for idx := 0; idx < b.N; idx++ {
		var r = UnsafeFlag(`async`)
		r.RenderTo(&out)
	}
}

func TestFlag(t *testing.T) {
	var out strings.Builder
	var r = Flag(`1 & 0`)
	var expected = `1 &amp; 0`

	if r.Type() != AttrNode {
		t.Errorf(`unexpected node type: got %s, expected %s`, r.Type().String(), AttrNode.String())
	}

	r.RenderTo(&out)

	if out.String() != expected {
		t.Errorf(`unexpected output: got %s, expected %s`, out.String(), expected)
	}
}

func BenchmarkFlag(b *testing.B) {
	var out strings.Builder

	b.ReportAllocs()
	b.ResetTimer()

	for idx := 0; idx < b.N; idx++ {
		var r = Flag(`1 & 0`)
		r.RenderTo(&out)
	}
}

func TestUnsafeAttr(t *testing.T) {
	var out strings.Builder
	var r = UnsafeAttr{`id`, `msg`}
	var expected = `id="msg"`

	if r.Type() != AttrNode {
		t.Errorf(`unexpected node type: got %s, expected %s`, r.Type().String(), AttrNode.String())
	}

	r.RenderTo(&out)

	if out.String() != expected {
		t.Errorf(`unexpected output: got %s, expected %s`, out.String(), expected)
	}
}

func BenchmarkUnsafeAttr(b *testing.B) {
	var out strings.Builder

	b.ReportAllocs()
	b.ResetTimer()

	for idx := 0; idx < b.N; idx++ {
		var r = UnsafeAttr{`id`, `msg`}
		r.RenderTo(&out)
	}
}

func TestAttr(t *testing.T) {
	var out strings.Builder
	var r = Attr{`id`, `1 & 0`}
	var expected = `id="1 &amp; 0"`

	if r.Type() != AttrNode {
		t.Errorf(`unexpected node type: got %s, expected %s`, r.Type().String(), AttrNode.String())
	}

	r.RenderTo(&out)

	if out.String() != expected {
		t.Errorf(`unexpected output: got %s, expected %s`, out.String(), expected)
	}
}

func BenchmarkAttr(b *testing.B) {
	var out strings.Builder

	b.ReportAllocs()
	b.ResetTimer()

	for idx := 0; idx < b.N; idx++ {
		var r = Attr{`id`, `1 & 0`}
		r.RenderTo(&out)
	}
}

func TestDynamicAttrFunc(t *testing.T) {
	var out strings.Builder
	var count = 0
	var r = DynamicAttrFunc(func(w Renderer) (err error) {
		count++
		return Attr{`id`, strconv.Itoa(count)}.RenderTo(w)
	})

	if r.Type() != AttrNode {
		t.Errorf(`unexpected node type: got %s, expected %s`, r.Type().String(), AttrNode.String())
	}

	for idx := 0; idx < 5; idx++ {
		r.RenderTo(&out)

		if out.String() != `id="`+strconv.Itoa(idx+1)+`"` {
			t.Errorf(`unexpected output: got %s, expected %s`, out.String(), strconv.Itoa(idx+1))
		}

		out.Reset()
	}
}

func BenchmarkDynamicAttrFunc(b *testing.B) {
	var out strings.Builder
	var count = 0

	b.ReportAllocs()
	b.ResetTimer()

	for idx := 0; idx < b.N; idx++ {
		var r = DynamicAttrFunc(func(w Renderer) (err error) {
			count++
			return Attr{`id`, strconv.Itoa(count)}.RenderTo(w)
		})
		r.RenderTo(&out)
	}
}
