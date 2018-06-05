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

	l := lex(input)
	for i, e := range expects {
		v := l.nextToken()
		if got, want := v, e; got != want {
			t.Error("got:", got, "want:", want, "i:", i)
		}
	}
}

func tokNewline() token {
	return token{tokenNewline, "\n"}
}

func tokModule() token {
	return token{tokenModule, "module"}
}

func tokRequire() token {
	return token{tokenRequire, "require"}
}

func tokExclude() token {
	return token{tokenExclude, "exclude"}
}

func tokReplace() token {
	return token{tokenReplace, "replace"}
}

func tokArrowFun() token {
	return token{tokenArrowFunction, "=>"}
}

func tokLeftParen() token {
	return token{tokenLeftParenthese, "("}
}

func tokRightParen() token {
	return token{tokenRightParenthese, ")"}
}

func tokString(s string) token {
	return token{tokenString, fmt.Sprintf("%q", s)}
}

func tokVersion(s string) token {
	return token{tokenVersion, s}
}

func tokEOF() token {
	return token{tokenEOF, ""}
}
