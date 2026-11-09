# JSON Parser in Go

A JSON parser implemented from scratch in Go, exploring tokenization, recursive-descent parsing, grammar validation, and error handling without relying on an existing JSON parsing library.

The project is written in Go and implements a command-line validator for JSON documents whose root value is an object. It was built to understand how JSON parsing works internally, rather than to provide a replacement for Go's standard JSON library or to claim complete JSON specification compliance.

Original challenge: [Coding Challenge #2 – Build a JSON Parser](https://codingchallenges.substack.com/p/coding-challenge-2?utm_source=chatgpt.com)

## Motivation and objectives

- Understand JSON parsing by implementing lexical and syntactic analysis separately.
- Explore recursive-descent parsing and JSON grammar validation.
- Build a small command-line utility with automated tests.

## Features and supported JSON constructs

The parser accepts a JSON object as the document's root value. Within objects and arrays, it supports:

| Feature | Support |
| --- | --- |
| Empty JSON objects | Yes |
| Objects with string keys | Yes |
| String, number, boolean, and `null` values | Yes |
| Empty arrays and arrays of primitive values | Yes |
| Nested arrays and nested objects | Yes |
| Arrays containing objects and objects containing arrays | Yes |
| JSON whitespace (`space`, tab, line feed, carriage return) | Yes |
| String escape and unescaped control-character validation | Yes |
| Number format validation | Yes |
| Trailing comma rejection | Yes |
| Multiple root values and trailing tokens rejection | Yes |
| Root-level arrays, strings, numbers, booleans, or `null` | No |

This is a learning project and has not been verified against a comprehensive JSON conformance suite. Its supported root type is intentionally limited to objects.

## Architecture and design

The validation flow is:

```text
Raw JSON document
        |
        v
Tokenizer -- recognizes tokens, checks lexical rules, records byte offsets
        |
        v
Parser -- validates object/array grammar using recursive descent
        |
        v
Validation result -- valid JSON or an error diagnostic
```

### Tokenizer

- Recognizes braces, brackets, colons, commas, strings, numbers, and the `true`, `false`, and `null` literals.
- Checks lexical rules, including JSON number formats, string escapes, control characters, and UTF-8 in strings.
- Tracks byte offsets for tokens and lexical errors.
- Skips only whitespace permitted by JSON.

### Parser

- Maintains a cursor through the token sequence.
- Parses objects and arrays, recursively handling nested values.
- Requires string object keys and rejects invalid token sequences.
- Ensures the complete token sequence is consumed, so a valid prefix is not accepted as a complete document.

### CLI

- Reads one input file and invokes validation.
- Prints `Valid JSON` for valid input or reports an invalid document on standard error.
- Uses exit codes to distinguish valid input, invalid JSON, and usage or file-reading errors.

## Project structure

```text
json-parser/
├── main.go
├── tokenize.go
├── parser.go
├── validate.go
├── main_test.go
├── go.mod
├── tests/
│   └── step1/
│       ├── invalid.json
│       └── valid.json
└── README.md
```

| File | Responsibility |
| --- | --- |
| `main.go` | CLI entry point, file handling, output, and exit codes |
| `tokenize.go` | Lexical analysis and tokenization |
| `parser.go` | Recursive-descent grammar validation |
| `validate.go` | Validation orchestration and JSON whitespace helper |
| `main_test.go` | Validator, tokenizer, token-boundary, and CLI tests |
| `go.mod` | Go module name and minimum Go version |
| `tests/step1/` | Sample valid and invalid JSON inputs |

## Installation and execution

### Prerequisites

- [Go 1.19 or later](https://go.dev/dl/) (the version declared in `go.mod`)
- Git

### Clone

```sh
git clone https://github.com/kushal-kamra/json-parser.git
cd json-parser
```

### Run

```sh
go run . path/to/input.json
```

### Build and execute

```sh
go build -o json-parser .
./json-parser path/to/input.json
```

The CLI accepts exactly one file path.

## Example usage

Save this as `input.json`:

```json
{
  "name": "Kushal",
  "age": 36,
  "active": true,
  "metadata": null,
  "skills": ["Go", "JavaScript"]
}
```

Run `./json-parser input.json`:

```text
Valid JSON
```

This input has a trailing comma and is invalid:

```json
{
  "name": "Kushal",
  "age": 36,
}
```

Run `./json-parser input.json` with that content:

```text
Invalid JSON: unexpected token "}" at byte 35; expected string object key
```

The diagnostic above is the CLI's actual output for the displayed input, including its line breaks.

## Exit codes

| Exit code | Meaning |
| --- | --- |
| `0` | Valid JSON |
| `1` | Invalid JSON |
| `2` | CLI usage or file-reading error |

## Testing and static analysis

Run the automated tests:

```sh
go test ./...
```

Run Go's static analysis checks:

```sh
go vet ./...
```

The tests cover valid and invalid JSON structures, strings and escapes, number formats, literal recognition, token boundaries, nested objects and arrays, missing values and separators, trailing commas, unexpected trailing tokens, and empty or malformed input. They also exercise CLI output and exit codes.

## Engineering decisions

| Decision | Rationale |
| --- | --- |
| Separate tokenizer and parser | Isolates lexical validation from grammar validation |
| Distinct primitive token types | Makes parser value handling explicit and extensible |
| Recursive-descent parsing | Maps nested JSON structures naturally to parsing functions |
| Token cursor | Provides predictable sequential parsing |
| Byte offsets | Makes diagnostics point to locations in the original input |
| Grammar-driven validation | Avoids relying on unrelated special-case checks |
| Full input consumption | Prevents accepting a valid prefix followed by extra tokens |

For example, `[1,]` is rejected because after the comma the array grammar requires another value. The parser does not need a separate trailing-comma rule.

## Limitations and future improvements

- Support non-object root values: arrays, strings, numbers, booleans, and `null`.
- Improve diagnostics with line and column numbers in addition to byte offsets.
- Add conformance tests against a comprehensive JSON test suite.
- Expand malformed UTF-8 and Unicode escape test coverage.
- Add fuzz testing for tokenizer and parser edge cases.

## Acknowledgements and references

- [John Crickett – Coding Challenge #2: Build a JSON Parser](https://codingchallenges.substack.com/p/coding-challenge-2?utm_source=chatgpt.com)
- [The Go Programming Language](https://go.dev/)
