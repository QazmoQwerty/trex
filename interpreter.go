package main

import (
	"strconv"
	"strings"
)

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

type PredeclaredDefinitionValue struct {
	fn func(Value, ListValue, Position) Value
}

func (this StringValue) String() string {
	return this.val
}

func (this ListValue) String() string {
	ret := ""
	for i, v := range this.vals {
		if i != 0 {
			ret += ", "
		}
		switch v.(type) {
		case ListValue:
			ret += "(" + v.String() + ")"
		default:
			ret += v.String()
		}
	}
	return ret
}

func (this PredeclaredDefinitionValue) String() string {
	return "<#Definition>"
}

func (this NullValue) String() string {
	return ""
}

func (this DefinitionValue) String() string {
	return "<#Definition>"
}

func assertParamsNum(expected int, list ListValue, pos Position) {
	if len(list.vals) != expected {
		panic(myErr{"incorrect parameter count.\n    have: " + strconv.Itoa(len(list.vals)) +
			"\n    want: " + strconv.Itoa(expected), pos, ERR_INTERPRETER})
	}
}

var predeclaredFuncs = map[string]func(Value, ListValue, Position) Value{
	"len": func(input Value, params ListValue, pos Position) Value {
		assertParamsNum(0, params, pos)
		switch v := input.(type) {
		case ListValue:
			return StringValue{strconv.Itoa(len(v.vals))}
		case StringValue:
			return StringValue{strconv.Itoa(len(v.val))}
		case NullValue:
			return StringValue{"0"}
		case DefinitionValue:
			return StringValue{"1"}
		}
		return nil
	},
	"split": func(input Value, params ListValue, pos Position) Value {
		assertParamsNum(1, params, pos)
		ret := ListValue{}
		for _, i := range strings.Split(input.String(), params.vals[0].String()) {
			ret.vals = append(ret.vals, StringValue{i})
		}
		return ret
	},
	"lines": func(input Value, params ListValue, pos Position) Value {
		assertParamsNum(0, params, pos)
		ret := ListValue{}
		for _, i := range strings.Split(input.String(), "\n") {
			ret.vals = append(ret.vals, StringValue{i})
		}
		return ret
	},
}

var definitions = []map[string]Definition{map[string]Definition{}}
var values = []map[string]Value{map[string]Value{}}

func (this Program) interpret(input Value) Value {
	enterBlock()

	if len(this.lines) == 1 {
		val := this.lines[0].interpret(input)
		exitBlock()
		return val
	}

	ret := StringValue{""}
	for i, n := range this.lines {
		s := n.interpret(input)

		switch s.(type) {
		case NullValue, *NullValue:
			break
		default:
			// ret.val = toString(s, input)
			ret.val += s.String()

			if i+1 != len(this.lines) {
				ret.val += "\n"
			}
		}
	}
	exitBlock()
	return ret
}

func (this Definition) interpret(input Value) Value {
	definitions[len(definitions)-1][this.id.id] = this
	return NullValue{}
}

func (this Literal) interpret(input Value) Value {
	return StringValue{this.value}
}

func (this EmptyExpression) interpret(input Value) Value {
	return NullValue{}
}

func (this Identifier) interpret(input Value) Value {
	if fn, ok := predeclaredFuncs[this.id]; ok {
		return PredeclaredDefinitionValue{fn}
	}
	for i := len(definitions) - 1; i >= 0; i-- {
		if val, ok := values[i][this.id]; ok {
			return val
		}
		if val, ok := definitions[i][this.id]; ok {
			return DefinitionValue{val}
		}
	}
	panic(myErr{"undefined identifier \"" + this.id + "\"", this.pos, ERR_INTERPRETER})
}

func createBoolValue(b bool) StringValue {
	if b {
		return StringValue{"1"}
	}
	return StringValue{""}
}

