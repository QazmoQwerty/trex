package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"reflect"

	"github.com/fatih/color"
	"github.com/peterh/liner"
)

const version = "0.2.3"
const gitlabLink = "gitlab.com/QazmoQwerty/trex"

var globals struct {
	liner   *liner.State
	showAst bool
}

func main() {
	ioSetup()
	defer ioExit()
	args := os.Args[1:]

	input := ""
	fileNames := []string{}
	redBold := color.New(color.FgRed).Add(color.Bold).PrintfFunc()

	globals.showAst = false

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
				ioExit()
			case "-ast":
				globals.showAst = true
			default:
				redBold("Error: ")
				println("Unknown flag \"" + arg + "\".")
				println("Usage: trex <input> <files> [arguments]")
				println("Try \"trex -h\" for more information.")
				ioExit()
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
		ioExit()
		return
	}

	if input[0] == '[' && input[len(input)-1] == ']' {
		input = input[1 : len(input)-1]
	} else {
		content, err := ioutil.ReadFile(input)
		if err != nil {
			redBold("Error: ")
			println("could not open file \"" + input + "\"")
			ioExit()
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
	for true {
		tokens := TokenQueue{}
		lexLine(&tokens, true)
		ast := parseProgram(&tokens, TT_EOF)
		if globals.showAst {
			println(printAst(ast).Print())
		}
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
