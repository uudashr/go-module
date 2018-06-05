package module

import (
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
		token{tokenModule, "module"},
		token{tokenString, `"my/thing"`},

		token{tokenRequire, "require"},
		token{tokenString, `"other/thing"`},
		token{tokenVersion, "v1.0.2"},

		token{tokenRequire, "require"},
		token{tokenString, `"new/thing"`},
		token{tokenVersion, "v2.3.4"},

		token{tokenExclude, "exclude"},
		token{tokenString, `"old/thing"`},
		token{tokenVersion, "v1.2.3"},

		token{tokenRequire, "require"},
		token{tokenLeftParenthese, "("},
		token{tokenString, `"future/thing"`}, token{tokenVersion, "v2.3.4"},
		token{tokenString, `"great/thing"`}, token{tokenVersion, "v1.2.3"},
		token{tokenRightParenthese, ")"},

		token{tokenReplace, "replace"},
		token{tokenString, `"bad/thing"`}, token{tokenVersion, "v1.4.5"},
		token{tokenArrowFunction, "=>"},
		token{tokenString, `"good/thing"`}, token{tokenVersion, "v1.4.5"},
		token{tokenEOF, ""},
	}

	l := lex(input)
	for _, e := range expects {
		v := l.nextToken()
		if got, want := v, e; got != want {
			t.Error("got:", got, "want:", want)
		}
	}

}
