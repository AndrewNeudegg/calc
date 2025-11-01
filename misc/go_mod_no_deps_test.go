package misc

import (
	"os"
	"strings"
	"testing"
)

func TestGoModHasNoDependencies(t *testing.T) {
	data, err := os.ReadFile("../go.mod")
	if err != nil {
		t.Fatalf("failed to read go.mod: %v", err)
	}
	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		if strings.HasPrefix(strings.TrimSpace(line), "require ") || strings.HasPrefix(strings.TrimSpace(line), "replace ") {
			t.Errorf("go.mod should not have dependencies, found: %s", line)
		}
	}
}
