package main

import (
	"strconv"
	"strings"
)

func lexLine(tokens *TokenQueue, isFirstLine bool) {
	prompt := ">>> "
	if isFirstLine {
		defer func() {
			if err := recover(); err != nil {
				lineCount++
				switch e := err.(type) {
				case error:
					printError(e)
				default:
					panic(err)
				}
				for tokens.size() > 0 {
					tokens.next()
				}
				tokens.pushBackEOF()
			}
		}()
	} else {
		prompt = "... "
	}
	line := readLine(prompt)
	if line == "\n" {
		lineCount++
		tokens.pushBackEOF()
		return
	}
	if line == "exit\n" || line == "quit\n" {
		ioExit()
		return
	} else if line == "help\n" || strings.HasPrefix(line, "help ") {
		showHelp(line)
		lineCount++
		tokens.pushBackEOF()
		return
	} else if line == "example\n" || strings.HasPrefix(line, "example ") {
		showExample(line)
		lineCount++
		tokens.pushBackEOF()
		return
	}
	switch line[len(line)-2] {
	case '\\':
		lex(line[:len(line)-2], tokens)
		lexLine(tokens, false)
	case '{':
		lex(line, tokens)
		for {
			line := readLine("... ")
			lex(line, tokens)
			if line == "}\n" {
				break
			}
		}
	default:
		lex(line, tokens)
	}
	if isFirstLine {
		tokens.pushBackEOF()
	}
}

func lexProgram(str string, tokens *TokenQueue) {
	lineCount = 1
	lex(str, tokens)
	tokens.pushBack(Token{TT_EOF, "", Position{lineCount, 0, 0}})
}

var lineCount = 1

