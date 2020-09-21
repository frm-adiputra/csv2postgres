package interpolation

import "github.com/frm-adiputra/csv2postgres/schema"

type Interpolator struct {
	BaseImportPath string
	TableSpecs     []*schema.Table
	ViewSpecs      []*schema.View
}
