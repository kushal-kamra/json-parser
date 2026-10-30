package main

import (
	"fmt"
	"unicode/utf8"
)

type tokenType uint8

const (
	tokenLeftBrace tokenType = iota
	tokenRightBrace
	tokenColon
	tokenComma
	tokenString
	tokenNumber
	tokenTrue
	tokenFalse
	tokenNull
	tokenLeftBracket
	tokenRightBracket
)

type token struct {
	kind   tokenType
	lexeme string
	offset int
}

func tokenize(input string) ([]token, error) {
	tokens := make([]token, 0)
	for i := 0; i < len(input); {
		if isJSONWhitespace(input[i]) {
			i++
			continue
		}

		switch input[i] {
		case '{':
			tokens = append(tokens, token{kind: tokenLeftBrace, lexeme: "{", offset: i})
			i++
		case '}':
			tokens = append(tokens, token{kind: tokenRightBrace, lexeme: "}", offset: i})
			i++
		case '[':
			tokens = append(tokens, token{kind: tokenLeftBracket, lexeme: "[", offset: i})
			i++
		case ']':
			tokens = append(tokens, token{kind: tokenRightBracket, lexeme: "]", offset: i})
			i++
		case ':':
			tokens = append(tokens, token{kind: tokenColon, lexeme: ":", offset: i})
			i++
		case ',':
			tokens = append(tokens, token{kind: tokenComma, lexeme: ",", offset: i})
			i++
		case '"':
			next, stringToken, err := tokenizeString(input, i)
			if err != nil {
				return nil, err
			}
			tokens = append(tokens, stringToken)
			i = next
		case 't':
			next, literalToken, err := tokenizeLiteral(input, i, "true", tokenTrue)
			if err != nil {
				return nil, err
			}
			tokens = append(tokens, literalToken)
			i = next
		case 'f':
			next, literalToken, err := tokenizeLiteral(input, i, "false", tokenFalse)
			if err != nil {
				return nil, err
			}
			tokens = append(tokens, literalToken)
			i = next
		case 'n':
			next, literalToken, err := tokenizeLiteral(input, i, "null", tokenNull)
			if err != nil {
				return nil, err
			}
			tokens = append(tokens, literalToken)
			i = next
		case '-', '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
			next, numberToken, err := tokenizeNumber(input, i)
			if err != nil {
				return nil, err
			}
			tokens = append(tokens, numberToken)
			i = next
		default:
			return nil, fmt.Errorf("unexpected character %q at byte %d", input[i], i)
		}
	}

	return tokens, nil
}

func tokenizeLiteral(input string, start int, literal string, kind tokenType) (int, token, error) {
	if len(input)-start < len(literal) || input[start:start+len(literal)] != literal {
		return 0, token{}, fmt.Errorf("invalid literal at byte %d", start)
	}

	end := start + len(literal)
	if end < len(input) && !isValueDelimiter(input[end]) {
		return 0, token{}, fmt.Errorf("invalid literal at byte %d", start)
	}
	return end, token{kind: kind, lexeme: literal, offset: start}, nil
}

func tokenizeNumber(input string, start int) (int, token, error) {
	i := start
	if input[i] == '-' {
		i++
		if i == len(input) {
			return 0, token{}, fmt.Errorf("incomplete number at byte %d", start)
		}
	}

	if input[i] == '0' {
		i++
		if i < len(input) && isDigit(input[i]) {
			return 0, token{}, fmt.Errorf("leading zero in number at byte %d", start)
		}
	} else if input[i] >= '1' && input[i] <= '9' {
		for i < len(input) && isDigit(input[i]) {
			i++
		}
	} else {
		return 0, token{}, fmt.Errorf("invalid number at byte %d", start)
	}

	if i < len(input) && input[i] == '.' {
		i++
		fractionStart := i
		for i < len(input) && isDigit(input[i]) {
			i++
		}
		if fractionStart == i {
			return 0, token{}, fmt.Errorf("incomplete number fraction at byte %d", start)
		}
	}

	if i < len(input) && (input[i] == 'e' || input[i] == 'E') {
		i++
		if i < len(input) && (input[i] == '+' || input[i] == '-') {
			i++
		}
		exponentStart := i
		for i < len(input) && isDigit(input[i]) {
			i++
		}
		if exponentStart == i {
			return 0, token{}, fmt.Errorf("incomplete number exponent at byte %d", start)
		}
	}

	if i < len(input) && !isValueDelimiter(input[i]) {
		return 0, token{}, fmt.Errorf("invalid number at byte %d", start)
	}

	return i, token{kind: tokenNumber, lexeme: input[start:i], offset: start}, nil
}

func isValueDelimiter(char byte) bool {
	return isJSONWhitespace(char) || char == ',' || char == '}' || char == ']'
}

func isDigit(char byte) bool {
	return char >= '0' && char <= '9'
}

func tokenizeString(input string, start int) (int, token, error) {
	for i := start + 1; i < len(input); {
		switch char := input[i]; {
		case char == '"':
			return i + 1, token{
				kind:   tokenString,
				lexeme: input[start : i+1],
				offset: start,
			}, nil
		case char == '\\':
			i++
			if i == len(input) {
				return 0, token{}, fmt.Errorf("unterminated escape sequence at byte %d", i-1)
			}

			switch input[i] {
			case '"', '\\', '/', 'b', 'f', 'n', 'r', 't':
				i++
			case 'u':
				if len(input)-i < 5 {
					return 0, token{}, fmt.Errorf("incomplete Unicode escape at byte %d", i-1)
				}
				for j := i + 1; j <= i+4; j++ {
					if !isHexDigit(input[j]) {
						return 0, token{}, fmt.Errorf("invalid Unicode escape at byte %d", i-1)
					}
				}
				i += 5
			default:
				return 0, token{}, fmt.Errorf("invalid escape sequence at byte %d", i-1)
			}
		case char < 0x20:
			return 0, token{}, fmt.Errorf("unescaped control character in string at byte %d", i)
		case char >= utf8.RuneSelf:
			_, size := utf8.DecodeRuneInString(input[i:])
			if size == 1 {
				return 0, token{}, fmt.Errorf("invalid UTF-8 in string at byte %d", i)
			}
			i += size
		default:
			i++
		}
	}

	return 0, token{}, fmt.Errorf("unterminated string at byte %d", start)
}

func isHexDigit(char byte) bool {
	return char >= '0' && char <= '9' ||
		char >= 'a' && char <= 'f' ||
		char >= 'A' && char <= 'F'
}
