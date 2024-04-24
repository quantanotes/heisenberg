package parser

type tokenKind int

const (
	tokenIdent tokenKind = iota
	tokenKeyword
	tokenType
	tokenNumber
	tokenString
	tokenEqual
	tokenLBrace
	tokenRBrace
	tokenComma
	tokenColon
	tokenNewLine
	tokenEOF
	tokenInvalid
)

var tokenNames = [...]string{
	tokenIdent:   "identifier",
	tokenKeyword: "keyword",
	tokenType:    "type",
	tokenNumber:  "number",
	tokenString:  "string",
	tokenEqual:   "equal",
	tokenLBrace:  "left brace",
	tokenRBrace:  "right brace",
	tokenComma:   "comma",
	tokenColon:   "colon",
	tokenNewLine: "new line",
	tokenEOF:     "end of file",
	tokenInvalid: "invalid",
}

type token struct {
	kind  tokenKind
	value string
	pos   int
}
