package module

import (
	"errors"
	"fmt"
)

// File represents the mod file.
type File struct {
	Module   string
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
func Parse(input string) (*File, error) {
	f := &File{}
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
	file  *File
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
	t := p.nextToken()
	if t.kind != tokenModule {
		return p.errorf("expect module declaration")
	}

	return parseModuleName
}

func parseModuleName(p *parser) parseFn {
	t := p.nextToken()
	if t.kind != tokenString {
		return p.errorf("expect module name")
	}

	p.file.Module = unquote(t.val)
	return parseVerb
}

func parseVerb(p *parser) parseFn {
	switch p.nextToken().kind {
	case tokenRequire:
		return parsePkgList(p.requirePkg)
	case tokenExclude:
		return parsePkgList(p.excludePkg)
	case tokenReplace:
		return parsePkgMapList(p.replacePkg)
	case tokenEOF:
		return nil
	default:
		return p.errorf("expect verb declaration")
	}
}

func parsePkgList(add func(pkg Package)) parseFn {
	return func(p *parser) parseFn {
		t := p.nextToken()
		if t.kind == tokenLeftParenthese {
			return parsePkgListElem(add)
		}

		pkg, err := readPkg(t, p)
		if err != nil {
			return p.error(err)
		}

		add(*pkg)
		return parseVerb
	}
}

func parsePkgListElem(add func(pkg Package)) parseFn {
	return func(p *parser) parseFn {
		t := p.nextToken()
		if t.kind == tokenRightParenthese {
			return parseVerb
		}

		pkg, err := readPkg(t, p)
		if err != nil {
			return p.error(err)
		}

		add(*pkg)
		return parsePkgListElem(add)
	}
}

func parsePkgMapList(add func(m PackageMap)) parseFn {
	return func(p *parser) parseFn {
		t := p.nextToken()
		if t.kind == tokenLeftParenthese {
			return parsePkgMapListElem(add)
		}

		pkgMap, err := readPkgMap(t, p)
		if err != nil {
			return p.error(err)
		}

		add(*pkgMap)
		return parseVerb
	}
}

func parsePkgMapListElem(add func(m PackageMap)) parseFn {
	return func(p *parser) parseFn {
		t := p.nextToken()
		if t.kind == tokenRightParenthese {
			return parseVerb
		}

		pkgMap, err := readPkgMap(t, p)
		if err != nil {
			return p.error(err)
		}

		add(*pkgMap)
		return parsePkgMapListElem(add)
	}
}

func readPkg(t token, p *parser) (*Package, error) {
	if t.kind != tokenString {
		return nil, errors.New("expect package declaration")
	}

	path := unquote(t.val)

	t = p.nextToken()
	if t.kind != tokenVersion {
		return nil, errors.New("expect package version")
	}

	return &Package{path, t.val}, nil
}

func readPkgMap(t token, p *parser) (*PackageMap, error) {
	old, err := readPkg(t, p)
	if err != nil {
		return nil, err
	}

	if p.nextToken().kind != tokenArrowFunction {
		return nil, errors.New("expect '=>'")
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
