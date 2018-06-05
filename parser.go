package module

import (
	"fmt"
)

// Module represents the mod file.
type Module struct {
	Name     string
	Requires []Package
	Excludes []Package
	Replaces []PackageMap
}

// PackageMap package mapping defintion.
type PackageMap struct {
	From Package
	To   Package
}

// Package represents the package info.
type Package struct {
	Path    string
	Version string
}

// Parse mod file.
func Parse(input string) (*Module, error) {
	f := &Module{}
	l := lex(input)
	p := &parser{lexer: l, file: f}

	for state := parseModule; state != nil; {
		state = state(p)
	}

	if p.err != nil {
		return nil, p.err
	}

	return f, nil
}

type parser struct {
	lexer *lexer
	file  *Module
	err   error
}

func (p *parser) nextToken() token {
	return p.lexer.nextToken()
}

func (p *parser) error(err error) parseFn {
	p.err = err
	return nil
}

func (p *parser) errorf(format string, args ...interface{}) parseFn {
	return p.error(fmt.Errorf(format, args...))
}

func (p *parser) requirePkg(pkg Package) {
	p.file.Requires = append(p.file.Requires, pkg)
}

func (p *parser) excludePkg(pkg Package) {
	p.file.Excludes = append(p.file.Excludes, pkg)
}

func (p *parser) replacePkg(m PackageMap) {
	p.file.Replaces = append(p.file.Replaces, m)
}

type parseFn func(p *parser) parseFn

func parseModule(p *parser) parseFn {
Loop:
	for {
		switch t := p.nextToken(); t.kind {
		case tokenNewline:
			// skip
		case tokenModule:
			break Loop
		default:
			return p.errorf("expect module declaration, got %s", t)
		}
	}

	return parseModuleName
}

func parseModuleName(p *parser) parseFn {
	t := p.nextToken()
	if t.kind != tokenString {
		return p.errorf("expect module name, got %s", t)
	}

	p.file.Name = unquote(t.val)

	if t = p.nextToken(); t.kind != tokenNewline {
		return p.errorf("expect newline, got %s", t)
	}
	return parseVerb
}

func parseVerb(p *parser) parseFn {
	switch t := p.nextToken(); t.kind {
	case tokenRequire:
		return parsePkgList(p.requirePkg)
	case tokenExclude:
		return parsePkgList(p.excludePkg)
	case tokenReplace:
		return parsePkgMapList(p.replacePkg)
	case tokenNewline:
		// ignore
		return parseVerb
	case tokenEOF:
		return nil
	default:
		return p.errorf("expect verb declaration, got %s", t)
	}
}

func parsePkgList(add func(pkg Package)) parseFn {
	return func(p *parser) parseFn {
		t := p.nextToken()
		if t.kind == tokenLeftParenthese {
			if t = p.nextToken(); t.kind != tokenNewline {
				return p.errorf("expect newline, got %s", t)
			}

			return parsePkgListElem(add)
		}

		pkg, err := readPkg(t, p)
		if err != nil {
			return p.error(err)
		}

		if t = p.nextToken(); t.kind != tokenNewline {
			return p.errorf("expect newline, got %s", t)
		}

		add(*pkg)
		return parseVerb
	}
}

func parsePkgListElem(add func(pkg Package)) parseFn {
	return func(p *parser) parseFn {
		t := p.nextToken()
		if t.kind == tokenRightParenthese {
			if t = p.nextToken(); t.kind != tokenNewline {
				return p.errorf("expect newline, got %s", t)
			}

			return parseVerb
		}

		pkg, err := readPkg(t, p)
		if err != nil {
			return p.error(err)
		}

		if t = p.nextToken(); t.kind != tokenNewline {
			return p.errorf("expect newline, got %s", t)
		}

		add(*pkg)
		return parsePkgListElem(add)
	}
}

func parsePkgMapList(add func(m PackageMap)) parseFn {
	return func(p *parser) parseFn {
		t := p.nextToken()
		if t.kind == tokenLeftParenthese {
			if t = p.nextToken(); t.kind != tokenNewline {
				return p.errorf("expect newline, got %s", t)
			}

			return parsePkgMapListElem(add)
		}

		pkgMap, err := readPkgMap(t, p)
		if err != nil {
			return p.error(err)
		}

		if t = p.nextToken(); t.kind != tokenNewline {
			return p.errorf("expect newline, got %s", t)
		}

		add(*pkgMap)
		return parseVerb
	}
}

func parsePkgMapListElem(add func(m PackageMap)) parseFn {
	return func(p *parser) parseFn {
		t := p.nextToken()
		if t.kind == tokenRightParenthese {
			if t = p.nextToken(); t.kind != tokenNewline {
				return p.errorf("expect newline, got %s", t)
			}

			return parseVerb
		}

		pkgMap, err := readPkgMap(t, p)
		if err != nil {
			return p.error(err)
		}

		if t = p.nextToken(); t.kind != tokenNewline {
			return p.errorf("expect newline, got %s", t)
		}

		add(*pkgMap)
		return parsePkgMapListElem(add)
	}
}

func readPkg(t token, p *parser) (*Package, error) {
	if t.kind != tokenString {
		return nil, fmt.Errorf("expect package declaration, got %s", t)
	}

	path := unquote(t.val)

	if t = p.nextToken(); t.kind != tokenVersion {
		return nil, fmt.Errorf("expect package version, got %s", t)
	}

	return &Package{path, t.val}, nil
}

func readPkgMap(t token, p *parser) (*PackageMap, error) {
	old, err := readPkg(t, p)
	if err != nil {
		return nil, err
	}

	if t := p.nextToken(); t.kind != tokenArrowFunction {
		return nil, fmt.Errorf("expect '=>', got %s", t)
	}

	new, err := readPkg(p.nextToken(), p)
	if err != nil {
		return nil, err
	}

	return &PackageMap{*old, *new}, nil
}

func unquote(s string) string {
	return s[1 : len(s)-1]
}

type tokenizer interface {
	nextToken() token
}
