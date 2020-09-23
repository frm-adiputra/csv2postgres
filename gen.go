package csv2postgres

import (
	"os"
	"path/filepath"

	"github.com/frm-adiputra/csv2postgres/interpolation"
	"github.com/frm-adiputra/csv2postgres/schema"
)

// Generator specifies the generator configurations.
type Generator struct {
	BaseImportPath string
	RootDir        string
	DefaultSchema  string
}

// Generate generates source codes based on spec
func (g Generator) Generate() error {
	defaultSchema := "public"
	if g.DefaultSchema != "" {
		defaultSchema = g.DefaultSchema
	}

	tableFiles, err := listYamlFiles("tables")
	if err != nil {
		return err
	}

	viewFiles, err := listYamlFiles("views")
	if err != nil {
		return err
	}

	tableSpecs, err := createTableSpecs(tableFiles, defaultSchema)
	if err != nil {
		return err
	}

	viewSpecs, err := createViewSpecs(viewFiles, defaultSchema)
	if err != nil {
		return err
	}

	rootDir := "."
	if g.RootDir != "" {
		rootDir = g.RootDir
	}

	i, err := interpolation.NewInterpolator(g.BaseImportPath, rootDir, defaultSchema, tableSpecs, viewSpecs)
	if err != nil {
		return err
	}

	err = i.Interpolate()
	if err != nil {
		return err
	}

	err = g.generateCommons(i)
	if err != nil {
		return err
	}

	err = g.generateTables(i.TablesData)
	if err != nil {
		return err
	}

	return nil
}

func listYamlFiles(rootPath string) ([]string, error) {
	files := make([]string, 0)
	err := filepath.Walk(rootPath,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			if info.Mode().IsRegular() && (filepath.Ext(path) == ".yaml" || filepath.Ext(path) == ".yml") {
				files = append(files, path)
			}
			return nil
		})
	if err != nil {
		return nil, err
	}
	return files, nil
}

func createTableSpecs(tableFiles []string, defaultSchema string) ([]*schema.Table, error) {
	specs := make([]*schema.Table, len(tableFiles))
	for i, specFile := range tableFiles {
		s, err := schema.NewTable(specFile, defaultSchema)
		if err != nil {
			return nil, err
		}
		specs[i] = s
	}
	return specs, nil
}

func createViewSpecs(viewFiles []string, defaultSchema string) ([]*schema.View, error) {
	specs := make([]*schema.View, len(viewFiles))
	for i, specFile := range viewFiles {
		s, err := schema.NewView(specFile, defaultSchema)
		if err != nil {
			return nil, err
		}
		specs[i] = s
	}
	return specs, nil
}

func (g Generator) generateCommons(itp *interpolation.Interpolator) error {
	err := execTemplate(
		filepath.Join(g.RootDir, generatedFilename("runner.go")),
		"runner.go", itp)
	if err != nil {
		return err
	}
	return nil
}

func (g Generator) generateTables(tds []*interpolation.TableData) error {
	for _, std := range tds {

		// create directory for package
		err := os.MkdirAll(std.PkgDir, 0777)
		if err != nil {
			return err
		}

		// generate source codes
		tmplNames := []string{
			"csvReader.go",
			"fieldProvider.go",
			"converter.go",
			"computer.go",
			"validator.go",
			"dbSync.go",
		}

		for _, tn := range tmplNames {
			err = execTemplate(
				filepath.Join(std.PkgDir, generatedFilename(tn)),
				tn, std)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
