package generator

import (
	"fmt"
	"path"
	"path/filepath"
	"strings"

	"github.com/frm-adiputra/csv2postgres/spec"
)

type templateData struct {
	Generator string
	Specs     []specTemplateData
	PkgVar    string
}

type specTemplateData struct {
	Generator         string
	Name              string
	ImportPath        string
	PkgDir            string
	PkgVar            string
	RootPkgDir        string
	RootPkgVar        string
	DataSource        string
	CSVSeparator      string
	Fields            []*fieldTemplateData
	ComputedFields    []*computedFieldTemplateData
	ComputePackage    string
	ComputePkgVar     string
	HasComputed       bool
	RequireSQLPkg     bool
	RequireStrconvPkg bool
	RequireTimePkg    bool
	ComputeFns        []*computeFnTemplateData
	HasValidation     bool
	Table             *spec.Table
}

type fieldTemplateData struct {
	*spec.Field
	GoType string
}

type computedFieldTemplateData struct {
	*spec.ComputedField
	GoType string
}

type computeFnTemplateData struct {
	Name         string
	ArgumentType string
	ReturnType   string
}

func newTemplateData(g Generator, specs []*spec.Spec) (*templateData, error) {
	csa := make([]specTemplateData, len(specs))
	errStr := "generating code for '%s': %w"
	for i, s := range specs {
		fields, err := createFieldsTemplateData(s.Fields)
		if err != nil {
			return nil, fmt.Errorf(errStr, s.SpecFile, err)
		}

		computedFields, err := createComputedFieldsTemplateData(s.ComputedFields)
		if err != nil {
			return nil, fmt.Errorf(errStr, s.SpecFile, err)
		}

		computeFns, err := computeFns(fields, computedFields)
		if err != nil {
			return nil, fmt.Errorf(errStr, s.SpecFile, err)
		}

		requireSQLPkg := requireSQLPkg(fields)
		requireStrconvPkg := requireStrconvPkg(fields)
		requireTimePkg := requireTimePkg(fields)
		hasValidation := hasValidation(s.Fields, s.ComputedFields)
		csa[i] = specTemplateData{
			Generator:         "github.com/frm-adiputra/csv2postgres",
			Name:              s.Name,
			ImportPath:        path.Join(g.BaseImportPath, "internal", strings.ToLower(s.Name)),
			PkgDir:            filepath.Join(g.OutDir, "internal", strings.ToLower(s.Name)),
			PkgVar:            strings.ToLower(s.Name),
			RootPkgDir:        g.OutDir,
			RootPkgVar:        strings.ToLower(path.Base(g.BaseImportPath)),
			DataSource:        s.Source,
			CSVSeparator:      s.Separator,
			Fields:            fields,
			ComputedFields:    computedFields,
			ComputePackage:    s.ComputePackage,
			ComputePkgVar:     path.Base(s.ComputePackage),
			HasComputed:       len(computeFns) != 0,
			RequireSQLPkg:     requireSQLPkg,
			RequireStrconvPkg: requireStrconvPkg,
			RequireTimePkg:    requireTimePkg,
			ComputeFns:        computeFns,
			HasValidation:     hasValidation,
			Table:             s.Table,
		}
	}
	return &templateData{
		Generator: "github.com/frm-adiputra/csv2postgres",
		Specs:     csa,
		PkgVar:    strings.ToLower(path.Base(g.BaseImportPath)),
	}, nil
}

func createFieldsTemplateData(fs []*spec.Field) ([]*fieldTemplateData, error) {
	a := make([]*fieldTemplateData, len(fs))

	for i, f := range fs {
		t, err := goType(f.Type, f.Required)
		if err != nil {
			return nil, err
		}
		a[i] = &fieldTemplateData{
			Field:  f,
			GoType: t,
		}
	}

	return a, nil
}

func createComputedFieldsTemplateData(fs []*spec.ComputedField) ([]*computedFieldTemplateData, error) {
	a := make([]*computedFieldTemplateData, len(fs))

	for i, f := range fs {
		t, err := goType(f.Type, f.Required)
		if err != nil {
			return nil, err
		}
		a[i] = &computedFieldTemplateData{
			ComputedField: f,
			GoType:        t,
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

func requireSQLPkg(fs []*fieldTemplateData) bool {
	for _, f := range fs {
		if strings.HasPrefix(f.GoType, "sql.") {
			return true
		}
	}
	return false
}

func requireTimePkg(fs []*fieldTemplateData) bool {
	for _, f := range fs {
		if f.GoType == "time.Time" {
			return true
		}
	}
	return false
}

func requireStrconvPkg(fs []*fieldTemplateData) bool {
	for _, f := range fs {
		switch f.GoType {
		case "bool", "float64", "int32", "int64":
			return true
		}
	}
	return false
}

func computeFns(fs []*fieldTemplateData, cfs []*computedFieldTemplateData) ([]*computeFnTemplateData, error) {
	a := make([]*computeFnTemplateData, 0)
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
			a = append(a, &computeFnTemplateData{
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
			a = append(a, &computeFnTemplateData{
				Name:         f.ComputeFn,
				ArgumentType: "map[string]interface{}",
				ReturnType:   f.GoType,
			})
		}
	}
	return a, nil
}

func hasValidation(fs []*spec.Field, cfs []*spec.ComputedField) bool {
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