package deps

import (
	"errors"
	"fmt"

	"github.com/yourbasic/graph"
)

type depsStore struct {
	l []string
	m map[string]int
}

func newDepsStore(l []string) depsStore {
	m := make(map[string]int, len(l))
	for i, v := range l {
		m[v] = i
	}
	return depsStore{l, m}
}

func (d depsStore) IndexOf(s string) (index int, ok bool) {
	index, ok = d.m[s]
	return index, ok
}

func (d depsStore) NameOF(index int) (name string, ok bool) {
	if index >= len(d.l) {
		return "", false
	}
	return d.l[index], true
}

type dependencies struct {
	N        int
	g        *graph.Mutable
	g1       *graph.Mutable // Transpose of g
	topSort  []int
	topSort1 []int
}

func newDependencies(n int) *dependencies {
	return &dependencies{
		N:  n,
		g:  graph.New(n),
		g1: graph.New(n),
	}
}

// DependsOn creates new dependency: a depends to b
func (d *dependencies) dependsOn(a, b int) error {
	if a < 0 || b < 0 || a >= d.N || b >= d.N {
		return errors.New("invalid a or b value")
	}
	d.g.AddCost(a, b, 1)
	d.g1.AddCost(b, a, 1)
	return nil
}

// Finalize finalizes dependency graph. Querying dependency order can only be
// calculated if it has been finalizes.
func (d *dependencies) finalize() error {
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

func (d dependencies) createOrderAll() (path []int) {
	return d.topSort1
}

// CreateOrder returns list of creation order.
func (d dependencies) createOrder(v int) (path []int) {
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

func (d dependencies) dropOrderAll() (path []int) {
	return d.topSort
}

// DropOrder returns list of deletion order.
func (d dependencies) dropOrder(v int) (path []int) {
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

// Reverse reverses array order.
func reverse(numbers []int) []int {
	for i := 0; i < len(numbers)/2; i++ {
		j := len(numbers) - i - 1
		numbers[i], numbers[j] = numbers[j], numbers[i]
	}
	return numbers
}

// Graph specifies dependencies graph.
type Graph struct {
	store depsStore
	deps  *dependencies
}

// NewGraph returns a new Graph.
func NewGraph(targets []string) *Graph {
	return &Graph{
		newDepsStore(targets),
		newDependencies(len(targets)),
	}
}

// DependsOn creates new dependency: a depends to b
func (d *Graph) DependsOn(a, b string) error {
	aIdx, ok := d.store.IndexOf(a)
	if !ok {
		return fmt.Errorf("unknown target: %s", a)
	}

	bIdx, ok := d.store.IndexOf(b)
	if !ok {
		return fmt.Errorf("unknown target: %s", b)
	}

	err := d.deps.dependsOn(aIdx, bIdx)
	if err != nil {
		return err
	}
	return nil
}

// Finalize finalizes dependency graph. Querying dependency order can only be
// calculated if it has been finalizes.
func (d *Graph) Finalize() error {
	return d.deps.finalize()
}

// CreateOrderAll returns list of creation order.
func (d Graph) CreateOrderAll() (path []string) {
	p := d.deps.createOrderAll()
	path = make([]string, len(p))
	for i, idx := range p {
		t, ok := d.store.NameOF(idx)
		if !ok {
			panic(fmt.Sprintf("invalid dependency index: %d", idx))
		}
		path[i] = t
	}
	return path
}

// CreateOrder returns list of creation order for target.
func (d Graph) CreateOrder(target string) (path []string, err error) {
	idx, ok := d.store.IndexOf(target)
	if !ok {
		return nil, fmt.Errorf("unknown target: %s", target)
	}

	p := d.deps.createOrder(idx)
	path = make([]string, len(p))
	for i, idx := range p {
		t, ok := d.store.NameOF(idx)
		if !ok {
			panic(fmt.Sprintf("invalid dependency index: %d", idx))
		}
		path[i] = t
	}
	return path, nil
}

// DropOrderAll returns list of deletion order.
func (d Graph) DropOrderAll() (path []string) {
	p := d.deps.dropOrderAll()
	path = make([]string, len(p))
	for i, idx := range p {
		t, ok := d.store.NameOF(idx)
		if !ok {
			panic(fmt.Sprintf("invalid dependency index: %d", idx))
		}
		path[i] = t
	}
	return path
}

// DropOrder returns list of deletion order for target.
func (d Graph) DropOrder(target string) (path []string, err error) {
	idx, ok := d.store.IndexOf(target)
	if !ok {
		return nil, fmt.Errorf("unknown target: %s", target)
	}

	p := d.deps.dropOrder(idx)
	path = make([]string, len(p))
	for i, idx := range p {
		t, ok := d.store.NameOF(idx)
		if !ok {
			panic(fmt.Sprintf("invalid dependency index: %d", idx))
		}
		path[i] = t
	}
	return path, nil
}
