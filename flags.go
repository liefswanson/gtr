package main

import (
	"flag"
	"os"

	"github.com/fatih/color"
)

////////////////////////////////////////////////////////////////////////////////
// flags
// to figure out which flags do what, run the program with the -help flag
type testFlags struct {
	compile            bool
	codegen            bool
	optimize           bool
	optimizeStandalone bool

	clean bool

	invertFlags bool
}

type viewFlags struct {
	run   bool
	asm   bool
	build bool

	diff bool
}

type createFlags struct {
	open  bool
	isAsm bool
}

type acceptFlags struct {
	asm bool

	all bool
}

////////////////////////////////////////////////////////////////////////////////
// parsers
func makeTestFlags(args []string) testFlags {
	flags := testFlags{}
	test := flag.NewFlagSet("test", flag.ExitOnError)
	test.BoolVar(&flags.clean, "clean", false,
		"Clean out the output directories before running tests")

	test.BoolVar(&flags.codegen, "codegen", false,
		"Generate asm but DON'T optimize")
	test.BoolVar(&flags.compile, "compile", false,
		"Generate asm and optimize")
	test.BoolVar(&flags.optimize, "optimize", false,
		"Optimize asm from -codegen.\n"+
			"\tDifferent to -compile:\n"+
			"\tCode is read back in, after being written to a file")
	test.BoolVar(&flags.optimizeStandalone, "optimizeStandalone", false,
		"Optimize asm written explicitly for testing")

	test.BoolVar(&flags.invertFlags, "invert", false,
		"Inverts all flags, making them subtractive instead of additive")

	test.Parse(args)
	return flags
}

func makeViewFlags(args []string) (viewFlags, string) {
	flags := viewFlags{}
	view := flag.NewFlagSet("view", flag.ExitOnError)
	view.BoolVar(&flags.diff, "diff", false,
		"view as a diff, instead of result and expectation separately")
	view.BoolVar(&flags.run, "run", false,
		"compare results of the run phase of testing")
	view.BoolVar(&flags.asm, "asm", false,
		"compare asm generated by building phase of testing")
	view.BoolVar(&flags.build, "build", false,
		"compare results of the build phase of testing")

	view.Parse(args)
	if len(view.Args()) == 0 {
		color.Magenta("No test was specified to view")
		os.Exit(1)
	}
	return flags, view.Arg(0)
}

func makeCreateFlags(args []string) (createFlags, string) {
	flags := createFlags{}
	create := flag.NewFlagSet("create", flag.ExitOnError)
	create.BoolVar(&flags.open, "open", false,
		"opens the test which was just created in your default editor")
	create.BoolVar(&flags.isAsm, "asm", false,
		"test specified is a .asm file, default to .pika otherwise")

	create.Parse(args)
	if len(create.Args()) == 0 {
		color.Magenta("No test was specified to create")
		os.Exit(1)
	}
	return flags, create.Arg(0)
}

func makeAcceptFlags(args []string) (acceptFlags, string) {
	flags := acceptFlags{}
	accept := flag.NewFlagSet("accept", flag.ExitOnError)

	accept.BoolVar(&flags.asm, "asm", false,
		"test specified is a .asm file, default to .pika otherwise")

	accept.BoolVar(&flags.all, "all", false,
		"move result folder to expect\n"+
			"\tsupersedes all other flags"+
			"\tsaves old expects to .expect folder as a backup"+
			"\tanything previously in .expect gets deleted")

	accept.Parse(args)
	if len(accept.Args()) == 0 && flags.all == false {
		color.Magenta("No test was specified to accept")
		os.Exit(1)
	}
	if len(accept.Args()) == 0 {
		return flags, ""
	}
	return flags, accept.Arg(0)
}
