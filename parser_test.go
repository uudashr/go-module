package module_test

import (
	"reflect"
	"testing"

	module "github.com/uudashr/go-module"
)

func TestParse(t *testing.T) {
	in := `
		module my/thing
		require other/thing v1.0.2
		require (
			new/thing v2.3.4
			other/new/thing v1.2.3
			indirect/thing v5.6.7 // indirect
		)
		exclude old/thing v1.2.3
		exclude (
			bad/thing v1.0.0
			new/bad/thing v2.2.3
		)
		replace bad/thing v1.4.5 => good/thing v1.4.5
		replace (
			bad/thing v1.0.0 => good/thing v1.0.0
			new/bad/thing v2.2.3 => new/good/thing v2.2.3
		)
		// refer: https://github.com/gin-gonic/gin/issues/1673
	`

	expectReqs := []module.Package{
		{Path: "other/thing", Version: "v1.0.2"},
		{Path: "new/thing", Version: "v2.3.4"},
		{Path: "other/new/thing", Version: "v1.2.3"},
		{Path: "indirect/thing", Version: "v5.6.7", Indirect: true},
	}

	expectExcl := []module.Package{
		{Path: "old/thing", Version: "v1.2.3"},
		{Path: "bad/thing", Version: "v1.0.0"},
		{Path: "new/bad/thing", Version: "v2.2.3"},
	}

	expectRepl := []module.PackageMap{
		{
			From: module.Package{Path: "bad/thing", Version: "v1.4.5"},
			To:   module.Package{Path: "good/thing", Version: "v1.4.5"},
		},
		{
			From: module.Package{Path: "bad/thing", Version: "v1.0.0"},
			To:   module.Package{Path: "good/thing", Version: "v1.0.0"},
		},
		{
			From: module.Package{Path: "new/bad/thing", Version: "v2.2.3"},
			To:   module.Package{Path: "new/good/thing", Version: "v2.2.3"},
		},
	}

	m, err := module.ParseInString(in)
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
