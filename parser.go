package main

func parseProgram(tokens *TokenChanManager, expected TokenType) *Program {
	prog := Program{nil, tokens.peek().pos}
	eatWS(tokens)
	eatToken(tokens, TT_TERMINATOR)
	defer func() {
		if err := recover(); err != nil {
			for !eatToken(tokens, TT_TERMINATOR) && !eatToken(tokens, expected) {
				tokens.next() // forward to the next line since this one is invalid
			}
			panic(err)
		}
	}()
	for !eatToken(tokens, expected) {
		node := parse(tokens, 0)
		prog.lines = append(prog.lines, node)
		if eatToken(tokens, expected) {
			return &prog
		}
		expectToken(tokens, TT_TERMINATOR)
	}
	return &prog
}

func eatWS(tokens *TokenChanManager) bool {
	return eatToken(tokens, TT_WHITESPACE)
}

func parse(tokens *TokenChanManager, lastPrecedence byte) Node {
	var left Node = nud(tokens)
	if left == nil {
		return nil
	}

	for {
		ateWS := eatWS(tokens)
		if precedence(tokens.peek()) <= lastPrecedence {
			if ateWS && tokens.peek().ty != TT_TERMINATOR {
				tokens.insertAtFront(Token{TT_WHITESPACE, " ", tokens.peek().pos})
			}
			return left
		}
		left = led(tokens, left, ateWS)
	}
}

func expectToken(tokens *TokenChanManager, ty TokenType) {
	if !eatToken(tokens, ty) {
		panic(myErr{"\"" + getOperatorByType(ty).str + "\" expected", tokens.peek().pos, ERR_PARSER})
	}
}

func eatToken(tokens *TokenChanManager, ty TokenType) bool {
	if tokens.peek().ty == ty {
		tokens.next()
		return true
	}
	return false
}

func parseExpressionList(tokens *TokenChanManager, prec byte) *ExpressionList {
	return convertToExpressionList(parse(tokens, prec))
}

func parseExpression(tokens *TokenChanManager, prec byte) Expression {
	return convertToExpression(parse(tokens, prec))
}

func parseIdentifierList(tokens *TokenChanManager) *IdentifierList {
	return convertToIdentifierList(parse(tokens, 0))
}

func convertToExpression(node Node) Expression {
	if node == nil {
		return nil
	}
	switch v := node.(type) {
	case Expression:
		return v
	default:
		panic(myErr{"expected an expression", v.getPosition(), ERR_PARSER})
	}
}

func convertToIdentifier(node Node) *Identifier {
	if node == nil {
		return nil
	}
	switch v := node.(type) {
	case *Identifier:
		return v
	default:
		panic(myErr{"expected an identifier", v.getPosition(), ERR_PARSER})
	}
}

func convertToIdentifierList(node Node) *IdentifierList {
	if node == nil {
		return nil
	}
	switch v := node.(type) {
	case *ExpressionList:
		ret := IdentifierList{[]Identifier{}, v.pos}
		for _, exp := range v.expressions {
			switch id := exp.(type) {
			case *Identifier:
				ret.identifiers = append(ret.identifiers, *id)
			default:
				panic(myErr{"expected an identifier", id.getPosition(), ERR_PARSER})
			}
		}
		return &ret
	case *IdentifierList:
		return v
	case *Identifier:
		return &IdentifierList{[]Identifier{*v}, v.pos}
	default:
		panic(myErr{"expected an identifier list", v.getPosition(), ERR_PARSER})
	}
}

func convertToExpressionList(node Node) *ExpressionList {
	if node == nil {
		return nil
	}
	switch v := node.(type) {
	case *ExpressionList:
		return v
	case Expression:
		return &ExpressionList{[]Expression{v}, v.getPosition()}
	default:
		panic(myErr{"expected an expression list", v.getPosition(), ERR_PARSER})
	}
}

func nud(tokens *TokenChanManager) Expression {
	eatWS(tokens)
	token := tokens.peek()
	switch tokens.peek().ty {
	case TT_IDENTIFIER:
		tokens.next()
		return &Identifier{token.data, token.pos}
	case TT_LITERAL:
		tokens.next()
		return &Literal{token.data, token.pos}
	case TT_PARENTHESIS_OPEN:
		tokens.next()
		inner := parseExpression(tokens, 0)
		expectToken(tokens, TT_PARENTHESIS_CLOSE)
		return inner
	case TT_SQUARE_BRACKETS_OPEN:
		tokens.next()
		node := Subscript{nil, nil, nil, nil, token.pos}
		node.idx1 = parseExpression(tokens, leftPrecedenceByOp(getOperatorByType(TT_COLON)))
		eatWS(tokens)
		if eatToken(tokens, TT_COLON) {
			node.idx2 = parseExpression(tokens, leftPrecedenceByOp(getOperatorByType(TT_COLON)))
			if node.idx2 == nil {
				node.idx2 = &Literal{"", tokens.peek().pos}
			}
		} else if eatToken(tokens, TT_SQUARE_BRACKETS_CLOSE) {
			return &node
		} else {
			panic(myErr{"expected ']' or ':'", tokens.peek().pos, ERR_PARSER})
		}
		eatWS(tokens)
		if eatToken(tokens, TT_COLON) {
			node.idx3 = parseExpression(tokens, leftPrecedenceByOp(getOperatorByType(TT_COLON)))
			if node.idx3 == nil {
				node.idx3 = &Literal{"", tokens.peek().pos}
			}
		} else if eatToken(tokens, TT_SQUARE_BRACKETS_CLOSE) {
			return &node
		} else {
			panic(myErr{"expected ']' or ':'", tokens.peek().pos, ERR_PARSER})
		}
		eatWS(tokens)
		expectToken(tokens, TT_SQUARE_BRACKETS_CLOSE)
		return &node
	}
	if isUnaryOperator(token.ty) {
		tokens.next()
		return &UnaryOperation{parseExpression(tokens, leftPrecedence(token)), getOperator(token.data), token.pos}
	}
	return nil
}

