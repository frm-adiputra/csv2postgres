package interpolation

import (
	"fmt"
	"sort"

	"github.com/frm-adiputra/csv2postgres/schema"
)

// Interpolator represents interpolator struct
type Interpolator struct {
	Generator  string
	ImportPath string
	TablesData []*TableData
	ViewsData  []*ViewData
	HasExport  bool

	CreateDepsAll []dependencyData
	DropDepsAll   []dependencyData

	pdeps         *projectDeps
	defaultSchema string
	nameToTarget  map[string]string
	targetToName  map[string]string
}

// NewInterpolator returns a new interpolator
func NewInterpolator(baseImportPath, rootDir, defaultSchema string, ts []*schema.Table, vs []*schema.View) (*Interpolator, error) {
	var err error

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

	nameToTarget, targetToName := generateTargetNames(lt, vt)
	pdeps, err := newProjectDeps(ts, vs, nameToTarget)
	if err != nil {
		return nil, err
	}

	ip := &Interpolator{
		Generator:     "github.com/frm-adiputra/csv2postgres",
		ImportPath:    baseImportPath,
		TablesData:    lt,
		ViewsData:     vt,
		CreateDepsAll: pdeps.createAllDependenciesData(false),
		DropDepsAll:   pdeps.createAllDependenciesData(true),
		HasExport:     pdeps.HasExport(),
		pdeps:         pdeps,
		defaultSchema: defaultSchema,
		nameToTarget:  nameToTarget,
		targetToName:  targetToName,
	}

	// ip.generateTargetNames()
	return ip, nil
}

// Interpolate do interpolation based on specification files provided
func (ip *Interpolator) Interpolate() error {
	err := ip.checkDuplicateName()
	if err != nil {
		return err
	}

	err = ip.linkDependencies()
	if err != nil {
		return err
	}
	return nil
}

func generateTargetNames(ts []*TableData, vs []*ViewData) (map[string]string, map[string]string) {
	nameToTarget := make(map[string]string)
	targetToName := make(map[string]string)
	i := 0
	for _, t := range ts {
		i++
		target := fmt.Sprintf("Target%03d", i)
		t.TargetName = target
		nameToTarget[t.RefName] = target
		targetToName[target] = t.RefName
	}

	for _, t := range vs {
		i++
		target := fmt.Sprintf("Target%03d", i)
		t.TargetName = target
		nameToTarget[t.RefName] = target
		targetToName[target] = t.RefName
	}
	return nameToTarget, targetToName
}

func (ip *Interpolator) generateTargetNames() {
	i := 0
	for _, t := range ip.TablesData {
		i++
		target := fmt.Sprintf("Target%03d", i)
		t.TargetName = target
		ip.nameToTarget[t.RefName] = target
		// ip.targetToName[target] = t.RefName
	}

	for _, t := range ip.ViewsData {
		i++
		target := fmt.Sprintf("Target%03d", i)
		t.TargetName = target
		ip.nameToTarget[t.RefName] = target
		// ip.targetToName[target] = t.RefName
	}
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
	fmt.Println(">>>> Link")
	var err error
	for _, t := range ip.TablesData {
		t.CreateDeps, err = ip.pdeps.createDependenciesData(t.RefName, false)
		if err != nil {
			return err
		}
		t.DropDeps, err = ip.pdeps.createDependenciesData(t.RefName, true)
		if err != nil {
			return err
		}

		t.CreateDepsIncludeTable = hasTableDep(t.CreateDeps)
		t.DropDepsIncludeTable = hasTableDep(t.DropDeps)
	}

	for _, v := range ip.ViewsData {
		v.CreateDeps, err = ip.pdeps.createDependenciesData(v.RefName, false)
		if err != nil {
			return err
		}
		v.DropDeps, err = ip.pdeps.createDependenciesData(v.RefName, true)
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

// DisplayTargets displays targets info to stdout
func (ip Interpolator) DisplayTargets() {
	l := make([]string, 0, len(ip.targetToName))
	for k := range ip.targetToName {
		l = append(l, k)
	}
	sort.Strings(l)
	for _, v := range l {
		fmt.Printf("%s: %s\n", v, ip.targetToName[v])
	}
}
