package main

import (
	"regexp"
	"sort"
	"strconv"
	"strings"
	"unicode"
)

var predeclaredFuncs = map[string]func(Value, ListValue, Position) Value{
	"len": func(input Value, params ListValue, pos Position) Value {
		assertParamsNum(0, params, pos)
		return StringValue{strconv.Itoa(len(input.String()))}
	},
	"count": func(input Value, params ListValue, pos Position) Value {
		assertParamsNum(0, params, pos)
		return StringValue{strconv.Itoa(len(valAsList(input).vals))}
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
		var min Value
		var minVal int
		for _, i := range valAsList(input).vals {
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
		var max Value
		var maxVal int
		for _, i := range valAsList(input).vals {
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
		ret := ListValue{}
		for _, i := range valAsList(input).vals {
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
		vals := valAsList(input).vals
		if len(vals) == 1 {
			count = strings.Count(input.String(), params.vals[0].String())
		} else {
			for _, i := range vals {
				if i.String() == params.vals[0].String() {
					count++
				}
			}
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
		// TODO - fix this to fit docs.txt
		return createBoolValue(len([]rune(input.String())) == 1 && unicode.IsUpper([]rune(input.String())[0]))
	},
	"isLower": func(input Value, params ListValue, pos Position) Value {
		assertParamsNum(0, params, pos)
		// TODO - fix this to fit docs.txt
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
		for _, i := range valAsList(input).vals {
			ret.val += i.String()
		}
		return ret
	},
	"fold": func(input Value, params ListValue, pos Position) Value {
		assertParamsNum(1, params, pos)
		v := valAsList(input)
		if len(v.vals) == 0 {
			return NullValue{}
		}
		ret := v.vals[len(v.vals)-1]
		for i := len(v.vals) - 2; i >= 0; i-- {
			list := ListValue{[]Value{v.vals[i], ret}}
			ret = callDefinition(params.vals[0], input, list, pos)
		}
		return ret
	},
	"foldr": func(input Value, params ListValue, pos Position) Value {
		assertParamsNum(1, params, pos)
		v := valAsList(input)
		if len(v.vals) == 0 {
			return NullValue{}
		}
		ret := v.vals[len(v.vals)-1]
		for i := len(v.vals) - 2; i >= 0; i-- {
			list := ListValue{[]Value{v.vals[i], ret}}
			ret = callDefinition(params.vals[0], input, list, pos)
		}
		return ret
	},
	"foldl": func(input Value, params ListValue, pos Position) Value {
		assertParamsNum(1, params, pos)
		v := valAsList(input)
		if len(v.vals) == 0 {
			return NullValue{}
		}
		ret := v.vals[0]
		for i := 1; i < len(v.vals); i++ {
			list := ListValue{[]Value{ret, v.vals[i]}}
			ret = callDefinition(params.vals[0], input, list, pos)
		}
		return ret
	},
	"sort": func(input Value, params ListValue, pos Position) Value {
		assertParamsNum(1, params, pos)
		v := valAsList(input)
		sort.SliceStable(v.vals, func(i, j int) bool {
			a := atoi(callDefinition(params.vals[0], v.vals[i], ListValue{}, pos).String(), pos)
			b := atoi(callDefinition(params.vals[0], v.vals[j], ListValue{}, pos).String(), pos)
			return a < b
		})
		return input
	},
	"reverse": func(input Value, params ListValue, pos Position) Value {
		assertParamsNum(0, params, pos)
		v := valAsList(input)
		if len(v.vals) == 1 {
			r := []rune(input.String())
			for i, j := 0, len(r)-1; i < len(r)/2; i, j = i+1, j-1 {
				r[i], r[j] = r[j], r[i]
			}
			return StringValue{string(r)}
		} else {
			for i := len(v.vals)/2 - 1; i >= 0; i-- {
				opp := len(v.vals) - 1 - i
				v.vals[i], v.vals[opp] = v.vals[opp], v.vals[i]
			}
			return v
		}
	},
	"replace": func(input Value, params ListValue, pos Position) Value {
		assertParamsNum(2, params, pos)
		return StringValue{strings.ReplaceAll(input.String(), params.vals[0].String(), params.vals[1].String())}
	},
	"bool": func(input Value, params ListValue, pos Position) Value {
		assertParamsNum(0, params, pos)
		if input.String() != "" {
			return StringValue{"true"}
		}
		return StringValue{"false"}
	},
}
