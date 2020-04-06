package main

import (
	"os"
	"strings"
	"unicode"

	"github.com/disiqueira/gotree"
	"github.com/fatih/color"
	"github.com/peterh/liner"
)

var allUserInput = []string{} // TODO - read history from liner
var lastLine = ""

func ioSetup() {
	globals.liner = liner.NewLiner()
	globals.liner.SetWordCompleter(wordCompleter)
}

func wordCompleter(line string, pos int) (string, []string, string) {
	if len(line) == 0 {
		return "", []string{"help", "exit", ""}, ""
	}

	low := pos - 1
	if low < 0 || low > len([]rune(line)) || !unicode.IsLetter([]rune(line)[low]) {
		return line[:pos], []string{}, line[pos:]
	}
	for low > 0 && unicode.IsLetter([]rune(line)[low-1]) {
		low--
	}
	high := pos - 1
	for high < len([]rune(line)) && unicode.IsLetter([]rune(line)[high]) {
		high++
	}
	word := line[low:high]
	toLower := strings.ToLower(word)

	completions := []string{}

	for k := range predeclaredFuncs {
		if strings.HasPrefix(strings.ToLower(k), toLower) {
			completions = append(completions, k)
		}
	}
	for i := len(definitions) - 1; i >= 0; i-- {
		for k := range definitions[i] {
			if strings.HasPrefix(strings.ToLower(k), toLower) {
				completions = append(completions, k)
			}
		}
	}
	for i := len(values) - 1; i >= 0; i-- {
		for k := range values[i] {
			if strings.HasPrefix(strings.ToLower(k), toLower) {
				completions = append(completions, k)
			}
		}
	}
	for _, s := range getWordOperators() {
		if strings.HasPrefix(strings.ToLower(s), toLower) {
			completions = append(completions, s)
		}
	}
	completions = append(completions, word)

	return line[:low], completions, line[high:]
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
