package main

const (
	bin  = "./bin"
	pika = "./tests/pika"
	asm  = "./tests/asm"

	expectDir = "./expect"
	resultDir = "./result"
	backupDir = "./.backup"

	pikaExt = ".pika"
	asmExt  = ".asm"
	asmoExt = ".asmo"
	txtExt  = ".txt"

	linuxOpen         = "xdg-open"
	macOpen           = "open"
	java              = "java"
	compilerName      = "pika-compiler.jar"
	codegeneratorName = "pika-codegen.jar"
	optimizerName     = "pika-optimizer.jar"
	wine              = "wine"

	loggingMessage = "logging.PikaLogger log"

	basicAsmFile  = "Halt\n"
	basicPikaFile = "exec {\n\n}\n"
)

// these are vars, but just as a technical restriction
// they should be considered constants
var (
	compiler      = []string{"-ea", "-jar", bin + "/" + compilerName}
	codegenerator = []string{"-ea", "-jar", bin + "/" + codegeneratorName}
	optimizer     = []string{"-ea", "-jar", bin + "/" + optimizerName}
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
