package commands

import (
	"testing"

	"github.com/andrewneudegg/calc/pkg/settings"
)

// TestQuitCommand tests that quit command sets the shouldQuit flag
func TestQuitCommand(t *testing.T) {
	h := New(settings.Default())
	
	// Initially should not quit
	if h.ShouldQuit() {
		t.Error("Expected ShouldQuit to be false initially")
	}
	
	// Execute quit command
	msg := h.Execute("quit", nil)
	
	// Should return empty message
	if msg != "" {
		t.Errorf("Expected empty message from quit, got: %s", msg)
	}
	
	// Should now be set to quit
	if !h.ShouldQuit() {
		t.Error("Expected ShouldQuit to be true after quit command")
	}
}

// TestQuitCommandVariants tests all quit command variants
func TestQuitCommandVariants(t *testing.T) {
	variants := []string{"quit", "exit", "q", "QUIT", "EXIT", "Q"}
	
	for _, variant := range variants {
		t.Run(variant, func(t *testing.T) {
			h := New(settings.Default())
			msg := h.Execute(variant, nil)
			
			if msg != "" {
				t.Errorf("Expected empty message from %s, got: %s", variant, msg)
			}
			
			if !h.ShouldQuit() {
				t.Errorf("Expected ShouldQuit to be true after %s command", variant)
			}
		})
	}
}
