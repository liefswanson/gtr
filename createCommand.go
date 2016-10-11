package main

import (
	"io/ioutil"

	"github.com/fatih/color"
)

func createCommand(flags createFlags, arg string) {
	var path string
	var contents []byte

	if flags.asm {
		filename := arg + asmExt
		path = buildPath(asm, filename)
		contents = []byte(basicAsmFile)
	} else {
		filename := arg + pikaExt
		path = buildPath(pika, filename)
		contents = []byte(basicPikaFile)
	}

	if !exists(path) {
		ioutil.WriteFile(path, contents, 0777)
	} else {
		color.Magenta(path + " already exists")
	}

	if flags.open {
		openDefaultEditor(path)
	}
}
