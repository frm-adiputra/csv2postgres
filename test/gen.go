// +build ignore
// This program generates codes for data processing.
// It must be invoked by running go generate

package main

import (
	"fmt"
	"os"

	"github.com/frm-adiputra/csv2postgres"
)

func main() {
	generateDemo()
}

func generateDemo() {
	g := csv2postgres.Generator{
		BaseImportPath: "github.com/frm-adiputra/csv2postgres/test/generated",
		RootDir:        "generated",
	}
	err := g.Generate()
	if err != nil {
		exitWithError(err)
	}
}

func exitWithError(err error) {
	fmt.Fprintln(os.Stderr, err)
	os.Exit(1)
}
