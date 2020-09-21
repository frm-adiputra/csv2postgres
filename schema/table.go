package schema

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/goccy/go-yaml"
)

// Field for Table
// It is allowed for multiple fields with different name to have a same original
// name.
type Field struct {
	Name string

	// Column is the name of column in CSV file.
	Column string

	// Type is PostgreSQL field type
	Type string

	// If Required is true then the field cannot be NULL
	Required bool

	// Required only for time type. For info about the format,
	// see layout in https://golang.org/pkg/time/#Parse
	TimeFormat string `yaml:"timeFormat"`

	// Set to true to exclude field from data table.
	// Set this to true if you need this field only for reference in
	// computedFields
	Exclude bool

	// ComputeFn is a function that will be used to compute the field's value
	// based on this field value alone. If you need to compute value based on
	// multiple field values then you must use the ComputedFields.
	// This function must be available in the package defined at HelperPackage.
	// This function signature must be `func (T) (T, error)` where T must be the
	// type of the field.
	ComputeFn string `yaml:"computeFn"`

	// Validation is array of validation rule
	Validation []string `yaml:",flow"`

	// Length is the length of string type. So it only be used if type is
	// string.
	// If Length not specified, then the database column will have TEXT type.
	// If specified, it will be VARCHAR(Length).
	Length int
}

// ComputedField specifies computed field configurations.
type ComputedField struct {
	Name       string
	Type       string
	ComputeFn  string `yaml:"computeFn"`
	Required   bool
	Exclude    bool
	Validation []string `yaml:",flow"`
}

// Table specifies how to process a CSV file
type Table struct {
	Name      string `yaml:"-"` // Must contain only letters and/or number. First character must be not a number
	SpecFile  string `yaml:"-"`
	Source    string
	Separator string

	// Fields defined here will be included in data table.
	// Fields not defined here will NOT included in data table.
	// Exceptions are for fields excluded explicitly (see Exclude in Field).
	Fields []*Field `yaml:",flow"`

	// Must refer to the import path of package that provides functions to be
	// used in computeFn
	ComputePackage string `yaml:"computePackage"`

	ComputedFields []*ComputedField `yaml:"computedFields,flow"`
	Table          string

	// DependsOn specifies other spec name that the table in this spec depends
	// on.
	DependsOn  []string `yaml:"dependsOn,flow"`
	Dependants []string `yaml:"-"`

	// Constraints sepecifies the table constraints. It must be written using
	// syntax like in the table constraints of create table PostgreSQL clause.
	Constraints []string `yaml:",flow"`
}

// NewTable creates a new table spec from a YAML file.
func NewTable(specFile string) (*Table, error) {
	f, err := os.Open(specFile)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	t := &Table{}
	d := yaml.NewDecoder(f)
	err = d.Decode(t)
	if err != nil {
		return nil, err
	}

	t.SpecFile = specFile
	ns := strings.Split(filepath.Base(specFile), ".")
	t.Name = strings.Join(ns[0:len(ns)-1], ".")

	err = t.validate()
	if err != nil {
		return nil, fmt.Errorf("%s: %s", specFile, err.Error())
	}

	return t, nil
}
