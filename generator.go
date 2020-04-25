package generator

import (
	"os"
	"path/filepath"

	"github.com/frm-adiputra/csv2postgres/spec"
)

// Generator specifies the generator configurations.
type Generator struct {
	BaseImportPath string
	OutDir         string
	Specs          []string
}

// Generate generates source codes based on spec
func (g Generator) Generate() error {
	specs, err := g.createSpecs()
	if err != nil {
		return err
	}

	td, err := newTemplateData(g, specs)
	if err != nil {
		return err
	}

	err = g.generateCommons(td)
	if err != nil {
		return err
	}

	err = g.generateBasedOnSpec(td.Specs)
	if err != nil {
		return err
	}
	return nil
}

func (g Generator) createSpecs() ([]*spec.Spec, error) {
	specs := make([]*spec.Spec, len(g.Specs))
	for i, specFile := range g.Specs {
		s, err := spec.NewSpec(specFile)
		if err != nil {
			return nil, err
		}
		specs[i] = s
	}
	return specs, nil
}

func (g Generator) generateCommons(td *templateData) error {
	err := execTemplate(
		filepath.Join(g.OutDir, generatedFilename("runner.go")),
		"runner.go", td)
	if err != nil {
		return err
	}
	return nil
}

func (g Generator) generateBasedOnSpec(specsTD []specTemplateData) error {
	for _, std := range specsTD {

		// create directory for package
		err := os.MkdirAll(std.PkgDir, 0777)
		if err != nil {
			return err
		}

		// generate source codes
		tmplNames := []string{
			// "data.go",
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

		err = execTemplate(
			filepath.Join(
				std.RootPkgDir,
				generatedFilename(lowerCaseFirst(std.Name)+"Targets.go")),
			"targets.go", std)
		if err != nil {
			return err
		}
	}
	return nil
}
