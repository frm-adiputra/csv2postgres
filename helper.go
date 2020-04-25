package csv2postgres

import (
	"fmt"
	"strings"
	"unicode"
)

func lowerCaseFirst(s string) string {
	for i, v := range s {
		return string(unicode.ToLower(v)) + s[i+1:]
	}
	return ""
}

func upperCaseFirst(s string) string {
	for i, v := range s {
		return string(unicode.ToUpper(v)) + s[i+1:]
	}
	return ""
}

func generatedFilename(s string) string {
	return fmt.Sprintf("g_%s", s)
}

func tableName(s string) (string, error) {
	a := strings.Split(s, ".")
	if len(a) != 2 {
		return "", fmt.Errorf("invalid schema.table name: %s", s)
	}
	return a[1], nil
}

func schemaName(s string) (string, error) {
	a := strings.Split(s, ".")
	if len(a) != 2 {
		return "", fmt.Errorf("invalid schema.table name: %s", s)
	}
	return a[0], nil
}
