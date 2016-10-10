package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"sync"

	"github.com/fatih/color"
)

////////////////////////////////////////////////////////////////////////////////
// convenience commands
func batchCodeGen(count int) {
	executeAll(count,
		pika, pikaExt,
		result.build.codegenerator, txtExt,
		result.asm.codegenerator,
		java, codegenerator)

	compareAllResults(count,
		result.build.codegenerator,
		expect.build.codegenerator,
		result.build.codegenerator)
}

func batchRunUnoptimized(count int) {
	executeAll(count,
		result.asm.codegenerator, asmExt,
		result.run.codegenerator, txtExt,
		"",
		wine, emulator)

	compareAllResults(count,
		result.run.codegenerator,
		expect.run.codegenerator,
		result.asm.codegenerator)
}

func batchOptimize(count int) {
	executeAll(count,
		result.asm.codegenerator, asmExt,
		result.build.optimizer, txtExt,
		result.asm.optimizer,
		java, optimizer)

	compareAllResults(count,
		result.build.optimizer,
		expect.build.optimizer,
		result.asm.codegenerator)
}

func batchRunOptimized(count int) {
	executeAll(count,
		result.asm.optimizer, asmoExt,
		result.run.optimizer, txtExt,
		"",
		wine, emulator)

	compareAllResults(count,
		result.run.optimizer,
		expect.run.optimizer,
		result.asm.optimizer)
}

func batchCompile(count int) {
	executeAll(count,
		pika, pikaExt,
		result.build.compiler, txtExt,
		result.asm.compiler,
		java, compiler)

	compareAllResults(count,
		result.build.compiler,
		expect.build.compiler,
		result.asm.compiler)
}

func batchRunCompiled(count int) {
	executeAll(count,
		result.asm.compiler, asmExt,
		result.run.compiler, txtExt,
		"",
		wine, emulator)

	compareAllResults(count,
		result.run.compiler,
		expect.run.compiler,
		result.asm.compiler)
}

func batchOptimizeStandalone(count int) {
	executeAll(count,
		asm, asmExt,
		result.build.optimizerStandalone, txtExt,
		result.asm.optimizerStandalone,
		java, optimizer)

	compareAllResults(count,
		result.build.optimizerStandalone,
		expect.build.optimizerStandalone,
		result.asm.optimizerStandalone)
}

func batchRunOptimizedStandalone(count int) {
	executeAll(count,
		result.asm.optimizerStandalone, asmoExt,
		result.run.optimizerStandalone, txtExt,
		"",
		wine, emulator)

	compareAllResults(count,
		result.run.optimizerStandalone,
		expect.run.optimizerStandalone,
		result.asm.optimizerStandalone)
}

////////////////////////////////////////////////////////////////////////////////
// execution
func executeAll(count int,
	inDir string, inExt string,
	outDir string, outExt string,
	targetDir string,
	cmd string, args []string) {
	files := getAllFiles(inDir)
	files = filterFiles(files, inExt)
	var wg sync.WaitGroup
	for i := 0; i <= count; i++ {
		start, end := measureSlice(len(files), count, i)
		filesSlice := files[start:end]
		go executeEach(filesSlice,
			inDir, inExt,
			outDir, outExt,
			targetDir,
			&wg, cmd, args)
		wg.Add(1)
	}
	wg.Wait()
}

func executeEach(files []os.FileInfo,
	inDir string, inExt string,
	outDir string, outExt string,
	targetDir string,
	wg *sync.WaitGroup,
	cmd string, args []string) {

	for _, file := range files {
		srcPath := buildPath(inDir, file.Name())
		targetPath := targetDir + "/"
		completeArgs := append(args, srcPath)
		if targetPath != "" {
			completeArgs = append(completeArgs, targetPath)
		}

		toWrite := execute(cmd, completeArgs)

		outputFilename := replaceExtension(file.Name(), outExt)
		outputFilename = buildPath(outDir, outputFilename)

		ioutil.WriteFile(outputFilename, toWrite, 0777)
	}
	wg.Done()
}

func execute(cmd string, args []string) []byte {
	task := exec.Command(cmd, args...)
	var stdout, stderr bytes.Buffer
	task.Stdout, task.Stderr = &stdout, &stderr
	if cmd == wine {
		task.Stderr = nil
	}
	task.Run()

	return append(stdout.Bytes(), stderr.Bytes()...)
}

////////////////////////////////////////////////////////////////////////////////
// comparison
func compareAllResults(count int,
	resultDir string, expectDir string, refDir string) {

	temp := getAllFiles(refDir)
	testFiles := filterOutFiles(temp, ".gitignore")

	results := make(chan testResult)

	for i := 0; i < count; i++ {
		start, end := measureSlice(len(testFiles), count, i)
		slice := testFiles[start:end]
		go compareEachResult(slice, results, resultDir, expectDir)
	}

	passed := 0
	failed := make([]string, 0, len(testFiles))
	for _ = range testFiles {
		test := <-results
		if test.result {
			passed++
		} else {
			failed = append(failed, test.name)
		}
	}
	color.Green("passed: [", passed, "/", len(testFiles), "]")
	if len(failed) != 0 {
		color.Set(color.FgRed)
		fmt.Print("failed:\n")
	}
	for _, name := range failed {
		fmt.Println(name)
	}
	color.Unset()
	// return passed, len(testFiles), failed
}
func compareEachResult(files []os.FileInfo, results chan testResult,
	resultDir string, expectDir string) {
	for _, file := range files {
		resultFileName := replaceExtension(file.Name(), txtExt)
		resultFilePath := buildPath(resultDir, resultFileName)

		expectFileName := replaceExtension(file.Name(), txtExt)
		expectFilePath := buildPath(expectDir, expectFileName)

		results <- testResult{
			name:   replaceExtension(file.Name(), ""),
			result: compareResult(resultFilePath, expectFilePath)}
	}
}

// the stripping of the logging lines might need to happen when writing
// rather than during reading as it does here
// it will allow the use of diff in that case
func compareResult(resultFilePath string, expectFilePath string) bool {

	if exists(expectFilePath) && exists(resultFilePath) {
		expectedRaw, err := ioutil.ReadFile(expectFilePath)
		crashOnError(err)
		expected := string(expectedRaw[:])

		resultRaw, err := ioutil.ReadFile(resultFilePath)
		crashOnError(err)
		result := string(resultRaw[:])

		expected = stripLines(expected, loggingMessage)
		result = stripLines(result, loggingMessage)

		return strings.Compare(result, expected) == 0
	} else if !exists(expectFilePath) && !exists(resultFilePath) {
		return true
	}
	return false
}

func measureSlice(size int, total int, current int) (int, int) {
	start := current * size / total
	end := (current + 1) * size / total
	if end > size {
		end = size
	}
	return start, end
}
