package interpolation

import (
	"github.com/frm-adiputra/csv2postgres/deps"
	"github.com/frm-adiputra/csv2postgres/schema"
)

type projectDeps struct {
	graph           *deps.Graph
	targetIsTable   map[string]bool
	targetHasExport map[string]bool
	nameToTarget    map[string]string
}

type dependencyData struct {
	RefName    string
	TargetName string
	Table      bool
	Export     bool
}

func newProjectDeps(ts []*schema.Table, vs []*schema.View, nameToTarget map[string]string) (*projectDeps, error) {
	targets := make([]string, 0, len(ts)+len(vs))
	targetIsTable := make(map[string]bool)
	targetHasExport := make(map[string]bool)
	for _, s := range ts {
		targets = append(targets, s.RefName)
		targetIsTable[s.RefName] = true
	}
	for _, v := range vs {
		targets = append(targets, v.RefName)
		if v.Export != "" {
			targetHasExport[v.RefName] = true
		}
	}

	d := deps.NewGraph(targets)

	for _, s := range ts {
		for _, ds := range s.DependsOn {
			err := d.DependsOn(s.RefName, ds)
			if err != nil {
				return nil, err
			}
		}
	}

	for _, s := range vs {
		for _, ds := range s.DependsOn {
			err := d.DependsOn(s.RefName, ds)
			if err != nil {
				return nil, err
			}
		}
	}

	err := d.Finalize()
	if err != nil {
		return nil, err
	}

	return &projectDeps{
		graph:           d,
		nameToTarget:    nameToTarget,
		targetIsTable:   targetIsTable,
		targetHasExport: targetHasExport,
	}, nil
}

func (p projectDeps) createDependenciesData(name string, drop bool) ([]dependencyData, error) {
	var deps []string
	var err error

	if drop {
		deps, err = p.graph.DropOrder(name)
	} else {
		deps, err = p.graph.CreateOrder(name)
	}

	if err != nil {
		return nil, err
	}

	l := make([]dependencyData, len(deps))
	for i, v := range deps {
		_, isTable := p.targetIsTable[v]
		_, isExport := p.targetHasExport[v]
		l[i] = dependencyData{
			RefName:    v,
			TargetName: p.nameToTarget[v],
			Table:      isTable,
			Export:     isExport,
		}
	}
	return l, nil
}

func (p projectDeps) createAllDependenciesData(drop bool) []dependencyData {
	var deps []string

	if drop {
		deps = p.graph.DropOrderAll()
	} else {
		deps = p.graph.CreateOrderAll()
	}

	l := make([]dependencyData, len(deps))
	for i, v := range deps {
		_, isTable := p.targetIsTable[v]
		_, isExport := p.targetHasExport[v]
		l[i] = dependencyData{
			RefName:    v,
			TargetName: p.nameToTarget[v],
			Table:      isTable,
			Export:     isExport,
		}
	}
	return l
}

func (p projectDeps) HasExport() bool {
	return len(p.targetHasExport) != 0
}
