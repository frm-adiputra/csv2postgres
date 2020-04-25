package pipeline

import "fmt"

// RecordReader is a type that specify record reader.
type RecordReader struct {
	Name          string
	RowReader     RowReader
	FieldProvider FieldProvider
	Converter     Converter
	Computer      Computer
	Validator     Validator
}

// Open opens the record reader for reading.
func (r *RecordReader) Open() error {
	err := r.RowReader.Open()
	if err != nil {
		return err
	}
	err = r.FieldProvider.Setup(r.RowReader.HeaderRow())
	if err != nil {
		return err
	}
	return nil
}

// Close closes the record reader, no further reading allowed.
func (r *RecordReader) Close() error {
	return r.RowReader.Close()
}

// ReadRecord reads a single record.
func (r *RecordReader) ReadRecord() (map[string]interface{}, error) {
	row, err := r.RowReader.ReadRow()
	if err != nil {
		return nil, err
	}

	fields, err := r.FieldProvider.ProvideField(row)
	if err != nil {
		return nil, err
	}

	fields, err = r.Converter.Convert(fields)
	if err != nil {
		return nil, err
	}

	fields, err = r.Computer.Compute(fields)
	if err != nil {
		return nil, err
	}

	err = r.Validator.Validate(fields)
	if err != nil {
		return nil, fmt.Errorf("%s record #%d: %w",
			r.Name, r.RowReader.RowCount(), err)
	}
	return fields, nil
}
