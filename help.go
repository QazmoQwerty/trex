package main

func showHelp(cmd string) {
	str := ""
	if len(cmd) >= 6 {
		str = cmd[5 : len(cmd)-1]
	}
	switch str {
	case "":
		println("Help still needs to be written")
		println("For now see the Language Specification: gitlab.com/QazmoQwerty/trex/-/blob/master/docs/trex-spec.md")
	case "exit":

		println("Exits the interpreter. Identical to \"quit\"")
	default:
		println("No help exists for name \"" + str + "\".")
	}
}

func showExample(cmd string) {
	str := ""
	if len(cmd) >= 9 {
		str = cmd[8 : len(cmd)-1]
	}
	switch str {
	case "":
		println("Enter \"example [name]\" to see an example for a particular definition.")
	case "example":
		println("\n--> example toUpper")
		println("[Example for the use of toUpper]\n")
	case "toUpper":
		println("\n--> toUpper 'Hello World'")
		println("HELLO WORLD\n")
	default:
		println("No example exists for name \"" + str + "\".")
	}
}
