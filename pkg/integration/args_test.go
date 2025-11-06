package integration

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestArgDirective(t *testing.T) {
	// Build the calc binary for testing
	calcBin := buildCalcBinary(t)
	defer os.Remove(calcBin)

	tests := []struct {
		name       string
		script     string
		args       []string
		wantOutput []string
		wantError  bool
	}{
		{
			name: "simple arg with CLI value",
			script: `:arg count "Enter count:"
result = count * 2
print("Result: {result}")`,
			args: []string{"--arg", "count=5"},
			wantOutput: []string{
				"Result: 10.00",
			},
		},
		{
			name: "multiple args",
			script: `:arg a "Enter a:"
:arg b "Enter b:"
sum = a + b
print("Sum: {sum}")`,
			args: []string{"--arg", "a=10", "--arg", "b=20"},
			wantOutput: []string{
				"Sum: 30.00",
			},
		},
		{
			name: "arg with units",
			script: `:arg distance "Distance:"
:arg time "Time:"
speed = distance / time
print("Speed: {speed}")`,
			args: []string{"--arg", "distance=100 km", "--arg", "time=2 hours"},
			wantOutput: []string{
				"Speed: 50.00 km/hours",
			},
		},
		{
			name: "arg with currency",
			script: `:arg amount "Amount:"
doubled = amount * 2
print("Doubled: {doubled}")`,
			args: []string{"--arg", "amount=50 usd"},
			wantOutput: []string{
				"Doubled: $100.00",
			},
		},
		{
			name: "arg with expression",
			script: `:arg value "Value:"
result = value + 10
print("Result: {result}")`,
			args: []string{"--arg", "value=5*3"},
			wantOutput: []string{
				"Result: 25.00",
			},
		},
		{
			name: "short flag -a",
			script: `:arg x "Enter x:"
result = x * 2
print("Result: {result}")`,
			args: []string{"-a", "x=7"},
			wantOutput: []string{
				"Result: 14.00",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create temp script file
			tmpFile := createTempScript(t, tt.script)
			defer os.Remove(tmpFile)

			// Build command arguments
			cmdArgs := []string{"-f", tmpFile}
			cmdArgs = append(cmdArgs, tt.args...)

			// Execute calc
			cmd := exec.Command(calcBin, cmdArgs...)
			var stdout, stderr bytes.Buffer
			cmd.Stdout = &stdout
			cmd.Stderr = &stderr

			err := cmd.Run()
			if tt.wantError && err == nil {
				t.Error("expected error but got none")
				return
			}
			if !tt.wantError && err != nil {
				t.Errorf("unexpected error: %v\nstderr: %s\nstdout: %s", err, stderr.String(), stdout.String())
				return
			}

			output := stdout.String()
			if output == "" {
				t.Errorf("got empty output. stderr: %s", stderr.String())
				return
			}
			for _, want := range tt.wantOutput {
				if !strings.Contains(output, want) {
					t.Errorf("output missing expected string %q\ngot: %s\nstderr: %s", want, output, stderr.String())
				}
			}
		})
	}
}

func TestArgFile(t *testing.T) {
	calcBin := buildCalcBinary(t)
	defer os.Remove(calcBin)

	script := `:arg x "Enter x:"
:arg y "Enter y:"
result = x + y
print("Result: {result}")`

	argsFile := `x=15
y=25
`

	// Create temp script and args files
	tmpScript := createTempScript(t, script)
	defer os.Remove(tmpScript)

	tmpArgs := createTempScript(t, argsFile)
	defer os.Remove(tmpArgs)

	// Execute calc with --arg-file
	cmd := exec.Command(calcBin, "-f", tmpScript, "--arg-file", tmpArgs)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		t.Errorf("unexpected error: %v\nstderr: %s", err, stderr.String())
		return
	}

	output := stdout.String()
	if !strings.Contains(output, "Result: 40.00") {
		t.Errorf("output missing expected result\ngot: %s", output)
	}
}

func TestArgFileWithOverride(t *testing.T) {
	calcBin := buildCalcBinary(t)
	defer os.Remove(calcBin)

	script := `:arg x "Enter x:"
:arg y "Enter y:"
result = x + y
print("Result: {result}")`

	argsFile := `x=10
y=20
`

	// Create temp script and args files
	tmpScript := createTempScript(t, script)
	defer os.Remove(tmpScript)

	tmpArgs := createTempScript(t, argsFile)
	defer os.Remove(tmpArgs)

	// Execute calc with --arg-file and CLI override
	cmd := exec.Command(calcBin, "-f", tmpScript, "--arg-file", tmpArgs, "--arg", "y=30")
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		t.Errorf("unexpected error: %v\nstderr: %s", err, stderr.String())
		return
	}

	output := stdout.String()
	// y should be 30 (overridden), not 20 from file
	if !strings.Contains(output, "Result: 40.00") {
		t.Errorf("output missing expected result (should use CLI override)\ngot: %s", output)
	}
}

// Helper functions

func buildCalcBinary(t *testing.T) string {
	t.Helper()
	
	// Get repo root
	repoRoot, err := filepath.Abs("../..")
	if err != nil {
		t.Fatalf("failed to get repo root: %v", err)
	}

	tmpBin := filepath.Join(t.TempDir(), "calc")
	cmd := exec.Command("go", "build", "-a", "-o", tmpBin, filepath.Join(repoRoot, "cmd", "calc"))
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		t.Fatalf("failed to build calc binary: %v\nstderr: %s", err, stderr.String())
	}
	return tmpBin
}

func createTempScript(t *testing.T, content string) string {
	t.Helper()
	
	tmpFile, err := os.CreateTemp("", "calc-test-*.calc")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	
	if _, err := tmpFile.WriteString(content); err != nil {
		tmpFile.Close()
		os.Remove(tmpFile.Name())
		t.Fatalf("failed to write temp file: %v", err)
	}
	
	if err := tmpFile.Close(); err != nil {
		os.Remove(tmpFile.Name())
		t.Fatalf("failed to close temp file: %v", err)
	}
	
	return tmpFile.Name()
}
