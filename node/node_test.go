package node

import (
	"html"
	"strings"
	"testing"
)

func TestRaw(t *testing.T) {
	var out strings.Builder
	var component = Raw(`text`)
	var expect = `text`

	if _, err := component.WriteNode(&out); err != nil {
		t.Errorf(`failed to write component string: %+v`, err)
		return
	}

	if out.String() != expect {
		t.Errorf(`unexpected value: %+v => %+v != %+v`, `text`, out.String(), expect)
	}
}

func TestText(t *testing.T) {
	var out strings.Builder
	var tests = [][2]string{
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

		raw := test[1]
		escaped := test[0]

		component := Text(raw)

		if _, err := component.WriteNode(&out); err != nil {
			t.Errorf(`failed to write component string: %+v`, err)
			continue
		}

		if out.String() != escaped {
			t.Errorf(`unexpected value: %+v => %+v != %+v`, raw, out.String(), escaped)
			continue
		}

		value := html.UnescapeString(escaped)
		if value != raw {
			t.Errorf(`escape string mismatch: %+v => %+v !=  %+v`, escaped, raw, value)
			continue
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

		_, err := component.WriteNode(&out)
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
		_, err := component.WriteNode(&out)
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
		el     NodeWriter
		expect string
	}{
		{Element(`hr`), `<hr />`},
		{Element(`h&r`), `<h&amp;r />`},
		{Element(`hr`, Attr(`id`, `msg`)), `<hr id="msg" />`},
		{Element(`hr`, Attr(`id`, `msg`), Attr(`class`, `hr`)), `<hr class="hr" id="msg" />`},
		{Element(`p`, Text(`hello`), Text(`, `), Text(`world!`)), `<p>hello, world!</p>`},
		{Element(`p&p`, Text(`hello`), Text(`, `), Text(`world!`)), `<p&amp;p>hello, world!</p&amp;p>`},
		{Element(`p`, Attr(`id`, `msg`), Attr(`class`, `highlight`), Text(`hello`), Text(`, `), Text(`world!`)), `<p class="highlight" id="msg">hello, world!</p>`},
		{Element(`p`, Element(`strong`, Text(`hi,`))), `<p><strong>hi,</strong></p>`},
	}

	for _, test := range tests {
		out.Reset()
		el := test.el
		expect := test.expect

		if _, err := el.WriteNode(&out); err != nil {
			t.Errorf(`failed to write component string: %+v`, err)
			continue
		}

		if out.String() != expect {
			t.Errorf(`unexpected value: %+v != %+v`, out.String(), expect)
		}
	}
}

func BenchmarkRaw(b *testing.B) {
	var out strings.Builder
	var component = Raw(`test`)

	b.ReportAllocs()
	b.ResetTimer()

	for idx := 0; idx < b.N; idx++ {
		component.WriteNode(&out)
	}
}

func BenchmarkText(b *testing.B) {
	var out strings.Builder
	var component = Text(`test`)

	b.ReportAllocs()
	b.ResetTimer()

	for idx := 0; idx < b.N; idx++ {
		component.WriteNode(&out)
	}
}

func BenchmarkAttr(b *testing.B) {
	var out strings.Builder
	var single = Attr(`crossorigin`)
	var pair = Attr(`id`, `msg`)

	b.ReportAllocs()
	b.ResetTimer()

	for idx := 0; idx < b.N; idx++ {
		single.WriteNode(&out)
		pair.WriteNode(&out)
	}
}

func BenchmarkElement(b *testing.B) {
	var out strings.Builder
	var p = Element(`p`, Element(`strong`, Text(`hi, `), Text(`this is test message.`)))

	b.ReportAllocs()
	b.ResetTimer()

	for idx := 0; idx < b.N; idx++ {
		p.WriteNode(&out)
	}
}

func BenchmarkRealCase(b *testing.B) {
	var out strings.Builder
	var t = Element(`html`,
		Element(`body`,
			Element(`h1`, Text(`title`)),
			Element(`p`, Text(`This is an example page!`)),
			Element(`ul`,
				Element(`li`, Text(`foo`),
					Element(`li`, Text(`bar`),
						Element(`li`, Text(`baz`)))))))

	b.ReportAllocs()
	b.ResetTimer()

	for idx := 0; idx < b.N; idx++ {
		t.WriteNode(&out)
	}
}
