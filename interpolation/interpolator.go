package interpolation

import (
	"fmt"

	"github.com/frm-adiputra/csv2postgres/schema"
)

// Interpolator represents interpolator struct
type Interpolator struct {
	tablesData []*TableData
	viewsData  []*ViewData
	pdeps      *projectDeps
}

// NewInterpolator returns a new interpolator
func NewInterpolator(baseImportPath, rootDir string, ts []*schema.Table, vs []*schema.View) (*Interpolator, error) {
	var err error
	pdeps, err := newProjectDeps(ts, vs)

	lt := make([]*TableData, len(ts))
	for i, t := range ts {
		lt[i], err = newTableData(t, baseImportPath, rootDir)
		if err != nil {
			return nil, err
		}
	}

	vt := make([]*ViewData, len(vs))
	for i, v := range vs {
		vt[i] = &ViewData{View: v}
	}

	return &Interpolator{
		tablesData: lt,
		viewsData:  vt,
		pdeps:      pdeps,
	}, nil
}

// Interpolate do interpolation based on specification files provided
func (ip *Interpolator) Interpolate() error {
	err := ip.checkDuplicateName()
	if err != nil {
		return err
	}
	return nil
}

func (ip Interpolator) checkDuplicateName() error {
	m := make(map[string]bool)
	for _, t := range ip.tablesData {
		_, found := m[t.FullName]
		if found {
			return fmt.Errorf("duplicate name '%s' in '%s'",
				t.FullName, t.SpecFile)
		}
		m[t.Name] = true
	}
	for _, v := range ip.viewsData {
		_, found := m[v.FullName]
		if found {
			return fmt.Errorf("duplicate name '%s' in '%s'",
				v.FullName, v.SpecFile)
		}
		m[v.Name] = true
	}
	return nil
}

func (ip *Interpolator) linkDependencies() error {
	var err error
	for _, t := range ip.tablesData {
		t.CreateDeps, err = ip.pdeps.createDependenciesData(t.FullName, false)
		if err != nil {
			return err
		}
		t.DropDeps, err = ip.pdeps.createDependenciesData(t.FullName, true)
		if err != nil {
			return err
		}
	}

	for _, v := range ip.viewsData {
		v.CreateDeps, err = ip.pdeps.createDependenciesData(v.FullName, false)
		if err != nil {
			return err
		}
		v.DropDeps, err = ip.pdeps.createDependenciesData(v.FullName, true)
		if err != nil {
			return err
		}
	}
	return nil
}
