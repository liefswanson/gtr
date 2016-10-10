package main

import (
	"fmt"
	"runtime"

	"github.com/fatih/color"
)

func testCommand(flags testFlags) {
	if flags.invertFlags {
		flags.compile = !flags.compile
		flags.codegen = !flags.codegen
		flags.optimize = !flags.optimize
		flags.optimizeStandalone = !flags.optimizeStandalone

		flags.clean = !flags.clean
	}

	cores := runtime.NumCPU()
	runtime.GOMAXPROCS(cores + 1)

	if flags.clean {
		fmt.Print("CLEANING...")
		cleanDirs()
		fmt.Println(" done")
	}
	if flags.codegen {
		color.Cyan("GENERATING CODE...")
		color.Yellow("building...")
		batchCodeGen(cores)
		color.Yellow("running...")
		batchRunUnoptimized(cores)
	}
	if flags.optimize {
		color.Cyan("OPTIMIZING...")
		color.Yellow("building...")
		batchOptimize(cores)
		color.Yellow("running...")
		batchRunOptimized(cores)
	}
	if flags.compile {
		color.Cyan("COMPILING...")
		color.Yellow("building...")
		batchCompile(cores)
		color.Yellow("running...")
		batchRunCompiled(cores)
	}
	if flags.optimizeStandalone {
		color.Cyan("OPTIMIZING STANDALONE ASM...")
		color.Yellow("building...")
		batchOptimizeStandalone(cores)
		color.Yellow("running...")
		batchRunOptimizedStandalone(cores)
	}
}
