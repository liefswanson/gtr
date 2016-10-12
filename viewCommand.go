package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
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
			os.Exit(1)
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
		output := makeDiff(expectPath, resultPath)
		printDiff(output)
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
