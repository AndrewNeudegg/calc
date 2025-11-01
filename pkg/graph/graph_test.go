package graph

import (
	"testing"
)

func TestAddNode(t *testing.T) {
	g := NewGraph()
	g.AddNode(1, "x = 10", nil)
	g.AddNode(2, "y = 20", nil)
	
	if _, exists := g.GetNode(1); !exists {
		t.Error("Node 1 should exist")
	}
	
	if _, exists := g.GetNode(2); !exists {
		t.Error("Node 2 should exist")
	}
}

func TestGetNode(t *testing.T) {
	g := NewGraph()
	g.AddNode(1, "x = 10", nil)
	
	node, exists := g.GetNode(1)
	if !exists {
		t.Fatal("Node should exist")
	}
	
	if node.Expression != "x = 10" {
		t.Errorf("Expected expression 'x = 10', got %q", node.Expression)
	}
}

func TestGetDependents(t *testing.T) {
	g := NewGraph()
	g.AddNode(1, "x = 10", nil)
	g.AddNode(2, "y = x + 5", []int{1})
	g.AddNode(3, "z = x * 2", []int{1})
	
	dependents := g.GetDependents(1)
	
	if len(dependents) != 2 {
		t.Errorf("Expected 2 dependents, got %d", len(dependents))
	}
	
	has2 := false
	has3 := false
	for _, dep := range dependents {
		if dep == 2 {
			has2 = true
		}
		if dep == 3 {
			has3 = true
		}
	}
	
	if !has2 || !has3 {
		t.Error("Dependents should include nodes 2 and 3")
	}
}

func TestHasCycle(t *testing.T) {
	tests := []struct {
		name      string
		setup     func(*Graph)
		wantCycle bool
	}{
		{
			name: "no cycle",
			setup: func(g *Graph) {
				g.AddNode(1, "a = 10", nil)
				g.AddNode(2, "b = a + 5", []int{1})
				g.AddNode(3, "c = b * 2", []int{2})
			},
			wantCycle: false,
		},
		{
			name: "simple cycle",
			setup: func(g *Graph) {
				g.AddNode(1, "a = b", []int{2})
				g.AddNode(2, "b = a", []int{1})
			},
			wantCycle: true,
		},
		{
			name: "long cycle",
			setup: func(g *Graph) {
				g.AddNode(1, "a = c", []int{3})
				g.AddNode(2, "b = a", []int{1})
				g.AddNode(3, "c = b", []int{2})
			},
			wantCycle: true,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := NewGraph()
			tt.setup(g)
			
			hasCycle := g.HasCycle()
			if hasCycle != tt.wantCycle {
				t.Errorf("HasCycle() = %v, want %v", hasCycle, tt.wantCycle)
			}
		})
	}
}

func TestTopologicalSort(t *testing.T) {
	g := NewGraph()
	g.AddNode(1, "a = 10", nil)
	g.AddNode(2, "b = a + 5", []int{1})
	g.AddNode(3, "c = b * 2", []int{2})
	
	order, err := g.TopologicalSort()
	if err != nil {
		t.Fatalf("TopologicalSort failed: %v", err)
	}
	
	if len(order) != 3 {
		t.Errorf("Expected 3 nodes in order, got %d", len(order))
	}
	
	// Find indices
	aIndex := -1
	bIndex := -1
	cIndex := -1
	
	for i, id := range order {
		switch id {
		case 1:
			aIndex = i
		case 2:
			bIndex = i
		case 3:
			cIndex = i
		}
	}
	
	// Check dependencies are satisfied (deps come before dependents)
	if aIndex > bIndex {
		t.Error("Node 1 (a) should come before node 2 (b)")
	}
	
	if bIndex > cIndex {
		t.Error("Node 2 (b) should come before node 3 (c)")
	}
}

func TestTopologicalSortWithCycle(t *testing.T) {
	g := NewGraph()
	g.AddNode(1, "a = b", []int{2})
	g.AddNode(2, "b = a", []int{1})
	
	_, err := g.TopologicalSort()
	
	if err == nil {
		t.Error("TopologicalSort should return error for graph with cycle")
	}
}

func TestClear(t *testing.T) {
	g := NewGraph()
	g.AddNode(1, "x = 10", nil)
	g.AddNode(2, "y = 20", nil)
	
	g.Clear()
	
	if _, exists := g.GetNode(1); exists {
		t.Error("Node 1 should not exist after clear")
	}
	
	if _, exists := g.GetNode(2); exists {
		t.Error("Node 2 should not exist after clear")
	}
}
