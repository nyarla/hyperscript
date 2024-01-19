package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
)

var (
	help bool
	out  string
	lang string
	kind string
)

type dataSet struct {
	tags  []string
	attrs []struct {
		name string
		attr string
	}
}

func init() {
	flag.BoolVar(&help, `h`, false, `show help`)
	flag.BoolVar(&help, `help`, false, `show help`)
	flag.StringVar(&kind, `kind`, `tags`, `the kind of output from lang: [tags|attrs]`)
	flag.StringVar(&lang, `lang`, `html`, `the markup langauge: [html]`)
	flag.StringVar(&out, `out`, ``, `output path to auto-genereated code`)

	flag.Parse()
}

func main() {
	if help {
		flag.PrintDefaults()
		os.Exit(0)
	}

	if lang == `` {
		log.Fatal(`-lang is empty. please set -lang flag.`)
	}

	if kind == `` {
		log.Fatal(`-kind is emptu. please set -kind flag.`)
	}

	if kind != `tags` && kind != `attrs` {
		log.Fatalf(`unsupported -kind: %s`, kind)
	}

	var (
		w       strings.Builder
		dataset dataSet
	)

	switch lang {
	case `html`:
		dataset = htmlDataSet
	default:
		log.Fatalf(`unsupported -lang: %s`, lang)
	}

	printSharedHeader(&w, kind)
	printCode(&w, kind, dataset)

	if err := os.WriteFile(out, []byte(w.String()), 0644); err != nil {
		log.Fatal(err)
	}
}

func printSharedHeader(w *strings.Builder, pkg string) {
	fmt.Fprintf(w, `// Code generated by internal/cmd/autogen/main.go; DO NOT EDIT
package %s

    import (
      "github.com/nyarla/hyperscript/node"
    )

    type NodeBuilder = node.NodeBuilder
`, pkg)

}

func printCode(w *strings.Builder, kind string, dataset dataSet) {
	if kind == `tags` {
		for _, name := range dataset.tags {
			fmt.Fprintf(w, "func %s (contains ...NodeBuilder) NodeBuilder { return node.Element(`%s`, contains...) }\n", name, strings.ToLower(name))
		}

		return
	}

	if kind == `attrs` {
		for _, attrs := range dataset.attrs {
			name := attrs.name
			attr := attrs.attr

			if attr == `` {
				attr = strings.ToLower(name)
			}

			fmt.Fprintf(w, "func %s (value string) NodeBuilder { return node.Attr(`%s`, value) }\n", name, attr)
		}

		return
	}
}

