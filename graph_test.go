package csv2postgres

import (
	"testing"

	"github.com/pkg/errors"
	"github.com/yourbasic/graph"
)

type DepsGraph struct {
	N        int
	g        *graph.Mutable
	g1       *graph.Mutable // Transpose of g
	topSort  []int
	topSort1 []int
}

func NewDeps(n int) *DepsGraph {
	return &DepsGraph{
		N:  n,
		g:  graph.New(n),
		g1: graph.New(n),
	}
}

// DependsOn creates new dependency: a depends to b
func (d *DepsGraph) DependsOn(a, b int) error {
	if a < 1 || b < 1 || a > d.N || b > d.N {
		return errors.New("invalid a or b value")
	}
	d.g.AddCost(a, b, 1)
	d.g1.AddCost(b, a, 1)
	return nil
}

func (d *DepsGraph) Finalize() error {
	// if !graph.Acyclic(d.g) {
	// 	return errors.New("must be acyclic dependencies")
	// }
	// ends := make([]int, 0)
	// for i := 1; i < d.End; i++ {
	// 	if d.g.Degree(i) == 0 {
	// 		ends = append(ends, i)
	// 	}
	// }

	// starts := make([]int, 0)
	// for i := 1; i < d.End; i++ {
	// 	if d.g1.Degree(i) == 0 {
	// 		starts = append(starts, i)
	// 	}
	// }

	// for _, i := range ends {
	// 	d.g.AddCost(i, d.End, 1)
	// 	d.g1.AddCost(d.Start, i, 1)
	// }

	// for _, i := range starts {
	// 	d.g.AddCost(d.Start, i, 1)
	// 	d.g1.AddCost(i, d.End, 1)
	// }

	ts, ok := graph.TopSort(d.g)
	if !ok {
		return errors.New("must be acyclic dependencies")
	}

	ts1, ok := graph.TopSort(d.g1)
	if !ok {
		return errors.New("must be acyclic dependencies")
	}

	d.topSort = ts
	d.topSort1 = ts1

	return nil
}

func (d *DepsGraph) CreateRoute(v int) (path []int) {
	m := make(map[int]bool)
	m[v] = true
	graph.BFS(d.g, v, func(v, w int, c int64) {
		m[w] = true
	})

	result := make([]int, 0, len(m))
	for _, i := range d.topSort {
		_, ok := m[i]
		if ok {
			result = append(result, i)
		}
	}
	return reverse(result)
}

func (d *DepsGraph) DropRoute(v int) (path []int) {
	m := make(map[int]bool)
	m[v] = true
	graph.BFS(d.g1, v, func(v, w int, c int64) {
		m[w] = true
	})

	result := make([]int, 0, len(m))
	for _, i := range d.topSort1 {
		_, ok := m[i]
		if ok {
			result = append(result, i)
		}
	}
	return reverse(result)
}

func reverse(numbers []int) []int {
	for i := 0; i < len(numbers)/2; i++ {
		j := len(numbers) - i - 1
		numbers[i], numbers[j] = numbers[j], numbers[i]
	}
	return numbers
}

func TestA(t *testing.T) {
	// Notes:
	// - Vertex 0 is start all target
	// - Last vertext is the end (6) connected to leaves
	// - all targets: use TopSort
	// - DependsOn:
	// - Dependants:

	deps := NewDeps(5)
	err := deps.DependsOn(1, 2)
	if err != nil {
		t.Error(err)
	}

	err = deps.DependsOn(1, 3)
	if err != nil {
		t.Error(err)
	}

	err = deps.DependsOn(1, 4)
	if err != nil {
		t.Error(err)
	}

	err = deps.DependsOn(2, 3)
	if err != nil {
		t.Error(err)
	}

	err = deps.DependsOn(2, 4)
	if err != nil {
		t.Error(err)
	}

	err = deps.DependsOn(4, 3)
	if err != nil {
		t.Error(err)
	}

	err = deps.Finalize()
	if err != nil {
		t.Error(err)
	}

	p := deps.CreateRoute(1)
	t.Errorf("CreateRoute 1: %+v", p)

	p = deps.CreateRoute(3)
	t.Errorf("CreateRoute 3: %+v", p)

	p = deps.DropRoute(1)
	t.Errorf("DropRoute 3: %+v", p)

	p = deps.DropRoute(3)
	t.Errorf("DropRoute 3: %+v", p)
}
