package main

type Node interface {
	getPosition() Position
	toString() string
	getChildren() []Node
	interpret(input string) (error, Value)
}

type Statement interface {
	Node
}

type Expression interface {
	Node
}

type BinaryOperation struct {
	left  Expression
	right Expression
	op    Operator
	pos   Position
}

func (this BinaryOperation) getPosition() Position {
	return this.pos
}

func (this BinaryOperation) toString() string {
	return this.op.str
}

func (this BinaryOperation) getChildren() []Node {
	return []Node{this.left, this.right}
}

type UnaryOperation struct {
	expression Expression
	op         Operator
	pos        Position
}

func (this UnaryOperation) getPosition() Position {
	return this.pos
}
func (this UnaryOperation) toString() string {
	return this.op.str
}

func (this UnaryOperation) getChildren() []Node {
	return []Node{this.expression}
}

type Conditional struct {
	thenBranch Expression
	elseBranch Expression
	condition  Expression
	pos        Position
}

func (this Conditional) getPosition() Position {
	return this.pos
}

func (this Conditional) toString() string {
	return "<if>"
}

func (this Conditional) getChildren() []Node {
	return []Node{this.condition, this.thenBranch, this.elseBranch}
}

type ForEach struct {
	exp Expression
	ids IdentifierList
	in  Expression
	pos Position
}

func (this ForEach) getPosition() Position {
	return this.pos
}

func (this ForEach) toString() string {
	return "<for>"
}

func (this ForEach) getChildren() []Node {
	return []Node{this.exp, this.ids, this.in}
}

type Identifier struct {
	id  string
	pos Position
}

func (this Identifier) getPosition() Position {
	return this.pos
}

func (this Identifier) toString() string {
	return this.id
}

func (this Identifier) getChildren() []Node {
	return nil
}

type ExpressionList struct {
	expressions []Expression
	pos         Position
}

func (this ExpressionList) getPosition() Position {
	return this.pos
}

func (this ExpressionList) toString() string {
	return "<expList>"
}

func (this ExpressionList) getChildren() []Node {
	arr := make([]Node, len(this.expressions))
	for i := 0; i < len(arr); i++ {
		arr[i] = this.expressions[i]
	}
	return arr
}

type FunctionCall struct {
	callee Expression
	params ExpressionList
	arg    Expression
	pos    Position
}

func (this FunctionCall) getPosition() Position {
	return this.pos
}

func (this FunctionCall) toString() string {
	return "<call>"
}

func (this FunctionCall) getChildren() []Node {
	return []Node{this.callee, this.params, this.arg}
}

type Literal struct {
	value string
	pos   Position
}

func (this Literal) getPosition() Position {
	return this.pos
}

func (this Literal) toString() string {
	return "\"" + this.value + "\""
}

func (this Literal) getChildren() []Node {
	return nil
}

type Subscript struct {
	expression Expression
	idx1       Expression
	idx2       Expression
	idx3       Expression
	pos        Position
}

func (this Subscript) getPosition() Position {
	return this.pos
}

func (this Subscript) toString() string {
	return "<[]>"
}

func (this Subscript) getChildren() []Node {
	switch {
	case this.idx2 == nil:
		return []Node{this.expression, this.idx1}
	case this.idx3 == nil:
		return []Node{this.expression, this.idx1, this.idx2}
	default:
		return []Node{this.expression, this.idx1, this.idx2, this.idx3}
	}

}

type IdentifierList struct {
	identifiers []Identifier
	pos         Position
}

func (this IdentifierList) getPosition() Position {
	return this.pos
}

func (this IdentifierList) toString() string {
	return "<idList>"
}

func (this IdentifierList) getChildren() []Node {
	arr := make([]Node, len(this.identifiers))
	for i := 0; i < len(arr); i++ {
		arr[i] = this.identifiers[i]
	}
	return arr
}

type Definition struct {
	id      Identifier
	params  IdentifierList
	content Program
	pos     Position
}

func (this Definition) getPosition() Position {
	return this.pos
}

func (this Definition) toString() string {
	return this.id.id + "=>"
}

func (this Definition) getChildren() []Node {
	var content Node = this.content
	if len(this.content.lines) == 1 {
		content = this.content.lines[0]
	}
	if len(this.params.identifiers) == 0 {
		return []Node{content}
	} else if len(this.params.identifiers) == 1 {
		return []Node{this.params.identifiers[0], content}
	}
	return []Node{this.params, content}
}

type Program struct {
	lines []Node
	pos   Position
}

func (this Program) getPosition() Position {
	return this.pos
}

func (this Program) toString() string {
	return "{}"
}

func (this Program) getChildren() []Node {
	return this.lines
}
