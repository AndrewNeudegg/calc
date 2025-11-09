package display

import (
	"os"
	"testing"
)

// TestIsATTY tests the TTY detection function
func TestIsATTY(t *testing.T) {
	// Test with stdin - in a test environment, this is typically not a TTY
	stdinFd := os.Stdin.Fd()
	// We don't assert true/false since it depends on test environment
	// Just ensure the function doesn't panic
	_ = isATTY(stdinFd)
	
	// Test with an invalid fd (should return false)
	invalidFd := uintptr(9999)
	if isATTY(invalidFd) {
		t.Error("Expected isATTY to return false for invalid fd")
	}
}

// TestRawModeEnableRestore tests enabling and restoring raw mode
func TestRawModeEnableRestore(t *testing.T) {
	// Skip if not running in a TTY environment
	if !isATTY(os.Stdin.Fd()) {
		t.Skip("Skipping test: not running in a TTY environment")
	}
	
	fd := int(os.Stdin.Fd())
	
	// Enable raw mode
	state, err := enableRawMode(fd)
	if err != nil {
		t.Fatalf("Failed to enable raw mode: %v", err)
	}
	
	// Ensure we got a valid state
	if state == nil {
		t.Fatal("Expected non-nil state from enableRawMode")
	}
	
	// Restore immediately
	restoreRawMode(fd, state)
	
	// Test that restore with nil state doesn't panic
	restoreRawMode(fd, nil)
}

// TestRawModeRestoreNilState tests that restoring with a nil state is safe
func TestRawModeRestoreNilState(t *testing.T) {
	// This should not panic
	restoreRawMode(int(os.Stdin.Fd()), nil)
}

// TestRawModeInvalidFd tests error handling with invalid file descriptor
func TestRawModeInvalidFd(t *testing.T) {
	invalidFd := -1
	_, err := enableRawMode(invalidFd)
	if err == nil {
		t.Error("Expected error when enabling raw mode on invalid fd")
	}
}
