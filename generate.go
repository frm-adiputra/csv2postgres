package csv2postgres

import (
	"github.com/frm-adiputra/csv2postgres/generator"
)

// Generate generates source codes based on spec.
// This function must be called from client gen.go file.
func Generate(baseImportPath, rootDir string) error {
	g := generator.Generator{
		BaseImportPath: baseImportPath,
		RootDir:        rootDir,
	}

	return g.Generate()
}
