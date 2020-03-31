package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"reflect"

	"github.com/peterh/liner"

	"github.com/disiqueira/gotree"
	"github.com/fatih/color"
)

const version = "0.1.1"
const gitlabLink = "gitlab.com/QazmoQwerty/trex"

var globals struct {
	liner *liner.State
}

func exitProgram() {
	globals.liner.Close()
	os.Exit(0)
}

func main() {
	globals.liner = liner.NewLiner()
	defer globals.liner.Close()
	args := os.Args[1:]

	input := ""
	fileNames := []string{}
	redBold := color.New(color.FgRed).Add(color.Bold).PrintfFunc()

	for _, arg := range args {
		if arg[0] == '-' {
			switch arg {
			case "-h":
				println("Usage: trex <input> <files> [flags]")
				println("    ")
				println("    input: Either a file text inside square brackets []")
				println("    files: Files to be run. If no files are specified trex will run in interpreter mode.")
				println("    flags: ")
				println("        -h (show this message)")
				exitProgram()
			default:
				redBold("Error: ")
				println("Unknown flag \"" + arg + "\".")
				println("Usage: trex <input> <files> [arguments]")
				println("Try \"trex -h\" for more information.")
				exitProgram()
			}
		} else {
			if input != "" {
				fileNames = append(fileNames, arg)
			} else {
				input = arg
			}
		}
	}

	if input == "" {
		redBold("Error: ")
		println("missing input string")
		println("Try \"trex -h\" for more information.")
		exitProgram()
		return
	}

	for _, arg := range args {
		println(arg)
	}

	if input[0] == '[' && input[len(input)-1] == ']' {
		input = input[1 : len(input)-1]
	} else {
		content, err := ioutil.ReadFile(input)
		if err != nil {
			redBold("Error: ")
			println("could not open file \"" + input + "\"")
			exitProgram()
		}
		input = string(content)
	}

	if len(fileNames) == 0 {
		startInterpreter(input)
	}
}

func startInterpreter(input string) {
	fmt.Printf("Trex %s (%s)\n", version, gitlabLink)
	fmt.Printf("Type \"help\" for help, \"exit\" to exit.\n")
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
