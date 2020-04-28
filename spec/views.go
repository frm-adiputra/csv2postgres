package spec

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/goccy/go-yaml"
)

// Views specifies a list of views
type Views struct {
	Views    []*View `yaml:",flow"`
	SpecFile string  `yaml:"-"`
}

// View specifies a database view
type View struct {
	Name      string `yaml:"-"`
	Schema    string
	DependsOn []string `yaml:"dependsOn"`

	// File defines the SQL file name. It should follow Golang identifier name
	// rules. The view name will be generated based on file name and converted
	// to snake case.
	File string
}

// NewViews creates a new views spec from a YAML file.
func NewViews(specFile string) (*Views, error) {
	f, err := os.Open(specFile)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	vs := &Views{}
	d := yaml.NewDecoder(f)
	err = d.Decode(vs)
	if err != nil {
		return nil, err
	}

	vs.SpecFile = specFile

	err = vs.validate()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", specFile, err)
	}

	vs.initValues()

	return vs, nil
}

func (vs *Views) initValues() {
	for _, v := range vs.Views {
		v.initValues()
	}
}

func (v *View) initValues() {
	// View name is its file name without file extension
	a := strings.Split(filepath.Base(v.File), ".")
	v.Name = strings.Join(a[0:len(a)-1], ".")
}

func (vs *Views) validate() error {
	for _, v := range vs.Views {
		err := v.validate()
		if err != nil {
			return err
		}
	}
	return nil
}

func (v *View) validate() error {
	if v.File == "" {
		return errors.New("file required")
	}
	return nil
}
