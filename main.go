package main

import (
	"fmt"

	"github.com/disiqueira/gotree"
)

func main() {
	const input = "abcd"

	const prog = `
//d=>(as(10))10

//f(11) 12

// #1 + 2
a(x, y) => x(10) << y 10
// 	x << y
// }
a(10, 20)
// a 10
// E3 {
    // "" if len < 2 else substr(0, 1) << substr(-2, -1)
// }

// [1:2:3]

// a(10, 20) [0:1]

// E3_v2 {
    // "" if len < 2 else [0:1] << [-2:-1]
// }

// E4 {
//     [0] << replace([0], '$') [1:]
// }
`

	tokens := make(chan Token)
	go lexProgram(prog, tokens)
	// for {
	// 	showToken(<-tokens)
	// }
	tokMan := createTokenChanManager(tokens)
	ast := parseProgram(&tokMan, TT_EOF)
	println(printAst(ast).Print())
	for _, n := range ast.lines {
		runLine(n, input)
	}
}

func runLine(node Node, input string) {
	defer func() {
		if err := recover(); err != nil {
			switch e := err.(type) {
			case error:
				printError(e)
			default:
				panic(err)
			}
		}
	}()
	val := node.interpret(input)
	switch node.(type) {
	case *Definition:
		break
	default:
		s := toString(val, input)
		println(s)
	}
}

func printError(err error) {
	switch e := err.(type) {
	case myErr:
		fmt.Printf("%d:%d:%d - %s\n", e.pos.line, e.pos.start, e.pos.end, e.msg)
	default:
		println(e.Error())
	}
}

func printAst(ast Node) gotree.Tree {
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
