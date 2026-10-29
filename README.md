# json-parser
Coding Challenge #2 - Build a JSON Parser

## Step 1 validator

`ValidateJSON(input string) (isValid bool, errorMessage string)` currently
validates only an empty JSON object, allowing JSON whitespace around the object
and between its braces. Successful validation returns an empty error message.
`go test ./...` covers valid `{}` input and the following invalid inputs:

empty input, arrays, non-whitespace object contents, multiple objects, trailing
content, and non-breaking space.
