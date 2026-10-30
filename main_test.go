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
		{input: `{"key":"value"}`, isValid: true},
		{input: `{"a":"1","b":"2"}`, isValid: true},
		{input: `{"age":36}`, isValid: true},
		{input: `{"active":true}`, isValid: true},
		{input: `{"deleted":false}`, isValid: true},
		{input: `{"metadata":null}`, isValid: true},
		{input: `{"ratio":-1.25e+2}`, isValid: true},
		{input: `{"key":"value","key-n":101,"key-o":{},"key-l":[]}`, isValid: true},
		{input: `{"nested":{"array":[1,true,null,"text",{},[]]}}`, isValid: true},
		{input: `{"arrays":[[],[{"key":{"nested":[]}}]]}`, isValid: true},
		{input: `{"key":"value","other":"text"}`, isValid: true},
		{input: `{"key":"a,b:c"}`, isValid: true},
		{input: `{"key":"escaped quote: \"ok\""}`, isValid: true},
		{input: "", isValid: false},
		{input: "[]", isValid: false},
		{input: "{ x }", isValid: false},
		{input: `{"age":}`, isValid: false},
		{input: `{"age":+36}`, isValid: false},
		{input: `{"age":01}`, isValid: false},
		{input: `{"active":True}`, isValid: false},
		{input: `{"active":truth}`, isValid: false},
		{input: `{"age":1.}`, isValid: false},
		{input: `{"age":1e}`, isValid: false},
		{input: `{"array":[1,]}`, isValid: false},
		{input: `{"array":[,1]}`, isValid: false},
		{input: `{"array":[1 2]}`, isValid: false},
		{input: `{"array":[{"key":1]}`, isValid: false},
		{input: `{"object":{"key":1,}}`, isValid: false},
		{input: `{"object":{1:2}}`, isValid: false},
		{input: `{"key" "value"}`, isValid: false},
		{input: `{"key":}`, isValid: false},
		{input: `{"key":"value",}`, isValid: false},
		{input: `{"key":"value" "other":"text"}`, isValid: false},
		{input: "{}{}", isValid: false},
		{input: `{"key":"value"}{}`, isValid: false},
		{input: `{,"key":"value"}`, isValid: false},
		{input: `{"key":"value",}`, isValid: false},
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

