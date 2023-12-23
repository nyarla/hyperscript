package node

import "testing"

func TestText(t *testing.T) {
	tests := [][2]string{
		{string(Raw(`hello`)), `hello`},
		{string(Text(`hello&workd`)), `hello&amp;workd`},
	}

	for _, test := range tests {
		if test[0] != test[1] {
			t.Errorf(`compare values mismatch: %v != %v`, test[0], test[1])
		}
	}
}

func TestAttr(t *testing.T) {
	tests := [][2]string{
		{string(Flag(`crossorigin`)), `crossorigin`},
		{string(Flag(`cross&origin`)), `cross&amp;origin`},

		{string(Pair(`id`, `foo`)), `id="foo"`},
		{string(Pair(`i&d`, `f&o`)), `i&amp;d="f&amp;o"`},
	}

	for _, test := range tests {
		if test[0] != test[1] {
			t.Errorf(`compare values mismatch: %v != %v`, test[0], test[1])
		}
	}
}

func TestNode(t *testing.T) {
	tests := [][2]string{
		{string(Node(`hr`)), `<hr />`},
		{string(Node(`hr`, Pair(`id`, `foo`))), `<hr id="foo" />`},
		{string(Node(`p`, Text(`hello`))), `<p >hello</p>`},
		{string(Node(`p`, Pair(`class`, `msg`), Text(`hi,`))), `<p class="msg">hi,</p>`},
	}

	for _, test := range tests {
		if test[0] != test[1] {
			t.Errorf(`compare values mismatch: %v != %v`, test[0], test[1])
		}
	}
}
