package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

func initDirs() {
	mkdirIfNotExist(bin)

	mkdirIfNotExist(pika)
	mkdirIfNotExist(asm)

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
		if !strings.Contains(file.Name(), ext) {
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

func stripLines(input string, containing string) string {
	lines := strings.Split(input, "\n")

	result := ""
	for _, line := range lines {
		if !strings.Contains(line, containing) {
			result += line + "\n"
		}
	}
	return result
}

func helpMessage() {
	fmt.Println("usage of gtr: gtr <command> <flags>* <target>?")
	fmt.Println()
	fmt.Println("gtr commands:")
	fmt.Println("test:\t\trun tests")
	fmt.Println("create:\t\tcreate a new test, requires test name as <target>")
	fmt.Println("view:\t\tview a specified test's results, requires test name as <target>")
	fmt.Println("accept:\t\taccept the current output of a test in the future, may require test name as <target>")
	fmt.Println("init:\t\tbuild the directory structure needed to run gtr in this directory")
	fmt.Println()
	fmt.Println("see gtr <command> --help for details on that command's flags")
}
