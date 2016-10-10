package main

import (
	"os"

	"github.com/fatih/color"
)

func main() {

	if len(os.Args) == 1 {
		helpMessage()
		os.Exit(0)
	}

	command := os.Args[1]
	args := os.Args[2:]
	switch command {
	case "test":
		flags := makeTestFlags(args)
		testCommand(flags)
	case "view":
		flags, target := makeViewFlags(args)
		viewCommand(flags, target)
	case "create":
		flags, target := makeCreateFlags(args)
		createCommand(flags, target)
	case "accept":
		flags, target := makeAcceptFlags(args)
		acceptCommand(flags, target)
	case "init":
		initDirs()
	case "help":
		fallthrough
	case "-help":
		fallthrough
	case "--help":
		helpMessage()
		os.Exit(0)
	default:
		color.Magenta("invalid command: " + "\"" + command + "\"")
		helpMessage()
		os.Exit(1)
	}
	os.Exit(0)
}
