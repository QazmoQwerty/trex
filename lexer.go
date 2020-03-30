package main

var allUserInput = []string{}
var lastLine = ""

func insertLine(line string) {
	allUserInput = append(allUserInput, line)
	lastLine = line
}

func readLine(prompt string) string {
	line, err := globals.liner.Prompt(prompt)
	if err != nil {
		panic(err)
	}
	insertLine(line)
	return line + "\n"
}

func showHelp() {
	println("Help still needs to be written")
	println("For now see the Language Specification: gitlab.com/QazmoQwerty/trex/-/blob/master/docs/trex-spec.md")
}

func lexLine(tokens chan Token, isFirstLine bool) {
	prompt := ">>> "
	if isFirstLine {

		defer func() {
			if err := recover(); err != nil {
				switch e := err.(type) {
				case error:
					printError(e)
				default:
					panic(err)
				}
				for len(tokens) > 0 {
					<-tokens
				}
				tokens <- Token{TT_EOF, "", Position{0, 0, 0}}
			}
		}()
	} else {
		prompt = "... "
	}
	line := readLine(prompt)
	if line == "exit\n" || line == "quit\n" {
		exitProgram()
		return
	} else if line == "help\n" {
		showHelp()
		tokens <- Token{TT_EOF, "", Position{0, 0, 0}}
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
		tokens <- Token{TT_EOF, "", Position{0, 0, 0}}
	}
}

func lexProgram(str string, tokens chan Token) {
	lex(str, tokens)
	tokens <- Token{TT_EOF, "", Position{0, 0, 0}}
}

var lineCount = 1

func lex(str string, tokens chan Token) {
	runes := []rune(str)
	// line := 1
	pos := 0
	idx := 0
	var beforeLast Token
	beforeLast.ty = TT_UNKNOWN
	var last Token
	last.ty = TT_UNKNOWN
	for idx < len(runes) {
		outputToken := true
		// tok := Token{CT_ILLEGAL, string(runes[idx]), Position{line, pos, pos + 1}}
		tok := Token{CT_ILLEGAL, string(runes[idx]), Position{lineCount, pos, pos + 1}}
		curr := runes[idx]
		idx++
		pos++
		switch runeType(curr) {
		case CT_WHITESPACE:
			tok.ty = TT_WHITESPACE
			if last.ty == TT_UNKNOWN || last.ty == TT_TERMINATOR || last.ty == TT_WHITESPACE {
				outputToken = false
			}
		case CT_NEWLINE:
			// line++
			lineCount++
			pos = 0
			tok.ty = TT_TERMINATOR
			lastTy := last.ty
			if lastTy == TT_WHITESPACE {
				lastTy = beforeLast.ty
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
			for idx < len(runes) && runeType(runes[idx]) == CT_DIGIT {
				tok.data += string(runes[idx])
				idx++
				pos++
			}
			tok.ty = TT_LITERAL
			tok.pos.end = pos
		case CT_LETTER:
			for idx < len(runes) && (runeType(runes[idx]) == CT_DIGIT || runeType(runes[idx]) == CT_LETTER) {
				tok.data += string(runes[idx])
				idx++
				pos++
			}
			if tok.data == "false" || tok.data == "true" {
				tok.ty = TT_LITERAL
			} else if isOperator(tok.data) {
				tok.ty = opType(tok.data)
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
			case TT_SINGLE_QUOTE, TT_DOUBLE_QUOTE:
				tok.ty = TT_LITERAL
				close := rune(tok.data[0])
				tok.data = ""

				for idx < len(runes) {
					if runes[idx] == close {
						count := 0
						for count < idx && runes[idx-count-1] == '\\' {
							count++
						}
						if count%2 == 0 {
							break
						}
					} else if runes[idx] == '\n' {
						// line++
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
						// line++
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
			// TODO - special character sequences like '\n'
			break
		case CT_ILLEGAL:
			break
		}
		if pos != 0 {
			tok.pos.end = pos
		}
		if outputToken {
			tokens <- tok
			beforeLast = last
			last = tok
		}
	}
}
