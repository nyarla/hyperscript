package node

import (
	"strconv"
	"strings"
	"testing"
)

func TestElement(t *testing.T) {
	var out strings.Builder
	var tests = []struct {
		el       Element
		expected string
	}{
		{Element{`hr`}, `<hr />`},

		{Element{`hr`, Attr{`id`, `msg`}}, `<hr id="msg" />`},
		{Element{`b`, Text(`hi`)}, `<b>hi</b>`},
		{Element{`p`, Element{`b`, Text(`hi`)}}, `<p><b>hi</b></p>`},
		{Element{`hr`, Attr{`id`, `sep`}, Attr{`class`, `primary`}}, `<hr class="primary" id="sep" />`},
		{Element{`p`, Text(`this `), Text(`is `), Text(`the `), Text(`test!`)}, `<p>this is the test!</p>`},
		{Element{`p`, Text(`this `), Attr{`id`, `msg`}, Text(`is `), Attr{`class`, `primary`}, Text(`the `), Text(`test!`)}, `<p class="primary" id="msg">this is the test!</p>`},
	}

	for _, test := range tests {
		test.el.RenderTo(&out)
		if out.String() != test.expected {
			t.Errorf(`unexpected output: got %q, expected %q`, out.String(), test.expected)
		}

		out.Reset()
	}
}

func BenchmarkElement(b *testing.B) {
	var out strings.Builder
	var r = Element{
		`body`,
		Element{`hr`},
		Element{`hr`, Attr{`id`, `msg`}},
		Element{`b`, Text(`hi`)},
		Element{`p`, Element{`b`, Text(`hi`)}},
		Element{`hr`, Attr{`id`, `sep`}, Attr{`class`, `primary`}},
		Element{`p`, Text(`this `), Text(`is `), Text(`the `), Text(`test!`)},
		Element{`p`, Text(`this `), Attr{`id`, `msg`}, Text(`is `), Attr{`class`, `primary`}, Text(`the `), Text(`test!`)},
	}

	b.ReportAllocs()
	b.ResetTimer()

	for idx := 0; idx < b.N; idx++ {
		r.RenderTo(&out)
	}
}

func TestDynamicElementFunc(t *testing.T) {
	var out strings.Builder
	var count = 0
	var r = DynamicElementFunc(func(w Renderer) (err error) {
		count++
		return Element{`p`, Text(strconv.Itoa(count))}.RenderTo(w)
	})

	if r.Type() != ElementNode {
		t.Errorf(`unexpected node type: got %s, expected %s`, r.Type().String(), AttrNode.String())
	}

	for idx := 0; idx < 5; idx++ {
		r.RenderTo(&out)

		if out.String() != `<p>`+strconv.Itoa(idx+1)+`</p>` {
			t.Errorf(`unexpected output: got %s, expected %s`, out.String(), strconv.Itoa(idx+1))
		}

		out.Reset()
	}
}

func BenchmarkDynamicElementFunc(b *testing.B) {
	var out strings.Builder
	var count = 0

	b.ReportAllocs()
	b.ResetTimer()

	for idx := 0; idx < b.N; idx++ {
		var r = DynamicElementFunc(func(w Renderer) (err error) {
			count++
			return Element{`p`, Text(strconv.Itoa(count))}.RenderTo(w)
		})
		r.RenderTo(&out)
	}
}
