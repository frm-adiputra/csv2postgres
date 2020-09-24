package csv2postgres

import (
	"fmt"
	"os"
	"text/template"

	"github.com/frm-adiputra/csv2postgres/internal/box"
	"github.com/iancoleman/strcase"
)

var (
	tmpl = template.New("")
)

func init() {
	tmpl.Option("missingkey=error")
	tmpl.Funcs(map[string]interface{}{
		"lowerCaseFirst": lowerCaseFirst,
		"upperCaseFirst": upperCaseFirst,
		"toPackageName":  toPackageName,
		"toExportedName": toExportedName,
		"readFile":       readFile,
		"toSnake":        strcase.ToSnake,
		// "hasTableDep":    hasTableDep,
	})
	addTemplate(tmpl, "runner.go", "/runner.go.tmpl")
	addTemplate(tmpl, "csvReader.go", "/csvReader.go.tmpl")
	addTemplate(tmpl, "fieldProvider.go", "/fieldProvider.go.tmpl")
	addTemplate(tmpl, "converter.go", "/converter.go.tmpl")
	addTemplate(tmpl, "computer.go", "/computer.go.tmpl")
	addTemplate(tmpl, "validator.go", "/validator.go.tmpl")
	addTemplate(tmpl, "dbSync.go", "/dbSync.go.tmpl")
	addTemplate(tmpl, "targets.go", "/targets.go.tmpl")
	addTemplate(tmpl, "view.go", "/view.go.tmpl")
	addTemplate(tmpl, "main.go", "/main.go.tmpl")
}

func addTemplate(t *template.Template, name, path string) {
	_, err := t.New(name).Parse(string(box.Get(path)))
	if err != nil {
		panic(err)
	}
}

func execTemplate(fileName, templateName string, data interface{}) error {
	fmt.Printf("Generating '%s' ...", fileName)
	f, err := os.Create(fileName)
	if err != nil {
		fmt.Println("[FAILED]")
		return err
	}
	defer f.Close()

	err = tmpl.ExecuteTemplate(f, templateName, data)
	if err != nil {
		fmt.Println("[FAILED]")
		return err
	}

	fmt.Println("[OK]")
	return nil
}
