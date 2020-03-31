package main

func showHelp(cmd string) {
	str := ""
	if len(cmd) >= 6 {
		str = cmd[5 : len(cmd)-1]
	}
	switch str {
	case "":
		println("Use \"help xxx\" to see help for a particular subject.")
		println("Or read the Language Specification: gitlab.com/QazmoQwerty/trex/-/blob/master/docs/trex-spec.md")
	case "exit":
		println("Exits the interpreter. Identical to \"quit\".")
	case "quit":
		println("Exits the interpreter. Identical to \"exit\".")
	case "toUpper":
		println("\n\"toUpper\":\n")
		println("Returns the input with all unicode letters mapped to their upper case.\n")
		println("Use \"example toUpper\" to see an example.\n")
	default:
		println("No help exists for \"" + str + "\".")
	}
}

func showExample(cmd string) {
	str := ""
	if len(cmd) >= 9 {
		str = cmd[8 : len(cmd)-1]
	}
	switch str {
	case "":
		println("Use \"example xxx\" to see an example for a particular subject.")
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
	default:
		println("No example exists for \"" + str + "\".")
	}
}
