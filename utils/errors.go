package utils

import (
	"fmt"
)

// FieldError records an error that happens in a field.
type FieldError struct {
	Field   string
	Message string
}

func (e FieldError) Error() string {
	return fmt.Sprintf("field '%s': %s", e.Field, e.Message)
}

// RecordError records an error that happens in a record.
type RecordError struct {
	Source    string
	RecordNum int64
	Message   string
}

func (e RecordError) Error() string {
	return fmt.Sprintf("record #%d in %s: %s",
		e.RecordNum, e.Source, e.Message)
}

// SourceFieldError records an error that happen when converting field value.
type SourceFieldError struct {
	Source    string
	RecordNum int64
	FieldError
}

func (e SourceFieldError) Error() string {
	return fmt.Sprintf("record #%d field '%s' in %s: %s",
		e.RecordNum, e.Field, e.Source, e.Message)
}

// SourceGenericError records a generic error that happen when processing csv file.
type SourceGenericError struct {
	Source    string
	RecordNum int64
	Message   string
}

func (e SourceGenericError) Error() string {
	return fmt.Sprintf("record #%d in %s: %s",
		e.RecordNum, e.Source, e.Message)
}

// ErrEmptyValue indicates an empty value error
func ErrEmptyValue(field string) FieldError {
	return FieldError{
		Field:   field,
		Message: "empty value not allowed",
	}
}
