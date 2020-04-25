package pipeline

import (
	"encoding/csv"
	"io"
	"os"
)

// CSVReader is implementation of RowReader using CSV as data source.
type CSVReader struct {
	FileName  string
	Separator rune

	f        *os.File
	rcsv     *csv.Reader
	header   []string
	rowCount int64
}

// Open opens CSV file for reading.
func (r *CSVReader) Open() error {
	f, err := os.Open(r.FileName)
	if err != nil {
		return err
	}
	r.f = f
	r.rowCount = 0
	rcsv := csv.NewReader(f)
	rcsv.Comma = r.Separator
	rcsv.TrimLeadingSpace = true
	r.rcsv = rcsv

	err = r.readHeader()
	if err != nil {
		return err
	}
	return nil
}

func (r *CSVReader) readHeader() error {
	rec, err := r.rcsv.Read()
	if err == io.EOF {
		return err
	}
	if err != nil {
		return err
	}

	r.header = rec
	return nil
}

// HeaderRow returns header row.
func (r *CSVReader) HeaderRow() []string { return r.header }

// Close closes CSV file.
func (r *CSVReader) Close() error {
	return r.f.Close()
}

// ReadRow reads a single row
func (r *CSVReader) ReadRow() ([]string, error) {
	a, err := r.rcsv.Read()
	if err != nil {
		return nil, err
	}
	r.rowCount++
	return a, nil
}

// RowCount returns the number of rows read.
func (r *CSVReader) RowCount() int64 { return r.rowCount }

// Source returns the source of rows.
func (r *CSVReader) Source() string { return r.FileName }
