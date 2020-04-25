package csv2postgres

import (
	"fmt"
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
