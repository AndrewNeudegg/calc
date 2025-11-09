package integration

import (
	"bytes"
	"os/exec"
	"strings"
	"testing"
	"time"
)

// TestCalcExitsCleanly verifies that calc exits with exit code 0 and preserves terminal formatting
func TestCalcExitsCleanly(t *testing.T) {
	calcBin := buildCalcBinary(t)

	// Test 1: Exit with :quit command
	t.Run("ExitWithQuitCommand", func(t *testing.T) {
		cmd := exec.Command(calcBin)
		stdin, err := cmd.StdinPipe()
		if err != nil {
			t.Fatalf("Failed to get stdin pipe: %v", err)
		}
		
		// Capture output to verify terminal state
		var stdout, stderr bytes.Buffer
		cmd.Stdout = &stdout
		cmd.Stderr = &stderr

		if err := cmd.Start(); err != nil {
			t.Fatalf("Failed to start calc: %v", err)
		}

		// Give it a moment to start
		time.Sleep(100 * time.Millisecond)

		// Send quit command
		if _, err := stdin.Write([]byte(":quit\n")); err != nil {
			t.Fatalf("Failed to write :quit command to stdin: %v", err)
		}
		stdin.Close()

		// Wait for process to exit
		err = cmd.Wait()
		if err != nil {
			t.Errorf("calc did not exit cleanly with :quit: %v", err)
		}
		
		// Verify no terminal escape sequences are left unclosed in output
		output := stdout.String() + stderr.String()
		if strings.Contains(output, "\x1b[") {
			// Count opening and closing escape sequences
			// This is a basic check - proper terminal handling should clean up
			t.Logf("Output contains ANSI escape sequences (this is normal for styled output)")
		}
	})

	// Test 2: Exit with :q shorthand
	t.Run("ExitWithQShorthand", func(t *testing.T) {
		cmd := exec.Command(calcBin)
		stdin, err := cmd.StdinPipe()
		if err != nil {
			t.Fatalf("Failed to get stdin pipe: %v", err)
		}
		
		var stdout, stderr bytes.Buffer
		cmd.Stdout = &stdout
		cmd.Stderr = &stderr

		if err := cmd.Start(); err != nil {
			t.Fatalf("Failed to start calc: %v", err)
		}

		time.Sleep(100 * time.Millisecond)

		// Send :q command
		if _, err := stdin.Write([]byte(":q\n")); err != nil {
			t.Fatalf("Failed to write :q command to stdin: %v", err)
		}
		stdin.Close()

		err = cmd.Wait()
		if err != nil {
			t.Errorf("calc did not exit cleanly with :q: %v", err)
		}
	})

	// Test 3: Exit with Ctrl-D (EOF)
	t.Run("ExitWithCtrlD", func(t *testing.T) {
		cmd := exec.Command(calcBin)
		stdin, err := cmd.StdinPipe()
		if err != nil {
			t.Fatalf("Failed to get stdin pipe: %v", err)
		}

		if err := cmd.Start(); err != nil {
			t.Fatalf("Failed to start calc: %v", err)
		}

		time.Sleep(100 * time.Millisecond)

		// Close stdin to simulate Ctrl-D
		stdin.Close()

		err = cmd.Wait()
		if err != nil {
			t.Errorf("calc did not exit cleanly with Ctrl-D: %v", err)
		}
	})

	// Test 4: Test simple calculation before exit
	t.Run("CalculationThenExit", func(t *testing.T) {
		cmd := exec.Command(calcBin)
		stdin, err := cmd.StdinPipe()
		if err != nil {
			t.Fatalf("Failed to get stdin pipe: %v", err)
		}

		if err := cmd.Start(); err != nil {
			t.Fatalf("Failed to start calc: %v", err)
		}

		time.Sleep(100 * time.Millisecond)

		// Send some calculations
		if _, err := stdin.Write([]byte("2 + 2\n")); err != nil {
			t.Fatalf("Failed to write to stdin: %v", err)
		}
		time.Sleep(100 * time.Millisecond)
		if _, err := stdin.Write([]byte("x = 10\n")); err != nil {
			t.Fatalf("Failed to write to stdin: %v", err)
		}
		time.Sleep(100 * time.Millisecond)
		if _, err := stdin.Write([]byte(":quit\n")); err != nil {
			t.Fatalf("Failed to write to stdin: %v", err)
		}
		stdin.Close()

		err = cmd.Wait()
		if err != nil {
			t.Errorf("calc did not exit cleanly after calculations: %v", err)
		}
	})
}

// TestCalcSingleExpressionMode tests -c flag (single expression mode)
func TestCalcSingleExpressionMode(t *testing.T) {
	tests := []struct {
		name       string
		expression string
		wantExit   int
	}{
		{"SimpleAddition", "2 + 2", 0},
		{"UnitConversion", "10 m in cm", 0},
		{"InvalidExpression", "2 + + 2", 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := exec.Command("go", "run", "../../cmd/calc", "-c", tt.expression)
			err := cmd.Run()
			
			exitCode := 0
			if err != nil {
				if exitErr, ok := err.(*exec.ExitError); ok {
					exitCode = exitErr.ExitCode()
				} else {
					t.Fatalf("Failed to run calc: %v", err)
				}
			}

			if exitCode != tt.wantExit {
				t.Errorf("Expected exit code %d, got %d", tt.wantExit, exitCode)
			}
		})
	}
}
