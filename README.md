# json-parser
Coding Challenge #2 - Build a JSON Parser

## Step 1 validator

`ValidateJSON(input string) (isValid bool, errorMessage string)` currently
validates only an empty JSON object, allowing JSON whitespace around the object
and between its braces. Successful validation returns an empty error message.
The `json-parser` CLI takes exactly one file path. It prints `Valid JSON` and
exits with status 0 for valid input. Invalid input is reported on stderr as
`Invalid JSON: <reason>` with status 1. Incorrect argument counts and file-read
errors are also reported on stderr, with status 2. `go test ./...` covers valid
`{}` input and the following invalid inputs:

empty input, arrays, non-whitespace object contents, multiple objects, trailing
content, and non-breaking space.

## Step 2 tokenizer

The tokenizer recognizes object braces, colons, commas, and quoted strings,
skipping only JSON-permitted whitespace. String escapes and control characters
are checked lexically. The parser and Step 2 object validation are not
implemented yet.
