package main

func ValidateJSON(input string) (isValid bool, errorMessage string) {
	start, end := 0, len(input)
	for start < end && isJSONWhitespace(input[start]) {
		start++
	}
	for start < end && isJSONWhitespace(input[end-1]) {
		end--
	}

	if start == end {
		return false, "input is empty"
	}
	if input[start] != '{' || input[end-1] != '}' {
		return false, "expected an empty JSON object"
	}
	for i := start + 1; i < end-1; i++ {
		if !isJSONWhitespace(input[i]) {
			return false, "expected only whitespace between object braces"
		}
	}

	return true, ""
}

func isJSONWhitespace(char byte) bool {
	return char == ' ' || char == '\t' || char == '\n' || char == '\r'
}
