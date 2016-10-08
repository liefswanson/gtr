package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"sync"
)

type testResult struct {
	name   string
	result bool
}

type testStorageDir struct {
	codegenerator       string
	compiler            string
	optimizer           string
	optimizerStandalone string
}

type testDirTree struct {
	build testStorageDir
	run   testStorageDir
	asm   testStorageDir
}

const (
	bin  = "bin"
	pika = "src/pika"
	asm  = "src/asm"

	pikaExt = ".pika"
	asmExt  = ".asm"
	asmoExt = ".asmo"
	txtExt  = ".txt"

	// resExt = ".res"
	// expExt = ".exp"

	java = "java"
	wine = "wine"
)

var (
	compiler      = []string{"-jar", "-ea ", bin + "/pika-compiler.jar"}
	codegenerator = []string{"-jar", "-ea ", bin + "/pika-codegen.jar"}
	optimizer     = []string{"-jar", "-ea ", bin + "/pika-optimizer.jar"}
	emulator      = []string{bin + "/ASMEmu.exe"}

	result = testDirTree{
		build: testStorageDir{
			codegenerator:       "result/build/codegenerator",
			compiler:            "result/build/compiler",
			optimizer:           "result/build/optimizer",
			optimizerStandalone: "result/build/optimizer-standalone",
		},
		run: testStorageDir{
			codegenerator:       "result/run/codegenerator",
			compiler:            "result/run/compiler",
			optimizer:           "result/run/optimizer",
			optimizerStandalone: "result/run/optimizer-standalone",
		},
		asm: testStorageDir{
			codegenerator:       "result/asm/codegenerator",
			compiler:            "result/asm/compiler",
			optimizer:           "result/asm/optimizer",
			optimizerStandalone: "result/asm/optimizer-standalone",
		},
	}

	expect = testDirTree{
		build: testStorageDir{
			codegenerator:       "expect/build/codegenerator",
			compiler:            "expect/build/compiler",
			optimizer:           "expect/build/optimizer",
			optimizerStandalone: "expect/build/optimizer-standalone",
		},
		run: testStorageDir{
			codegenerator:       "expect/run/codegenerator",
			compiler:            "expect/run/compiler",
			optimizer:           "expect/run/optimizer",
			optimizerStandalone: "expect/run/optimizer-standalone",
		},
		asm: testStorageDir{
			codegenerator:       "expect/asm/codegenerator",
			compiler:            "expect/asm/compiler",
			optimizer:           "expect/asm/optimizer",
			optimizerStandalone: "expect/asm/optimizer-standalone",
		},
	}
)

func main() {
	var (
		compile            bool
		codegen            bool
		optimize           bool
		optimizeStandalone bool
		//reoptimize         bool

		clean      bool
		initialize bool

		invertFlags bool
	)

	flag.BoolVar(&initialize, "init", false,
		"make the directory structure required for running tests")
	flag.BoolVar(&clean, "clean", false,
		"clean out the output directories before running tests")

	flag.BoolVar(&codegen, "codegen", false,
		"generate asm but DON'T optimize")
	flag.BoolVar(&compile, "compile", false,
		"generate asm and optimize")
	flag.BoolVar(&optimize, "optimize", false,
		"optimize asm from --codegen.\n"+
			"Different to --compile:\n"+
			"Code is read back in, after being written to a file")
	flag.BoolVar(&optimizeStandalone, "optimizeStandalone", false,
		"optimize asm written explicitly for testing")
	// flag.BoolVar(&reoptimize, "reoptimize", false,
	// 	"run all optimized code through the optimizer again")

	flag.BoolVar(&invertFlags, "invertFlags", false,
		"inverts all flags, making them subtractive instead of additive")
	flag.Parse()

	cores := runtime.NumCPU()
	runtime.GOMAXPROCS(cores + 1)

	if invertFlags {
		compile = !compile
		codegen = !codegen
		optimize = !optimize
		optimizeStandalone = !optimizeStandalone
		// reoptimize = !reoptimize

		clean = !clean
		initialize = !initialize
	}

	if initialize {
		fmt.Print("INITIALIZING... ")
		initDirs()
		fmt.Println(" done")
	}
	if clean {
		fmt.Print("CLEANING... ")
		cleanDirs()
		fmt.Println(" done")
	}
	if codegen {
		fmt.Println("GENERATING CODE...")
		batchCodeGen(cores)
		fmt.Println("RUNNING...")
		batchRunUnoptimized(cores)
	}
	if optimize {
		fmt.Println("OPTIMIZING...")
		batchOptimize(cores)
		fmt.Println("RUNNING...")
		batchRunOptimized(cores)
	}
	if compile {
		fmt.Println("GENERATING CODE + OPTIMIZING...")
		batchCompile(cores)
		fmt.Println("RUNNING...")
		batchRunCompiled(cores)
	}
	if optimizeStandalone {
		fmt.Println("OPTIMIZING HAND WRITTEN ASM...")
		batchOptimizeStandalone(cores)
		fmt.Println("RUNNING...")
		batchRunOptimizedStandalone(cores)
	}

	os.Exit(0)
}

