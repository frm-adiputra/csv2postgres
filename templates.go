package generator

import (
	"bytes"
	"os"
	"text/template"

	"github.com/markbates/pkger"
)

var (
	tmpl = template.New("")
)

func init() {
	tmpl.Option("missingkey=error")
	tmpl.Funcs(map[string]interface{}{
		"lowerCaseFirst": lowerCaseFirst,
		"upperCaseFirst": upperCaseFirst,
		"tableName":      tableName,
		"schemaName":     schemaName,
	})
	addTemplate(tmpl, "runner.go", "/generator/templates/runner.go.tmpl")
	addTemplate(tmpl, "csvReader.go", "/generator/templates/csvReader.go.tmpl")
	addTemplate(tmpl, "fieldProvider.go", "/generator/templates/fieldProvider.go.tmpl")
	addTemplate(tmpl, "converter.go", "/generator/templates/converter.go.tmpl")
	addTemplate(tmpl, "computer.go", "/generator/templates/computer.go.tmpl")
	addTemplate(tmpl, "validator.go", "/generator/templates/validator.go.tmpl")
	addTemplate(tmpl, "dbSync.go", "/generator/templates/dbSync.go.tmpl")
	addTemplate(tmpl, "targets.go", "/generator/templates/targets.go.tmpl")
}

func addTemplate(t *template.Template, name, path string) {
	f, err := pkger.Open(path)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	buf := new(bytes.Buffer)
	buf.ReadFrom(f)
	_, err = t.New(name).Parse(buf.String())
	if err != nil {
		panic(err)
	}
}

func execTemplate(fileName, templateName string, data interface{}) error {
	f, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer f.Close()

	err = tmpl.ExecuteTemplate(f, templateName, data)
	if err != nil {
		return err
	}

	return nil
}
