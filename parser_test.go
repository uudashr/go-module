package module_test

import (
	"reflect"
	"testing"

	module "github.com/uudashr/go-module"
)

func TestParse(t *testing.T) {
	in := `
		module "my/thing"
		require "other/thing" v1.0.2
		require (
			"new/thing" v2.3.4
			"other/new/thing" v1.2.3
		)
		exclude "old/thing" v1.2.3
		exclude (
			"bad/thing" v1.0.0
			"new/bad/thing" v2.2.3
		)
		replace "bad/thing" v1.4.5 => "good/thing" v1.4.5
		replace (
			"bad/thing" v1.0.0 => "good/thing" v1.0.0
			"new/bad/thing" v2.2.3 => "new/good/thing" v2.2.3
		)
	`

	expectReqs := []module.Package{
		module.Package{Path: "other/thing", Version: "v1.0.2"},
		module.Package{Path: "new/thing", Version: "v2.3.4"},
		module.Package{Path: "other/new/thing", Version: "v1.2.3"},
	}

	expectExcl := []module.Package{
		module.Package{Path: "old/thing", Version: "v1.2.3"},
		module.Package{Path: "bad/thing", Version: "v1.0.0"},
		module.Package{Path: "new/bad/thing", Version: "v2.2.3"},
	}

	expectRepl := []module.PackageMap{
		module.PackageMap{
			From: module.Package{Path: "bad/thing", Version: "v1.4.5"},
			To:   module.Package{Path: "good/thing", Version: "v1.4.5"},
		},
		module.PackageMap{
			From: module.Package{Path: "bad/thing", Version: "v1.0.0"},
			To:   module.Package{Path: "good/thing", Version: "v1.0.0"},
		},
		module.PackageMap{
			From: module.Package{Path: "new/bad/thing", Version: "v2.2.3"},
			To:   module.Package{Path: "new/good/thing", Version: "v2.2.3"},
		},
	}

	m, err := module.Parse(in)
	if err != nil {
		t.Fatal(err)
	}

	if got, want := m.Requires, expectReqs; !reflect.DeepEqual(got, want) {
		t.Error("got:", got, "want:", want)
	}

	if got, want := m.Excludes, expectExcl; !reflect.DeepEqual(got, want) {
		t.Error("got:", got, "want:", want)
	}

	if got, want := m.Replaces, expectRepl; !reflect.DeepEqual(got, want) {
		t.Error("got:", got, "want:", want)
	}
}

func assertRequirePackage(t *testing.T, m *module.Module, pkgPath, pkgVer string) {
	t.Helper()
	for _, req := range m.Requires {
		if req.Path == pkgPath && req.Version == pkgVer {
			return
		}
	}
	t.Errorf("expect package %q %s", pkgPath, pkgVer)
}
