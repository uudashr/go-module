package module_test

import (
	"fmt"

	module "github.com/uudashr/go-module"
)

func ExampleParse() {
	in := `
		module "my/thing"
		require "other/thing" v1.0.2
		require "new/thing" v2.3.4
		exclude "old/thing" v1.2.3
		replace "bad/thing" v1.4.5 => "good/thing" v1.4.5
	`
	m, err := module.ParseInString(in)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Module %q\n", m.Name)
	for _, pkg := range m.Requires {
		fmt.Printf("Require %q %s\n", pkg.Path, pkg.Version)
	}

	for _, pkg := range m.Excludes {
		fmt.Printf("Exclude %q %s\n", pkg.Path, pkg.Version)
	}

	for _, rep := range m.Replaces {
		from, to := rep.From, rep.To
		fmt.Printf("Replace %q %s with %q %s\n", from.Path, from.Version, to.Path, to.Version)
	}

	// Output:
	// Module "my/thing"
	// Require "other/thing" v1.0.2
	// Require "new/thing" v2.3.4
	// Exclude "old/thing" v1.2.3
	// Replace "bad/thing" v1.4.5 with "good/thing" v1.4.5
}
