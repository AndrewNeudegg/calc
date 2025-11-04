package display

import "testing"

func TestCtrlCTipContainsGuidance(t *testing.T) {
	tip := ctrlCTip()
	if tip == "" {
		t.Fatalf("ctrlCTip should not be empty")
	}
	if !containsAll(tip, []string{"Ctrl-C", "Ctrl-D", ":help"}) {
		t.Fatalf("ctrlCTip missing expected guidance, got: %q", tip)
	}
}

func containsAll(s string, parts []string) bool {
	for _, p := range parts {
		if !contains(s, p) {
			return false
		}
	}
	return true
}

func contains(s, sub string) bool { return len(s) >= len(sub) && (indexOf(s, sub) >= 0) }

// simple substring search to avoid pulling additional packages
func indexOf(s, sub string) int {
	// naive search is fine for a tiny literal
	n, m := len(s), len(sub)
	if m == 0 {
		return 0
	}
	for i := 0; i+m <= n; i++ {
		if s[i:i+m] == sub {
			return i
		}
	}
	return -1
}
