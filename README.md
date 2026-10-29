# json-parser
Coding Challenge #2 - Build a JSON Parser

## Step 1 validator

`ValidateJSON(input string) (isValid bool, errorMessage string)` currently
validates only an empty JSON object, allowing JSON whitespace around the object
and between its braces. Successful validation returns an empty error message.
The `json-parser` CLI takes exactly one file path. It prints `Valid JSON` and
exits with status 0 for valid input, or prints `Invalid JSON` and exits with
status 1 for invalid input. Incorrect argument counts and file-read errors are
reported on stderr with status 2. `go test ./...` covers valid `{}` input and
the following invalid inputs:

empty input, arrays, non-whitespace object contents, multiple objects, trailing
content, and non-breaking space.
