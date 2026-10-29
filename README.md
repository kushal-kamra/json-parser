# json-parser
Coding Challenge #2 - Build a JSON Parser

## Step 1 validator

`ValidateJSON(input string) (isValid bool, errorMessage string)` validates JSON
objects whose keys and values are strings, including empty objects. It allows
JSON whitespace around and within the object. Successful validation returns an
empty error message.
The `json-parser` CLI takes exactly one file path. It prints `Valid JSON` and
exits with status 0 for valid input. Invalid input is reported on stderr as
`Invalid JSON: <reason>` with status 1. Incorrect argument counts and file-read
errors are also reported on stderr, with status 2. Tests cover empty and
string-valued objects, malformed object syntax, tokenization errors, and
trailing content.

## Step 2 tokenizer and parser

The tokenizer recognizes object braces, colons, commas, and quoted strings,
skipping only JSON-permitted whitespace. String escapes and control characters
are checked lexically. The parser consumes the complete token sequence and
validates objects containing string keys and string values.
