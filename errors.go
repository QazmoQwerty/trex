package main

type myErr struct {
	msg string
	pos Position
	ty  ErrorType
}

func (err myErr) Error() string {
	return err.msg
}

type ErrorType int

const (
	ERR_GENERAL = iota
	ERR_LEXER
	ERR_PARSER
	ERR_INTERPRETER
)