func atoi(str string) int {
	ret, err := strconv.Atoi(str)
	if err != nil {
		panic(err)
	}
	return ret
}

func (this BinaryOperation) interpret(input Value) Value {
	left := this.left.interpret(input)

	right := this.right.interpret(input)

	// leftStr := toString(left, input)
	leftStr := left.String()

	// rightStr := toString(right, input)
	rightStr := right.String()

	switch this.op.ty {
	case TT_STRING_ADD:
		return StringValue{leftStr + rightStr}
	case TT_EQUAL:
		return createBoolValue(leftStr == rightStr)
	case TT_NOT_EQUAL:
		return createBoolValue(leftStr != rightStr)
	case TT_AND:
		return createBoolValue(leftStr != "" && rightStr != "")
	case TT_OR:
		return createBoolValue(leftStr != "" || rightStr != "")
	case TT_ADD:
		switch l := left.(type) {
		case ListValue:
			switch r := right.(type) {
			case ListValue:
				list := ListValue{l.vals}
				for _, i := range r.vals {
					list.vals = append(list.vals, i)
				}
				return list
			case NullValue:
				return l
			default:
				return ListValue{append(l.vals, r)}
			}
		case NullValue:
			return right
		}
		switch r := right.(type) {
		case ListValue:
			switch l := left.(type) {
			case ListValue:
				list := ListValue{l.vals}
				for _, i := range r.vals {
					list.vals = append(list.vals, i)
				}
				return list
			case NullValue:
				return l
			default:
				list := ListValue{[]Value{l}}
				for _, i := range r.vals {
					list.vals = append(list.vals, i)
				}
				return list
			}
		case NullValue:
			return left
		}
		return StringValue{strconv.Itoa(atoi(leftStr) + atoi(rightStr))}
	case TT_SUB:
		return StringValue{strconv.Itoa(atoi(leftStr) - atoi(rightStr))}
	case TT_DIV:
		return StringValue{strconv.Itoa(atoi(leftStr) / atoi(rightStr))}
	case TT_MUL:
		return StringValue{strconv.Itoa(atoi(leftStr) * atoi(rightStr))}
	case TT_MOD:
		return StringValue{strconv.Itoa(atoi(leftStr) % atoi(rightStr))}
	case TT_RANGE:
		low := atoi(leftStr)
		high := atoi(rightStr)
		list := ListValue{}
		for i := low; i < high; i++ {
			list.vals = append(list.vals, StringValue{strconv.Itoa(i)})
		}
		return list
	default:
		panic(myErr{"unimplemented binary operator \"" + this.op.str + "\"", this.pos, ERR_INTERPRETER})
	}
}

func (this UnaryOperation) interpret(input Value) Value {
	val := this.expression.interpret(input)

	if this.op.ty == TT_INDIRECTION {
		return val
	}

	// str := toString(val, input)
	str := val.String()

	switch this.op.ty {
	case TT_NOT:
		return createBoolValue(str == "")
	case TT_ADD:
		return StringValue{str}
	case TT_SUB:
		return StringValue{strconv.Itoa(-atoi(str))}
	default:
		panic(myErr{"unimplemented unary operator \"" + this.op.str + "\"", this.pos, ERR_INTERPRETER})
	}
}

func (this Conditional) interpret(input Value) Value {
	left := this.condition.interpret(input)
	if left.String() != "" {
		return this.thenBranch.interpret(input)
	}
	return this.elseBranch.interpret(input)
}

func enterBlock() {
	definitions = append(definitions, map[string]Definition{})
	values = append(values, map[string]Value{})
}

func exitBlock() {
	definitions = definitions[:len(definitions)-1]
	values = values[:len(values)-1]
}

func valToList(arr Value, input Value) []Value {
	list := []Value{}
	switch t := arr.(type) {
	default:
		break
	case ListValue:
		list = t.vals
		break
	case StringValue:
		for _, s := range t.val {
			list = append(list, StringValue{string(s)})
		}
		break
	}
	return list
}

