package module

import (
	"testing"
)

func TestLex(t *testing.T) {
	input := `
		module my/thing
		go 1.12
		require other/thing v1.0.2
		require new/thing v2.3.4
		exclude old/thing v1.2.3
		require (
			future/thing v2.3.4
			great/thing v1.2.3
			indirect/thing v5.6.7 // indirect
		)
		replace bad/thing v1.4.5 => good/thing v1.4.5
		// refer: https://github.com/gin-gonic/gin/issues/1673
	`
	// input := `
	// 	// refer: https://github.com/gin-gonic/gin/issues/1673
	// `
	expects := []token{
		tokNewline(),

		tokModule(),
		tokNakedVal("my/thing"), tokNewline(),

		tokGo(),
		tokNakedVal("1.12"), tokNewline(),

		tokRequire(),
		tokNakedVal("other/thing"), tokNakedVal("v1.0.2"), tokNewline(),

		tokRequire(),
		tokNakedVal("new/thing"), tokNakedVal("v2.3.4"), tokNewline(),

		tokExclude(),
		tokNakedVal("old/thing"), tokNakedVal("v1.2.3"), tokNewline(),

		tokRequire(),
		tokLeftParen(), tokNewline(),
		tokNakedVal("future/thing"), tokNakedVal("v2.3.4"), tokNewline(),
		tokNakedVal("great/thing"), tokNakedVal("v1.2.3"), tokNewline(),
		tokNakedVal("indirect/thing"), tokNakedVal("v5.6.7"), tokComment(), tokIndirectComment(), tokNewline(),
		tokRightParen(), tokNewline(),

		tokReplace(),
		tokNakedVal("bad/thing"), tokNakedVal("v1.4.5"),
		tokArrowFun(),
		tokNakedVal("good/thing"), tokNakedVal("v1.4.5"), tokNewline(),

		tokComment(), tokNewline(),

		tokEOF(),
	}
	// expects := []token{
	// 	tokNewline(),
	// 	tokComment(), tokNewline(),
	// 	tokEOF(),
	// }

	l := lexInString(input)
	for i, e := range expects {
		v := l.nextToken()
		if got, want := v, e; got != want {
			t.Errorf("got %d %s", got.kind, got.val)
			t.Errorf("want %d %s", want.kind, want.val)
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

func tokGo() token {
	return token{kind: tokenGo, val: "go"}
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

func tokComment() token {
	return token{kind: tokenComment, val: "//"}
}

func tokIndirectComment() token {
	return token{kind: tokenIndirectComment, val: "indirect"}
}

func tokNakedVal(s string) token {
	return token{kind: tokenNakedVal, val: s}
}

func tokEOF() token {
	return token{kind: tokenEOF, val: ""}
}
