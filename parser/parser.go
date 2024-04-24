package parser

import (
	"fmt"
	"heisenberg/lang"
)

type parser struct {
	tokens []token
	pos    int
}

func Parse(input string) ([]lang.Ast, error) {
	tokens, err := tokenise(input)
	if err != nil {
		return nil, err
	}

	var result []lang.Ast
	p := newParser(tokens)

	for p.current().kind != tokenEOF {
		clause, err := p.parseClause()
		if err != nil {
			return nil, err
		}
		result = append(result, clause)
		p.skipNewline()
	}

	return result, nil
}

func newParser(tokens []token) parser {
	return parser{tokens: tokens}
}

func (p *parser) parseClause() (lang.Ast, error) {
	kw, err := p.parseKeyword()
	if err != nil {
		return nil, err
	}

	switch lang.Keyword(kw) {
	case lang.KeywordFrom:
		return p.parseFromClause()
	case lang.KeywordCreate:
		return p.parseCreateClause()
	case lang.KeywordSelect:
		return p.parseSelectClause()
	default:
		return nil, fmt.Errorf("unbound symbol")
	}
}

func (p *parser) parseFromClause() (lang.Ast, error) {
	table, err := p.parseIdentifier()
	if err != nil {
		return nil, err
	}
	return lang.From{Table: table}, nil
}

func (p *parser) parseCreateClause() (lang.Ast, error) {
	kw, err := p.parseKeyword()
	if err != nil {
		return nil, err
	}

	switch lang.Keyword(kw) {
	case lang.KeywordTable:
		return p.parseCreateTableClause()
	default:
		return nil, fmt.Errorf("cannot create %s", kw)
	}
}

func (p *parser) parseSelectClause() (lang.Ast, error) {
	columns, err := p.parseCommaSeparatedList()
	if err != nil {
		return nil, err
	}
	return lang.Select{Columns: columns}, nil
}

func (p *parser) parseCreateTableClause() (lang.Ast, error) {
	columns := []lang.ColumnDefinition{}
	toks, err := p.parseManyTokens(tokenIdent, tokenLBrace)
	if err != nil {
		return nil, err
	}
	name := toks[0]

	for p.current().kind != tokenRBrace {
		column, err := p.parseColumnDefinition()
		if err != nil {
			return nil, err
		}
		columns = append(columns, column)
		if _, err := p.parseAnyToken(tokenComma); err != nil {
			if p.current().kind != tokenRBrace {
				return nil, err
			}
		}
	}

	if _, err := p.parseAnyToken(tokenRBrace); err != nil {
		return nil, err
	}

	return lang.CreateTable{Name: name, Columns: columns}, nil
}

func (p *parser) parseColumnDefinition() (lang.ColumnDefinition, error) {
	var def lang.ColumnDefinition
	for p.next().kind != tokenColon {
		// Parse constraints here
	}
	toks, err := p.parseManySkippingNewLine(tokenIdent, tokenColon, tokenType)
	if err != nil {
		return lang.ColumnDefinition{}, err
	}
	def.Name = toks[0]
	def.Type = lang.TypeKeywords[lang.Keyword(toks[2])]
	return def, nil
}

func (p *parser) parseAnyToken(kind tokenKind) (string, error) {
	if err := p.expect(kind); err != nil {
		return "", err
	}
	value := p.current().value
	p.pos++
	return value, nil
}

func (p *parser) parseManyTokens(kinds ...tokenKind) ([]string, error) {
	var result []string
	for _, kind := range kinds {
		value, err := p.parseAnyToken(kind)
		if err != nil {
			return nil, err
		}
		result = append(result, value)
	}
	return result, nil
}

func (p *parser) parseManySkippingNewLine(kinds ...tokenKind) ([]string, error) {
	var result []string
	for _, kind := range kinds {
		p.skipNewline()
		value, err := p.parseAnyToken(kind)
		if err != nil {
			return nil, err
		}
		result = append(result, value)
		p.skipNewline()
	}
	return result, nil
}

func (p *parser) parseIdentifier() (string, error) {
	return p.parseAnyToken(tokenIdent)
}

func (p *parser) parseKeyword() (string, error) {
	return p.parseAnyToken(tokenKeyword)
}

func (p *parser) parseType() (string, error) {
	return p.parseAnyToken(tokenType)
}

func (p *parser) parseCommaSeparatedList() ([]string, error) {
	var list []string
	for p.pos < len(p.tokens) {
		ident, err := p.parseIdentifier()
		if err != nil {
			return list, nil
		}
		list = append(list, ident)
		if p.current().kind != tokenComma {
			break
		}
		p.pos++
	}
	return list, nil
}

func (p *parser) skipNewline() {
	for p.current().kind == tokenNewLine {
		p.pos++
	}
}

func (p *parser) current() token {
	return p.tokens[p.pos]
}

func (p *parser) next() token {
	return p.tokens[p.pos+1]
}

func (p *parser) expect(kind tokenKind) error {
	if p.pos >= len(p.tokens) || p.tokens[p.pos].kind != kind {
		return fmt.Errorf("expected: %s, got: %s, at position: %d", tokenNames[kind], p.tokens[p.pos].value, p.pos)
	}
	return nil
}
