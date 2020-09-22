package interpolation

import (
	"github.com/frm-adiputra/csv2postgres/deps"
	"github.com/frm-adiputra/csv2postgres/schema"
)

type projectDeps struct {
	graph           *deps.Graph
	targetIsTable   map[string]bool
	targetHasExport map[string]bool
}

type dependencyData struct {
	Name   string
	Table  bool
	Export bool
}

func newProjectDeps(ts []*schema.Table, vs []*schema.View) (*projectDeps, error) {
	targets := make([]string, 0, len(ts)+len(vs))
	targetIsTable := make(map[string]bool)
	targetHasExport := make(map[string]bool)
	for _, s := range ts {
		targets = append(targets, s.Name)
		targetIsTable[s.Name] = true
	}
	for _, v := range vs {
		targets = append(targets, v.Name)
		if v.Export != "" {
			targetHasExport[v.Name] = true
		}
	}

	d := deps.NewGraph(targets)

	for _, s := range ts {
		for _, ds := range s.DependsOn {
			err := d.DependsOn(s.Name, ds)
			if err != nil {
				return nil, err
			}
		}
	}

	for _, s := range vs {
		for _, ds := range s.DependsOn {
			err := d.DependsOn(s.Name, ds)
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
			Name:   v,
			Table:  isTable,
			Export: isExport,
		}
	}
	return l, nil
}