func lex(str string, tokens *TokenQueue) {
	runes := []rune(str)
	pos := 0
	idx := 0
	for idx < len(runes) {
		outputToken := true
		tok := Token{CT_ILLEGAL, string(runes[idx]), Position{lineCount, pos, pos + 1}}
		curr := runes[idx]
		idx++
		pos++
		switch runeType(curr) {
		case CT_WHITESPACE:
			tok.ty = TT_WHITESPACE
			if tokens.peekBack().ty == TT_UNKNOWN || tokens.peekBack().ty == TT_TERMINATOR || tokens.peekBack().ty == TT_WHITESPACE {
				outputToken = false
			}
		case CT_NEWLINE:
			lineCount++
			pos = 0
			tok.ty = TT_TERMINATOR
			tok.data = "|"
			lastTy := tokens.peekBack().ty
			if lastTy == TT_WHITESPACE {
				tokens.popBack()
				lastTy = tokens.peekBack().ty
			}
			switch lastTy {
			case TT_LITERAL, TT_IDENTIFIER, TT_PARENTHESIS_CLOSE, TT_CURLY_BRACES_CLOSE, TT_SQUARE_BRACKETS_CLOSE:
				tok.ty = TT_TERMINATOR
			default:
				tok.ty = TT_WHITESPACE
			}
			if lastTy == TT_UNKNOWN || lastTy == TT_TERMINATOR || (tok.ty == TT_WHITESPACE && lastTy == TT_WHITESPACE) {
				outputToken = false
			}
		case CT_DIGIT:
			tok.ty = TT_LITERAL
			if idx < len(runes) && curr == '0' && (runes[idx] == 'x' || runes[idx] == 'X') {
				idx++
				tok.data = "0x"
				for idx < len(runes) && (runeType(runes[idx]) == CT_DIGIT || runeType(runes[idx]) == CT_LETTER) {
					tok.data += string(runes[idx])
					idx++
					pos++
				}
				tok.data = strconv.Itoa(atoi(tok.data, tok.pos))
			} else {
				for idx < len(runes) && runeType(runes[idx]) == CT_DIGIT {
					tok.data += string(runes[idx])
					idx++
					pos++
				}
			}
			tok.pos.end = pos
		case CT_LETTER:
			for idx < len(runes) && (runeType(runes[idx]) == CT_DIGIT || runeType(runes[idx]) == CT_LETTER) {
				tok.data += string(runes[idx])
				idx++
				pos++
			}
			if tok.data == "false" || tok.data == "true" {
				tok.ty = TT_LITERAL
				if tok.data == "false" {
					tok.data = ""
				} else {
					tok.data = "1"
				}
			} else if isOperator(tok.data) {
				tok.ty = opType(tok.data)
				if tok.ty == TT_IN && tokens.peekBack().ty == TT_WHITESPACE && tokens.peekBeforeBack().ty == TT_NOT {
					tokens.popBack()
					tok = tokens.popBack()
					tok.ty = TT_NOT_IN
					tok.data = "not in"
				}
			} else {
				tok.ty = TT_IDENTIFIER
			}
			tok.pos.end = pos
		case CT_OPERATOR:
			for idx < len(runes) && runeType(runes[idx]) == CT_OPERATOR {
				newOp := tok.data + string(runes[idx])
				if isOperator(newOp) || !isOperator(tok.data) {
					tok.data = newOp
					idx++
					pos++
				} else {
					break
				}
			}
			tok.ty = opType(tok.data)
			switch opType(tok.data) {
			case TT_UNKNOWN:
				tok.pos.end = pos
				panic(myErr{"Operator \"" + tok.data + "\" does not exist.", tok.pos, ERR_LEXER})
			case TT_SINGLE_QUOTE, TT_DOUBLE_QUOTE, TT_TICK_QUOTE:
				tok.ty = TT_LITERAL
				close := rune(tok.data[0])
				tok.data = ""

				for idx < len(runes) && runes[idx] != close {
					if runes[idx] == '\n' {
						lineCount++
						pos = 0
					}
					tok.data += string(runes[idx])
					idx++
					pos++
				}
				idx++
				pos++
			case TT_SINGLE_LINE_COMMENT:
				for idx < len(runes) && runes[idx] != '\n' {
					idx++
					pos++
				}
				outputToken = false
			case TT_MULTI_LINE_COMMENT_OPEN:
				for idx+1 < len(runes) && !(runes[idx] == '*' && runes[idx+1] == '/') {
					if runes[idx] == '\n' {
						lineCount++
						pos = 0
					} else {
						pos++
					}
					idx++
				}
				idx += 2
				pos += 2
				outputToken = false
			}
		case CT_ESCAPE:
			idx--
			pos--
			tok.data = ""
			tok.ty = TT_LITERAL
			for idx < len(runes) && runes[idx] == '\\' {
				startPos := pos
				idx++
				pos++
				c := runes[idx]
				idx++
				pos++
				switch c {
				case 'n':
					tok.data += "\n"
				case 't':
					tok.data += "\t"
				case 'x':
					str = "0x"
					for runeType(runes[idx]) == CT_DIGIT || runeType(runes[idx]) == CT_LETTER {
						str += string(runes[idx])
						pos++
						idx++
					}
					tok.pos.end = pos
					tok.data += string(atoi(str, Position{lineCount, startPos, pos}))
					break
				default:
					if '0' <= c && c <= '9' {
						str = string(c)
						for runeType(runes[idx]) == CT_DIGIT || runeType(runes[idx]) == CT_LETTER {
							str += string(runes[idx])
							pos++
							idx++
						}
						tok.pos.end = pos
						tok.data += string(atoi(str, Position{lineCount, startPos, pos}))
					} else {
						tok.pos.end = pos
						panic(myErr{"Invalid escape sequence.", Position{lineCount, startPos, pos}, ERR_LEXER})
					}
				}
			}

			break
		case CT_ILLEGAL:
			panic(myErr{`unknown character ` + strconv.QuoteRune(curr), tok.pos, ERR_LEXER})
		}
		if pos != 0 {
			tok.pos.end = pos
		}
		if outputToken {
			tokens.pushBack(tok)
		}
	}
}
