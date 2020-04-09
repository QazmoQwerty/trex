package main

type Operator struct {
	ty            TokenType
	str           string
	associativity bool
	precedence    byte
	isBinary      bool
}

func precedence(token Token) byte {
	switch token.ty {
	case TT_IDENTIFIER, TT_LITERAL, TT_WHITESPACE:
		return 135
	default:
		return getOperator(token.data).precedence
	}
}

func leftPrecedence(token Token) byte {
	op := getOperator(token.data)
	if op.associativity == RIGHT_TO_LEFT {
		return op.precedence - 1
	}
	return op.precedence
}

func leftPrecedenceByTy(ty TokenType) byte {
	return leftPrecedenceByOp(getOperatorByType(ty))
}

func leftPrecedenceByOp(op Operator) byte {
	if op.associativity == RIGHT_TO_LEFT {
		return op.precedence - 1
	}
	return op.precedence
}

func isUnaryOperator(ty TokenType) bool {
	return ty == TT_NOT || ty == TT_INDIRECTION || ty == TT_SUB || ty == TT_ADD || ty == TT_ANON_DEFINE
}

const functionPrecedence = 110

func getOperator(str string) Operator {
	switch str {
	case "(":
		return Operator{TT_PARENTHESIS_OPEN, str, LEFT_TO_RIGHT, 140, false}
	case ")":
		return Operator{TT_PARENTHESIS_CLOSE, str, LEFT_TO_RIGHT, 0, false}
	case "{":
		return Operator{TT_CURLY_BRACES_OPEN, str, LEFT_TO_RIGHT, 135, false}
	case "}":
		return Operator{TT_CURLY_BRACES_CLOSE, str, LEFT_TO_RIGHT, 0, false}
	case "[":
		return Operator{TT_SQUARE_BRACKETS_OPEN, str, LEFT_TO_RIGHT, 140, false}
	case "]":
		return Operator{TT_SQUARE_BRACKETS_CLOSE, str, LEFT_TO_RIGHT, 0, false}
	case "#":
		return Operator{TT_INDIRECTION, str, LEFT_TO_RIGHT, 130, false}
	case "not":
		return Operator{TT_NOT, str, RIGHT_TO_LEFT, 120, false}
	case "/":
		return Operator{TT_DIV, str, LEFT_TO_RIGHT, 100, true}
	case "%":
		return Operator{TT_MOD, str, LEFT_TO_RIGHT, 100, true}
	case "*":
		return Operator{TT_MUL, str, LEFT_TO_RIGHT, 100, true}
	case "**":
		return Operator{TT_STRING_MUL, str, LEFT_TO_RIGHT, 100, true}
	case "..":
		return Operator{TT_RANGE, str, LEFT_TO_RIGHT, 95, true}
	case "+":
		return Operator{TT_ADD, str, LEFT_TO_RIGHT, 90, true}
	case "-":
		return Operator{TT_SUB, str, LEFT_TO_RIGHT, 90, true}
	case "<<":
		return Operator{TT_STRING_ADD, str, LEFT_TO_RIGHT, 90, true}
	case ".<.":
		return Operator{TT_LEXICAL_SMALLER, str, LEFT_TO_RIGHT, 85, true}
	case ".<=.":
		return Operator{TT_LEXICAL_SMALLER_EQUAL, str, LEFT_TO_RIGHT, 85, true}
	case ".>.":
		return Operator{TT_LEXICAL_GREATER, str, LEFT_TO_RIGHT, 85, true}
	case ".>=.":
		return Operator{TT_LEXICAL_GREATER_EQUAL, str, LEFT_TO_RIGHT, 85, true}
	case "<":
		return Operator{TT_SMALLER, str, LEFT_TO_RIGHT, 80, true}
	case "<=":
		return Operator{TT_SMALLER_EQUAL, str, LEFT_TO_RIGHT, 80, true}
	case ">":
		return Operator{TT_GREATER, str, LEFT_TO_RIGHT, 80, true}
	case ">=":
		return Operator{TT_GREATER_EQUAL, str, LEFT_TO_RIGHT, 80, true}
	case "=":
		return Operator{TT_EQUAL, str, LEFT_TO_RIGHT, 70, true}
	case "!=":
		return Operator{TT_NOT_EQUAL, str, LEFT_TO_RIGHT, 70, true}
	case "and":
		return Operator{TT_AND, str, LEFT_TO_RIGHT, 60, true}
	case "or":
		return Operator{TT_OR, str, LEFT_TO_RIGHT, 50, true}
	case ",":
		return Operator{TT_COMMA, str, RIGHT_TO_LEFT, 40, false}
	case "where":
		return Operator{TT_WHERE, str, LEFT_TO_RIGHT, 35, true}
	case "if":
		return Operator{TT_IF, str, RIGHT_TO_LEFT, 30, false}
	case ":":
		return Operator{TT_COLON, str, LEFT_TO_RIGHT, 20, false}
	case "->":
		return Operator{TT_ANON_DEFINE, str, LEFT_TO_RIGHT, 15, true}
	case "=>":
		return Operator{TT_DEFINE, str, LEFT_TO_RIGHT, 10, true}
	case "|":
		return Operator{TT_TERMINATOR, str, false, 0, false}
	case "//":
		return Operator{TT_SINGLE_LINE_COMMENT, str, false, 0, false}
	case "/*":
		return Operator{TT_MULTI_LINE_COMMENT_OPEN, str, false, 0, false}
	case "*/":
		return Operator{TT_MULTI_LINE_COMMENT_CLOSE, str, false, 0, false}
	case "'":
		return Operator{TT_SINGLE_QUOTE, str, false, 0, false}
	case "\"":
		return Operator{TT_DOUBLE_QUOTE, str, false, 0, false}
	case "else":
		return Operator{TT_ELSE, str, false, 0, false}
	case "for":
		return Operator{TT_FOR, str, false, 35, false}
	case "in":
		return Operator{TT_IN, str, false, 45, true}
	case "not in":
		return Operator{TT_NOT_IN, str, false, 45, true}
	default:
		return Operator{TT_UNKNOWN, str, false, 0, false}
	}
}

