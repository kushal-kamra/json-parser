package main

func ValidateJSON(input string) (isValid bool, errorMessage string) {
	tokens, err := tokenize(input)
	if err != nil {
		return false, err.Error()
	}

	if err := parseJSON(tokens); err != nil {
		return false, err.Error()
	}
	return true, ""
}

func isJSONWhitespace(char byte) bool {
	return char == ' ' || char == '\t' || char == '\n' || char == '\r'
}
