package node

import (
	"html"
	"strings"
	"testing"
)

func TestUnsafe(t *testing.T) {
	var out strings.Builder
	var tests = []struct{ in, expect string }{
		{`text`, `text`},
	}

	for _, test := range tests {
		out.Reset()

		component := Unsafe(test.in)

		if _, err := component.BuildNode(&out); err != nil {
			t.Errorf(`failed to render node by component: %+v`, err)
			continue
		}

		if out.String() != test.expect {
			t.Errorf(`unexpected value: (%+v) => %+v != %+v`, test.in, out.String(), test.expect)
		}
	}
}

func TestSafe(t *testing.T) {
	var out strings.Builder
	var tests = []struct{ expect, in string }{
		{`&amp;`, `&`},
		{`&gt;`, `>`},
		{`&lt;`, `<`},
		{`&#34;`, `"`},
		{`&#39;`, `'`},
		{`&#96;`, "`"},
		{`&#123;`, `{`},
		{`&#125;`, `}`},
	}

	for _, test := range tests {
		out.Reset()

		component := Safe(test.in)

		if _, err := component.BuildNode(&out); err != nil {
			t.Errorf(`failed to render node by component: %+v`, err)
			continue
		}

		if out.String() != test.expect {
			t.Errorf(`unexpected value: (%+v) => %+v != %+v`, test.in, out.String(), test.expect)
			continue
		}

		unescaped := html.UnescapeString(out.String())
		if test.in != unescaped {
			t.Errorf(`malformed html escape: (%+v) => %+v != %+v`, test.expect, test.in, unescaped)
		}
	}
}

func TestAttr(t *testing.T) {
	var out strings.Builder

	var single = []struct{ in, expect string }{
		{`crossorigin`, `crossorigin`},
		{`cross&origin`, `cross&amp;origin`},
	}

	for _, test := range single {
		out.Reset()

		component := Attr(test.in)

		if _, err := component.BuildNode(&out); err != nil {
			t.Errorf(`failed to render node by component: %+v`, err)
			continue
		}

		if out.String() != test.expect {
			t.Errorf(`unexpected value: (%+v) => %+v != %+v`, test.in, out.String(), test.expect)
		}
	}

	var multiple = []struct {
		in     []string
		expect string
	}{
		{[]string{`id`, `msg`}, `id="msg"`},
		{[]string{`id`, `m&g`}, `id="m&amp;g"`},
		{[]string{`i&d`, `msg`}, `i&amp;d="msg"`},
		{[]string{`i&d`, `m&g`}, `i&amp;d="m&amp;g"`},

		{[]string{`id`, `msg`, `class`, `highlight`}, `id="msg" class="highlight"`},
	}

	for _, test := range multiple {
		out.Reset()

		component := Attr(test.in...)

		if _, err := component.BuildNode(&out); err != nil {
			t.Errorf(`failed to render node by component: %+v`, err)
			continue
		}

		if out.String() != test.expect {
			t.Errorf(`unexpected value: (%+v) => %+v != %+v`, test.in, out.String(), test.expect)
		}
	}

}

func TestElement(t *testing.T) {
	var out strings.Builder
	var tests = []struct {
		component NodeBuilder
		expect    string
	}{
		// only tag
		{Element(`hr`), `<hr />`},

		// tag with attrs
		{Element(`hr`, Attr(`id`, `sep`)), `<hr id="sep" />`},
		{Element(`hr`, Attr(`id`, `sep`, `class`, `mark`), Attr(`data-bind`, `foo`)), `<hr data-bind="foo" id="sep" class="mark" />`},

		// tag with contents
		{Element(`p`, Unsafe(`hello `), Element(`strong`, Unsafe(`nyarla`))), `<p>hello <strong>nyarla</strong></p>`},
		{Element(`p`, Element(`mark`, Unsafe(`hello`)), Unsafe(` `), Element(`strong`, Unsafe(`nyarla`))), `<p><mark>hello</mark> <strong>nyarla</strong></p>`},

		// tag with attrs and contents
		{Element(`p`, Attr(`id`, `msg`), Unsafe(`hello `), Element(`strong`, Unsafe(`nyarla`))), `<p id="msg">hello <strong>nyarla</strong></p>`},
	}

	for _, test := range tests {
		out.Reset()

		if _, err := test.component.BuildNode(&out); err != nil {
			t.Errorf(`failed to render node by component: %+v`, err)
			continue
		}

		if out.String() != test.expect {
			t.Errorf(`unexpected value: %+v != %+v`, out.String(), test.expect)
		}
	}
}

func BenchmarkUnsafe(b *testing.B) {
	var out strings.Builder
	var component = Unsafe(`test`)

	b.ReportAllocs()
	b.ResetTimer()

	for idx := 0; idx < b.N; idx++ {
		component.BuildNode(&out)
	}
}

func BenchmarkSafe(b *testing.B) {
	var out strings.Builder
	var component = Safe(`test`)

	b.ReportAllocs()
	b.ResetTimer()

	for idx := 0; idx < b.N; idx++ {
		component.BuildNode(&out)
	}
}

func BenchmarkAttr(b *testing.B) {
	var out strings.Builder
	var single = Attr(`crossorigin`)
	var pair = Attr(`id`, `msg`)

	b.ReportAllocs()
	b.ResetTimer()

	for idx := 0; idx < b.N; idx++ {
		single.BuildNode(&out)
		pair.BuildNode(&out)
	}
}

func BenchmarkElement(b *testing.B) {
	var out strings.Builder
	var p = Element(`p`, Element(`strong`, Safe(`hi, `), Safe(`this is test message.`)))

	b.ReportAllocs()
	b.ResetTimer()

	for idx := 0; idx < b.N; idx++ {
		p.BuildNode(&out)
	}
}

func BenchmarkRealCase(b *testing.B) {
	var out strings.Builder
	var t = Element(`html`,
		Element(`body`,
			Element(`h1`, Safe(`title`)),
			Element(`p`, Safe(`This is an example page!`)),
			Element(`ul`,
				Element(`li`, Safe(`foo`)),
				Element(`li`, Safe(`bar`)),
				Element(`li`, Safe(`baz`)),
				Element(`hr`))))

	b.ReportAllocs()
	b.ResetTimer()

	for idx := 0; idx < b.N; idx++ {
		t.BuildNode(&out)
	}
}
