package main

import (
	"os"
	"os/exec"

	"github.com/fatih/color"
)

func acceptCommand(flags acceptFlags, testname string) {
	if flags.all {
		acceptAll()
		return
	}

	var path string
	if flags.asm {
		filename := testname + asmExt
		path = buildPath(asmDir, filename)
	} else {
		filename := testname + pikaExt
		path = buildPath(pikaDir, filename)
	}

	if !exists(path) {
		color.Magenta(path + " does not exist")
		os.Exit(1)
		return // not necessary, just to be explicit
	}

	if flags.asm {
		acceptOptimizedStandalone(testname)
	} else {
		acceptCodeGen(testname)
		acceptOptimized(testname)
		acceptCompile(testname)
	}
}

func acceptAll() {
	os.RemoveAll(backupDir)
	// os.Rename(expectDir, backupDir)
	// os.Rename(resultDir, expectDir)
	// initResultDirs()
	backup()
	acceptAllPika()
	acceptAllStandalone()
}

func acceptAllPika() {
	files := getAllFiles(pikaDir)
	files = filterOutFiles(files, ".gitignore")

	for _, file := range files {
		testname := replaceExtension(file.Name(), "")
		acceptCodeGen(testname)
		acceptOptimized(testname)
		acceptCompile(testname)
	}
}

func acceptAllStandalone() {
	files := getAllFiles(asmDir)
	files = filterOutFiles(files, ".gitignore")

	for _, file := range files {
		testname := replaceExtension(file.Name(), "")
		acceptOptimizedStandalone(testname)
	}
}

func backup() {
	exec.Command("cp", "-rf", expectDir, backupDir).Run()
}

func acceptTest(testname, oldPath, newPath, ext string) {
	old := buildPath(oldPath, testname+ext)
	new := buildPath(newPath, testname+ext)
	moveIfExists(old, new)
}

func acceptCompile(testname string) {
	acceptTest(testname, result.build.compiler, expect.build.compiler, txtExt)
	acceptTest(testname, result.run.compiler, expect.run.compiler, txtExt)
	acceptTest(testname, result.asm.compiler, expect.asm.compiler, asmExt)

	acceptTest(testname, result.buildo.compiler, expect.buildo.compiler, txtExt)
	acceptTest(testname, result.asmo.compiler, expect.asmo.compiler, asmoExt)
}

func acceptCodeGen(testname string) {
	acceptTest(testname, result.build.codegenerator, expect.build.codegenerator, txtExt)
	acceptTest(testname, result.run.codegenerator, expect.run.codegenerator, txtExt)
	acceptTest(testname, result.asm.codegenerator, expect.asm.codegenerator, asmExt)
}

func acceptOptimized(testname string) {
	acceptTest(testname, result.build.optimizer, expect.build.optimizer, txtExt)
	acceptTest(testname, result.run.optimizer, expect.run.optimizer, txtExt)
	acceptTest(testname, result.asm.optimizer, expect.asm.optimizer, asmoExt)

	acceptTest(testname, result.buildo.optimizer, expect.buildo.optimizer, txtExt)
	acceptTest(testname, result.asmo.optimizer, expect.asmo.optimizer, asmoExt)
}

func acceptOptimizedStandalone(testname string) {
	acceptTest(testname, result.build.optimizerStandalone, expect.build.optimizerStandalone, txtExt)
	acceptTest(testname, result.run.optimizerStandalone, expect.run.optimizerStandalone, txtExt)
	acceptTest(testname, result.asm.optimizerStandalone, expect.asm.optimizerStandalone, asmoExt)

	acceptTest(testname, result.buildo.optimizerStandalone, expect.buildo.optimizerStandalone, txtExt)
	acceptTest(testname, result.asmo.optimizerStandalone, expect.asmo.optimizerStandalone, asmoExt)
}
