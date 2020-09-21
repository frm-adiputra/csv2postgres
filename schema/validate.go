package schema

import (
	"errors"
	"fmt"
	"strings"
)

var (
	errSeparator = errors.New("separator must be a single character")
	fieldGoTypes = []string{
		"bool",
		"float64",
		"int32",
		"int64",
		"string",
		"time",
	}
	fieldTypes = []string{
		"bigint",
		"boolean",
		"cidr",
		"date",
		"double precision",
		"inet",
		"integer",
		"json",
		"macaddr",
		"real",
		"smallint",
		"text",
		"uuid",
	}
	fieldTypesPrefix = []string{
		"bit",       // bit [ (n) ] | bit varying [ (n) ]
		"character", // character [ (n) ] | character varying [ (n) ]
		"char",      // char [ (n) ]
		"varchar",   // varchar [ (n) ]
		"time",      // time [ (p) ] [ without time zone ] | time [ (p) ] with time zone
		"timestamp", // timestamp [ (p) ] [ without time zone ] | timestamp [ (p) ] with time zone
	}
)

func (s *Table) validate() error {
	if s.Name == "" {
		return errors.New("name is required")
	}

	if len(s.Separator) != 1 {
		return errSeparator
	}

	for _, f := range s.Fields {
		err := validateField(f)
		if err != nil {
			return err
		}
	}

	for _, f := range s.ComputedFields {
		err := validateComputedField(f)
		if err != nil {
			return err
		}
	}

	err := s.checkDuplicateFieldNames()
	if err != nil {
		return err
	}

	if s.Table == "" {
		return errors.New("table name required")
	}

	return nil
}

func (s *Table) checkDuplicateFieldNames() error {
	m := make(map[string]bool)

	for _, f := range s.Fields {
		_, found := m[f.Name]
		if found {
			return fmt.Errorf("duplicate field name: '%s'", f.Name)
		}
		m[f.Name] = true
	}

	for _, f := range s.ComputedFields {
		_, found := m[f.Name]
		if found {
			return fmt.Errorf("duplicate computed field name: '%s'", f.Name)
		}
		m[f.Name] = true
	}
	return nil
}

func validateField(f *Field) error {
	if f.Name == "" {
		return errors.New("name is required")
	}

	if f.Column == "" {
		f.Column = f.Name
	}

	if f.Type == "" {
		return fmt.Errorf("validating field '%s': type required", f.Name)
	}

	if !validFieldType(f.Type) {
		return fmt.Errorf("validating field '%s': invalid type '%s'", f.Name, f.Type)
	}

	if f.Type == "time" && f.TimeFormat == "" {
		return fmt.Errorf("validating field '%s': timeFormat must not be empty", f.Name)
	}

	return nil
}

func validFieldType(ft string) bool {
	for _, t := range fieldTypes {
		if t == ft {
			return true
		}
	}
	for _, t := range fieldTypesPrefix {
		if strings.HasPrefix(ft, t) {
			return true
		}
	}
	return false
}

func validateComputedField(f *ComputedField) error {
	if f.Name == "" {
		return errors.New("name is required")
	}

	if f.Type == "" {
		return fmt.Errorf("validating computed field '%s': type required", f.Name)
	}

	if f.ComputeFn == "" {
		return fmt.Errorf(
			"validating computed field '%s': computeFn is required", f.Name)
	}

	if !validFieldType(f.Type) {
		return fmt.Errorf("validating computed field '%s': invalid type '%s'",
			f.Name, f.Type)
	}

	return nil
}
