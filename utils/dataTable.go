package utils

import (
	"database/sql"
	"fmt"
	"strconv"
	"time"
)

// CreateHeader creates header from columns.
// Using fieldColumnMap to filter only requested field and also rename it.
func CreateHeader(columns []string, fieldColumnMap map[string]string) (map[string]int, error) {
	m := make(map[string]int, len(fieldColumnMap))
	for n, c := range fieldColumnMap {
		for i, h := range columns {
			if c == h {
				m[n] = i
			}
		}

		_, ok := m[n]
		if !ok {
			return nil, fmt.Errorf("column '%s' for field '%s' not found", c, n)
		}
	}

	return m, nil
}

// ProvideField provides row to be accessed by field names using map.
// The returned map will only provide the specified fields.
func ProvideField(fieldColumnMap map[string]int, row []string) (map[string]interface{}, error) {
	m := make(map[string]interface{}, len(fieldColumnMap))
	for k, v := range fieldColumnMap {
		if v >= len(row) {
			return nil, fmt.Errorf("invalid index for field '%s'", k)
		}
		m[k] = row[v]
	}
	return m, nil
}

// StringToNullBool converts string to sql.NullBool
func StringToNullBool(s string) (sql.NullBool, error) {
	if s == "" {
		return sql.NullBool{}, nil
	}

	v, err := strconv.ParseBool(s)
	if err != nil {
		return sql.NullBool{}, err
	}

	return sql.NullBool{
		Bool:  v,
		Valid: true,
	}, nil
}

// StringToNullFloat64 converts string to sql.NullFloat64
func StringToNullFloat64(s string) (sql.NullFloat64, error) {
	if s == "" {
		return sql.NullFloat64{}, nil
	}

	v, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return sql.NullFloat64{}, err
	}

	return sql.NullFloat64{
		Float64: v,
		Valid:   true,
	}, nil
}

// StringToNullInt32 converts string to sql.NullInt32
func StringToNullInt32(s string) (sql.NullInt32, error) {
	if s == "" {
		return sql.NullInt32{}, nil
	}

	v, err := strconv.ParseInt(s, 10, 32)
	if err != nil {
		return sql.NullInt32{}, err
	}

	return sql.NullInt32{
		Int32: int32(v),
		Valid: true,
	}, nil
}

// StringToNullInt64 converts string to sql.NullInt64
func StringToNullInt64(s string) (sql.NullInt64, error) {
	if s == "" {
		return sql.NullInt64{}, nil
	}

	v, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return sql.NullInt64{}, err
	}

	return sql.NullInt64{
		Int64: v,
		Valid: true,
	}, nil
}

// StringToNullString converts string to sql.NullString
func StringToNullString(s string) (sql.NullString, error) {
	if s == "" {
		return sql.NullString{}, nil
	}

	return sql.NullString{
		String: s,
		Valid:  true,
	}, nil
}

// StringToTime converts string to time.Time
func StringToTime(s, layout string) (time.Time, error) {
	return time.Parse(s, layout)
}

// StringToNullTime converts string to sql.NullTime
func StringToNullTime(layout, s string) (sql.NullTime, error) {
	if s == "" {
		return sql.NullTime{}, nil
	}

	v, err := time.Parse(s, layout)
	if err != nil {
		return sql.NullTime{}, err
	}

	return sql.NullTime{
		Time:  v,
		Valid: true,
	}, nil
}
