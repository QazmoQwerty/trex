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

func callDefinition(callee Value, input Value, params ListValue, pos Position) Value {
	switch def := callee.(type) {
	case PredeclaredDefinitionValue:
		return def.fn(input, params, pos)
	case DefinitionValue:
		enterBlock()
		if len(params.vals) != len(def.def.params.identifiers) {
			panic(myErr{"incorrect parameter count\n    have: " + strconv.Itoa(len(params.vals)) +
				"\n    want: " + strconv.Itoa(len(def.def.params.identifiers)), pos, ERR_INTERPRETER})
		}
		for i := 0; i < len(params.vals); i++ {
			id := Identifier{def.def.params.identifiers[i].id, def.def.pos}
			values[len(values)-1][id.id] = params.vals[i]
		}
		ret := def.def.content.interpret(input)
		exitBlock()
		return ret
	default:
		panic(myErr{"cannot call non-definition value", pos, ERR_INTERPRETER})
	}
}

// var definitions = []map[string]Definition{map[string]Definition{}}
// var values = []map[string]Value{map[string]Value{}}

var definitions = []map[string]Definition{{}}
var values = []map[string]Value{{}}

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
	if len(definitions) == 1 {
		globals.liner.RegisterFunction(this.id.id)
	}
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

func atoi(str string, pos Position) int {
	i, err := strconv.ParseInt(str, 0, 32)
	ret := int(i)
	if err != nil {
		if len(str) > 30 {
			panic(myErr{strconv.QuoteToGraphic(str[:30]) + `... cannot be converted to a number\n Note: full value was not shown due to length.`, pos, ERR_INTERPRETER})
		} else {
			panic(myErr{strconv.QuoteToGraphic(str) + ` cannot be converted to a number`, pos, ERR_INTERPRETER})
		}
	}
	return ret
}

func (this BinaryOperation) interpret(input Value) Value {
	left := this.left.interpret(input)
	right := this.right.interpret(input)
	leftPos := this.left.getPosition()
	rightPos := this.right.getPosition()

	switch this.op.ty {
	case TT_IN:
		switch r := right.(type) {
		case ListValue:
			subStr := left.String()
			for _, i := range r.vals {
				if strings.Contains(i.String(), subStr) {
					return createBoolValue(true)
				}
			}
			return createBoolValue(false)
		case StringValue:
			return createBoolValue(strings.Contains(right.String(), left.String()))
		default:
			return createBoolValue(false)
		}
	case TT_NOT_IN:
		switch r := right.(type) {
		case ListValue:
			subStr := left.String()
			for _, i := range r.vals {
				if strings.Contains(i.String(), subStr) {
					return createBoolValue(false)
				}
			}
			return createBoolValue(true)
		case StringValue:
			return createBoolValue(!strings.Contains(right.String(), left.String()))
		default:
			return createBoolValue(true)
		}
	case TT_STRING_ADD:
		return StringValue{left.String() + right.String()}
	case TT_STRING_MUL:
		return StringValue{strings.Repeat(left.String(), atoi(right.String(), this.right.getPosition()))}
	case TT_EQUAL:
		return createBoolValue(left.String() == right.String())
	case TT_NOT_EQUAL:
		return createBoolValue(left.String() != right.String())
	case TT_AND:
		return createBoolValue(left.String() != "" && right.String() != "")
	case TT_OR:
		return createBoolValue(left.String() != "" || right.String() != "")
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
		return StringValue{strconv.Itoa(atoi(left.String(), leftPos) + atoi(right.String(), rightPos))}
	case TT_SUB:
		return StringValue{strconv.Itoa(atoi(left.String(), leftPos) - atoi(right.String(), rightPos))}
	case TT_DIV:
		return StringValue{strconv.Itoa(atoi(left.String(), leftPos) / atoi(right.String(), rightPos))}
	case TT_MUL:
		return StringValue{strconv.Itoa(atoi(left.String(), leftPos) * atoi(right.String(), rightPos))}
	case TT_MOD:
		return StringValue{strconv.Itoa(atoi(left.String(), leftPos) % atoi(right.String(), rightPos))}
	case TT_RANGE:
		low := atoi(left.String(), leftPos)
		high := atoi(right.String(), rightPos)
		if low > high {
			list := ListValue{make([]Value, low-high)}
			for i := 0; i < low-high; i++ {
				list.vals[i] = StringValue{strconv.Itoa(low - i - 1)}
			}
			return list
		} else {
			list := ListValue{make([]Value, high-low)}
			for i := 0; i+low < high; i++ {
				list.vals[i] = StringValue{strconv.Itoa(i + low)}
			}
			return list
		}
	case TT_SMALLER:
		return createBoolValue(atoi(left.String(), leftPos) < atoi(right.String(), rightPos))
	case TT_SMALLER_EQUAL:
		return createBoolValue(atoi(left.String(), leftPos) <= atoi(right.String(), rightPos))
	case TT_GREATER:
		return createBoolValue(atoi(left.String(), leftPos) > atoi(right.String(), rightPos))
	case TT_GREATER_EQUAL:
		return createBoolValue(atoi(left.String(), leftPos) >= atoi(right.String(), rightPos))
	case TT_LEXICAL_SMALLER:
		return createBoolValue(strings.Compare(left.String(), right.String()) < 0)
	case TT_LEXICAL_SMALLER_EQUAL:
		return createBoolValue(strings.Compare(left.String(), right.String()) <= 0)
	case TT_LEXICAL_GREATER:
		return createBoolValue(strings.Compare(left.String(), right.String()) > 0)
	case TT_LEXICAL_GREATER_EQUAL:
		return createBoolValue(strings.Compare(left.String(), right.String()) >= 0)
	default:
		panic(myErr{"unimplemented binary operator \"" + this.op.str + "\"", this.pos, ERR_INTERPRETER})
	}
}

