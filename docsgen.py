import os

inputfile = "docs/docs.txt"
mdfile = "docs/builtin-defs.md"
gofile = "help.go"
txtfile = "docs/docs.txt"

items = []

class Item:
    def __init__(self, name, explanation, example):
        self.name = name
        self.explanation = explanation
        self.example = example
    def __str__(self):
        return "## " + self.name + "\n" + self.explanation + "\n" + self.example

with open(inputfile) as docs:
    for s in docs.read().split("## "):
        if s == "":
            continue
        item = Item("", "", "")
        stage = 0
        for line in s.splitlines():
            if line == "":
                continue
            if stage == 0:
                item.name = line
                stage += 1
            elif stage == 1:
                if line.startswith("-->"):
                    item.example += line
                    stage += 1
                else:
                    if item.explanation != "":
                        item.explanation += "\n"
                    item.explanation += line
            else:
                item.example += "\n" + line
        items += [item]

items.sort(key=lambda i: i.name)

txtstr = ""

for i in items:
    txtstr += "## " + i.name + "\n"
    for line in i.explanation.splitlines():
        txtstr += line + "\n"
    for line in i.example.splitlines():
        txtstr += line + "\n"
    txtstr += "\n"

with open(txtfile, 'w+') as out:
    out.write(txtstr)
print('outputted txt file "' + txtfile + '"')

mdstr = "# Trex Built-In Definitions\n\n## Table of Contents:\n\n"

for i in range(len(items)):
    mdstr += str(i+1) + '. [' + items[i].name + '](#' + items[i].name + ')\n'
mdstr += '\n'

for i in items:
    mdstr += "## " + i.name + "\n\n"
    for line in i.explanation.splitlines():
        mdstr += line + "\n\n"
    mdstr += "```\n"
    for line in i.example.splitlines():
        mdstr += line + "\n"
    mdstr += "```\n\n"

with open(mdfile, 'w+') as out:
    out.write(mdstr)
print('outputted md file "' + mdfile + '"!')

examplefunc = """func showExample(cmd string) {
	s := ""
	if len(cmd) >= 9 {
		s = cmd[8 : len(cmd)-1]
	}
	switch s {
	default:
		globals.outputColor.Println(`No example exists for "` + s + `".`)
	case "":
		globals.outputColor.Println(`Try "example xxx" to see an example for a particular subject.`)
	case "example":
		globals.outputColor.Println(`
--> example example
[Do you really need to see this?]
`)
	case "help":
		globals.outputColor.Println(`
--> help example
[help for for how to use the example command]
`)
	case "quit":
		globals.outputColor.Println(`
--> quit
[trex will exit]
`)
	case "exit":
		globals.outputColor.Println(`
--> exit
[trex will exit]
`)
"""
for i in items:
    examplefunc += '\tcase "' + i.name + '":\n\t\tglobals.outputColor.Println(`\n' + i.example + '\n`)\n'
examplefunc += "\t}\n}\n"

helpfunc = """func showHelp(cmd string) {
	s := ""
	if len(cmd) >= 6 {
		s = cmd[5 : len(cmd)-1]
	}
	switch s {
	default:
		globals.outputColor.Println(`No help exists for "` + s + `".`)
	case "":
		globals.outputColor.Println(`Use "help xxx" to see help for a particular subject.
Or read the Language Specification: gitlab.com/QazmoQwerty/trex/-/blob/master/docs/trex-spec.md`)
	case "example":
		globals.outputColor.Println(`Use "example xxx" to an example of a particular subject.`)
	case "help":
		globals.outputColor.Println(`Use "help xxx" to see help for a particular subject.`)
	case "exit":
		globals.outputColor.Println(`
"exit":
Exits the interpreter. Identical to "quit".
Input: none
Parameters: none
`)
	case "quit":
		globals.outputColor.Println(`
"quit":
Exits the interpreter. Identical to "exit".
Input: a list
Parameters: none
`)
"""
for i in items:
    helpfunc += '\tcase "' + i.name + '":\n\t\tglobals.outputColor.Println(`\n' + \
                '"' + i.name + '":\n' + i.explanation + \
                '\nTip: try "example ' + i.name + '" to see an example.' + '\n`)\n'
helpfunc += "\t}\n}\n"

with open(gofile, 'w+') as out:
    out.write("package main\n\n// generated by " + os.path.basename(__file__) + "\n\n" + helpfunc + "\n" + examplefunc)
print('outputted go file "' + gofile + '"!')