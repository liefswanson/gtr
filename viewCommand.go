package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"

	"github.com/fatih/color"
)

func viewCommand(flags viewFlags, testname string) {
	if flags.asm {
		phase := "asm"
		color.Cyan("ASM...")
		viewOutput(phase, flags.set, testname, flags.asm, flags.diff)
	}
	if flags.build {
		phase := "build"
		color.Cyan("BUILD...")
		viewOutput(phase, flags.set, testname, flags.asm, flags.diff)
	}
	if flags.run {
		phase := "run"
		color.Cyan("RUN...")
		viewOutput(phase, flags.set, testname, flags.asm, flags.diff)
	}
}

// TODO need to take into account asmo files
// TODO need to strip timestamps from outputs when writing to file
//      not when reading them in
func viewOutput(phase, set, testname string, asm bool, diff bool) {
	resultPath := buildPath(resultDir, phase, set, testname+txtExt)
	expectPath := buildPath(expectDir, phase, set, testname+txtExt)
	if asm {
		resultPath = buildPath(resultDir, phase, set, testname+asmExt)
		expectPath = buildPath(expectDir, phase, set, testname+asmExt)
	}

	if !exists(resultPath) {
		color.Magenta("there is no result set for " + testname)
		color.Magenta(resultPath + " does not exist")
		os.Exit(1)
	}

	if !exists(expectPath) {
		color.Magenta("there is no expectation set for " + testname)
		color.Magenta(expectPath + " does not exist")
		os.Exit(1)
	}

	if diff {
		color.Yellow("diff...")
		gitDiff := exec.Command("git", "diff", "--no-index", resultPath, expectPath)
		var stdOut bytes.Buffer
		gitDiff.Stdout = &stdOut
		gitDiff.Run()
		// TODO need to recolor the output of git diff
		fmt.Print(string(stdOut.Bytes()))
	} else {
		exp, err := ioutil.ReadFile(expectPath)
		crashOnError(err)
		res, err := ioutil.ReadFile(resultPath)
		crashOnError(err)
		color.Yellow("expect...")
		fmt.Print(string(exp))
		color.Yellow("result...")
		fmt.Print(string(res))
	}
}
