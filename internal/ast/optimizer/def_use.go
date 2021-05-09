package optimizer

import "github.com/ryan-berger/jpl/internal/ast"

type defUse struct {
	graph        map[string]map[ast.Node]bool
	reverseGraph map[ast.Node]map[string]bool
	parent       *defUse
	children     map[*ast.Function]*defUse
}

func (d *defUse) recordDef(key string) {
	d.graph[key] = make(map[ast.Node]bool)
}

func (d *defUse) recordUse(key string, n ast.Node) {
	g, ok := d.graph[key]
	if !ok {
		if d.parent == nil {
			panic("oopsies")
		}
		d.parent.recordUse(key, n)
		return
	}
	g[n] = true
	if _, ok := d.reverseGraph[n]; !ok {
		d.reverseGraph[n] = make(map[string]bool)
	}
	d.reverseGraph[n][key] = true
}

func (d *defUse) getUses(key string) []ast.Node {
	var nodes []ast.Node
	for n := range d.graph[key] {
		nodes = append(nodes, n)
	}

	for _, child := range d.children {
		nodes = append(nodes, child.getUses(key)...)
	}
	return nodes
}

func (d *defUse) clearUse(node ast.Node) {
	for ident := range d.reverseGraph[node] {
		delete(d.graph[ident], node)
	}
}

func (d *defUse) getIdentifiers(node ast.Node) []string {
	var idents []string
	for n := range d.reverseGraph[node] {
		idents = append(idents, n)
	}

	for _, child := range d.children {
		idents = append(idents, child.getIdentifiers(node)...)
	}
	return idents
}
