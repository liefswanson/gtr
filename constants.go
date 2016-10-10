package main

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
	bin  = "./bin"
	pika = "./tests/pika"
	asm  = "./tests/asm"

	pikaExt = ".pika"
	asmExt  = ".asm"
	asmoExt = ".asmo"
	txtExt  = ".txt"

	java = "java"
	wine = "wine"
)

var (
	compiler      = []string{"-ea", "-jar", bin + "/pika-compiler.jar"}
	codegenerator = []string{"-ea", "-jar", bin + "/pika-codegen.jar"}
	optimizer     = []string{"-ea", "-jar", bin + "/pika-optimizer.jar"}
	emulator      = []string{bin + "/ASMEmu.exe"}

	result = testDirTree{
		build: testStorageDir{
			codegenerator:       "./result/build/codegenerator",
			compiler:            "./result/build/compiler",
			optimizer:           "./result/build/optimizer",
			optimizerStandalone: "./result/build/optimizer-standalone",
		},
		run: testStorageDir{
			codegenerator:       "./result/run/codegenerator",
			compiler:            "./result/run/compiler",
			optimizer:           "./result/run/optimizer",
			optimizerStandalone: "./result/run/optimizer-standalone",
		},
		asm: testStorageDir{
			codegenerator:       "./result/asm/codegenerator",
			compiler:            "./result/asm/compiler",
			optimizer:           "./result/asm/optimizer",
			optimizerStandalone: "./result/asm/optimizer-standalone",
		},
	}

	expect = testDirTree{
		build: testStorageDir{
			codegenerator:       "./expect/build/codegenerator",
			compiler:            "./expect/build/compiler",
			optimizer:           "./expect/build/optimizer",
			optimizerStandalone: "./expect/build/optimizer-standalone",
		},
		run: testStorageDir{
			codegenerator:       "./expect/run/codegenerator",
			compiler:            "./expect/run/compiler",
			optimizer:           "./expect/run/optimizer",
			optimizerStandalone: "./expect/run/optimizer-standalone",
		},
		asm: testStorageDir{
			codegenerator:       "./expect/asm/codegenerator",
			compiler:            "./expect/asm/compiler",
			optimizer:           "./expect/asm/optimizer",
			optimizerStandalone: "./expect/asm/optimizer-standalone",
		},
	}
)