func TestTokenize(t *testing.T) {
	tests := []struct {
		name       string
		input      string
		wantTokens []token
		wantErr    bool
	}{
		{
			name:  "punctuation and strings",
			input: ` { "key" : ["value", {}] , "other": "text" } `,
			wantTokens: []token{
				{kind: tokenLeftBrace, lexeme: "{", offset: 1},
				{kind: tokenString, lexeme: `"key"`, offset: 3},
				{kind: tokenColon, lexeme: ":", offset: 9},
				{kind: tokenLeftBracket, lexeme: "[", offset: 11},
				{kind: tokenString, lexeme: `"value"`, offset: 12},
				{kind: tokenComma, lexeme: ",", offset: 19},
				{kind: tokenLeftBrace, lexeme: "{", offset: 21},
				{kind: tokenRightBrace, lexeme: "}", offset: 22},
				{kind: tokenRightBracket, lexeme: "]", offset: 23},
				{kind: tokenComma, lexeme: ",", offset: 25},
				{kind: tokenString, lexeme: `"other"`, offset: 27},
				{kind: tokenColon, lexeme: ":", offset: 34},
				{kind: tokenString, lexeme: `"text"`, offset: 36},
				{kind: tokenRightBrace, lexeme: "}", offset: 43},
			},
		},
		{
			name:       "only permitted whitespace",
			input:      " \t\n\r",
			wantTokens: []token{},
		},
		{
			name:  "escaped string",
			input: `"quote: \" slash: \\ unicode: \u0041"`,
			wantTokens: []token{
				{kind: tokenString, lexeme: `"quote: \" slash: \\ unicode: \u0041"`, offset: 0},
			},
		},
		{
			name:  "punctuation inside string",
			input: `"a,b:c"`,
			wantTokens: []token{
				{kind: tokenString, lexeme: `"a,b:c"`, offset: 0},
			},
		},
		{
			name:  "escaped quotes inside string",
			input: `"escaped quote: \"ok\""`,
			wantTokens: []token{
				{kind: tokenString, lexeme: `"escaped quote: \"ok\""`, offset: 0},
			},
		},
		{
			name:  "number and literals",
			input: `36 -1.25e+2 true false null`,
			wantTokens: []token{
				{kind: tokenNumber, lexeme: "36", offset: 0},
				{kind: tokenNumber, lexeme: "-1.25e+2", offset: 3},
				{kind: tokenTrue, lexeme: "true", offset: 12},
				{kind: tokenFalse, lexeme: "false", offset: 17},
				{kind: tokenNull, lexeme: "null", offset: 23},
			},
		},
		{name: "unquoted text", input: `{ key: "value" }`, wantErr: true},
		{name: "unterminated string", input: `"value`, wantErr: true},
		{name: "trailing backslash", input: `"trailing backslash\`, wantErr: true},
		{name: "invalid escape", input: `"value\q"`, wantErr: true},
		{name: "incomplete Unicode escape", input: `"\u12"`, wantErr: true},
		{name: "invalid Unicode hex digit", input: `"\u12G4"`, wantErr: true},
		{name: "unescaped control character", input: "\"line\nbreak\"", wantErr: true},
		{name: "non-breaking space outside string", input: "{\u00a0}", wantErr: true},
		{name: "invalid UTF-8 in string", input: string([]byte{'"', 0xff, '"'}), wantErr: true},
		{name: "leading zero in number", input: "01", wantErr: true},
		{name: "incomplete fraction", input: "1.", wantErr: true},
		{name: "incomplete exponent", input: "1e+", wantErr: true},
		{name: "invalid true literal", input: "True", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tokenize(tt.input)
			if (err != nil) != tt.wantErr {
				t.Fatalf("tokenize(%q) error = %v, wantErr %t", tt.input, err, tt.wantErr)
			}
			if err != nil {
				return
			}
			if len(got) != len(tt.wantTokens) {
				t.Fatalf("tokenize(%q) returned %d tokens, want %d", tt.input, len(got), len(tt.wantTokens))
			}
			for i := range got {
				if got[i] != tt.wantTokens[i] {
					t.Errorf("token %d = %#v, want %#v", i, got[i], tt.wantTokens[i])
				}
			}
		})
	}
}

func TestPrimitiveValueDelimiters(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		parse    func(string) (int, token, error)
		wantEnd  int
		wantKind tokenType
		wantErr  bool
	}{
		{
			name:     "literal followed by closing brace",
			input:    "true}",
			wantEnd:  4,
			wantKind: tokenTrue,
			parse: func(input string) (int, token, error) {
				return tokenizeLiteral(input, 0, "true", tokenTrue)
			},
		},
		{
			name:     "literal followed by comma",
			input:    "true,",
			wantEnd:  4,
			wantKind: tokenTrue,
			parse: func(input string) (int, token, error) {
				return tokenizeLiteral(input, 0, "true", tokenTrue)
			},
		},
		{
			name:     "literal followed by whitespace",
			input:    "true ",
			wantEnd:  4,
			wantKind: tokenTrue,
			parse: func(input string) (int, token, error) {
				return tokenizeLiteral(input, 0, "true", tokenTrue)
			},
		},
		{
			name:     "literal at end of input",
			input:    "true",
			wantEnd:  4,
			wantKind: tokenTrue,
			parse: func(input string) (int, token, error) {
				return tokenizeLiteral(input, 0, "true", tokenTrue)
			},
		},
		{
			name:    "literal followed by text",
			input:   "trueValue",
			parse:   func(input string) (int, token, error) { return tokenizeLiteral(input, 0, "true", tokenTrue) },
			wantErr: true,
		},
		{
			name:     "number followed by closing bracket",
			input:    "12]",
			wantEnd:  2,
			wantKind: tokenNumber,
			parse: func(input string) (int, token, error) {
				return tokenizeNumber(input, 0)
			},
		},
		{
			name:     "number at end of input",
			input:    "12",
			wantEnd:  2,
			wantKind: tokenNumber,
			parse:    func(input string) (int, token, error) { return tokenizeNumber(input, 0) },
		},
		{
			name:    "number followed by text",
			input:   "12abc",
			parse:   func(input string) (int, token, error) { return tokenizeNumber(input, 0) },
			wantErr: true,
		},
		{
			name:    "number followed by second decimal point",
			input:   "12.3.4",
			parse:   func(input string) (int, token, error) { return tokenizeNumber(input, 0) },
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			end, got, err := tt.parse(tt.input)
			if (err != nil) != tt.wantErr {
				t.Fatalf("parse(%q) error = %v, wantErr %t", tt.input, err, tt.wantErr)
			}
			if err == nil && (end != tt.wantEnd || got.kind != tt.wantKind || got.lexeme != tt.input[:end]) {
				t.Errorf("parse(%q) = (%d, %#v), want end %d and token kind %d", tt.input, end, got, tt.wantEnd, tt.wantKind)
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
		{name: "invalid JSON", args: []string{invalidObjectPath}, wantStderrPart: "Invalid JSON:", wantExit: 1},
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
