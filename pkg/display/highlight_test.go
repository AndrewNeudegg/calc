package display

import (
	"strings"
	"testing"
)

func TestHighlighterColorsCommand(t *testing.T) {
	th := DefaultTheme()
	h := NewHighlighter(th)
	out := h.Colorize(":help")
	if !strings.Contains(out, th.Command) || !strings.Contains(out, th.Reset) {
		t.Fatalf("expected command to be wrapped in command color and reset: %q", out)
	}
}

func TestHighlighterColorsNumbersAndKeywordsAndUnits(t *testing.T) {
	th := DefaultTheme()
	h := NewHighlighter(th)
	s := "2 km in m"
	out := h.Colorize(s)
	// Expect presence of number, unit, keyword styles at least once
	if !strings.Contains(out, th.Number) {
		t.Fatalf("expected number color present: %q", out)
	}
	if !strings.Contains(out, th.Unit) {
		t.Fatalf("expected unit color present: %q", out)
	}
	if !strings.Contains(out, th.Keyword) {
		t.Fatalf("expected keyword color present: %q", out)
	}
}
