package csv2postgres

import (
	"fmt"
	"path"
	"path/filepath"
	"strings"

	"github.com/frm-adiputra/csv2postgres/spec"
)

type templateData struct {
	Generator     string
	Specs         []specTemplateData
	PkgVar        string
	CreateDepsAll []dependencyData
	DropDepsAll   []dependencyData
	Views         []viewTemplateData
	ImportPath    string
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
	Table             *tableTemplateData
	DependsOn         []string
	CreateDeps        []dependencyData
	DropDeps          []dependencyData
	Constraints       []string
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

type tableTemplateData struct {
	Name           string
	QuotedFullName string
	SchemaName     string
	TableName      string
}

type viewTemplateData struct {
	*spec.View
	CreateDeps []dependencyData
	DropDeps   []dependencyData
}

type dependencyData struct {
	Name  string
	Table bool
}

func newTemplateData(g Generator, specs []*spec.Spec, vs *spec.Views) (*templateData, error) {
	err := checkDuplicateNames(specs, vs)
	if err != nil {
		return nil, err
	}

	err = checkDuplicateTableNames(specs)
	if err != nil {
		return nil, err
	}

	deps, targetIsTable, err := linkDependencies(specs, vs)
	if err != nil {
		return nil, err
	}

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

		tableData, err := createTableTemplateData(s.Table)
		if err != nil {
			return nil, fmt.Errorf(errStr, s.SpecFile, err)
		}

		cDeps, err := deps.CreateOrder(s.Name)
		if err != nil {
			return nil, fmt.Errorf(errStr, s.SpecFile, err)
		}
		createDeps := createDependenciesData(cDeps, targetIsTable)

		dDeps, err := deps.DropOrder(s.Name)
		if err != nil {
			return nil, fmt.Errorf(errStr, s.SpecFile, err)
		}
		dropDeps := createDependenciesData(dDeps, targetIsTable)

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
			Table:             tableData,
			DependsOn:         s.DependsOn,
			CreateDeps:        createDeps,
			DropDeps:          dropDeps,
			Constraints:       s.Constraints,
		}
	}

	vsd, err := createViewsTemplateData(vs, deps, targetIsTable)
	if err != nil {
		return nil, fmt.Errorf(errStr, vs.SpecFile, err)
	}

	return &templateData{
		Generator:     "github.com/frm-adiputra/csv2postgres",
		Specs:         csa,
		PkgVar:        strings.ToLower(path.Base(g.BaseImportPath)),
		CreateDepsAll: createDependenciesData(deps.CreateOrderAll(), targetIsTable),
		DropDepsAll:   createDependenciesData(deps.DropOrderAll(), targetIsTable),
		Views:         vsd,
		ImportPath:    g.BaseImportPath,
	}, nil
}

func createDependenciesData(deps []string, targetIsTable map[string]bool) []dependencyData {
	l := make([]dependencyData, len(deps))
	for i, v := range deps {
		_, isTable := targetIsTable[v]
		l[i] = dependencyData{
			Name:  v,
			Table: isTable,
		}
	}
	return l
}

func checkDuplicateNames(specs []*spec.Spec, vs *spec.Views) error {
	m := make(map[string]bool)
	for _, s := range specs {
		_, found := m[s.Name]
		if found {
			return fmt.Errorf("duplicate name '%s' in '%s'",
				s.Name, s.SpecFile)
		}
		m[s.Name] = true
	}
	for _, v := range vs.Views {
		_, found := m[v.Name]
		if found {
			return fmt.Errorf("duplicate name '%s' in '%s'",
				v.Name, vs.SpecFile)
		}
		m[v.Name] = true
	}
	return nil
}

func linkDependencies(specs []*spec.Spec, vs *spec.Views) (DepsGraph, map[string]bool, error) {
	targets := make([]string, 0, len(specs)+len(vs.Views))
	targetIsTable := make(map[string]bool)
	for _, s := range specs {
		targets = append(targets, s.Name)
		targetIsTable[s.Name] = true
	}
	for _, v := range vs.Views {
		targets = append(targets, v.Name)
	}

	d := NewDepsGraph(targets)

	for _, s := range specs {
		for _, ds := range s.DependsOn {
			err := d.DependsOn(s.Name, ds)
			if err != nil {
				return DepsGraph{}, nil, err
			}
		}
	}

	for _, s := range vs.Views {
		for _, ds := range s.DependsOn {
			err := d.DependsOn(s.Name, ds)
			if err != nil {
				return DepsGraph{}, nil, err
			}
		}
	}

	err := d.Finalize()
	if err != nil {
		return DepsGraph{}, nil, err
	}

	return *d, targetIsTable, nil
}

func checkDuplicateTableNames(specs []*spec.Spec) error {
	m := make(map[string]bool)
	for _, s := range specs {
		_, found := m[s.Table]
		if found {
			return fmt.Errorf("duplicate table name '%s' in '%s'",
				s.Table, s.SpecFile)
		}
		m[s.Table] = true
	}
	return nil
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

func createTableTemplateData(table string) (*tableTemplateData, error) {
	quotedFullName, schemaName, tableName, err := parseTableName(table)
	if err != nil {
		return nil, err
	}
	return &tableTemplateData{
		Name:           table,
		QuotedFullName: quotedFullName,
		SchemaName:     schemaName,
		TableName:      tableName,
	}, nil
}

func parseTableName(s string) (quotedFullName, schemaName, tableName string, err error) {
	a := strings.Split(s, ".")
	if len(a) > 2 {
		return "", "", "", fmt.Errorf("invalid table name: %s", s)
	}

	if len(a) == 2 {
		schemaName = a[0]
		tableName = a[1]
		quotedFullName = fmt.Sprintf(`"%s"."%s"`, a[0], a[1])
	} else if len(a) == 1 {
		tableName = a[0]
		quotedFullName = fmt.Sprintf(`"%s"`, a[0])
	}
	return quotedFullName, schemaName, tableName, nil
}

func createViewsTemplateData(vs *spec.Views, deps DepsGraph, targetIsTable map[string]bool) ([]viewTemplateData, error) {
	l := make([]viewTemplateData, 0, len(vs.Views))
	for _, v := range vs.Views {
		cDeps, err := deps.CreateOrder(v.Name)
		if err != nil {
			return nil, err
		}
		createDeps := createDependenciesData(cDeps, targetIsTable)
		dDeps, err := deps.DropOrder(v.Name)
		if err != nil {
			return nil, err
		}
		dropDeps := createDependenciesData(dDeps, targetIsTable)
		l = append(l, viewTemplateData{
			View:       v,
			CreateDeps: createDeps,
			DropDeps:   dropDeps,
		})
	}
	return l, nil
}
