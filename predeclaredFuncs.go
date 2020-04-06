package main

import (
	"regexp"
	"strconv"
	"strings"
	"unicode"
)

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
	"join": func(input Value, params ListValue, pos Position) Value {
		assertParamsNum(0, params, pos)
		ret := StringValue{}
		vals := input.(ListValue).vals
		for _, i := range vals {
			ret.val += i.String()
		}
		return ret
	},
	"fold":  foldRFunc,
	"foldr": foldRFunc,
	"foldl": foldLFunc,
}

func foldRFunc(input Value, params ListValue, pos Position) Value {
	assertParamsNum(1, params, pos)
	v := input.(ListValue)
	if len(v.vals) == 1 {
		return v.vals[0]
	}
	list := ListValue{[]Value{v.vals[0], foldRFunc(ListValue{v.vals[1:]}, params, pos)}}
	return callDefinition(params.vals[0], input, list, pos)
}

func foldLFunc(input Value, params ListValue, pos Position) Value {
	assertParamsNum(1, params, pos)
	v := input.(ListValue)
	if len(v.vals) == 1 {
		return v.vals[0]
	}
	list := ListValue{[]Value{v.vals[len(v.vals)-1], foldLFunc(ListValue{v.vals[:len(v.vals)-1]}, params, pos)}}
	return callDefinition(params.vals[len(v.vals)-1], input, list, pos)
}