func getWordOperators() []string {
	return []string{
		"else", "for", "in", "and", "if", "where", "or", "not", "exit", "help", "quit", "example",
	}
}

func getOperatorByType(op TokenType) Operator {
	switch op {
	case TT_PARENTHESIS_OPEN:
		return Operator{TT_PARENTHESIS_OPEN, "(", LEFT_TO_RIGHT, 140, true}
	case TT_PARENTHESIS_CLOSE:
		return Operator{TT_PARENTHESIS_CLOSE, ")", LEFT_TO_RIGHT, 0, false}
	case TT_CURLY_BRACES_OPEN:
		return Operator{TT_CURLY_BRACES_OPEN, "{", LEFT_TO_RIGHT, 135, true}
	case TT_CURLY_BRACES_CLOSE:
		return Operator{TT_CURLY_BRACES_CLOSE, "}", LEFT_TO_RIGHT, 0, false}
	case TT_SQUARE_BRACKETS_OPEN:
		return Operator{TT_SQUARE_BRACKETS_OPEN, "[", LEFT_TO_RIGHT, 140, true}
	case TT_SQUARE_BRACKETS_CLOSE:
		return Operator{TT_SQUARE_BRACKETS_CLOSE, "]", LEFT_TO_RIGHT, 0, false}
	case TT_INDIRECTION:
		return Operator{TT_INDIRECTION, "#", LEFT_TO_RIGHT, 130, false}
	case TT_NOT:
		return Operator{TT_NOT, "not", RIGHT_TO_LEFT, 120, false}
	case TT_DIV:
		return Operator{TT_DIV, "/", LEFT_TO_RIGHT, 100, true}
	case TT_MOD:
		return Operator{TT_MOD, "%", LEFT_TO_RIGHT, 100, true}
	case TT_MUL:
		return Operator{TT_MUL, "*", LEFT_TO_RIGHT, 100, true}
	case TT_STRING_MUL:
		return Operator{TT_STRING_MUL, "**", LEFT_TO_RIGHT, 100, true}
	case TT_RANGE:
		return Operator{TT_RANGE, "..", LEFT_TO_RIGHT, 95, true}
	case TT_ADD:
		return Operator{TT_ADD, "+", LEFT_TO_RIGHT, 90, true}
	case TT_SUB:
		return Operator{TT_SUB, "-", LEFT_TO_RIGHT, 90, true}
	case TT_STRING_ADD:
		return Operator{TT_STRING_ADD, "<<", LEFT_TO_RIGHT, 90, true}
	case TT_LEXICAL_SMALLER:
		return Operator{TT_SMALLER, ".<.", LEFT_TO_RIGHT, 85, true}
	case TT_LEXICAL_SMALLER_EQUAL:
		return Operator{TT_SMALLER_EQUAL, ".<=.", LEFT_TO_RIGHT, 85, true}
	case TT_LEXICAL_GREATER:
		return Operator{TT_GREATER, ".>.", LEFT_TO_RIGHT, 85, true}
	case TT_LEXICAL_GREATER_EQUAL:
		return Operator{TT_GREATER_EQUAL, ".>=.", LEFT_TO_RIGHT, 85, true}
	case TT_SMALLER:
		return Operator{TT_SMALLER, "<", LEFT_TO_RIGHT, 80, true}
	case TT_SMALLER_EQUAL:
		return Operator{TT_SMALLER_EQUAL, "<=", LEFT_TO_RIGHT, 80, true}
	case TT_GREATER:
		return Operator{TT_GREATER, ">", LEFT_TO_RIGHT, 80, true}
	case TT_GREATER_EQUAL:
		return Operator{TT_GREATER_EQUAL, ">=", LEFT_TO_RIGHT, 80, true}
	case TT_EQUAL:
		return Operator{TT_EQUAL, "=", LEFT_TO_RIGHT, 70, true}
	case TT_NOT_EQUAL:
		return Operator{TT_NOT_EQUAL, "!=", LEFT_TO_RIGHT, 70, true}
	case TT_AND:
		return Operator{TT_AND, "and", LEFT_TO_RIGHT, 60, true}
	case TT_OR:
		return Operator{TT_OR, "or", LEFT_TO_RIGHT, 50, true}
	case TT_COMMA:
		return Operator{TT_COMMA, ",", RIGHT_TO_LEFT, 40, false}
	case TT_WHERE:
		return Operator{TT_WHERE, "where", LEFT_TO_RIGHT, 35, true}
	case TT_IF:
		return Operator{TT_IF, "if", RIGHT_TO_LEFT, 30, false}
	case TT_COLON:
		return Operator{TT_COLON, ":", LEFT_TO_RIGHT, 20, false}
	case TT_ANON_DEFINE:
		return Operator{TT_ANON_DEFINE, "->", LEFT_TO_RIGHT, 15, true}
	case TT_DEFINE:
		return Operator{TT_DEFINE, "=>", LEFT_TO_RIGHT, 10, true}
	case TT_TERMINATOR:
		return Operator{TT_TERMINATOR, "|", false, 0, false}
	case TT_SINGLE_LINE_COMMENT:
		return Operator{TT_SINGLE_LINE_COMMENT, "//", false, 0, false}
	case TT_MULTI_LINE_COMMENT_OPEN:
		return Operator{TT_MULTI_LINE_COMMENT_OPEN, "/*", false, 0, false}
	case TT_MULTI_LINE_COMMENT_CLOSE:
		return Operator{TT_MULTI_LINE_COMMENT_CLOSE, "*/", false, 0, false}
	case TT_SINGLE_QUOTE:
		return Operator{TT_SINGLE_QUOTE, "'", false, 0, false}
	case TT_DOUBLE_QUOTE:
		return Operator{TT_DOUBLE_QUOTE, "\"", false, 0, false}
	case TT_ELSE:
		return Operator{TT_ELSE, "else", false, 0, false}
	case TT_FOR:
		return Operator{TT_FOR, "for", false, 32, false}
	case TT_NOT_IN:
		return Operator{TT_NOT_IN, "not in", false, 45, true}
	case TT_IN:
		return Operator{TT_IN, "in", false, 45, true}
	default:
		return Operator{TT_UNKNOWN, "", false, 0, false}
	}
}

