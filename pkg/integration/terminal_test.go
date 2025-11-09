package integration

import (
	"os/exec"
	"testing"
	"time"
)

// TestCalcExitsCleanly verifies that calc exits cleanly without breaking terminal formatting
func TestCalcExitsCleanly(t *testing.T) {
	// Build the calc binary
	buildCmd := exec.Command("go", "build", "-o", "/tmp/calc_test", "../../cmd/calc")
	if err := buildCmd.Run(); err != nil {
		t.Fatalf("Failed to build calc: %v", err)
	}

	// Test 1: Exit with :quit command
	t.Run("ExitWithQuitCommand", func(t *testing.T) {
		cmd := exec.Command("/tmp/calc_test")
		stdin, err := cmd.StdinPipe()
		if err != nil {
			t.Fatalf("Failed to get stdin pipe: %v", err)
		}

		if err := cmd.Start(); err != nil {
			t.Fatalf("Failed to start calc: %v", err)
		}

		// Give it a moment to start
		time.Sleep(100 * time.Millisecond)

		// Send quit command
		stdin.Write([]byte(":quit\n"))
		stdin.Close()

		// Wait for process to exit
		err = cmd.Wait()
		if err != nil {
			t.Errorf("calc did not exit cleanly with :quit: %v", err)
		}
	})

	// Test 2: Exit with :q shorthand
	t.Run("ExitWithQShorthand", func(t *testing.T) {
		cmd := exec.Command("/tmp/calc_test")
		stdin, err := cmd.StdinPipe()
		if err != nil {
			t.Fatalf("Failed to get stdin pipe: %v", err)
		}

		if err := cmd.Start(); err != nil {
			t.Fatalf("Failed to start calc: %v", err)
		}

		time.Sleep(100 * time.Millisecond)

		// Send :q command
		stdin.Write([]byte(":q\n"))
		stdin.Close()

		err = cmd.Wait()
		if err != nil {
			t.Errorf("calc did not exit cleanly with :q: %v", err)
		}
	})

	// Test 3: Exit with Ctrl-D (EOF)
	t.Run("ExitWithCtrlD", func(t *testing.T) {
		cmd := exec.Command("/tmp/calc_test")
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
		cmd := exec.Command("/tmp/calc_test")
		stdin, err := cmd.StdinPipe()
		if err != nil {
			t.Fatalf("Failed to get stdin pipe: %v", err)
		}

		if err := cmd.Start(); err != nil {
			t.Fatalf("Failed to start calc: %v", err)
		}

		time.Sleep(100 * time.Millisecond)

		// Send some calculations
		stdin.Write([]byte("2 + 2\n"))
		time.Sleep(100 * time.Millisecond)
		stdin.Write([]byte("x = 10\n"))
		time.Sleep(100 * time.Millisecond)
		stdin.Write([]byte(":quit\n"))
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
