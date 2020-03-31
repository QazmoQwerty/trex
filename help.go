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
		println(`Use "help xxx" to see help for a particular subject.
Or read the Language Specification: gitlab.com/QazmoQwerty/trex/-/blob/master/docs/trex-spec.md`)
	case "example":
		println(`Use "help xxx" to see help for a particular subject.`)
	case "help":
		println(`Use "help xxx" to see help for a particular subject.`)
	case "exit":
		println(`
"exit":
Exits the interpreter. Identical to "quit".
Input: none
Parameters: none
`)
	case "quit":
		println(`
"quit":
Exits the interpreter. Identical to "exit".
Input: a list
Parameters: none
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
Input: a string
Parameters: none
Tip: try "example toUpper" to see an example.
`)
	case "toLower":
		println(`
"toLower":
Returns the input with all unicode letters mapped to their lower case.
Input: a string
Parameters: none
Tip: try "example toLower" to see an example.
`)
	case "len":
		println(`
"len":
Returns the length of a given string.
Input: a string
Parameters: none
Tip: try "example len" to see an example.
`)
	case "count":
		println(`
"count":
Returns the number of values in a given list.
Input: a list
Parameters: none
Tip: Try "example count" to see an example.
`)
	case "split":
		println(`
"split":
TODO - explanation for "split"
Input: a string
Parameters: 1
Tip: try "example split" to see an example.
`)
	case "lines":
		println(`
"lines":
Splits a given string into lines.
Input: a list
Parameters: none
Tip: try "example lines" to see an example."
`)
	case "words":
		println(`
"words":
Splits a given string into words.
Input: a string
Parameters: none
Tip: try "example words" to see an example.
`)
	case "chars":
		println(`
"chars":
Splits a given string into a list of single characters.
Input: a string
Parameters: none
Tip: try "example chars" to see an example.
`)
	case "min":
		println(`
"min":
Finds the smallest value in a list based on a specified order.
Input: a list
Parameters: 1
 - the definition by which to order the values, which must return a value convertible to a number.
Tip: try "example min" to see an example.
`)
	case "max":
		println(`
"max":
Finds the largest value in a list based on a specified order.
Input: a list
Parameters: 1
 - the definition by which to order the values, which must return a value convertible to a number..
Tip: try "example max" to see an example.
`)
	case "unique":
		println(`
Returns a list of all unique values in a given list.
Input: a list
Parameters: none
Tip: try "example unique" to see an example.
`)
	case "numOccurs":
		println(`
Returns the number of times a value occurs inside a given list or string.
Input: a list or string.
Parameters: 1
 - the value to count occurences of
Tip: try "example numOccurs" to see an example.
`)
	case "sort":
		println(`
Sorts a list based on a specified order.
Input: a list
Parameters: 1
 - the definition by which to order the values, which must return a value convertible to a number..
Tip: try "example sort" to see an example.

NOTE: currently unimplemented
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
		println(`Try "example xxx" to see an example for a particular subject.`)
	case "example":
		println(`
--> example example")
println("[Do you really need to see this?]
`)
	case "help":
		println(`
--> help example")
[help for for how to use the example command]
`)
	case "quit":
		println(`
--> quit
[trex will exit]
`)
	case "exit":
		println(`
--> exit
[trex will exit]
`)
	case "toLower":
		println(`
--> toLower "Hello World"
hello world
`)
	case "toUpper":
		println(`
--> toUpper "Hello World"
HELLO WORLD
`)
	case "min":
		println(`
--> []
word
another
foo
--> min(#len)
foo
`)
	case "max":
		println(`
--> []
word
another
foo
--> max(#len)
another
`)
	case "numOccurs":
		println(`
--> numOccurs('fo') 'foobafo'
2
`)
	case "chars":
		println(`
--> chars 12343
1, 2, 3, 4, 3
`)
	case "words":
		println(`
--> foo => "this is a sentence"
--> words foo
this, is, a, sentence
`)
	case "lines":
		println(`
--> []
one
two
three
--> lines
one, two, three
`)
	case "count":
		println(`
--> lines
one, two, three
--> count lines
3
`)
	case "len":
		println(`
--> len "example"
7
`)
	case "sort":
		println(`
--> words
one, three, four
--> sort(#len) words
one, four, three
`)
	default:
		println("No example exists for \"" + str + "\".")
	}
}
