package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"reflect"

	"github.com/fatih/color"
	"gitlab.com/QazmoQwerty/go-liner-highlight"
)

const gitlabLink = "gitlab.com/QazmoQwerty/trex"

var globals struct {
	liner                      *liner.State
	showAst                    bool
	showLex                    bool
	forceInterpret             bool
	interpreterSyntaxHighlight bool
	codeFile                   string
	errorColor                 *color.Color
	outputColor                *color.Color
}

func main() {
	args := os.Args[1:]

	input := ""
	fileNames := []string{}
	globals.errorColor = color.New(color.FgHiRed)
	globals.outputColor = color.New()
	globals.codeFile = ""
	globals.showAst = false
	globals.showLex = false
	globals.interpreterSyntaxHighlight = false
	globals.forceInterpret = false

	for _, arg := range args {
		if arg[0] == '-' {
			switch arg {
			case "-h":
				println(`Usage: trex <input> <files> [flags]
	input: Either a file text inside square brackets []
	files: Files to be run. If no files are specified trex will run in interpreter mode.
	flags: 
		-h (show this message)
		-i (run interpreter after code files have ben executed)
		-v (show version)
		-hl (turn on syntax highlighting in the interpreter)

	debug flags:
		-lex (show output of the lexer)
		-ast (show output of the parser)`)
				ioExit()
			case "-ast":
				globals.showAst = true
			case "-lex":
				globals.showLex = true
			case "-hl":
				globals.interpreterSyntaxHighlight = true
			case "-i":
				globals.forceInterpret = true
			case "-v":
				fmt.Printf("Trex %s (%s)\n", version, gitlabLink)
				ioExit()
			default:
				globals.errorColor.Print("Error:")
				println(" Unknown flag \"" + arg + "\".")
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

	ioSetup()
	defer ioExit()

	if globals.interpreterSyntaxHighlight {
		globals.outputColor = color.New(color.FgHiBlack)
	}

	if input == "" {
		globals.errorColor.Print("Error:")
		println(" missing input string")
		println("Try \"trex -h\" for more information.")
		ioExit()
		return
	}
	if input != "" {
		if input[0] == '[' && input[len(input)-1] == ']' {
			input = input[1 : len(input)-1]
		} else {
			content, err := ioutil.ReadFile(input)
			if err != nil {
				globals.errorColor.Print("Error:")
				println(" could not open file \"" + input + "\"")
				ioExit()
			}
			input = string(content)
		}
	}

	if len(fileNames) == 0 {
		startInterpreter(input)
	} else {
		for _, f := range fileNames {
			interpretFile(input, f)
		}
		if globals.forceInterpret {
			startInterpreter(input)
		}
	}
}

func interpretFile(input string, file string) {
	defer recoverer()
	globals.codeFile = file
	content, err := ioutil.ReadFile(file)
	if err != nil {
		println("op" + file)
		globals.errorColor.Print("Error:")
		println(" could not open file \"" + file + "\"")
		ioExit()
	}
	tokens := TokenQueue{}
	lexProgram(string(content), &tokens)
	if globals.showLex {
		for _, tok := range tokens.tokens {
			showToken(tok)
		}
	}
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

func startInterpreter(input string) {
	fmt.Printf("Trex %s (%s)\n", version, gitlabLink)
	fmt.Printf("Type \"help\" for help, \"exit\" to exit.\n")
	lineCount = 1
	globals.codeFile = ""
	for true {
		tokens := TokenQueue{}
		lexLine(&tokens, true)
		if globals.showLex {
			for _, tok := range tokens.tokens {
				showToken(tok)
			}
		}
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
		globals.outputColor.Print(val.String())
		println()
		// println(val.String())
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
