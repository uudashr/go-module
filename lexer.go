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
	tokenError tokenKind = iota
	tokenEOF

	// operators
	tokenArrowFunction // =>

	// delimiters
	tokenLeftParenthese  // (
	tokenRightParenthese // )

	// literals
	tokenString // with double quote
	tokenVersion

	// keywords
	tokenModule
	tokenRequire
	tokenExclude
	tokenReplace
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
		case t := <-l.tokens:
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

func (l *lexer) peek() rune {
	r := l.next()
	l.backup()
	return r
}

func (l *lexer) accept(valid string) bool {
	if strings.ContainsRune(valid, l.next()) {
		return true
	}

	l.backup()
	return false
}

func (l *lexer) acceptRun(valid string) {
	for strings.ContainsRune(valid, l.next()) {
		// read another
	}
	l.backup()
}

func (l *lexer) acceptWhile(fn func(r rune) bool) {
	for fn(l.next()) {
		// read another
	}
	l.backup()
}

func (l *lexer) emit(t tokenKind) {
	i := token{t, l.input[l.start:l.pos]}
	l.tokens <- i
	l.width = l.pos
}

func (l *lexer) errorf(format string, args ...interface{}) lexFn {
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
			// ignore all whitespace
			l.ignore()
		case r == 'v':
			return lexVersion
		case r == '"':
			return lexString
		case r == '(':
			l.emit(tokenLeftParenthese)
			return lexFile
		case r == ')':
			l.emit(tokenRightParenthese)
			return lexFile
		case r == '=':
			if l.next() != '>' {
				return l.errorf("expect =>")
			}
			l.emit(tokenArrowFunction)
			return lexFile
		case isAlphaLower(r):
			return lexKeyword
		case r == eof:
			l.emit(tokenEOF)
			return nil
		default:
			return l.errorf("expecting valid keyword while lexFile, got %q", string(r))
		}
	}
}

func lexKeyword(l *lexer) lexFn {
	for {
		switch r := l.next(); {
		case isAlphaLower(r):
			// absorb
		default:
			if r != ' ' {
				return l.errorf("unexpected character %q while lexKeyword", string(r))
			}

			l.backup()
			word := l.input[l.start:l.pos]
			t, ok := key[word]
			if !ok {
				return l.errorf("invalid keyword %q while lexKeyword", word)
			}

			l.emit(t)
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
		case r == eof || isEOL(r):
			return l.errorf("unterminated string, got %s", string(r))
		default:
			// absorp
		}
	}
}

func lexVersion(l *lexer) lexFn {
	for {
		switch r := l.next(); {
		case r == ' ' || r == eof || isEOL(r):
			l.backup()
			l.emit(tokenVersion)
			return lexFile
		default:
			// absorb
		}
	}
}

func isEOL(r rune) bool {
	return strings.ContainsRune("\n\r", r)
}

func isWhiteSpace(r rune) bool {
	return unicode.IsSpace(r)
}

func isAlphaLower(r rune) bool {
	return unicode.IsLetter(r) && unicode.IsLower(r)
}
