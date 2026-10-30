package main

import "fmt"

type parser struct {
	tokens []token
	cursor int
}

func parseJSON(tokens []token) error {
	p := parser{tokens: tokens}
	if len(tokens) == 0 {
		return fmt.Errorf("input is empty")
	}

	if err := p.parseObject(); err != nil {
		return err
	}
	if p.cursor != len(p.tokens) {
		return p.unexpectedToken("end of input")
	}

	return nil
}

func (p *parser) parseObject() error {
	if err := p.expect(tokenLeftBrace, "{"); err != nil {
		return err
	}
	if p.consume(tokenRightBrace) {
		return nil
	}

	for {
		if err := p.expect(tokenString, "string object key"); err != nil {
			return err
		}
		if err := p.expect(tokenColon, ":"); err != nil {
			return err
		}
		if err := p.parseValue(); err != nil {
			return err
		}

		if p.consume(tokenRightBrace) {
			return nil
		}
		if err := p.expect(tokenComma, ", or }"); err != nil {
			return err
		}
	}
}

func (p *parser) parseValue() error {
	if p.cursor == len(p.tokens) {
		return fmt.Errorf("unexpected end of input; expected object value")
	}

	switch p.tokens[p.cursor].kind {
	case tokenString, tokenNumber, tokenTrue, tokenFalse, tokenNull:
		p.cursor++
		return nil
	default:
		return p.unexpectedToken("object value")
	}
}

func (p *parser) expect(kind tokenType, expected string) error {
	if p.cursor == len(p.tokens) {
		return fmt.Errorf("unexpected end of input; expected %s", expected)
	}
	if p.tokens[p.cursor].kind != kind {
		return p.unexpectedToken(expected)
	}

	p.cursor++
	return nil
}

func (p *parser) consume(kind tokenType) bool {
	if p.cursor < len(p.tokens) && p.tokens[p.cursor].kind == kind {
		p.cursor++
		return true
	}
	return false
}

func (p *parser) unexpectedToken(expected string) error {
	current := p.tokens[p.cursor]
	return fmt.Errorf("unexpected token %q at byte %d; expected %s", current.lexeme, current.offset, expected)
}
