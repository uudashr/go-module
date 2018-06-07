package module

import (
	"fmt"
	"strings"
	"unicode"
	"unicode/utf8"
)

const eof = rune(-1)

type tokenKind int

const (
	// special kind
	tokenError tokenKind = iota // error
	tokenEOF                    // end of file

	// operators
	tokenMapFun // "=>"

	// delimiters
	tokenLeftParen  // "("
	tokenRightParen // ")"
	tokenNewline    // "\n"

	// literals
	tokenString  // wrapped by double quote
	tokenVersion // semver prefixed with "v" ex: v2.1.6

	// keywords
	tokenModule  // module
	tokenRequire // require
	tokenExclude // exclude
	tokenReplace // replace
)

var key = map[string]tokenKind{
	"module":  tokenModule,
	"require": tokenRequire,
	"exclude": tokenExclude,
	"replace": tokenReplace,
}

type token struct {
	kind tokenKind
	val  string
}

func (t token) String() string {
	switch t.kind {
	case tokenEOF:
		return "EOF"
	case tokenError:
		return t.val
	case tokenNewline:
		return "newline"
	}

	if len(t.val) > 10 {
		return fmt.Sprintf("%.10q...", t.val)
	}
	return fmt.Sprintf("%q", t.val)
}

type lexFn func(l *lexer) lexFn

type lexer struct {
	start  int        // start position of the token
	pos    int        // current read position of the input
	width  int        // width of the last runes read from the input
	input  string     // the input string being scanned
	tokens chan token // the scanned tokens
	state  lexFn      // the current state of lexer
}

func lex(input string) *lexer {
	l := &lexer{
		input:  input,
		tokens: make(chan token, 2),
		state:  lexFile,
	}
	return l
}

func (l *lexer) nextToken() token {
	for {
		select {
		case t, ok := <-l.tokens:
			if !ok {
				return token{tokenError, "no more token"}
			}

			if t.kind == tokenEOF {
				close(l.tokens)
			}
			return t
		default:
			l.state = l.state(l)
		}
	}
}

func (l *lexer) next() (r rune) {
	if l.pos >= len(l.input) {
		l.width = 0
		return eof
	}

	r, l.width = utf8.DecodeRuneInString(l.input[l.pos:])
	l.pos += l.width
	return r
}

func (l *lexer) backup() {
	l.pos -= l.width
}

func (l *lexer) val() string {
	return l.input[l.start:l.pos]
}

func (l *lexer) emit(kind tokenKind) {
	i := token{kind, l.val()}
	l.tokens <- i
	l.start = l.pos
}

func (l *lexer) emitErrorf(format string, args ...interface{}) lexFn {
	l.tokens <- token{tokenError, fmt.Sprintf(format, args...)}
	return nil
}

func (l *lexer) ignore() {
	l.start = l.pos
}

func lexFile(l *lexer) lexFn {
	for {
		switch r := l.next(); {
		case isWhiteSpace(r):
			l.ignore()
		case r == '\n':
			l.emit(tokenNewline)
			return lexFile
		case r == 'v':
			return lexVersion
		case r == '"':
			return lexString
		case r == '(':
			l.emit(tokenLeftParen)
			return lexFile
		case r == ')':
			l.emit(tokenRightParen)
			return lexFile
		case r == '=':
			if l.next() != '>' {
				return l.emitErrorf("expect => got %q", string(r))
			}

			l.emit(tokenMapFun)
			return lexFile
		case isAlphaLower(r):
			return lexKeyword
		case r == eof:
			l.ignore()
			l.emit(tokenEOF)
			return nil
		default:
			return l.emitErrorf("expecting valid keyword while lexFile, got %q", string(r))
		}
	}
}

func lexKeyword(l *lexer) lexFn {
	for {
		switch r := l.next(); {
		case isAlphaLower(r):
			// absorb
		default:
			l.backup()
			word := l.val()
			kind, ok := key[word]
			if !ok {
				return l.emitErrorf("invalid keyword %q while lexKeyword", word)
			}

			l.emit(kind)
			return lexFile
		}
	}
}

func lexString(l *lexer) lexFn {
	for {
		switch r := l.next(); {
		case r == '"':
			l.emit(tokenString)
			return lexFile
		case r == '\n', r == eof:
			return l.emitErrorf("unterminated string, got %s", string(r))
		case r == '\\':
			r = l.next()
			if !(r == 't' || r == '\\') {
				return l.emitErrorf(`invalid escape char \%s`, string(r))
			}
			fallthrough
		default:
			// absorp
		}
	}
}

func lexVersion(l *lexer) lexFn {
	for {
		switch l.next() {
		case ' ', '\n', eof:
			l.backup()
			l.emit(tokenVersion)
			return lexFile
		default:
			// absorb
		}
	}
}

func isWhiteSpace(r rune) bool {
	return strings.ContainsRune(" \t", r)
}

func isAlphaLower(r rune) bool {
	return unicode.IsLetter(r) && unicode.IsLower(r)
}