func (this UnaryOperation) interpret(input Value) Value {
	val := this.expression.interpret(input)

	if this.op.ty == TT_INDIRECTION {
		return val
	}

	str := val.String()

	switch this.op.ty {
	case TT_NOT:
		return createBoolValue(str == "")
	case TT_ADD:
		return StringValue{str}
	case TT_SUB:
		return StringValue{strconv.Itoa(-atoi(str, this.expression.getPosition()))}
	default:
		panic(myErr{"unimplemented unary operator \"" + this.op.str + "\"", this.pos, ERR_INTERPRETER})
	}
}

func (this AnonDefinition) interpret(input Value) Value {
	return DefinitionValue{
		Definition{
			Identifier{"", this.pos},
			this.ids,
			Program{[]Node{this.exp}, this.pos},
			this.pos,
		},
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

func valAsList(val Value) ListValue {
	switch v := val.(type) {
	case ListValue:
		return v
	case NullValue:
		return ListValue{}
	default:
		return ListValue{[]Value{val}}
	}
}

func valToList(arr Value) []Value {
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

func (this Comprehension) runComprehension(input Value, idx int, list []Value) ListValue {
	ret := ListValue{}
	enterBlock()
	switch len(this.fors) - idx {
	case 0:
		break
	case 1:
		for _, v := range list {
			values[len(values)-1][this.fors[idx].id.id] = v
			if this.where == nil || this.where.interpret(input).String() != "" {
				ret.vals = append(ret.vals, this.exp.interpret(input))
			}
		}
	default:
		for _, v := range list {
			values[len(values)-1][this.fors[idx].id.id] = v
			ret.vals = append(ret.vals, this.runComprehension(input, idx+1, list).vals...)
		}
	}
	exitBlock()
	return ret
}

func (this Comprehension) interpret(input Value) Value {
	return this.runComprehension(input, 0, valToList(this.fors[0].exp.interpret(input)))
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

func assertInRange(idx, len int, pos Position) {
	if idx < 0 || idx >= len {
		panic(myErr{
			"list index out of range [" + strconv.Itoa(idx) + "] with length " + strconv.Itoa(len),
			pos,
			ERR_INTERPRETER,
		})
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
	vals := valToList(val)

	if this.idx2 == nil && this.idx3 == nil {
		idx := atoi(this.idx1.interpret(input).String(), this.idx1.getPosition())
		if idx < 0 {
			idx += len(vals)
		}
		assertInRange(idx, len(vals), this.idx1.getPosition())
		return vals[idx]
	}

	lowStr := this.idx1.interpret(input).String()
	highStr := this.idx2.interpret(input).String()
	low, high := 0, len(vals)
	if lowStr != "" {
		low = atoi(lowStr, this.idx1.getPosition())
		if low < 0 {
			low += len(vals)
		}
		assertInRange(low, len(vals), this.idx1.getPosition())
	}
	if highStr != "" {
		high = atoi(highStr, this.idx2.getPosition())
		if high < 0 {
			high += len(vals)
		}
		assertInRange(high, len(vals), this.idx2.getPosition())
	}

	if this.idx3 == nil {
		switch t := val.(type) {
		case ListValue:
			return ListValue{vals[low:high]}
		case StringValue:
			return StringValue{t.String()[low:high]}
		}
	}

	stepStr := this.idx3.interpret(input).String()
	step := 1
	if stepStr != "" {
		step = atoi(stepStr, this.idx3.getPosition())
	}

	if step == 0 {
		panic(myErr{"slice step index cannot be zero", this.idx3.getPosition(), ERR_INTERPRETER})
	}
	switch t := val.(type) { // TODO - make this more efficient by preallocating memory
	case ListValue:
		newVals := ListValue{}
		if step > 0 {
			for i := low; i < high; i++ {
				if (i-low)%step == 0 {
					newVals.vals = append(newVals.vals, vals[i])
				}
			}
		} else {
			for i := high - 1; i >= low; i-- {
				if (len(vals)-1-i+low)%step == 0 {
					newVals.vals = append(newVals.vals, vals[i])
				}
			}
		}
		return newVals
	case StringValue:
		newStr := StringValue{}
		str := t.String()
		if step > 0 {
			for i := low; i < high; i++ {
				if i%step == 0 {
					newStr.val += string(str[i])
				}
			}
		} else {
			for i := high - 1; i >= low; i-- {
				if (len(vals)-1-i)%step == 0 {
					newStr.val += string(str[i])
				}
			}
		}
		return newStr
	}

	panic(myErr{"Third indices are not supported yet.", this.pos, ERR_INTERPRETER})
}

func (this IdentifierList) interpret(input Value) Value {
	list := ListValue{}
	for _, n := range this.identifiers {
		list.vals = append(list.vals, n.interpret(input))
	}
	return list
}
