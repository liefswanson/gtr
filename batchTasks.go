package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"sort"
	"sync"

	"github.com/fatih/color"
)

////////////////////////////////////////////////////////////////////////////////
// convenience commands
func batchCodeGen(count int) {
	executeAll(count,
		pikaDir, pikaExt,
		result.build.codegenerator, txtExt,
		result.asm.codegenerator,
		java, codegeneratorArgs)

	compareAllResults(count,
		result.build.codegenerator, txtExt,
		expect.build.codegenerator, txtExt,
		pikaDir)
}

func batchRunUnoptimized(count int) {
	executeAll(count,
		result.asm.codegenerator, asmExt,
		result.run.codegenerator, txtExt,
		"",
		wine, emulatorArgs)

	compareAllResults(count,
		result.run.codegenerator, txtExt,
		expect.run.codegenerator, txtExt,
		result.asm.codegenerator)
}

func batchOptimize(count int) {
	executeAll(count,
		result.asm.codegenerator, asmExt,
		result.build.optimizer, txtExt,
		result.asm.optimizer,
		java, optimizerArgs)

	compareAllResults(count,
		result.build.optimizer, txtExt,
		expect.build.optimizer, txtExt,
		result.asm.codegenerator)
}

// TODO reoptimize
func batchReoptimizeOptimize(count int) {
	executeAll(count,
		result.asm.optimizer, asmoExt,
		result.buildo.optimizer, txtExt,
		result.asmo.optimizer,
		java, optimizerArgs)

	compareAllResults(count,
		result.asmo.optimizer, asmoExt,
		expect.asmo.optimizer, asmoExt,
		result.asm.optimizer)
}

func batchRunOptimized(count int) {
	executeAll(count,
		result.asm.optimizer, asmoExt,
		result.run.optimizer, txtExt,
		"",
		wine, emulatorArgs)

	compareAllResults(count,
		result.run.optimizer, txtExt,
		expect.run.optimizer, txtExt,
		result.asm.optimizer)
}

func batchCompile(count int) {
	executeAll(count,
		pikaDir, pikaExt,
		result.build.compiler, txtExt,
		result.asm.compiler,
		java, compilerArgs)

	compareAllResults(count,
		result.build.compiler, txtExt,
		expect.build.compiler, txtExt,
		pikaDir)
}

// TODO
func batchReoptimizeCompile(count int) {
	executeAll(count,
		result.asm.compiler, asmExt,
		result.buildo.compiler, txtExt,
		result.asmo.compiler,
		java, optimizerArgs)

	compareAllResults(count,
		result.asmo.compiler, asmoExt,
		expect.asmo.compiler, asmoExt,
		result.asm.compiler)
}

func batchRunCompiled(count int) {
	executeAll(count,
		result.asm.compiler, asmExt,
		result.run.compiler, txtExt,
		"",
		wine, emulatorArgs)

	compareAllResults(count,
		result.run.compiler, txtExt,
		expect.run.compiler, txtExt,
		result.asm.compiler)
}

// TODO run once in the beginning to check that you have the same result
// TODO add several directories for extra run phases, that check against results

func batchOptimizeStandalone(count int) {
	executeAll(count,
		asmDir, asmExt,
		result.build.optimizerStandalone, txtExt,
		result.asm.optimizerStandalone,
		java, optimizerArgs)

	compareAllResults(count,
		result.build.optimizerStandalone, txtExt,
		expect.build.optimizerStandalone, txtExt,
		asmDir)
}

// TODO
func batchReoptimizeOptimizeStandalone(count int) {
	executeAll(count,
		result.asm.optimizerStandalone, asmoExt,
		result.buildo.optimizerStandalone, txtExt,
		result.asmo.optimizerStandalone,
		java, optimizerArgs)

	compareAllResults(count,
		result.asmo.optimizerStandalone, asmoExt,
		expect.asm.optimizerStandalone, asmExt,
		result.asm.optimizerStandalone)
}

func batchRunOptimizedStandalone(count int) {
	executeAll(count,
		result.asm.optimizerStandalone, asmoExt,
		result.run.optimizerStandalone, txtExt,
		"",
		wine, emulatorArgs)

	compareAllResults(count,
		result.run.optimizerStandalone, txtExt,
		expect.run.optimizerStandalone, txtExt,
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
	wg.Add(count)
	for i := 0; i < count; i++ {
		start, end := measureSlice(len(files), count, i)
		filesSlice := files[start:end]
		go executeEach(filesSlice,
			inDir, inExt,
			outDir, outExt,
			targetDir,
			&wg, cmd, args)
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

		bytesToWrite := execute(cmd, completeArgs)

		outputFilename := replaceExtension(file.Name(), outExt)
		outputFilename = buildPath(outDir, outputFilename)

		toWrite := string(bytesToWrite)
		toWrite = stripLines(toWrite, loggingMessage)

		bytesToWrite = []byte(toWrite)
		ioutil.WriteFile(outputFilename, bytesToWrite, 0777)
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
	resultDir, resExt,
	expectDir, expExt,
	refDir string) {

	testFiles := getAllFiles(refDir)
	testFiles = filterOutFiles(testFiles, ".gitignore")
	testFiles = filterOutFiles(testFiles, ".directory")

	results := make(chan testResult)

	for i := 0; i < count; i++ {
		start, end := measureSlice(len(testFiles), count, i)
		slice := testFiles[start:end]
		go compareEachResult(slice, results, resultDir, resExt, expectDir, expExt)
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
	sort.Strings(failed)

	green := color.New(color.FgGreen)
	total := len(testFiles)
	green.Println("passed: [", passed, "/", total, "]")

	if len(failed) == 0 {
		return
	}

	color.Set(color.FgRed)
	fmt.Println("failed:")
	for _, testName := range failed {
		fmt.Println(testName)
	}
	color.Unset()
}

func compareEachResult(files []os.FileInfo, results chan testResult,
	resultDir, resExt,
	expectDir, expExt string) {
	for _, file := range files {
		resultFileName := replaceExtension(file.Name(), resExt)
		resultFilePath := buildPath(resultDir, resultFileName)

		expectFileName := replaceExtension(file.Name(), resExt)
		expectFilePath := buildPath(expectDir, expectFileName)

		results <- testResult{
			name:   replaceExtension(file.Name(), ""),
			result: compareResult(resultFilePath, expectFilePath)}
	}
}

func compareResult(resultFilePath string, expectFilePath string) bool {

	if exists(expectFilePath) && exists(resultFilePath) {
		expectRaw, err := ioutil.ReadFile(expectFilePath)
		crashOnError(err)
		expect := string(expectRaw)

		resultRaw, err := ioutil.ReadFile(resultFilePath)
		crashOnError(err)
		result := string(resultRaw)

		return result == expect
	} else if !exists(expectFilePath) && !exists(resultFilePath) {
		return true
	}
	return false
}

func measureSlice(length int, slices int, sliceNum int) (int, int) {
	start := sliceNum * length / slices
	end := (sliceNum + 1) * length / slices
	if end > length {
		end = length
	}
	return start, end
}
