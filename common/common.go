package common

import (
	"fmt"
	"path/filepath"
	"strings"
)

// Names represents entity names
type Names struct {
	Name        string `yaml:"-"`
	SchemaName  string `yaml:"-"`
	FullName    string `yaml:"-"`
	SQLFullName string `yaml:"-"`
}

// NewNames create names based on filename
func NewNames(filename string) (*Names, error) {
	baseName := filepath.Base(filename)
	n := strings.Split(baseName, ".")

	if len(n) != 3 {
		return nil, fmt.Errorf(
			"%s: Table/views file spec must be in format 'schema_name.table_or_view_name.yaml'",
			filename)
	}

	schemaName := n[0]
	name := n[1]

	return &Names{
		SchemaName:  schemaName,
		Name:        name,
		FullName:    fmt.Sprintf("%s.%s", schemaName, name),
		SQLFullName: fmt.Sprintf(`"%s"."%s"`, schemaName, name),
	}, nil
}
