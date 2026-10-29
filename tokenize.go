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
		default:
			return nil, fmt.Errorf("unexpected character %q at byte %d", input[i], i)
		}
	}

	return tokens, nil
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
