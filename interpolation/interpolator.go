package interpolation

import (
	"fmt"

	"github.com/frm-adiputra/csv2postgres/schema"
)

// Interpolator represents interpolator struct
type Interpolator struct {
	Generator     string
	TablesData    []*TableData
	ViewsData     []*ViewData
	pdeps         *projectDeps
	defaultSchema string
}

// NewInterpolator returns a new interpolator
func NewInterpolator(baseImportPath, rootDir, defaultSchema string, ts []*schema.Table, vs []*schema.View) (*Interpolator, error) {
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
		Generator:     "github.com/frm-adiputra/csv2postgres",
		TablesData:    lt,
		ViewsData:     vt,
		pdeps:         pdeps,
		defaultSchema: defaultSchema,
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
	for _, t := range ip.TablesData {
		_, found := m[t.FullName]
		if found {
			return fmt.Errorf("duplicate name '%s' in '%s'",
				t.FullName, t.SpecFile)
		}
		m[t.Name] = true
	}
	for _, v := range ip.ViewsData {
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
	for _, t := range ip.TablesData {
		t.CreateDeps, err = ip.pdeps.createDependenciesData(t.FullName, false)
		if err != nil {
			return err
		}
		t.DropDeps, err = ip.pdeps.createDependenciesData(t.FullName, true)
		if err != nil {
			return err
		}

		t.CreateDepsIncludeTable = hasTableDep(t.CreateDeps)
		t.DropDepsIncludeTable = hasTableDep(t.DropDeps)
	}

	for _, v := range ip.ViewsData {
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

func hasTableDep(l []dependencyData) bool {
	for _, v := range l {
		if v.Table {
			return true
		}
	}
	return false
}
