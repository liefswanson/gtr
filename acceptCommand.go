package main

import (
	"os"

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
	os.Rename(expectDir, backupDir)
	os.Rename(resultDir, expectDir)
	initResultDirs()
}

func acceptCompile(testname string) {
	new := buildPath(result.build.compiler, testname+txtExt)
	old := buildPath(expect.build.compiler, testname+txtExt)
	moveIfExists(old, new)

	new = buildPath(result.asm.compiler, testname+asmExt)
	old = buildPath(expect.asm.compiler, testname+asmExt)
	moveIfExists(old, new)

	new = buildPath(result.run.compiler, testname+txtExt)
	old = buildPath(expect.run.compiler, testname+txtExt)
	moveIfExists(old, new)
}

func acceptCodeGen(testname string) {
	new := buildPath(result.build.codegenerator, testname+txtExt)
	old := buildPath(expect.build.codegenerator, testname+txtExt)
	moveIfExists(old, new)

	new = buildPath(result.asm.codegenerator, testname+asmExt)
	old = buildPath(expect.asm.codegenerator, testname+asmExt)
	moveIfExists(old, new)

	new = buildPath(result.run.codegenerator, testname+txtExt)
	old = buildPath(expect.run.codegenerator, testname+txtExt)
	moveIfExists(old, new)
}

func acceptOptimized(testname string) {
	new := buildPath(result.build.optimizer, testname+txtExt)
	old := buildPath(expect.build.optimizer, testname+txtExt)
	moveIfExists(old, new)

	new = buildPath(result.asm.optimizer, testname+asmoExt)
	old = buildPath(expect.asm.optimizer, testname+asmoExt)
	moveIfExists(old, new)

	new = buildPath(result.run.optimizer, testname+txtExt)
	old = buildPath(expect.run.optimizer, testname+txtExt)
	moveIfExists(old, new)
}

func acceptOptimizedStandalone(testname string) {
	new := buildPath(result.build.optimizerStandalone, testname+txtExt)
	old := buildPath(expect.build.optimizerStandalone, testname+txtExt)
	moveIfExists(old, new)

	new = buildPath(result.asm.optimizerStandalone, testname+asmoExt)
	old = buildPath(expect.asm.optimizerStandalone, testname+asmoExt)
	moveIfExists(old, new)

	new = buildPath(result.run.optimizerStandalone, testname+txtExt)
	old = buildPath(expect.run.optimizerStandalone, testname+txtExt)
	moveIfExists(old, new)
}
