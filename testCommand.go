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
		flags.reoptimize = !flags.reoptimize

		flags.clean = !flags.clean
	}

	runtime.GOMAXPROCS(flags.threads)

	if flags.clean {
		fmt.Print("CLEANING...")
		cleanResultDirs()
		fmt.Println(" done")
	}
	if flags.codegen {
		color.Cyan("GENERATING CODE...")
		color.Yellow("building...")
		batchCodeGen(flags.threads)
		color.Yellow("running...")
		batchRunUnoptimized(flags.threads)
	}
	if flags.optimize {
		color.Cyan("OPTIMIZING...")
		color.Yellow("building...")
		batchOptimize(flags.threads)
		color.Yellow("running...")
		batchRunOptimized(flags.threads)
		if flags.reoptimize {
			color.Yellow("reoptimizing...")
			batchReoptimizeOptimize(flags.threads)
		}
	}
	if flags.compile {
		color.Cyan("COMPILING...")
		color.Yellow("building...")
		batchCompile(flags.threads)
		color.Yellow("running...")
		batchRunCompiled(flags.threads)
		if flags.reoptimize {
			color.Yellow("reoptimizing...")
			batchReoptimizeCompile(flags.threads)
		}
	}
	if flags.optimizeStandalone {
		color.Cyan("OPTIMIZING STANDALONE ASM...")
		color.Yellow("building...")
		batchOptimizeStandalone(flags.threads)
		color.Yellow("running...")
		batchRunOptimizedStandalone(flags.threads)
		if flags.reoptimize {
			color.Yellow("reoptimizing...")
			batchOptimizeStandalone(flags.threads)
		}
	}
}
