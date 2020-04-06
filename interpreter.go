package main

import (
	"regexp"
	"strconv"
	"strings"
	"unicode"
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
		case DefinitionValue, PredeclaredDefinitionValue:
			return StringValue{"1"}
		}
		return nil
	},
	"count": func(input Value, params ListValue, pos Position) Value {
		assertParamsNum(0, params, pos)
		switch v := input.(type) {
		case ListValue:
			return StringValue{strconv.Itoa(len(v.vals))}
		case StringValue, DefinitionValue, PredeclaredDefinitionValue:
			return StringValue{"1"}
		case NullValue:
			return StringValue{"0"}
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
	"words": func(input Value, params ListValue, pos Position) Value {
		assertParamsNum(0, params, pos)
		ret := ListValue{}
		for _, i := range strings.Fields(input.String()) {
			ret.vals = append(ret.vals, StringValue{i})
		}
		return ret
	},
	"chars": func(input Value, params ListValue, pos Position) Value {
		assertParamsNum(0, params, pos)
		str := []rune(input.String())
		ret := ListValue{make([]Value, len(str))}
		for i, c := range str {
			a := StringValue{string(c)}
			ret.vals[i] = a
		}
		return ret
	},
	"min": func(input Value, params ListValue, pos Position) Value {
		assertParamsNum(1, params, pos)
		list := input.(ListValue)
		var min Value
		var minVal int
		for _, i := range list.vals {
			currVal := atoi(callDefinition(params.vals[0], i, ListValue{}, pos).String(), pos)
			if min == nil || currVal < minVal {
				min = i
				minVal = currVal
			}
		}
		return min
	},
	"max": func(input Value, params ListValue, pos Position) Value {
		assertParamsNum(1, params, pos)
		list := input.(ListValue)
		var max Value
		var maxVal int
		for _, i := range list.vals {
			currVal := atoi(callDefinition(params.vals[0], i, ListValue{}, pos).String(), pos)
			if max == nil || currVal > maxVal {
				max = i
				maxVal = currVal
			}
		}
		return max
	},
	"unique": func(input Value, params ListValue, pos Position) Value {
		assertParamsNum(0, params, pos)
		list := input.(ListValue)
		ret := ListValue{}
		for _, i := range list.vals {
			b := true
			for _, j := range ret.vals {
				if j.String() == i.String() {
					b = false
					break
				}
			}
			if b {
				ret.vals = append(ret.vals, i)
			}
		}
		return ret
	},
	"numOccurs": func(input Value, params ListValue, pos Position) Value {
		assertParamsNum(1, params, pos)
		count := 0
		switch v := input.(type) {
		case ListValue:
			for _, i := range v.vals {
				if i.String() == params.vals[0].String() {
					count++
				}
			}
		case StringValue:
			count = strings.Count(v.val, params.vals[0].String())
		}
		return StringValue{strconv.Itoa(count)}
	},
	"toUpper": func(input Value, params ListValue, pos Position) Value {
		assertParamsNum(0, params, pos)
		return StringValue{strings.ToUpper(input.String())}
	},
	"toLower": func(input Value, params ListValue, pos Position) Value {
		assertParamsNum(0, params, pos)
		return StringValue{strings.ToLower(input.String())}
	},
	"isLetter": func(input Value, params ListValue, pos Position) Value {
		assertParamsNum(0, params, pos)
		return createBoolValue(len([]rune(input.String())) == 1 && unicode.IsLetter([]rune(input.String())[0]))
	},
	"isUpper": func(input Value, params ListValue, pos Position) Value {
		assertParamsNum(0, params, pos)
		return createBoolValue(len([]rune(input.String())) == 1 && unicode.IsUpper([]rune(input.String())[0]))
	},
	"isLower": func(input Value, params ListValue, pos Position) Value {
		assertParamsNum(0, params, pos)
		return createBoolValue(len([]rune(input.String())) == 1 && unicode.IsLower([]rune(input.String())[0]))
	},
	"isDigit": func(input Value, params ListValue, pos Position) Value {
		assertParamsNum(0, params, pos)
		return createBoolValue(len([]rune(input.String())) == 1 && unicode.IsDigit([]rune(input.String())[0]))
	},
	"ascii": func(input Value, params ListValue, pos Position) Value {
		assertParamsNum(0, params, pos)
		vals := ListValue{}
		for _, i := range []rune(input.String()) {
			vals.vals = append(vals.vals, StringValue{strconv.Itoa(int(i))})
		}
		if len(vals.vals) == 1 {
			return vals.vals[0]
		}
		return vals
	},
	"matches": func(input Value, params ListValue, pos Position) Value {
		assertParamsNum(1, params, pos)
		r := regexp.MustCompile(params.vals[0].String())
		matches := r.FindAllString(input.String(), -1)
		assertParamsNum(1, params, pos)
		ret := ListValue{make([]Value, len(matches))}
		for i := 0; i < len(matches); i++ {
			ret.vals[i] = StringValue{matches[i]}
		}
		return ret
	},
	"hasMatch": func(input Value, params ListValue, pos Position) Value {
		assertParamsNum(1, params, pos)
		r := regexp.MustCompile(params.vals[0].String())
		return createBoolValue(r.MatchString(input.String()))
	},
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

func atoi(str string, pos Position) int {
	ret, err := strconv.Atoi(str)
	if err != nil {
		if len(str) > 30 {
			panic(myErr{`"` + str[:30] + `"... cannot be converted to a number (full value was not shown due to length)`, pos, ERR_INTERPRETER})
		} else {
			panic(myErr{`"` + str + `" cannot be converted to a number`, pos, ERR_INTERPRETER})
		}
	}
	return ret
}

func (this BinaryOperation) interpret(input Value) Value {
	left := this.left.interpret(input)
	right := this.right.interpret(input)
	leftStr := left.String()
	rightStr := right.String()
	leftPos := this.left.getPosition()
	rightPos := this.right.getPosition()

	switch this.op.ty {
	case TT_STRING_ADD:
		return StringValue{leftStr + rightStr}
	case TT_STRING_MUL:
		return StringValue{strings.Repeat(leftStr, atoi(rightStr, this.right.getPosition()))}
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
		return StringValue{strconv.Itoa(atoi(leftStr, leftPos) + atoi(rightStr, rightPos))}
	case TT_SUB:
		return StringValue{strconv.Itoa(atoi(leftStr, leftPos) - atoi(rightStr, rightPos))}
	case TT_DIV:
		return StringValue{strconv.Itoa(atoi(leftStr, leftPos) / atoi(rightStr, rightPos))}
	case TT_MUL:
		return StringValue{strconv.Itoa(atoi(leftStr, leftPos) * atoi(rightStr, rightPos))}
	case TT_MOD:
		return StringValue{strconv.Itoa(atoi(leftStr, leftPos) % atoi(rightStr, rightPos))}
	case TT_RANGE:
		low := atoi(leftStr, leftPos)
		high := atoi(rightStr, rightPos)
		list := ListValue{}
		for i := low; i < high; i++ {
			list.vals = append(list.vals, StringValue{strconv.Itoa(i)})
		}
		return list
	case TT_SMALLER:
		return createBoolValue(atoi(leftStr, leftPos) < atoi(rightStr, rightPos))
	case TT_SMALLER_EQUAL:
		return createBoolValue(atoi(leftStr, leftPos) <= atoi(rightStr, rightPos))
	case TT_GREATER:
		return createBoolValue(atoi(leftStr, leftPos) > atoi(rightStr, rightPos))
	case TT_GREATER_EQUAL:
		return createBoolValue(atoi(leftStr, leftPos) >= atoi(rightStr, rightPos))
	case TT_LEXICAL_SMALLER:
		return createBoolValue(strings.Compare(leftStr, rightStr) < 0)
	case TT_LEXICAL_SMALLER_EQUAL:
		return createBoolValue(strings.Compare(leftStr, rightStr) <= 0)
	case TT_LEXICAL_GREATER:
		return createBoolValue(strings.Compare(leftStr, rightStr) > 0)
	case TT_LEXICAL_GREATER_EQUAL:
		return createBoolValue(strings.Compare(leftStr, rightStr) >= 0)
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
		idx := atoi(this.idx1.interpret(input).String(), this.idx1.getPosition())
		if idx < 0 {
			idx += len(vals)
		}
		return vals[idx]
	}
	if this.idx3 == nil {
		lowStr := this.idx1.interpret(input).String()
		highStr := this.idx2.interpret(input).String()
		low, high := 0, len(vals)
		if lowStr != "" {
			low = atoi(lowStr, this.idx1.getPosition())
			if low < 0 {
				low += len(vals)
			}
		}
		if highStr != "" {
			high = atoi(highStr, this.idx2.getPosition())
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
	panic(myErr{"Third indices are not supported yet.", this.pos, ERR_INTERPRETER})
}

func (this IdentifierList) interpret(input Value) Value {
	list := ListValue{}
	for _, n := range this.identifiers {
		list.vals = append(list.vals, n.interpret(input))
	}
	return list
}
