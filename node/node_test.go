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

	var single = [][2]string{
		{`crossorigin`, `crossorigin`},
		{`cross&origin`, `cross&amp;origin`},
	}

	for _, test := range single {
		out.Reset()

		from := test[0]
		to := test[1]

		component := Attr(from)

		_, err := component.BuildNode(&out)
		if err != nil {
			t.Errorf(`failed to write component string: %+v`, err)
			continue
		}

		if out.String() != to {
			t.Errorf(`unexpected value: %+v => %+v != %+v`, from, out.String(), to)
			continue
		}
	}

	var pairs = [][3]string{
		{`id`, `msg`, `id="msg"`},
		{`i&d`, `msg`, `i&amp;d="msg"`},
		{`id`, `m&g`, `id="m&amp;g"`},
		{`i&d`, `m&g`, `i&amp;d="m&amp;g"`},
	}

	for _, test := range pairs {
		out.Reset()

		k := test[0]
		v := test[1]
		expect := test[2]

		component := Attr(k, v)
		_, err := component.BuildNode(&out)
		if err != nil {
			t.Errorf(`failed to write component string: %+v`, err)
			continue
		}

		if out.String() != expect {
			t.Errorf(`unexpected value: (%+v, %+v) => %+v != %+v`, k, v, out.String(), expect)
			continue
		}
	}
}

func TestElement(t *testing.T) {
	var out strings.Builder
	var tests = []struct {
		el     NodeBuilder
		expect string
	}{
		{Element(`hr`), `<hr />`},
		{Element(`hr`, Attr(`id`, `msg`)), `<hr id="msg" />`},
		{Element(`hr`, Attr(`id`, `msg`), Attr(`class`, `hr`)), `<hr class="hr" id="msg" />`},
		{Element(`p`, Safe(`hello`), Safe(`, `), Safe(`world!`)), `<p>hello, world!</p>`},
		{Element(`p`, Attr(`id`, `msg`), Attr(`class`, `highlight`), Safe(`hello`), Safe(`, `), Safe(`world!`)), `<p class="highlight" id="msg">hello, world!</p>`},
		{Element(`p`, Element(`strong`, Safe(`hi,`))), `<p><strong>hi,</strong></p>`},
	}

	for _, test := range tests {
		out.Reset()
		el := test.el
		expect := test.expect

		if _, err := el.BuildNode(&out); err != nil {
			t.Errorf(`failed to write component string: %+v`, err)
			continue
		}

		if out.String() != expect {
			t.Errorf(`unexpected value: %+v != %+v`, out.String(), expect)
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
