package csv2postgres

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
	addTemplate(tmpl, "runner.go", "/templates/runner.go.tmpl")
	addTemplate(tmpl, "csvReader.go", "/templates/csvReader.go.tmpl")
	addTemplate(tmpl, "fieldProvider.go", "/templates/fieldProvider.go.tmpl")
	addTemplate(tmpl, "converter.go", "/templates/converter.go.tmpl")
	addTemplate(tmpl, "computer.go", "/templates/computer.go.tmpl")
	addTemplate(tmpl, "validator.go", "/templates/validator.go.tmpl")
	addTemplate(tmpl, "dbSync.go", "/templates/dbSync.go.tmpl")
	addTemplate(tmpl, "targets.go", "/templates/targets.go.tmpl")
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
