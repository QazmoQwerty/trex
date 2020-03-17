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
	// errs := make(chan error)
	go lexProgram(prog, tokens)
	// for {
	// 	showToken(<-tokens)
	// }
	tokMan := createTokenChanManager(tokens)
	err, ast := parseProgram(&tokMan, TT_EOF)
	println(printAst(ast).Print())
	if err != nil {
		switch e := err.(type) {
		case myErr:
			fmt.Printf("%d:%d:%d - %s\n", e.pos.line, e.pos.start, e.pos.end, e.msg)
		default:
			println(e.Error())
		}
	} else {
		for _, n := range ast.lines {
			err, val := n.interpret(input)
			if err != nil {
				switch e := err.(type) {
				case myErr:
					fmt.Printf("%d:%d:%d - %s\n", e.pos.line, e.pos.start, e.pos.end, e.msg)
				default:
					println(e.Error())
				}
			} else {
				switch n.(type) {
				case *Definition:
					break
				default:
					err, s := toString(val, input)
					if err != nil {
						printError(err)
					} else {
						println(s)
					}
				}
			}
		}
		// fmt.Println(printAst(ast).Print())
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
