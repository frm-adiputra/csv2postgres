package common

import (
	"fmt"
	"path/filepath"
	"strings"
)

// Names represents entity names
type Names struct {
	RefName     string `yaml:"-"` // Name from filename, to be used on dependencies
	TargetName  string `yaml:"-"` // Name to be used on generated code
	Name        string `yaml:"-"`
	SchemaName  string `yaml:"-"`
	FullName    string `yaml:"-"`
	SQLFullName string `yaml:"-"`
}

// NewNames create names based on filename
func NewNames(filename, defaultSchema string) (*Names, error) {
	baseName := filepath.Base(filename)
	n := strings.Split(baseName, ".")

	var schemaName, name, refName string

	if len(n) == 2 {
		schemaName = defaultSchema
		name = n[0]
		refName = name
	} else if len(n) == 3 {
		schemaName = n[0]
		name = n[1]
		refName = fmt.Sprintf("%s.%s", schemaName, name)
	} else {
		return nil, fmt.Errorf(
			"%s: Table/views file spec must be in format 'schema_name.table_or_view_name.yaml' or 'table_or_view_name.yaml'",
			filename)
	}

	return &Names{
		RefName:     refName,
		SchemaName:  schemaName,
		Name:        name,
		FullName:    fmt.Sprintf("%s.%s", schemaName, name),
		SQLFullName: fmt.Sprintf(`"%s"."%s"`, schemaName, name),
	}, nil
}
