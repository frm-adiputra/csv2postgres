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
}

// Generate generates source codes based on spec
func (g Generator) Generate() error {
	tableFiles, err := listYamlFiles("tables")
	if err != nil {
		return err
	}

	viewFiles, err := listYamlFiles("views")
	if err != nil {
		return err
	}

	tableSpecs, err := createTableSpecs(tableFiles)
	if err != nil {
		return err
	}

	viewSpecs, err := createViewSpecs(viewFiles)
	if err != nil {
		return err
	}

	rootDir := "."
	if g.RootDir != "" {
		rootDir = g.RootDir
	}

	i, err := interpolation.NewInterpolator(g.BaseImportPath, rootDir, tableSpecs, viewSpecs)
	if err != nil {
		return err
	}

	err = i.Interpolate()
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

func createTableSpecs(tableFiles []string) ([]*schema.Table, error) {
	specs := make([]*schema.Table, len(tableFiles))
	for i, specFile := range tableFiles {
		s, err := schema.NewTable(specFile)
		if err != nil {
			return nil, err
		}
		specs[i] = s
	}
	return specs, nil
}

func createViewSpecs(viewFiles []string) ([]*schema.View, error) {
	specs := make([]*schema.View, len(viewFiles))
	for i, specFile := range viewFiles {
		s, err := schema.NewView(specFile)
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