package main

func parseProgram(tokens *TokenChanManager, expected TokenType) (error, *Program) {
	prog := Program{nil, tokens.peek().pos}
	eatWS(tokens)
	skipTerminator(tokens)
	for !eatToken(tokens, expected) {
		err, node := parse(tokens, 0)
		if err != nil {
			for !eatToken(tokens, TT_TERMINATOR) && !eatToken(tokens, expected) {
				tokens.next() // forward to the next line since this one is invalid
			}
			return err, nil
		}
		prog.lines = append(prog.lines, node)
		if eatToken(tokens, expected) {
			return nil, &prog
		}
		expectToken(tokens, TT_TERMINATOR)
	}
	return nil, &prog
}

func skipTerminator(tokens *TokenChanManager) {
	eatToken(tokens, TT_TERMINATOR)
}

func eatWS(tokens *TokenChanManager) bool {
	return eatToken(tokens, TT_WHITESPACE)
}

func parse(tokens *TokenChanManager, lastPrecedence byte) (error, Node) {
	var left Node
	err, left := nud(tokens)
	if err != nil {
		return err, left
	}
	if left == nil {
		return nil, nil
	}

	for {
		ateWS := eatWS(tokens)
		if precedence(tokens.peek()) <= lastPrecedence {
			if ateWS && tokens.peek().ty != TT_TERMINATOR {
				tokens.insertAtFront(Token{TT_WHITESPACE, " ", tokens.peek().pos})
			}
			return nil, left
		}

		err, left = led(tokens, left, ateWS)
		if err != nil {
			return err, nil
		}
	}
	// for precedence(tokens.peek()) > lastPrecedence {
	// 	err, left = led(tokens, left)
	// 	if err != nil {
	// 		return err, nil
	// 	}
	// }
	// return nil, left
}

func expectToken(tokens *TokenChanManager, ty TokenType) error {
	if !eatToken(tokens, ty) {
		return myErr{"\"" + getOperatorByType(ty).str + "\" expected", tokens.peek().pos, ERR_PARSER}
	}
	return nil
}

func eatToken(tokens *TokenChanManager, ty TokenType) bool {
	if tokens.peek().ty == ty {
		tokens.next()
		return true
	}
	return false
}

func parseExpressionList(tokens *TokenChanManager, prec byte) (error, *ExpressionList) {
	err, exp := parse(tokens, prec)
	if err != nil {
		return err, nil
	}
	if exp == nil {
		return nil, nil
	}
	err, ret := convertToExpressionList(exp)
	if err != nil {
		return err, nil
	}
	return nil, ret
}

func parseExpression(tokens *TokenChanManager, prec byte) (error, Expression) {
	err, node := parse(tokens, prec)
	if err != nil {
		return err, nil
	}
	if node == nil {
		return nil, nil
	}
	switch v := node.(type) {
	case Expression:
		return nil, v
	default:
		return myErr{"expected an expression", v.getPosition(), ERR_PARSER}, nil
	}
}

func parseIdentifierList(tokens *TokenChanManager) (error, *IdentifierList) {
	err, node := parse(tokens, 0)
	if err != nil {
		return err, nil
	}
	err, ret := convertToIdentifierList(node)
	if err != nil {
		return err, nil
	}
	return nil, ret
}

func convertToExpression(node Node) (error, Expression) {
	switch v := node.(type) {
	case Expression:
		return nil, v
	default:
		return myErr{"expected an expression", v.getPosition(), ERR_PARSER}, nil
	}
}

func convertToIdentifier(node Node) (error, *Identifier) {
	switch v := node.(type) {
	case *Identifier:
		return nil, v
	default:
		return myErr{"expected an identifier", v.getPosition(), ERR_PARSER}, nil
	}
}

func convertToIdentifierList(node Node) (error, *IdentifierList) {
	switch v := node.(type) {
	case *ExpressionList:
		ret := IdentifierList{[]Identifier{}, v.pos}
		for _, exp := range v.expressions {
			switch id := exp.(type) {
			case *Identifier:
				ret.identifiers = append(ret.identifiers, *id)
			default:
				return myErr{"expected an identifier", id.getPosition(), ERR_PARSER}, nil
			}
		}
		return nil, &ret
	case *IdentifierList:
		return nil, v
	case *Identifier:
		return nil, &IdentifierList{[]Identifier{*v}, v.pos}
	default:
		return myErr{"expected an identifier list", v.getPosition(), ERR_PARSER}, nil
	}
}