func opType(str string) TokenType {
	return getOperator(str).ty
}

func isOperator(str string) bool {
	return opType(str) != TT_UNKNOWN
}

const (
	RIGHT_TO_LEFT      = true
	LEFT_TO_RIGHT      = false
	ASSOCIATIVITY_NONE = false
)

const (
	TT_UNKNOWN = iota
	TT_EOF
	TT_IDENTIFIER
	TT_WHITESPACE
	TT_TERMINATOR
	TT_LITERAL
	TT_OPERATOR
	TT_SINGLE_QUOTE
	TT_DOUBLE_QUOTE
	TT_TICK_QUOTE
	TT_SINGLE_LINE_COMMENT
	TT_MULTI_LINE_COMMENT_OPEN
	TT_MULTI_LINE_COMMENT_CLOSE
	TT_EQUAL
	TT_SMALLER
	TT_SMALLER_EQUAL
	TT_GREATER
	TT_NOT_EQUAL
	TT_GREATER_EQUAL
	TT_LEXICAL_SMALLER
	TT_LEXICAL_SMALLER_EQUAL
	TT_LEXICAL_GREATER
	TT_LEXICAL_GREATER_EQUAL
	TT_NOT
	TT_AND
	TT_OR
	TT_IF
	TT_ELSE
	TT_FOR
	TT_IN
	TT_NOT_IN
	TT_PARENTHESIS_OPEN
	TT_PARENTHESIS_CLOSE
	TT_CURLY_BRACES_OPEN
	TT_CURLY_BRACES_CLOSE
	TT_SQUARE_BRACKETS_OPEN
	TT_SQUARE_BRACKETS_CLOSE
	TT_COLON
	TT_COMMA
	TT_ADD
	TT_SUB
	TT_RANGE
	TT_MUL
	TT_DIV
	TT_STRING_ADD
	TT_STRING_MUL
	TT_MOD
	TT_WHERE
	TT_INDIRECTION
	TT_DEFINE
	TT_ANON_DEFINE
)
