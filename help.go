package main

func showHelp(cmd string) {
	str := ""
	if len(cmd) >= 6 {
		str = cmd[5 : len(cmd)-1]
	}
	switch str {
	default:
		println(`No help exists for "` + str + `".`)
	case "":
		println(`Use \"help xxx\" to see help for a particular subject.
Or read the Language Specification: gitlab.com/QazmoQwerty/trex/-/blob/master/docs/trex-spec.md`)
	case "exit":
		println(`
"exit":
Exits the interpreter. Identical to "quit".
Expected parameters: none
`)
	case "quit":
		println(`
"quit":
Exits the interpreter. Identical to "exit".
Expected parameters: none
`)
	case "program", "Program":
		println(`
Programs:
A program consists of definitions and expressions. Each program gets an input value (called the argument) and outputs a Value.
Entering a top-level expression will cause the program to output it's value. Entering another top-level expression will cause the program to output a newline + the expression.
Only a program which outputs only one value may output a non-string value.
Tip: try "example program" to see an example.
`)
	case "toUpper":
		println(`
"toUpper":
Returns the input with all unicode letters mapped to their upper case.
Expected parameters: none
Tip: try "example toUpper" to see an example.
`)
	case "toLower":
		println(`
"toLower":
Returns the input with all unicode letters mapped to their lower case.
Expected parameters: none
Tip: try "example toLower" to see an example.
`)
	case "len":
		println(`
"len":
Returns the length of a given string.
Expected parameters: none
Tip: try "example len" to see an example.
`)
	case "count":
		println(`
"count":
Returns the number of values in a given list.
Expected parameters: none
Tip: Try "example count" to see an example.
`)
	case "split":
		println(`
"split":
TODO - explanation for "split"
Expected number of parameters: 1
Tip: try "example split" to see an example.
`)
	case "lines":
		println(`
"lines":
Splits a given string into lines.")
Expected parameters: none")
Tip: try "example lines" to see an example."
`)
	case "words":
		println(`
"words":
Splits a given string into words.
Expected parameters: none
Tip: try "example words" to see an example.
`)
	case "chars":
		println(`
"chars":
Splits a given string into a list of single characters.
Expected parameters: none
Tip: try "example chars" to see an example.
`)
	case "min":
		println(`
"min":
Finds the smallest value in a list based on a specified order.
Expected parameters: 1
 - the definition by which to order the values.

Tip: try "example min" to see an example.
`)
	case "max":
		println(`
"max":
Finds the largest value in a list based on a specified order.
Expected parameters: 1
 - the definition by which to order the values.

Tip: try "example max" to see an example.
`)
	}
}

func showExample(cmd string) {
	str := ""
	if len(cmd) >= 9 {
		str = cmd[8 : len(cmd)-1]
	}
	switch str {
	case "":
		println("Try \"example xxx\" to see an example for a particular subject.")
	case "example":
		println("\n--> example example")
		println("[Do you really need to see this?]\n")
	case "quit":
		println("\nquit\n[trex will exit]\n")
	case "exit":
		println("\nexit\n[trex will exit]\n")
	case "toLower":
		println("\n--> toLower 'Hello World'")
		println("hello world\n")
	case "toUpper":
		println("\n--> toUpper 'Hello World'")
		println("HELLO WORLD\n")
	case "min":
		println("\n--> []\nword\nanother\n12\nfoo")
		println("--> min(#len)\n12\n")
	default:
		println("No example exists for \"" + str + "\".")
	}
}
