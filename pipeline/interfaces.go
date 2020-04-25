package pipeline

import "database/sql"

// RowReader is the interface that wraps the functionality of row reader.
type RowReader interface {
	Open() error
	ReadRow() ([]string, error)
	HeaderRow() []string
	Close() error

	// RowCount returns the number of rows read.
	RowCount() int64

	// Source returns the source of rows.
	Source() string
}

// FieldProvider is the interface that wraps the functionality of accessing
// fields.
type FieldProvider interface {
	Setup(header []string) error
	ProvideField([]string) (map[string]interface{}, error)
}

// Converter is the interface that wraps the functionality of converting field
// value
type Converter interface {
	Convert(map[string]interface{}) (map[string]interface{}, error)
}

// Computer is the interface that wraps the functionality of computing field
// value
type Computer interface {
	Compute(map[string]interface{}) (map[string]interface{}, error)
}

// Validator is the interface that wraps the functionality of validating field
// value
type Validator interface {
	Validate(map[string]interface{}) error
}

// DBSynchronizer synchronize specs and data with database table.
type DBSynchronizer interface {
	// Name returns the table's name
	Name() string

	// DependsOn returns other tables that this table depends on
	DependsOn() []string

	// Create table
	Create(*sql.DB) error

	// Delete all rows from table
	Delete(*sql.DB) error

	// Drop table
	Drop(*sql.DB) error

	// Fill rows
	Fill(*sql.DB) error
}