var htmlDataSet = dataSet{
	// Ref. https://developer.mozilla.org/en-US/docs/Web/HTML/Element
	tags: []string{
		// Root Element
		`HTML`,

		// Document metadata
		`Base`,
		`Head`,
		`Link`,
		`Meta`,
		`Style`,
		`Title`,

		// Sectioning root
		`Body`,

		// Content sectioning
		`Address`,
		`Article`,
		`Aside`,
		`Footer`,
		`H1`,
		`H2`,
		`H3`,
		`H4`,
		`H5`,
		`H6`,
		`Header`,
		`Hgroup`,
		`Main`,
		`Nav`,
		`Search`,
		`Section`,

		// Text content
		`Blockquote`,
		`Dd`,
		`Div`,
		`Dl`,
		`Dt`,
		`Figcaption`,
		`Figure`,
		`Hr`,
		`Li`,
		`Menu`,
		`Ol`,
		`P`,
		`Pre`,
		`Ul`,

		// Inline text semantics
		`A`,
		`Abbr`,
		`B`,
		`Bdi`,
		`Bdo`,
		`Br`,
		`Cite`,
		`Code`,
		`Data`,
		`Dfn`,
		`Em`,
		`I`,
		`Kbd`,
		`Mark`,
		`Q`,
		`Rp`,
		`Rt`,
		`Ruby`,
		`S`,
		`Samp`,
		`Small`,
		`Span`,
		`Strong`,
		`Sub`,
		`Sup`,
		`Time`,
		`U`,
		`Var`,
		`Wbr`,

		// Image and multimedia
		`Area`,
		`Audio`,
		`Img`,
		`Map`,
		`Track`,
		`Video`,

		// Embedded content
		`Embed`,
		`Iframe`,
		`Object`,
		`Picture`,
		`Portal`,
		`Source`,

		// Svg and MathML
		`SVG`,
		`Math`,

		// Scripting
		`Canvas`,
		`Noscript`,
		`Script`,

		// Demarcating edits
		`Del`,
		`Ins`,

		// Table content
		`Caption`,
		`Col`,
		`Colgroup`,
		`Table`,
		`Tbody`,
		`Td`,
		`Tfoot`,
		`Th`,
		`Thead`,
		`Tr`,

		// Forms
		`Button`,
		`Datalist`,
		`Fieldset`,
		`Form`,
		`Input`,
		`Label`,
		`Legend`,
		`Meter`,
		`Optgroup`,
		`Option`,
		`Output`,
		`Progress`,
		`Select`,
		`Textarea`,

		// Intractive elements
		`Details`,
		`Dialog`,
		`Summary`,

		// Web Components
		`Slot`,
		`Template`,
	},
	// Ref. https://developer.mozilla.org/en-US/docs/Web/HTML/Attributes
	attrs: []struct {
		name string
		attr string
	}{
		{`Accpet`, ``},
		{`AcceptCharset`, `accept-charset`},
		{`AccessKey`, ``},
		{`Action`, ``},
		{`Allow`, ``},
		{`Alt`, ``},
		{`Async`, ``},
		{`AutoCapitalize`, ``},
		{`AutoComplete`, ``},
		{`AutoPlayer`, ``},
		{`Bufferd`, ``},
		{`Capture`, ``},
		{`Charset`, ``},
		{`Checked`, ``},
		{`Cite`, ``},
		{`Cols`, ``},
		{`Colspan`, ``},
		{`Content`, ``},
		{`ContentEditable`, ``},
		{`Controls`, ``},
		{`Coords`, ``},
		{`CrossOrigin`, ``},
		{`Data`, ``},
		{`DateTime`, ``},
		{`Decoding`, ``},
		{`Default`, ``},
		{`Defer`, ``},
		{`Dir`, ``},
		{`Dirname`, ``},
		{`Disabled`, ``},
		{`Download`, ``},
		{`Draggable`, ``},
		{`Enctype`, ``},
		{`For`, ``},
		{`Form`, ``},
		{`FormAction`, ``},
		{`FormEncType`, ``},
		{`FormMethod`, ``},
		{`FormNoValdiate`, ``},
		{`FormTarget`, ``},
		{`Headers`, ``},
		{`Height`, ``},
		{`Hidden`, ``},
		{`High`, ``},
		{`HrefLang`, ``},
		{`HttpEquiv`, ``},
		{`Id`, ``},
		{`Integrity`, ``},
		{`InputMode`, ``},
		{`IsMap`, ``},
		{`Itemprop`, ``},
		{`Kind`, ``},
		{`Label`, ``},
		{`Lang`, ``},
		{`List`, ``},
		{`Loop`, ``},
		{`Low`, ``},
		{`Max`, ``},
		{`MaxLength`, ``},
		{`MinLength`, ``},
		{`Media`, ``},
		{`Method`, ``},
		{`Multiple`, ``},
		{`Muted`, ``},
		{`Name`, ``},
		{`NoValidate`, ``},
		{`Open`, ``},
		{`Pattern`, ``},
		{`Ping`, ``},
		{`PlaceHolder`, ``},
		{`PlaysInline`, ``},
		{`Poster`, ``},
		{`Preload`, ``},
		{`Readonly`, ``},
		{`RefererPolicy`, ``},
		{`Rel`, ``},
		{`Required`, ``},
		{`Reversed`, ``},
		{`Role`, ``},
		{`Sandbox`, ``},
		{`Scope`, ``},
		{`Selected`, ``},
		{`Shape`, ``},
		{`Size`, ``},
		{`Sizes`, ``},
		{`Slot`, ``},
		{`Span`, ``},
		{`SpellCheck`, ``},
		{`Src`, ``},
		{`SrcDoc`, ``},
		{`SrcLang`, ``},
		{`SrcSet`, ``},
		{`Start`, ``},
		{`Step`, ``},
		{`Style`, ``},
		{`TabIndex`, ``},
		{`Target`, ``},
		{`Title`, ``},
		{`Translate`, ``},
		{`Type`, ``},
		{`UseMap`, ``},
		{`Value`, ``},
		{`Width`, ``},
		{`Wrap`, ``},
	},
}
