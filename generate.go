package csv2postgres

import (
	"github.com/frm-adiputra/csv2postgres/generator"
)

// Generate generates source codes based on spec.
// This function must be called from client gen.go file.
func Generate(baseImportPath string) error {
	err := generator.Generate(baseImportPath)
	if err != nil {
		return err
	}
	return nil
}