func (this Comprehension) interpret(input Value) Value {
	list := valToList(this.fors[0].exp.interpret(input), input)
	ret := ListValue{}
	enterBlock()
	for _, v := range list {
		values[len(values)-1][this.fors[0].id.id] = v
		// if this.where == nil || toString(this.where.interpret(input), input) != "" {
		if this.where == nil || this.where.interpret(input).String() != "" {
			ret.vals = append(ret.vals, this.exp.interpret(input))
		}
	}
	exitBlock()
	return ret
}

func (this ExpressionList) interpret(input Value) Value {
	list := ListValue{}
	for _, n := range this.expressions {
		val := n.interpret(input)

		list.vals = append(list.vals, val)
	}
	return list
}

func (this FunctionCall) interpret(input Value) Value {
	val := this.callee.interpret(input)

	switch def := val.(type) {
	default:
		if this.arg == nil && len(this.params.expressions) == 0 {
			return def
		}
		panic(myErr{"cannot call non-definition value", this.pos, ERR_INTERPRETER})
	case PredeclaredDefinitionValue:
		params := ListValue{}
		for _, exp := range this.params.expressions {
			params.vals = append(params.vals, exp.interpret(input))
		}
		var inputVal Value
		if this.arg == nil {
			inputVal = input
		} else {
			inputVal = this.arg.interpret(input)
		}
		return def.fn(inputVal, params, this.pos)
	case DefinitionValue:
		enterBlock()
		if len(this.params.expressions) != len(def.def.params.identifiers) {
			panic(myErr{"incorrect parameter count\n    have: " + strconv.Itoa(len(this.params.expressions)) +
				"\n    want: " + strconv.Itoa(len(def.def.params.identifiers)), this.pos, ERR_INTERPRETER})
		}
		for i := 0; i < len(this.params.expressions); i++ {
			val := this.params.expressions[i].interpret(input)
			id := Identifier{def.def.params.identifiers[i].id, def.def.pos}
			values[len(values)-1][id.id] = val
		}
		var inputVal Value
		if this.arg == nil {
			inputVal = input
		} else {
			inputVal = this.arg.interpret(input)
		}
		ret := def.def.content.interpret(inputVal)
		exitBlock()
		return ret
	}

}

func (this Subscript) interpret(input Value) Value {
	var val Value
	if this.expression == nil {
		val = input
	} else {
		val = this.expression.interpret(input)
	}
	if this.idx1 == nil && this.idx2 == nil && this.idx3 == nil {
		return val
	}
	vals := valToList(val, input)

	if this.idx2 == nil && this.idx3 == nil {
		// idx := atoi(toString(this.idx1.interpret(input), input))
		idx := atoi(this.idx1.interpret(input).String())
		if idx < 0 {
			idx += len(vals)
		}
		return vals[idx]
	}
	if this.idx3 == nil {
		// lowStr := toString(this.idx1.interpret(input), input)
		// highStr := toString(this.idx2.interpret(input), input)
		lowStr := this.idx1.interpret(input).String()
		highStr := this.idx2.interpret(input).String()
		low, high := 0, len(vals)
		if lowStr != "" {
			low = atoi(lowStr)
			if low < 0 {
				low += len(vals)
			}
		}
		if highStr != "" {
			high = atoi(highStr)
			if high < 0 {
				high += len(vals)
			}
		}

		switch t := val.(type) {
		case ListValue:
			return ListValue{vals[low:high]}
		case StringValue:
			return StringValue{t.String()[low:high]}
		}
	}
	panic(myErr{"unimplemented interpret method7", this.pos, ERR_INTERPRETER})
}

func (this IdentifierList) interpret(input Value) Value {
	list := ListValue{}
	for _, n := range this.identifiers {
		val := n.interpret(input)

		list.vals = append(list.vals, val)
	}
	return list
}
