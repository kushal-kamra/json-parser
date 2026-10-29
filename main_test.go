package main

import (
	"bytes"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestValidateJSON(t *testing.T) {
	tests := []struct {
		input   string
		isValid bool
	}{
		{input: "{}", isValid: true},
		{input: " \t{\n\r} ", isValid: true},
		{input: "", isValid: false},
		{input: "[]", isValid: false},
		{input: "{ x }", isValid: false},
		{input: "{}{}", isValid: false},
		{input: "{} trailing", isValid: false},
		{input: "{\u00a0}", isValid: false},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			isValid, errorMessage := ValidateJSON(tt.input)
			if isValid != tt.isValid {
				t.Errorf("ValidateJSON(%q) validity = %t, want %t", tt.input, isValid, tt.isValid)
			}
			if tt.isValid && errorMessage != "" {
				t.Errorf("ValidateJSON(%q) error message = %q, want empty", tt.input, errorMessage)
			}
			if !tt.isValid && errorMessage == "" {
				t.Errorf("ValidateJSON(%q) returned no error message for invalid input", tt.input)
			}
		})
	}
}

func TestCLI(t *testing.T) {
	tempDir := t.TempDir()
	binaryPath := filepath.Join(tempDir, "json-parser")
	build := exec.Command("go", "build", "-o", binaryPath, ".")
	if output, err := build.CombinedOutput(); err != nil {
		t.Fatalf("build CLI: %v\n%s", err, output)
	}

	validPath := filepath.Join(tempDir, "valid.json")
	if err := os.WriteFile(validPath, []byte("{}"), 0600); err != nil {
		t.Fatalf("write valid input: %v", err)
	}
	invalidPath := filepath.Join(tempDir, "invalid.json")
	if err := os.WriteFile(invalidPath, []byte("{"), 0600); err != nil {
		t.Fatalf("write invalid input: %v", err)
	}
	invalidObjectPath := filepath.Join(tempDir, "invalid-object.json")
	if err := os.WriteFile(invalidObjectPath, []byte("{ x }"), 0600); err != nil {
		t.Fatalf("write invalid object input: %v", err)
	}

	tests := []struct {
		name           string
		args           []string
		wantStdout     string
		wantStderr     string
		wantStderrPart string
		wantExit       int
	}{
		{name: "valid JSON", args: []string{validPath}, wantStdout: "Valid JSON\n", wantExit: 0},
		{name: "invalid JSON", args: []string{invalidObjectPath}, wantStderr: "Invalid JSON: expected only whitespace between object braces\n", wantExit: 1},
		{name: "no file path", wantStderr: "Usage: json-parser <file>\n", wantExit: 2},
		{name: "multiple file paths", args: []string{validPath, invalidPath}, wantStderr: "Usage: json-parser <file>\n", wantExit: 2},
		{name: "file cannot be read", args: []string{filepath.Join(tempDir, "missing.json")}, wantStderrPart: "Error reading file:", wantExit: 2},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := exec.Command(binaryPath, tt.args...)
			var stdout, stderr bytes.Buffer
			cmd.Stdout = &stdout
			cmd.Stderr = &stderr

			err := cmd.Run()
			exitCode := 0
			if err != nil {
				var exitErr *exec.ExitError
				if !errors.As(err, &exitErr) {
					t.Fatalf("run CLI: %v", err)
				}
				exitCode = exitErr.ExitCode()
			}

			if stdout.String() != tt.wantStdout {
				t.Errorf("stdout = %q, want %q", stdout.String(), tt.wantStdout)
			}
			if tt.wantStderrPart != "" {
				if !strings.Contains(stderr.String(), tt.wantStderrPart) {
					t.Errorf("stderr = %q, want it to contain %q", stderr.String(), tt.wantStderrPart)
				}
			} else if stderr.String() != tt.wantStderr {
				t.Errorf("stderr = %q, want %q", stderr.String(), tt.wantStderr)
			}
			if exitCode != tt.wantExit {
				t.Errorf("exit status = %d, want %d", exitCode, tt.wantExit)
			}
		})
	}
}
