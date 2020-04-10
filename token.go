package main

import (
	"fmt"
	"strconv"
)

type RuneType int

const (
	CT_ILLEGAL = iota
	CT_WHITESPACE
	CT_NEWLINE
	CT_OPERATOR
	CT_LETTER
	CT_DIGIT
	CT_ESCAPE
)

func runeType(r rune) RuneType {
	if ('a' <= r && r <= 'z') || ('A' <= r && r <= 'Z') || r == '_' {
		return CT_LETTER
	}
	if '0' <= r && r <= '9' {
		return CT_DIGIT
	}
	switch r {
	case '=', '!', '<', '>', '#', '+', '\'', '"', '-', '*', '/', '%', ':', '|', '[', '{', '(', ']', '}', ')', ',', '.':
		return CT_OPERATOR
	case '\\':
		return CT_ESCAPE
	case '\n':
		return CT_NEWLINE
	case ' ', '\t', '\r':
		return CT_WHITESPACE
	default:
		return CT_ILLEGAL
	}
}

type TokenType int

type Token struct {
	ty   TokenType
	data string
	pos  Position
}

type Position struct {
	line  int
	start int
	end   int
}

func showToken(tok Token) {
	fmt.Printf("[ type %d | %d:%d:%d - \t%s]\n", tok.ty, tok.pos.line, tok.pos.start, tok.pos.end, strconv.QuoteToGraphic(tok.data))
}
