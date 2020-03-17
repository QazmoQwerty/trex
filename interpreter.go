package main

import "strconv"

type Value interface {
	String() string
}

type ListValue struct {
	vals []Value
}

type StringValue struct {
	val string
}

type NullValue struct {
}

type DefinitionValue struct {
	def Definition
}

func (this StringValue) String() string {
	return this.val
}

func (this ListValue) String() string {
	ret := "["
	for i, v := range this.vals {
		if i != 0 {
			ret += ", "
		}
		ret += v.String()
	}
	ret += "]"
	return ret
}

func (this NullValue) String() string {
	return ""
}

func (this DefinitionValue) String() string {
	return "<#Definition>"
}

var definitions = []map[string]Definition{map[string]Definition{}}

func (this Program) interpret(input string) Value {
	ret := StringValue{""}
	definitions = append(definitions, map[string]Definition{})
	for i, n := range this.lines {
		s := n.interpret(input)

		switch s.(type) {
		case NullValue, *NullValue:
			break
		default:
			ret.val = toString(s, input)

			if i+1 != len(this.lines) {
				ret.val += "\n"
			}
		}
	}
	definitions = definitions[:len(definitions)-1]
	return ret
}

func toString(val Value, input string) string {
	switch t := val.(type) {
	case *NullValue:
		return ""
	case DefinitionValue:
		v := t.def.content.interpret(input)
		return v.String()
	default:
		return val.String()
	}
}

func (this Definition) interpret(input string) Value {
	definitions[len(definitions)-1][this.id.id] = this
	return NullValue{}
}

func (this Literal) interpret(input string) Value {
	return StringValue{this.value}
}

func (this Identifier) interpret(input string) Value {
	for i := len(definitions) - 1; i >= 0; i-- {
		if val, ok := definitions[i][this.id]; ok {
			// return val.content.interpret(input)
			return DefinitionValue{val}
		}
	}
	panic(myErr{"undefined identifier \"" + this.id + "\"", this.pos, ERR_INTERPRETER})
}

func (this BinaryOperation) interpret(input string) Value {
	left := this.left.interpret(input)

	right := this.right.interpret(input)

	leftStr := toString(left, input)

	rightStr := toString(right, input)

	switch this.op.ty {
	case TT_STRING_ADD:
		return StringValue{leftStr + rightStr}
	default:
		panic(myErr{"unimplemented binary operator \"" + this.op.str + "\"", this.pos, ERR_INTERPRETER})
	}
}

func (this UnaryOperation) interpret(input string) Value {
	panic(myErr{"unimplemented interpret method 'unaryoperation'", this.pos, ERR_INTERPRETER})
}

func (this Conditional) interpret(input string) Value {
	left := this.condition.interpret(input)

	if left.String() != "" { // 'true'
		ret := this.thenBranch.interpret(input)

		return ret
	}
	ret := this.elseBranch.interpret(input)

	return ret
}

func (this ForEach) interpret(input string) Value {
	panic(myErr{"unimplemented interpret method4", this.pos, ERR_INTERPRETER})
}

func (this ExpressionList) interpret(input string) Value {
	list := ListValue{}
	for _, n := range this.expressions {
		val := n.interpret(input)

		list.vals = append(list.vals, val)
	}
	return list
}

func (this FunctionCall) interpret(input string) Value {
	val := this.callee.interpret(input)

	switch def := val.(type) {
	default:
		panic(myErr{"cannot call non-definition value", this.pos, ERR_INTERPRETER})
	case DefinitionValue:
		definitions = append(definitions, map[string]Definition{})

		if len(this.params.expressions) != len(def.def.params.identifiers) {
			panic(myErr{"incorrect parameter count\ncount is: " + strconv.Itoa(len(this.params.expressions)) +
				"\nshould be: " + strconv.Itoa(len(def.def.params.identifiers)), this.pos, ERR_INTERPRETER})
		}

		for i := 0; i < len(this.params.expressions); i++ {
			val := this.params.expressions[i].interpret(input)

			str := toString(val, input)

			id := Identifier{def.def.params.identifiers[i].id, def.def.pos}
			prog := Program{[]Node{Literal{str, def.def.pos}}, def.def.pos}
			param := Definition{id, IdentifierList{}, prog, def.def.pos}
			definitions[len(definitions)-1][id.id] = param
		}

		inputStr := ""

		if this.arg == nil {
			inputStr = input
		} else {
			exp := this.arg.interpret(input)

			inputStr = toString(exp, input)

		}

		ret := def.def.content.interpret(inputStr)

		definitions = definitions[:len(definitions)-1]
		return ret
	}

}

func (this Subscript) interpret(input string) Value {
	panic(myErr{"unimplemented interpret method7", this.pos, ERR_INTERPRETER})
}

func (this IdentifierList) interpret(input string) Value {
	list := ListValue{}
	for _, n := range this.identifiers {
		val := n.interpret(input)

		list.vals = append(list.vals, val)
	}
	return list
}
