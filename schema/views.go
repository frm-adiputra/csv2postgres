package schema

import (
	"errors"
	"fmt"
	"os"

	"github.com/frm-adiputra/csv2postgres/common"
	"github.com/goccy/go-yaml"
)

// View specifies a database view
type View struct {
	*common.Names `yaml:"-"`
	SpecFile      string `yaml:"-"`

	DependsOn []string `yaml:"dependsOn"`

	// If set, the result of view will be exported to CSV file defined.
	Export string

	// The SQL for view
	SQL string `yaml:"sql"`
}

// NewView creates a new view spec from a YAML file.
func NewView(specFile, defaultSchema string) (*View, error) {
	f, err := os.Open(specFile)
	if err != nil {
		return nil, fmt.Errorf("%s: %s", specFile, err.Error())
	}
	defer f.Close()

	t := &View{}
	d := yaml.NewDecoder(f)
	err = d.Decode(t)
	if err != nil {
		return nil, fmt.Errorf("%s: %s", specFile, err.Error())
	}

	t.SpecFile = specFile
	t.Names, err = common.NewNames(specFile, defaultSchema)
	if err != nil {
		return nil, fmt.Errorf("%s: %s", specFile, err.Error())
	}

	err = t.validate()
	if err != nil {
		return nil, fmt.Errorf("%s: %s", specFile, err.Error())
	}

	return t, nil
}

func (v *View) validate() error {
	if v.SQL == "" {
		return errors.New("sql cannot be empty")
	}
	return nil
}