func led(tokens *TokenChanManager, node Node, ateWS bool) Node {
	token := tokens.peek()
	left := convertToExpression(node)
	switch token.ty {
	case TT_IF:
		tokens.next()
		cond := parseExpression(tokens, 0)
		eatWS(tokens)
		expectToken(tokens, TT_ELSE)
		elseB := parseExpression(tokens, 0)
		pos := Position{left.getPosition().line, left.getPosition().start, elseB.getPosition().end}
		return &Conditional{left, elseB, cond, pos}
	case TT_COMMA:
		tokens.next()
		right := parseExpression(tokens, 0)
		list := ExpressionList{nil, left.getPosition()}
		list.expressions = append(list.expressions, left)
		switch r := right.(type) {
		case *ExpressionList:
			for _, i := range r.expressions {
				list.expressions = append(list.expressions, i)
			}
		default:
			list.expressions = append(list.expressions, r)
		}
		return &list
	case TT_DEFINE:
		tokens.next()
		var def Definition
		id := convertToIdentifier(left)
		def.id = *id
		def.pos = id.pos
		prog := parseExpression(tokens, 0)
		def.content = Program{[]Node{prog}, prog.getPosition()}
		def.pos.end = def.content.pos.end
		return &def
	case TT_CURLY_BRACES_OPEN:
		tokens.next()
		return &Definition{
			*convertToIdentifier(left),
			IdentifierList{},
			*parseProgram(tokens, TT_CURLY_BRACES_CLOSE),
			left.getPosition(),
		}
	case TT_SQUARE_BRACKETS_OPEN:
		if ateWS {
			break
		}
		tokens.next()
		first := parseExpression(tokens, leftPrecedenceByOp(getOperatorByType(TT_COLON)))
		node := Subscript{left, first, nil, nil, token.pos}
		eatWS(tokens)
		if eatToken(tokens, TT_COLON) {
			node.idx2 = parseExpression(tokens, leftPrecedenceByOp(getOperatorByType(TT_COLON)))
			if node.idx2 == nil {
				node.idx2 = &Literal{"", tokens.peek().pos}
			}
		} else if eatToken(tokens, TT_SQUARE_BRACKETS_CLOSE) {
			return &node
		} else {
			panic(myErr{"expected ']' or ':'", tokens.peek().pos, ERR_PARSER})
		}
		eatWS(tokens)
		if eatToken(tokens, TT_COLON) {
			node.idx3 = parseExpression(tokens, leftPrecedenceByOp(getOperatorByType(TT_COLON)))
			if node.idx3 == nil {
				node.idx3 = &Literal{"", tokens.peek().pos}
			}
		} else if eatToken(tokens, TT_SQUARE_BRACKETS_CLOSE) {
			return &node
		} else {
			panic(myErr{"expected ']' or ':'", tokens.peek().pos, ERR_PARSER})
		}
		eatWS(tokens)
		expectToken(tokens, TT_SQUARE_BRACKETS_CLOSE)
		return &node
	}
	if getOperator(token.data).isBinary {
		tokens.next()
		return &BinaryOperation{left, parseExpression(tokens, leftPrecedence(token)), getOperator(token.data), token.pos}
	}

	right := parseExpression(tokens, 255)
	if ateWS {
		return &FunctionCall{left, ExpressionList{}, right, left.getPosition()}
	}
	ateWS = eatWS(tokens)
	if eatToken(tokens, TT_DEFINE) {
		id := convertToIdentifier(left)
		params := convertToIdentifierList(right)
		exp := parseExpression(tokens, 0)
		prog := Program{[]Node{exp}, exp.getPosition()}
		return &Definition{*id, *params, prog, id.pos}
	}
	if eatToken(tokens, TT_CURLY_BRACES_OPEN) {
		return &Definition{
			*convertToIdentifier(left),
			*convertToIdentifierList(right),
			*parseProgram(tokens, TT_CURLY_BRACES_CLOSE),
			left.getPosition(),
		}
	}
	if !ateWS {
		return &FunctionCall{left, *convertToExpressionList(right), nil, left.getPosition()}
	}
	return &FunctionCall{left, *convertToExpressionList(right), parseExpression(tokens, functionPrecedence), left.getPosition()}
}
