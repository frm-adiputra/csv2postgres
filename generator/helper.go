package generator

import (
	"fmt"
	"io/ioutil"
	"strings"
	"unicode"

	"github.com/iancoleman/strcase"
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

func toPackageName(s string) string {
	return strings.ToLower(strcase.ToCamel(s))
}

func toExportedName(s string) string {
	return strcase.ToCamel(s)
}

func readFile(p string) (string, error) {
	b, err := ioutil.ReadFile(p)
	if err != nil {
		return "", err
	}
	return string(b), nil
}
