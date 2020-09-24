package interpolation

import "github.com/frm-adiputra/csv2postgres/schema"

// ViewData represents data for view
type ViewData struct {
	*schema.View
	CreateDeps             []dependencyData
	DropDeps               []dependencyData
	CreateDepsIncludeTable bool
	DropDepsIncludeTable   bool
}
