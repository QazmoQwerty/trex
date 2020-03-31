package main

import (
	"os"

	"github.com/disiqueira/gotree"
	"github.com/fatih/color"
	"github.com/peterh/liner"
)

var allUserInput = []string{} // TODO - read history from liner
var lastLine = ""

func ioSetup() {
	globals.liner = liner.NewLiner()
}

func ioExit() {
	globals.liner.Close()
	os.Exit(0)
}

func insertLine(line string) {
	allUserInput = append(allUserInput, line)
	lastLine = line
}

func readLine(prompt string) string {
	line, err := globals.liner.Prompt(prompt)
	if err != nil {
		panic(err)
	}
	globals.liner.AppendHistory(line)
	insertLine(line)
	return line + "\n"
}

func showHelp() {
	println("Help still needs to be written")
	println("For now see the Language Specification: gitlab.com/QazmoQwerty/trex/-/blob/master/docs/trex-spec.md")
}

func printError(err error) {
	switch e := err.(type) {
	case myErr:
		whiteBold := color.New(color.FgWhite).Add(color.Bold).PrintfFunc()
		redBold := color.New(color.FgRed).Add(color.Bold).PrintfFunc()
		if e.pos.line == lineCount-1 {
			print("    ")
		} else {
			whiteBold("At line %d\n", e.pos.line)
			print("    " + allUserInput[e.pos.line-1] + "\n    ")
		}
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
