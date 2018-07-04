package module

import (
	"fmt"
	"testing"
)

func TestLex(t *testing.T) {
	input := `
		module "my/thing"
		require "other/thing" v1.0.2
		require "new/thing" v2.3.4
		exclude "old/thing" v1.2.3
		require (
			"future/thing" v2.3.4
			"great/thing" v1.2.3
		)
		replace "bad/thing" v1.4.5 => "good/thing" v1.4.5
	`
	expects := []token{
		tokNewline(),

		tokModule(),
		tokString("my/thing"), tokNewline(),

		tokRequire(),
		tokString("other/thing"), tokVersion("v1.0.2"), tokNewline(),

		tokRequire(),
		tokString("new/thing"), tokVersion("v2.3.4"), tokNewline(),

		tokExclude(),
		tokString("old/thing"), tokVersion("v1.2.3"), tokNewline(),

		tokRequire(),
		tokLeftParen(), tokNewline(),
		tokString("future/thing"), tokVersion("v2.3.4"), tokNewline(),
		tokString("great/thing"), tokVersion("v1.2.3"), tokNewline(),
		tokRightParen(), tokNewline(),

		tokReplace(),
		tokString("bad/thing"), tokVersion("v1.4.5"),
		tokArrowFun(),
		tokString("good/thing"), tokVersion("v1.4.5"), tokNewline(),

		tokEOF(),
	}

	l := lexInString(input)
	for i, e := range expects {
		v := l.nextToken()
		if got, want := v, e; got != want {
			t.Error("got:", got, "want:", want, "i:", i)
		}
	}
}

func tokNewline() token {
	return token{kind: tokenNewline, val: "\n"}
}

func tokModule() token {
	return token{kind: tokenModule, val: "module"}
}

func tokRequire() token {
	return token{kind: tokenRequire, val: "require"}
}

func tokExclude() token {
	return token{kind: tokenExclude, val: "exclude"}
}

func tokReplace() token {
	return token{kind: tokenReplace, val: "replace"}
}

func tokArrowFun() token {
	return token{kind: tokenMapFun, val: "=>"}
}

func tokLeftParen() token {
	return token{kind: tokenLeftParen, val: "("}
}

func tokRightParen() token {
	return token{kind: tokenRightParen, val: ")"}
}

func tokString(s string) token {
	return token{kind: tokenString, val: fmt.Sprintf("%q", s)}
}

func tokVersion(s string) token {
	return token{kind: tokenVersion, val: s}
}

func tokEOF() token {
	return token{kind: tokenEOF, val: ""}
}