func initDirs() {
	mkdirIfNotExist("bin")

	mkdirIfNotExist("src/asm")
	mkdirIfNotExist("src/pika")

	// expect
	mkdirIfNotExist(expect.asm.codegenerator)
	mkdirIfNotExist(expect.asm.compiler)
	mkdirIfNotExist(expect.asm.optimizer)
	mkdirIfNotExist(expect.asm.optimizerStandalone)

	mkdirIfNotExist(expect.build.codegenerator)
	mkdirIfNotExist(expect.build.compiler)
	mkdirIfNotExist(expect.build.optimizer)
	mkdirIfNotExist(expect.build.optimizerStandalone)

	mkdirIfNotExist(expect.run.codegenerator)
	mkdirIfNotExist(expect.run.compiler)
	mkdirIfNotExist(expect.run.optimizer)
	mkdirIfNotExist(expect.run.optimizerStandalone)

	// result
	mkdirIfNotExist(result.asm.codegenerator)
	mkdirIfNotExist(result.asm.compiler)
	mkdirIfNotExist(result.asm.optimizer)
	mkdirIfNotExist(result.asm.optimizerStandalone)

	mkdirIfNotExist(result.build.codegenerator)
	mkdirIfNotExist(result.build.compiler)
	mkdirIfNotExist(result.build.optimizer)
	mkdirIfNotExist(result.build.optimizerStandalone)

	mkdirIfNotExist(result.run.codegenerator)
	mkdirIfNotExist(result.run.compiler)
	mkdirIfNotExist(result.run.optimizer)
	mkdirIfNotExist(result.run.optimizerStandalone)
}

func cleanDirs() {
	cleanDir(result.asm.codegenerator)
	cleanDir(result.asm.compiler)
	cleanDir(result.asm.optimizer)
	cleanDir(result.asm.optimizerStandalone)

	cleanDir(result.build.codegenerator)
	cleanDir(result.build.compiler)
	cleanDir(result.build.optimizer)
	cleanDir(result.build.optimizerStandalone)

	cleanDir(result.run.codegenerator)
	cleanDir(result.run.compiler)
	cleanDir(result.run.optimizer)
	cleanDir(result.run.optimizerStandalone)
}

func cleanDir(dir string) {
	files := getAllFiles(dir)
	files = filterOutFiles(files, ".gitignore")
	for _, file := range files {
		if strings.Compare(file.Name(), ".gitignore") != 0 {
			err := os.Remove(buildPath(dir, file.Name()))
			crashOnError(err)
		}
	}
}

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
func measureSlice(size int, total int, current int) (int, int) {
	start := current * size / total
	end := (current + 1) * size / total
	if end > size {
		end = size
	}
	return start, end
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
	task.Run()

	return append(stdout.Bytes(), stderr.Bytes()...)
}

func batchCodeGen(count int) {
	executeAll(count,
		pika, pikaExt,
		result.build.codegenerator, txtExt,
		result.asm.codegenerator,
		java, codegenerator)

	compareAllResults(count,
		result.build.codegenerator,
		expect.build.codegenerator,
		pika)
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

func compareAllResults(count int,
	resultDir string, expectDir string, refDir string) {

	temp := getAllFiles(refDir)
	testFiles := filterOutFiles(temp, ".gitignore")

	results := make(chan testResult)

	for i := 0; i < count; i++ {
		start := i * len(testFiles) / count
		end := (i + 1) * len(testFiles) / count
		slice := testFiles[start:end]
		go compareEachResult(slice, results, resultDir, expectDir)
	}

	passed := 0
	failed := ""
	for _ = range testFiles {
		test := <-results
		if test.result {
			passed++
		} else {
			failed += test.name + "\n"
		}
	}
	fmt.Println("passed: [", passed, "/", len(testFiles), "]")
	if failed != "" {
		fmt.Print("failed:\n", failed)
	}

	// TODO
	// log failures to file
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
func compareResult(resultFilePath string, expectFilePath string) bool {

	if exists(expectFilePath) && exists(resultFilePath) {
		expectedRaw, err := ioutil.ReadFile(expectFilePath)
		crashOnError(err)
		expected := string(expectedRaw[:])

		resultRaw, err := ioutil.ReadFile(resultFilePath)
		crashOnError(err)
		result := string(resultRaw[:])

		expectedLines := strings.Split(expected, "\n")
		resultLines := strings.Split(result, "\n")
		expected, result = "", ""

		for _, line := range expectedLines {
			if !strings.Contains(line, "logging.PikaLogger log") {
				expected += line
			}
		}

		for _, line := range resultLines {
			if !strings.Contains(line, "logging.PikaLogger log") {
				result += line
			}
		}

		return strings.Compare(result, expected) == 0
	} else if !exists(expectFilePath) && !exists(resultFilePath) {
		return true
	}
	return false
}

func mkdirIfNotExist(path string) {
	if !exists(path) {
		os.MkdirAll(path, 0777)
	}
}
func exists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func getAllFiles(dir string) []os.FileInfo {
	files, err := ioutil.ReadDir(dir)
	crashOnError(err)
	return files
}

func buildPath(parts ...string) string {
	if len(parts) == 0 {
		return "./"
	}
	return strings.Join(parts, "/")
}

func replaceExtension(filename string, ext string) string {
	prefix := strings.Split(filename, ".")[0]
	return prefix + ext
}

func crashOnError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func filterOutFiles(in []os.FileInfo, ext string) []os.FileInfo {
	files := make([]os.FileInfo, 0, len(in))
	for _, file := range in {
		if strings.Contains(file.Name(), ext) {
			files = append(files, file)
		}
	}
	return files
}

func filterFiles(in []os.FileInfo, ext string) []os.FileInfo {
	files := make([]os.FileInfo, 0, len(in))
	for _, file := range in {
		if strings.Contains(file.Name(), ext) {
			files = append(files, file)
		}
	}
	return files
}
