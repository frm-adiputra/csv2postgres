package interpolation

import (
	"fmt"
	"path"
	"strings"

	"github.com/frm-adiputra/csv2postgres/common"
	"github.com/frm-adiputra/csv2postgres/schema"
)

// TableData represents interpolation result for table
type TableData struct {
	*common.Names

	SpecFile   string
	ImportPath string
	PkgVar     string

	DataSource     string
	CSVSeparator   string
	ComputePackage string
	ComputePkgVar  string

	Fields         []*FieldData
	ComputedFields []*ComputedFieldData
	ComputeFns     []*ComputeFnData
	Constraints    []string

	RequireSQLPkg     bool
	RequireStrconvPkg bool
	RequireTimePkg    bool
	HasValidation     bool
	HasComputed       bool

	DependsOn              []string
	CreateDeps             []dependencyData
	DropDeps               []dependencyData
	CreateDepsIncludeTable bool
	DropDepsIncludeTable   bool
}

// FieldData represents interpolation result for table's field
type FieldData struct {
	*schema.Field
	GoType string
}

// ComputedFieldData represents interpolation result for table's compute field
type ComputedFieldData struct {
	*schema.ComputedField
	GoType string
}

// ComputeFnData represents interpolation result for table's compute function
type ComputeFnData struct {
	Name         string
	ArgumentType string
	ReturnType   string
}

func newTableData(ts *schema.Table, baseImportPath, rootDir string) (*TableData, error) {
	errFmt := "fail processing '%s': %w"

	fields, err := newFieldsData(ts.Fields)
	if err != nil {
		return nil, fmt.Errorf(errFmt, ts.SpecFile, err)
	}

	computedFields, err := newComputedFieldsData(ts.ComputedFields)
	if err != nil {
		return nil, fmt.Errorf(errFmt, ts.SpecFile, err)
	}

	computeFns, err := newComputeFnsData(fields, computedFields)
	if err != nil {
		return nil, fmt.Errorf(errFmt, ts.SpecFile, err)
	}

	return &TableData{
		Names: ts.Names,

		SpecFile:   ts.SpecFile,
		ImportPath: path.Join(baseImportPath, "internal", strings.ToLower(ts.Name)),
		PkgVar:     strings.ToLower(ts.Name),

		DataSource:     ts.Source,
		CSVSeparator:   ts.Separator,
		ComputePackage: ts.ComputePackage,
		ComputePkgVar:  path.Base(ts.ComputePackage),

		Fields:         fields,
		ComputedFields: computedFields,
		ComputeFns:     computeFns,
		HasComputed:    len(computeFns) != 0,
		Constraints:    ts.Constraints,

		RequireSQLPkg:     requireSQLPkg(fields),
		RequireStrconvPkg: requireStrconvPkg(fields),
		RequireTimePkg:    requireTimePkg(fields),
		HasValidation:     hasValidation(ts.Fields, ts.ComputedFields),
	}, nil
}

func newFieldsData(fs []*schema.Field) ([]*FieldData, error) {
	a := make([]*FieldData, len(fs))

	for i, f := range fs {
		t, err := goType(f.Type, f.Required)
		if err != nil {
			return nil, err
		}
		a[i] = &FieldData{
			Field:  f,
			GoType: t,
		}
	}

	return a, nil
}

func newComputedFieldsData(fs []*schema.ComputedField) ([]*ComputedFieldData, error) {
	a := make([]*ComputedFieldData, len(fs))

	for i, f := range fs {
		t, err := goType(f.Type, f.Required)
		if err != nil {
			return nil, err
		}
		a[i] = &ComputedFieldData{
			ComputedField: f,
			GoType:        t,
		}
	}

	return a, nil
}

func newComputeFnsData(fs []*FieldData, cfs []*ComputedFieldData) ([]*ComputeFnData, error) {
	a := make([]*ComputeFnData, 0)
	m := make(map[string]string)
	for _, f := range fs {
		if f.ComputeFn != "" {
			e, found := m[f.ComputeFn]
			if found && e != f.GoType {
				return nil, fmt.Errorf(
					"field with duplicated computeFn and different type: '%s'",
					f.Name)
			}
			if found {
				continue
			}
			m[f.ComputeFn] = f.GoType
			a = append(a, &ComputeFnData{
				Name:         f.ComputeFn,
				ArgumentType: f.GoType,
				ReturnType:   f.GoType,
			})
		}
	}

	for _, f := range cfs {
		if f.ComputeFn != "" {
			e, found := m[f.ComputeFn]
			if found && e != "map[string]interface{}" {
				return nil, fmt.Errorf(
					"computed field with duplicated computeFn and different type: '%s'",
					f.Name)
			}
			if found {
				continue
			}
			m[f.ComputeFn] = f.GoType
			a = append(a, &ComputeFnData{
				Name:         f.ComputeFn,
				ArgumentType: "map[string]interface{}",
				ReturnType:   f.GoType,
			})
		}
	}
	return a, nil
}

func goType(fieldType string, required bool) (string, error) {
	baseType := "unknown"
	switch {
	case fieldType == "boolean":
		baseType = "bool"
	case fieldType == "smallint", fieldType == "integer":
		baseType = "int32"
	case fieldType == "bigint":
		baseType = "int64"
	case fieldType == "double precision", fieldType == "real":
		baseType = "float64"
	case fieldType == "cidr", fieldType == "inet", fieldType == "json", fieldType == "macaddr", fieldType == "text", fieldType == "uuid":
		baseType = "string"
	case fieldType == "date":
		baseType = "time"
	case strings.HasPrefix(fieldType, "bit"):
		baseType = "string"
	case strings.HasPrefix(fieldType, "character"):
		baseType = "string"
	case strings.HasPrefix(fieldType, "char"):
		baseType = "string"
	case strings.HasPrefix(fieldType, "varchar"):
		baseType = "string"
	case strings.HasPrefix(fieldType, "time"):
		baseType = "time"
	case strings.HasPrefix(fieldType, "timestamp"):
		baseType = "time"
	}

	switch baseType {
	case "bool":
		if required {
			return "bool", nil
		}
		return "sql.NullBool", nil
	case "float64":
		if required {
			return "float64", nil
		}
		return "sql.NullFloat64", nil
	case "int32":
		if required {
			return "int32", nil
		}
		return "sql.NullInt32", nil
	case "int64":
		if required {
			return "int64", nil
		}
		return "sql.NullInt64", nil
	case "string":
		if required {
			return "string", nil
		}
		return "sql.NullString", nil
	case "time":
		if required {
			return "time.Time", nil
		}
		return "sql.NullTime", nil
	}

	return "", fmt.Errorf("invalid type: %s", fieldType)
}

func requireSQLPkg(fs []*FieldData) bool {
	for _, f := range fs {
		if strings.HasPrefix(f.GoType, "sql.") {
			return true
		}
	}
	return false
}

func requireTimePkg(fs []*FieldData) bool {
	for _, f := range fs {
		if f.GoType == "time.Time" {
			return true
		}
	}
	return false
}

func requireStrconvPkg(fs []*FieldData) bool {
	for _, f := range fs {
		switch f.GoType {
		case "bool", "float64", "int32", "int64":
			return true
		}
	}
	return false
}

func hasValidation(fs []*schema.Field, cfs []*schema.ComputedField) bool {
	for _, f := range fs {
		if len(f.Validation) != 0 {
			return true
		}
	}

	for _, f := range cfs {
		if len(f.Validation) != 0 {
			return true
		}
	}

	return false
}