func convertToExpressionList(node Node) (error, *ExpressionList) {
	switch v := node.(type) {
	case *ExpressionList:
		return nil, v
	case Expression:
		return nil, &ExpressionList{[]Expression{v}, v.getPosition()}
	default:
		return myErr{"expected an expression list", v.getPosition(), ERR_PARSER}, nil
	}
}

func nud(tokens *TokenChanManager) (error, Expression) {
	eatWS(tokens)
	token := tokens.peek()
	switch tokens.peek().ty {
	case TT_IDENTIFIER:
		tokens.next()
		node := Identifier{token.data, token.pos}
		return nil, &node
	case TT_LITERAL:
		tokens.next()
		node := Literal{token.data, token.pos}
		return nil, &node
	case TT_PARENTHESIS_OPEN:
		tokens.next()
		err, inner := parseExpression(tokens, 0)
		if err != nil {
			return err, nil
		}
		err = expectToken(tokens, TT_PARENTHESIS_CLOSE)
		if err != nil {
			return err, nil
		}
		return nil, inner
	case TT_SQUARE_BRACKETS_OPEN:
		tokens.next()

		err, first := parseExpression(tokens, leftPrecedenceByOp(getOperatorByType(TT_COLON)))
		if err != nil {
			return err, nil
		}
		node := Subscript{nil, first, nil, nil, token.pos}
		eatWS(tokens)
		if eatToken(tokens, TT_COLON) {
			err, node.idx2 = parseExpression(tokens, leftPrecedenceByOp(getOperatorByType(TT_COLON)))
			if err != nil {
				return err, nil
			}
			if node.idx2 == nil {
				node.idx2 = &Literal{"", tokens.peek().pos}
			}
		} else if eatToken(tokens, TT_SQUARE_BRACKETS_CLOSE) {
			return nil, &node
		} else {
			return myErr{"expected ']' or ':'", tokens.peek().pos, ERR_PARSER}, nil
		}
		eatWS(tokens)
		if eatToken(tokens, TT_COLON) {
			err, node.idx3 = parseExpression(tokens, leftPrecedenceByOp(getOperatorByType(TT_COLON)))
			if err != nil {
				return err, nil
			}
			if node.idx3 == nil {
				node.idx3 = &Literal{"", tokens.peek().pos}
			}
		} else if eatToken(tokens, TT_SQUARE_BRACKETS_CLOSE) {
			return nil, &node
		} else {
			return myErr{"expected ']' or ':'", tokens.peek().pos, ERR_PARSER}, nil
		}
		eatWS(tokens)
		err = expectToken(tokens, TT_SQUARE_BRACKETS_CLOSE)
		if err != nil {
			return err, nil
		}
		return nil, &node
	}
	if isUnaryOperator(token.ty) {
		tokens.next()
		op := getOperator(token.data)
		err, exp := parseExpression(tokens, leftPrecedence(token))
		if err != nil {
			return err, nil
		}
		ret := UnaryOperation{exp, op, token.pos}
		return nil, &ret
	}
	return nil, nil
}

