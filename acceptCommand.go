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
		path = buildPath(asm, filename)
	} else {
		filename := testname + pikaExt
		path = buildPath(pika, filename)
	}

	if exists(path) {
		if flags.asm {
			acceptOptimizedStandalone(testname)
		} else {
			acceptCodeGen(testname)
			acceptOptimized(testname)
			acceptCompile(testname)
		}
	} else {
		color.Magenta(path + " does not exist")
	}
}

func acceptAll() {
	os.RemoveAll("./.expect")
	os.Rename("./expect/", "./.expect/")
	os.Rename("./result", "./expect")
	initResultDirs()
}

func acceptCompile(testname string) {
	old := buildPath(result.build.compiler, testname+txtExt)
	new := buildPath(expect.build.compiler, testname+txtExt)
	moveIfExists(old, new)

	new = buildPath(result.asm.compiler, testname+asmExt)
	old = buildPath(expect.asm.compiler, testname+asmExt)
	moveIfExists(old, new)

	new = buildPath(result.run.compiler, testname+txtExt)
	old = buildPath(expect.run.compiler, testname+txtExt)
	moveIfExists(old, new)
}

func acceptCodeGen(testname string) {
	old := buildPath(result.build.codegenerator, testname+txtExt)
	new := buildPath(expect.build.codegenerator, testname+txtExt)
	moveIfExists(old, new)

	new = buildPath(result.asm.codegenerator, testname+asmExt)
	old = buildPath(expect.asm.codegenerator, testname+asmExt)
	moveIfExists(old, new)

	new = buildPath(result.run.codegenerator, testname+txtExt)
	old = buildPath(expect.run.codegenerator, testname+txtExt)
	moveIfExists(old, new)
}

func acceptOptimized(testname string) {
	old := buildPath(result.build.optimizer, testname+txtExt)
	new := buildPath(expect.build.optimizer, testname+txtExt)
	moveIfExists(old, new)

	new = buildPath(result.asm.optimizer, testname+asmoExt)
	old = buildPath(expect.asm.optimizer, testname+asmoExt)
	moveIfExists(old, new)

	new = buildPath(result.run.optimizer, testname+txtExt)
	old = buildPath(expect.run.optimizer, testname+txtExt)
	moveIfExists(old, new)
}

func acceptOptimizedStandalone(testname string) {
	old := buildPath(result.build.optimizerStandalone, testname+txtExt)
	new := buildPath(expect.build.optimizerStandalone, testname+txtExt)
	moveIfExists(old, new)

	new = buildPath(result.asm.optimizerStandalone, testname+asmoExt)
	old = buildPath(expect.asm.optimizerStandalone, testname+asmoExt)
	moveIfExists(old, new)

	new = buildPath(result.run.optimizerStandalone, testname+txtExt)
	old = buildPath(expect.run.optimizerStandalone, testname+txtExt)
	moveIfExists(old, new)
}
