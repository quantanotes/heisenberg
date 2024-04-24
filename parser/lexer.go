package parser

import (
	"fmt"
	"heisenberg/lang"
	"strings"
	"unicode"
	"unicode/utf8"
)

type lexer struct {
	input string
	pos   int
	width int
}

func tokenise(input string) ([]token, error) {
	lexer := newLexer(input)
	var tokens []token
	for {
		tok, err := lexer.Next()
		if err != nil {
			return nil, err
		}
		tokens = append(tokens, tok)
		if tok.kind == tokenEOF {
			break
		}
	}
	return tokens, nil
}

func newLexer(input string) lexer {
	return lexer{input: input}
}

func (l *lexer) Next() (token, error) {
	if l.pos >= len(l.input) {
		return token{kind: tokenEOF}, nil
	}

	r, _ := utf8.DecodeRuneInString(l.input[l.pos:])
	switch {
	case isSpace(r):
		l.consumeWhitespace()
		return l.Next()
	case unicode.IsLetter(r):
		return l.scanIdentifier()
	case unicode.IsDigit(r):
		return l.scanNumber()
	case r == '"':
		return l.scanString()
	default:
		return l.scanPunctuation()
	}
}

func (l *lexer) scanIdentifier() (token, error) {
	start := l.pos
	for l.pos < len(l.input) {
		r, _ := utf8.DecodeRuneInString(l.input[l.pos:])
		if !unicode.IsLetter(r) && !unicode.IsDigit(r) && r != '_' {
			break
		}
		l.width = utf8.RuneLen(r)
		l.pos += l.width
	}
	value := l.input[start:l.pos]
	return token{kind: l.identify(value), value: value, pos: start}, nil
}

func (l *lexer) scanNumber() (token, error) {
	start := l.pos
	for l.pos < len(l.input) {
		r, _ := utf8.DecodeRuneInString(l.input[l.pos:])
		if !unicode.IsDigit(r) {
			break
		}
		l.width = utf8.RuneLen(r)
		l.pos += l.width
	}
	value := l.input[start:l.pos]
	return token{kind: tokenNumber, value: value, pos: start}, nil
}

func (l *lexer) scanString() (token, error) {
	start := l.pos
	l.pos++
	for l.pos < len(l.input) {
		r, _ := utf8.DecodeRuneInString(l.input[l.pos:])
		if r == '"' {
			l.pos++
			break
		}
		if r == '\\' {
			l.pos += 2
		} else {
			l.width = utf8.RuneLen(r)
			l.pos += l.width
		}
	}
	value := l.input[start:l.pos]
	return token{kind: tokenString, value: value, pos: start}, nil
}

func (l *lexer) scanPunctuation() (token, error) {
	switch l.input[l.pos] {
	case '=':
		return l.consumeToken(tokenEqual)
	case '{':
		return l.consumeToken(tokenLBrace)
	case '}':
		return l.consumeToken(tokenRBrace)
	case ':':
		return l.consumeToken(tokenColon)
	case ',':
		return l.consumeToken(tokenComma)
	case '\n':
		return l.consumeToken(tokenNewLine)
	default:
		return token{}, fmt.Errorf("unexpected token at position %d: %c", l.pos, l.input[l.pos])
	}
}

func (l *lexer) consumeWhitespace() {
	start := l.pos
	for l.pos < len(l.input) && isSpace(rune(l.input[l.pos])) {
		l.width = 1
		l.pos++
	}
	l.width = l.pos - start
}

func (l *lexer) consumeToken(kind tokenKind) (token, error) {
	value := l.input[l.pos : l.pos+1]
	token := token{kind: kind, value: value, pos: l.pos}
	l.pos++
	return token, nil
}

func (l *lexer) identify(value string) tokenKind {
	if _, ok := lang.TypeKeywords[lang.Keyword(strings.ToLower(value))]; ok {
		return tokenType
	}
	if _, ok := lang.Keywords[lang.Keyword(strings.ToLower(value))]; ok {
		return tokenKeyword
	}
	return tokenIdent
}

func isSpace(r rune) bool {
	return unicode.IsSpace(r) && r != '\n'
}