func led(tokens *TokenChanManager, node Node, ateWS bool) (error, Node) {
	token := tokens.peek()
	err, left := convertToExpression(node)
	if err != nil {
		return err, nil
	}
	switch token.ty {
	case TT_IF:
		tokens.next()
		err, cond := parseExpression(tokens, 0)
		if err != nil {
			return err, nil
		}
		eatWS(tokens)
		err = expectToken(tokens, TT_ELSE)
		if err != nil {
			return err, nil
		}
		err, elseB := parseExpression(tokens, 0)
		if err != nil {
			return err, nil
		}
		pos := Position{left.getPosition().line, left.getPosition().start, elseB.getPosition().end}
		return nil, &Conditional{left, elseB, cond, pos}
	case TT_COMMA:
		tokens.next()
		err, right := parseExpression(tokens, 0)
		if err != nil {
			return err, nil
		}
		list := ExpressionList{nil, left.getPosition()}
		list.expressions = append(list.expressions, left)
		switch r := right.(type) {
		case ExpressionList:
			for _, i := range r.expressions {
				list.expressions = append(list.expressions, i)
			}
		default:
			list.expressions = append(list.expressions, r)
		}
		return nil, &list
	case TT_DEFINE:
		tokens.next()
		var def Definition
		err, id := convertToIdentifier(left)
		if err != nil {
			return err, nil
		}
		def.id = *id
		def.pos = id.pos
		err, prog := parseExpression(tokens, 0)
		if err != nil {
			return err, nil
		}
		def.content = Program{[]Node{prog}, prog.getPosition()}
		def.pos.end = def.content.pos.end
		return nil, &def
	case TT_CURLY_BRACES_OPEN:
		tokens.next()
		err, id := convertToIdentifier(left)
		if err != nil {
			return err, nil
		}
		err, prog := parseProgram(tokens, TT_CURLY_BRACES_CLOSE)
		if err != nil {
			return err, nil
		}
		return nil, &Definition{*id, IdentifierList{}, *prog, id.pos}
	case TT_SQUARE_BRACKETS_OPEN:
		if ateWS {
			break
		}

		tokens.next()
		err, first := parseExpression(tokens, leftPrecedenceByOp(getOperatorByType(TT_COLON)))
		if err != nil {
			return err, nil
		}
		node := Subscript{left, first, nil, nil, token.pos}
		eatWS(tokens)
		if eatToken(tokens, TT_COLON) {
			err, node.idx2 = parseExpression(tokens, leftPrecedenceByOp(getOperatorByType(TT_COLON)))
			if err != nil {
				return err, nil
			}
			if node.idx2 == nil {
				node.idx2 = &Literal{"", tokens.peek().pos}
			}
		} else if eatToken(tokens, TT_SQUARE_BRACKETS_CLOSE) {
			return nil, &node
		} else {
			return myErr{"expected ']' or ':'", tokens.peek().pos, ERR_PARSER}, nil
		}
		eatWS(tokens)
		if eatToken(tokens, TT_COLON) {
			err, node.idx3 = parseExpression(tokens, leftPrecedenceByOp(getOperatorByType(TT_COLON)))
			if err != nil {
				return err, nil
			}
			if node.idx3 == nil {
				node.idx3 = &Literal{"", tokens.peek().pos}
			}
		} else if eatToken(tokens, TT_SQUARE_BRACKETS_CLOSE) {
			return nil, &node
		} else {
			return myErr{"expected ']' or ':'", tokens.peek().pos, ERR_PARSER}, nil
		}
		eatWS(tokens)
		err = expectToken(tokens, TT_SQUARE_BRACKETS_CLOSE)
		if err != nil {
			return err, nil
		}
		return nil, &node
	}
	if getOperator(token.data).isBinary {
		tokens.next()
		node := BinaryOperation{left, nil, getOperator(token.data), token.pos}
		var err error
		err, node.right = parseExpression(tokens, leftPrecedence(token))
		if err != nil {
			return err, nil
		}
		return nil, &node
	}

	err, right := parseExpression(tokens, 255)
	if err != nil {
		return err, nil
	}
	if ateWS {
		return nil, &FunctionCall{left, ExpressionList{}, right, left.getPosition()}
	}
	ateWS = eatWS(tokens)
	if eatToken(tokens, TT_DEFINE) {
		err, id := convertToIdentifier(left)
		if err != nil {
			return err, nil
		}
		err, params := convertToIdentifierList(right)
		if err != nil {
			return err, nil
		}
		err, exp := parseExpression(tokens, 0)
		if err != nil {
			return err, nil
		}
		prog := Program{[]Node{exp}, exp.getPosition()}
		return nil, &Definition{*id, *params, prog, id.pos}
	}
	if eatToken(tokens, TT_CURLY_BRACES_OPEN) {
		err, id := convertToIdentifier(left)
		if err != nil {
			return err, nil
		}
		err, params := convertToIdentifierList(right)
		if err != nil {
			return err, nil
		}
		err, prog := parseProgram(tokens, TT_CURLY_BRACES_CLOSE)
		if err != nil {
			return err, nil
		}
		return nil, &Definition{*id, *params, *prog, id.pos}
	}
	err, params := convertToExpressionList(right)
	if err != nil {
		return err, nil
	}
	if !ateWS {
		return nil, &FunctionCall{left, *params, nil, left.getPosition()}
	}
	err, arg := parseExpression(tokens, functionPrecedence)
	if err != nil {
		return err, nil
	}
	return nil, &FunctionCall{left, *params, arg, left.getPosition()}
}
