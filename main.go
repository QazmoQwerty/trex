package main

import (
	"os"
	"reflect"

	"github.com/disiqueira/gotree"
	"github.com/fatih/color"
)

func main() {
	args := os.Args[1:]

	input := ""
	fileNames := []string{}

	for _, arg := range args {
		if arg[0] == '-' { // Special arguments
			// TODO - what arguments do we need?
		} else {
			if input != "" {
				fileNames = append(fileNames, arg)
			} else {
				input = arg
			}
		}
	}

	if input == "" {
		input = "abcd"
		// println("ERROR: missing input string")
		// return
	}

	if len(fileNames) == 0 {
		startInterpreter(input)
	}

	// const input = "abcd"

	// 	const prog = `

	// // a => []
	// // b(n) => n << []
	// // a b(10)

	// sum(a, b) => a + b
	// len => 1 + len [1:] if [] else 0
	// count => len
	// fold(f) => () if len = 0 else f([0], fold(#f) [1:])
	// numOccurs(n) => 0 if not [] else (1 if [0] = n else 0) + numOccurs(n) [1:]
	// toBool => "true" if [] else "false"

	// "numOccurs:"
	// numOccurs(4) (4, 2, 4, 4, 6)
	// "fold by sum:"
	// fold(#sum) (1, 2, 3, 4)

	// isPrime(n) => count (i for i in 2..n where n % i = 0) = 0

	// primes(n) => i for i in 0..n where isPrime(i)

	// "toBool:"
	// toBool isPrime(13)

	// "primes:"
	// primes(100)

	// merge(a, b) => (a, b) if len a = 1 and len b = 1 else a + b

	// stringAdd(a, b) => a << b
	// collapse => fold(#stringAdd)

	// has(n) => count (i for i in [] where i = n) != 0

	// unique => () if len = 0 else merge(([-1] if numOccurs([-1]) = 1 else ()), (unique [:-1]))

	// collapse (1, 2, 3, 4, 5)

	// foo => 12344321

	// "foo:"
	// foo

	// "unique:"
	// unique foo

	// "collapse unique:"
	// collapse unique foo

	// `

	// 	tokens := make(chan Token)
	// 	go lexProgram(prog, tokens)
	// 	// for {
	// 	// 	showToken(<-tokens)
	// 	// }
	// 	tokMan := createTokenChanManager(tokens)
	// 	ast := parseProgram(&tokMan, TT_EOF)
	// 	println(printAst(ast).Print())
	// 	if !isNil(ast) {
	// 		for _, n := range ast.lines {
	// 			runLine(n, StringValue{input})
	// 		}
	// 	}
}

func startInterpreter(input string) {
	tokens := make(chan Token)
	for true {
		go lexLine(tokens, true)
		tokMan := createTokenChanManager(tokens)
		ast := parseProgram(&tokMan, TT_EOF)
		// println(printAst(ast).Print())
		if !isNil(ast) {
			for _, n := range ast.lines {
				runLine(n, StringValue{input})
			}
		}
	}
}

func recoverer() {
	if err := recover(); err != nil {
		switch e := err.(type) {
		case error:
			printError(e)
		default:
			panic(err)
		}
	}
}

func runLine(node Node, input Value) {
	defer recoverer()
	val := node.interpret(input)
	switch node.(type) {
	case Definition:
		break
	default:
		println(val.String())
		// println(toString(val, input))
	}
}

func printError(err error) {
	switch e := err.(type) {
	case myErr:
		whiteBold := color.New(color.FgWhite).Add(color.Bold).PrintfFunc()
		redBold := color.New(color.FgRed).Add(color.Bold).PrintfFunc()
		whiteBold("At line %d\n", e.pos.line)
		print("    " + allUserInput[e.pos.line-1] + "    ")
		for i := 0; i < e.pos.start; i++ {
			if allUserInput[e.pos.line-1][i] == '\t' {
				print("\t")
			} else {
				print(" ")
			}
		}
		for i := e.pos.start; i < e.pos.end; i++ {
			redBold("^")
		}
		redBold("\nError: ")
		println(e.msg)
	default:
		println(e.Error())
	}
}

func isNil(i interface{}) bool {
	if i == nil {
		return true
	}
	switch reflect.TypeOf(i).Kind() {
	case reflect.Ptr, reflect.Map, reflect.Array, reflect.Chan, reflect.Slice:
		return reflect.ValueOf(i).IsNil()
	}
	return false
}

func printAst(ast Node) gotree.Tree {
	if isNil(ast) {
		return gotree.New("{}")
	}
	tree := gotree.New(ast.toString())
	for _, n := range ast.getChildren() {
		if n == nil {
			tree.Add("")
		} else {
			tree.AddTree(printAst(n))
		}
	}
	return tree
}
