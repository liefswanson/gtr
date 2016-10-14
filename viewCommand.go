package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os/exec"
	"strings"

	"github.com/fatih/color"
)

func viewCommand(flags viewFlags, testname string) {
	if flags.asm {
		phase := asm
		color.Cyan("ASM...")
		viewOutput(phase, flags.testSet, testname, flags.asm, flags.diff)
	}
	if flags.build {
		phase := build
		color.Cyan("BUILD...")
		viewOutput(phase, flags.testSet, testname, flags.asm, flags.diff)
	}
	if flags.run {
		phase := run
		color.Cyan("RUN...")
		viewOutput(phase, flags.testSet, testname, flags.asm, flags.diff)
	}
	if flags.asmo {
		if flags.testSet == "codegenerator" {
			color.Magenta("there is no reoptimize phase for the codegenerator")
			return
		}
		phase := asmo
		color.Cyan("ASMO...")
		viewOutput(phase, flags.testSet, testname, flags.asm, flags.diff)
	}
}

// TODO need to take into account asmo files
func viewOutput(phase, testSet, testname string, isAsm bool, diff bool) {
	var resultPath, expectPath string
	if phase == asm || phase == asmo {
		if testSet == "optimizer" || testSet == "optimizer-standalone" ||
			phase == asmo {
			resultPath = buildPath(resultDir, phase, testSet, testname+asmoExt)
			expectPath = buildPath(expectDir, phase, testSet, testname+asmoExt)
		} else {
			resultPath = buildPath(resultDir, phase, testSet, testname+asmExt)
			expectPath = buildPath(expectDir, phase, testSet, testname+asmExt)
		}
	} else {
		resultPath = buildPath(resultDir, phase, testSet, testname+txtExt)
		expectPath = buildPath(expectDir, phase, testSet, testname+txtExt)
	}

	if diff {
		color.Yellow("diff...")
		if exists(resultPath) && exists(expectPath) {
			output := makeDiff(expectPath, resultPath)
			printDiff(output)
		}

		if !exists(expectPath) {
			color.Magenta("there is no expectation set for " + testname)
			color.Magenta(expectPath + " does not exist")
		}

		if !exists(resultPath) {
			color.Magenta("there is no result set for " + testname)
			color.Magenta(resultPath + " does not exist")
		}

	} else {
		color.Yellow("expect...")
		if !exists(expectPath) {
			color.Magenta("there is no expectation set for " + testname)
			color.Magenta(expectPath + " does not exist")
		} else {
			exp, err := ioutil.ReadFile(expectPath)
			crashOnError(err)

			fmt.Print(string(exp))
		}

		color.Yellow("result...")
		if !exists(resultPath) {
			color.Magenta("there is no result set for " + testname)
			color.Magenta(resultPath + " does not exist")
		} else {
			res, err := ioutil.ReadFile(resultPath)
			crashOnError(err)

			fmt.Print(string(res))

		}

	}
}

func makeDiff(expectPath, resultPath string) string {
	gitDiff := exec.Command("git", "diff", "--no-index", expectPath, resultPath)
	var stdOut bytes.Buffer
	gitDiff.Stdout = &stdOut
	gitDiff.Run()
	return string(stdOut.Bytes())
}

func printDiff(diff string) {
	lines := strings.Split(diff, "\n")
	if len(lines) < 4 {
		return
	}
	for _, line := range lines[2:] {
		if len(line) == 0 {
			continue
		}
		r := line[0]
		lightBlue := color.New(color.FgHiBlue)
		switch r {
		case '+':
			color.Green(line)
		case '-':
			color.Red(line)
		case '@':
			lightBlue.Println(line)
		default:
			fmt.Println(line)
		}
	}
}
