package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
	"unicode"

	"github.com/disiqueira/gotree"
	"github.com/fatih/color"
	"gitlab.com/QazmoQwerty/go-liner-highlight"
)

var allUserInput = []string{}
var lastLine = ""

var trexKeywords = []string{
	"if", "else", "for", "in", "from", "not", "or", "and",
}

//Note: we should put the longest operators first.
var trexOperators = []string{
	"=", "!", "<", ">", "#", "+", "-", "*", "/", "%", ":",
	"|", "[", "{", "(", "]", "}", ")", ",", ".",
}

var colors = map[liner.Category](*color.Color){
	liner.NumberType:   color.New(color.FgHiRed),
	liner.StringType:   color.New(color.FgHiRed),
	liner.KeywordType:  color.New(color.FgHiBlue),
	liner.CommentType:  color.New(color.FgHiGreen),
	liner.FunctionType: color.New(color.FgHiYellow),
	liner.OperatorType: color.New(),
	liner.IdentType:    color.New(),
}

func ioSetup() {
	globals.liner = liner.NewLiner()
	globals.liner.SetWordCompleter(wordCompleter)
	globals.liner.SetSyntaxHighlight(globals.interpreterSyntaxHighlight)
	globals.liner.RegisterOperators(trexOperators)
	globals.liner.RegisterKeywords(trexKeywords)
	globals.liner.RegisterColors(colors)
	for k := range predeclaredFuncs {
		globals.liner.RegisterFunction(k)
	}
	globals.liner.RegisterFunctions([]string{"exit", "quit", "help", "example"})
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
	if globals.liner != nil {
		globals.liner.Close()
	}
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
		// blueBold := color.New(color.FgBlue).Add(color.Bold).PrintfFunc()
		line := ""
		if globals.codeFile != "" {
			var err error
			line, err = getLineOfFile(globals.codeFile, e.pos.line)
			if err != nil {
				println(err.Error())
				return
			}
			redBold(" --> ")
			whiteBold("line %d\n", e.pos.line)
			redBold("  | \n  | ")
			print(line)
			if !strings.HasSuffix(line, "\n") {
				println()
			}
			redBold("  | ")
		} else {
			if e.pos.line <= 0 || len(allUserInput) < e.pos.line {
				println("an internal error occurred...")
				return
			}
			if e.pos.line == lineCount-1 {
				print("    ")
				line = allUserInput[e.pos.line-1]
			} else {
				redBold(" --> ")
				whiteBold("line %d\n", e.pos.line)
				line = allUserInput[e.pos.line-1]
				redBold("  | \n  | ")
				print(line)
				if !strings.HasSuffix(line, "\n") {
					println()
				}
				redBold("  | ")
			}
		}
		for i := 0; i < e.pos.start && i < len(line); i++ {
			if line[i] == '\t' {
				print("\t")
			} else {
				print(" ")
			}
		}
		for i := e.pos.start; i < e.pos.end; i++ {
			redBold("^")
		}
		println()
		globals.errorColor.Print("Error:")
		println(" " + e.msg)
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

func getLineOfFile(fn string, n int) (string, error) {
	if n < 1 {
		return "", fmt.Errorf("invalid request: line %d", n)
	}
	f, err := os.Open(fn)
	if err != nil {
		return "", err
	}
	defer f.Close()
	bf := bufio.NewReader(f)
	var line string
	for lnum := 0; lnum < n; lnum++ {
		// line, err = bf.ReadString('\n')
		line, err = bf.ReadString('\n')
		if err == io.EOF {
			if lnum == n-1 && line != "" {
				return line, nil
			}
			switch lnum {
			case 0:
				return "", errors.New("no lines in file")
			case 1:
				return "", errors.New("only 1 line")
			default:
				return "", fmt.Errorf("only %d lines", lnum)
			}
		}
		if err != nil {
			return "", err
		}
	}
	if line == "" {
		return "", fmt.Errorf("line %d empty", n)
	}
	return line, nil
}
