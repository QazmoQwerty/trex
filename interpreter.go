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

func (this Program) interpret(input string) (error, Value) {
	ret := StringValue{""}
	definitions = append(definitions, map[string]Definition{})
	for i, n := range this.lines {
		err, s := n.interpret(input)
		if err != nil {
			return err, nil
		}
		switch s.(type) {
		case NullValue, *NullValue:
			break
		default:
			err, ret.val = toString(s, input)
			if err != nil {
				return err, nil
			}
			if i+1 != len(this.lines) {
				ret.val += "\n"
			}
		}
	}
	definitions = definitions[:len(definitions)-1]
	return nil, ret
}

func toString(val Value, input string) (error, string) {
	switch t := val.(type) {
	case *NullValue:
		return nil, ""
	case DefinitionValue:
		err, v := t.def.content.interpret(input)
		if err != nil {
			return err, ""
		}
		return nil, v.String()
	default:
		return nil, val.String()
	}
}

func (this Definition) interpret(input string) (error, Value) {
	definitions[len(definitions)-1][this.id.id] = this
	return nil, NullValue{}
}

func (this Literal) interpret(input string) (error, Value) {
	return nil, StringValue{this.value}
}

func (this Identifier) interpret(input string) (error, Value) {
	for i := len(definitions) - 1; i >= 0; i-- {
		if val, ok := definitions[i][this.id]; ok {
			// return val.content.interpret(input)
			return nil, DefinitionValue{val}
		}
	}
	return myErr{"undefined identifier \"" + this.id + "\"", this.pos, ERR_INTERPRETER}, nil
}

func (this BinaryOperation) interpret(input string) (error, Value) {
	err, left := this.left.interpret(input)
	if err != nil {
		return err, nil
	}
	err, right := this.right.interpret(input)
	if err != nil {
		return err, nil
	}
	err, leftStr := toString(left, input)
	if err != nil {
		return err, nil
	}
	err, rightStr := toString(right, input)
	if err != nil {
		return err, nil
	}
	switch this.op.ty {
	case TT_STRING_ADD:
		return nil, StringValue{leftStr + rightStr}
	default:
		return myErr{"unimplemented binary operator \"" + this.op.str + "\"", this.pos, ERR_INTERPRETER}, nil
	}
}

func (this UnaryOperation) interpret(input string) (error, Value) {
	return myErr{"unimplemented interpret method 'unaryoperation'", this.pos, ERR_INTERPRETER}, nil
}

func (this Conditional) interpret(input string) (error, Value) {
	err, left := this.condition.interpret(input)
	if err != nil {
		return err, nil
	}
	if left.String() != "" { // 'true'
		err, ret := this.thenBranch.interpret(input)
		if err != nil {
			return err, nil
		}
		return nil, ret
	}
	err, ret := this.elseBranch.interpret(input)
	if err != nil {
		return err, nil
	}
	return nil, ret
}

func (this ForEach) interpret(input string) (error, Value) {
	return myErr{"unimplemented interpret method4", this.pos, ERR_INTERPRETER}, nil
}

func (this ExpressionList) interpret(input string) (error, Value) {
	list := ListValue{}
	for _, n := range this.expressions {
		err, val := n.interpret(input)
		if err != nil {
			return err, nil
		}
		list.vals = append(list.vals, val)
	}
	return nil, list
}

func (this FunctionCall) interpret(input string) (error, Value) {
	err, val := this.callee.interpret(input)
	if err != nil {
		return err, nil
	}

	switch def := val.(type) {
	default:
		return myErr{"cannot call non-definition value", this.pos, ERR_INTERPRETER}, nil
	case DefinitionValue:
		definitions = append(definitions, map[string]Definition{})

		if len(this.params.expressions) != len(def.def.params.identifiers) {
			return myErr{"incorrect parameter count\ncount is: " + strconv.Itoa(len(this.params.expressions)) + 
			"\nshould be: " + strconv.Itoa(len(def.def.params.identifiers)), this.pos, ERR_INTERPRETER}, nil
		}

		for i := 0; i < len(this.params.expressions); i++ {
			err, val := this.params.expressions[i].interpret(input)
			if err != nil {
				return err, nil
			}
			err, str := toString(val, input)
			if err != nil {
				return err, nil
			}
			id := Identifier{def.def.params.identifiers[i].id, def.def.pos}
			prog := Program{[]Node{Literal{str, def.def.pos}}, def.def.pos}
			param := Definition{id, IdentifierList{}, prog, def.def.pos}
			definitions[len(definitions)-1][id.id] = param
		}

		inputStr := ""

		if this.arg == nil {
			inputStr = input
		} else {
			err, exp := this.arg.interpret(input)
			if err != nil {
				return err, nil
			}
			err, inputStr = toString(exp, input)
			if err != nil {
				return err, nil
			}
		}

		err, ret := def.def.content.interpret(inputStr)
		if err != nil {
			return err, nil
		}
		definitions = definitions[:len(definitions)-1]
		return nil, ret
	}

}

func (this Subscript) interpret(input string) (error, Value) {
	return myErr{"unimplemented interpret method7", this.pos, ERR_INTERPRETER}, nil
}

func (this IdentifierList) interpret(input string) (error, Value) {
	list := ListValue{}
	for _, n := range this.identifiers {
		err, val := n.interpret(input)
		if err != nil {
			return err, nil
		}
		list.vals = append(list.vals, val)
	}
	return nil, list
}
