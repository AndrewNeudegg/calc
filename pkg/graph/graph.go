package graph

import (
	"fmt"
)

// Node represents a line in the calculation graph.
type Node struct {
	ID           int
	Expression   string
	Dependencies []int
}

// Graph manages dependencies between calculation lines.
type Graph struct {
	nodes map[int]*Node
}

// NewGraph creates a new dependency graph.
func NewGraph() *Graph {
	return &Graph{
		nodes: make(map[int]*Node),
	}
}

// AddNode adds a node to the graph.
func (g *Graph) AddNode(id int, expr string, deps []int) {
	g.nodes[id] = &Node{
		ID:           id,
		Expression:   expr,
		Dependencies: deps,
	}
}

// GetNode retrieves a node by ID.
func (g *Graph) GetNode(id int) (*Node, bool) {
	node, ok := g.nodes[id]
	return node, ok
}

// GetDependents returns all nodes that depend on the given node.
func (g *Graph) GetDependents(id int) []int {
	var dependents []int
	for nid, node := range g.nodes {
		for _, dep := range node.Dependencies {
			if dep == id {
				dependents = append(dependents, nid)
				break
			}
		}
	}
	return dependents
}

// TopologicalSort returns nodes in evaluation order.
func (g *Graph) TopologicalSort() ([]int, error) {
	visited := make(map[int]bool)
	tempMark := make(map[int]bool)
	var result []int
	
	var visit func(int) error
	visit = func(id int) error {
		if tempMark[id] {
			return fmt.Errorf("circular dependency detected involving line %d", id)
		}
		
		if visited[id] {
			return nil
		}
		
		tempMark[id] = true
		
		node, ok := g.nodes[id]
		if ok {
			for _, dep := range node.Dependencies {
				if err := visit(dep); err != nil {
					return err
				}
			}
		}
		
		tempMark[id] = false
		visited[id] = true
		result = append(result, id)
		
		return nil
	}
	
	for id := range g.nodes {
		if !visited[id] {
			if err := visit(id); err != nil {
				return nil, err
			}
		}
	}
	
	return result, nil
}

// HasCycle checks if the graph contains a circular dependency.
func (g *Graph) HasCycle() bool {
	_, err := g.TopologicalSort()
	return err != nil
}

// Clear removes all nodes from the graph.
func (g *Graph) Clear() {
	g.nodes = make(map[int]*Node)
}
