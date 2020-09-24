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
	g := csv2postgres.Generator{
		BaseImportPath: "github.com/frm-adiputra/csv2postgres/test/generated",
		RootDir:        "generated",
	}
	if err := g.Generate(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
