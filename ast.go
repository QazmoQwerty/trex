package main

type Node interface {
	getPosition() Position
	toString() string
	getChildren() []Node
	interpret(input Value) Value
}

type Statement interface {
	Node
}

type Expression interface {
	Node
}

type AnonDefinition struct {
	ids IdentifierList
	exp Expression
	pos Position
}

func (node AnonDefinition) getPosition() Position {
	return node.pos
}

func (node AnonDefinition) toString() string {
	return "->"
}

func (node AnonDefinition) getChildren() []Node {
	return []Node{node.ids, node.exp}
}

type BinaryOperation struct {
	left  Expression
	right Expression
	op    Operator
	pos   Position
}

func (node BinaryOperation) getPosition() Position {
	return node.pos
}

func (node BinaryOperation) toString() string {
	return node.op.str
}

func (node BinaryOperation) getChildren() []Node {
	return []Node{node.left, node.right}
}

type EmptyExpression struct {
	pos Position
}

func (node EmptyExpression) getPosition() Position {
	return node.pos
}
func (node EmptyExpression) toString() string {
	return "()"
}

func (node EmptyExpression) getChildren() []Node {
	return []Node{}
}

type UnaryOperation struct {
	expression Expression
	op         Operator
	pos        Position
}

func (node UnaryOperation) getPosition() Position {
	return node.pos
}
func (node UnaryOperation) toString() string {
	return node.op.str
}

func (node UnaryOperation) getChildren() []Node {
	return []Node{node.expression}
}

type Conditional struct {
	thenBranch Expression
	elseBranch Expression
	condition  Expression
	pos        Position
}

func (node Conditional) getPosition() Position {
	return node.pos
}

func (node Conditional) toString() string {
	return "<if>"
}

func (node Conditional) getChildren() []Node {
	return []Node{node.condition, node.thenBranch, node.elseBranch}
}

type Comprehension struct {
	exp   Expression
	fors  []ForClause
	where Expression
	pos   Position
}

type ForClause struct {
	id  Identifier
	exp Expression
}

func (node Comprehension) getPosition() Position {
	return node.pos
}

func (node Comprehension) toString() string {
	return "<for>"
}

func (node Comprehension) getChildren() []Node {
	ret := []Node{node.exp}
	for _, i := range node.fors {
		ret = append(ret, i.id)
		ret = append(ret, i.exp)
	}
	ret = append(ret, node.where)
	return ret
}

type Identifier struct {
	id  string
	pos Position
}

func (node Identifier) getPosition() Position {
	return node.pos
}

func (node Identifier) toString() string {
	return node.id
}

func (node Identifier) getChildren() []Node {
	return nil
}

type ExpressionList struct {
	expressions []Expression
	pos         Position
}

func (node ExpressionList) getPosition() Position {
	return node.pos
}

func (node ExpressionList) toString() string {
	return "<expList>"
}

func (node ExpressionList) getChildren() []Node {
	arr := make([]Node, len(node.expressions))
	for i := 0; i < len(arr); i++ {
		arr[i] = node.expressions[i]
	}
	return arr
}

type FunctionCall struct {
	callee Expression
	params ExpressionList
	arg    Expression
	pos    Position
}

func (node FunctionCall) getPosition() Position {
	return node.pos
}

func (node FunctionCall) toString() string {
	return "<call>"
}

func (node FunctionCall) getChildren() []Node {
	return []Node{node.callee, node.params, node.arg}
}

type Literal struct {
	value string
	pos   Position
}

func (node Literal) getPosition() Position {
	return node.pos
}

func (node Literal) toString() string {
	return "\"" + node.value + "\""
}

func (node Literal) getChildren() []Node {
	return nil
}

type Subscript struct {
	expression Expression
	idx1       Expression
	idx2       Expression
	idx3       Expression
	pos        Position
}

func (node Subscript) getPosition() Position {
	return node.pos
}

func (node Subscript) toString() string {
	return "[]"
}

func (node Subscript) getChildren() []Node {
	switch {
	case node.idx2 == nil:
		return []Node{node.expression, node.idx1}
	case node.idx3 == nil:
		return []Node{node.expression, node.idx1, node.idx2}
	default:
		return []Node{node.expression, node.idx1, node.idx2, node.idx3}
	}

}

type IdentifierList struct {
	identifiers []Identifier
	pos         Position
}

func (node IdentifierList) getPosition() Position {
	return node.pos
}

func (node IdentifierList) toString() string {
	return "<idList>"
}

func (node IdentifierList) getChildren() []Node {
	arr := make([]Node, len(node.identifiers))
	for i := 0; i < len(arr); i++ {
		arr[i] = node.identifiers[i]
	}
	return arr
}

type Definition struct {
	id      Identifier
	params  IdentifierList
	content Program
	pos     Position
}

func (node Definition) getPosition() Position {
	return node.pos
}

func (node Definition) toString() string {
	return node.id.id + "=>"
}

func (node Definition) getChildren() []Node {
	var content Node = node.content
	if len(node.content.lines) == 1 {
		content = node.content.lines[0]
	}
	if len(node.params.identifiers) == 0 {
		return []Node{content}
	} else if len(node.params.identifiers) == 1 {
		return []Node{node.params.identifiers[0], content}
	}
	return []Node{node.params, content}
}

type Program struct {
	lines []Node
	pos   Position
}

func (node Program) getPosition() Position {
	return node.pos
}

func (node Program) toString() string {
	return "{}"
}

func (node Program) getChildren() []Node {
	return node.lines
}
