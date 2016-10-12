package main

const (
	binDir  = "./bin"
	pikaDir = "./tests/pika"
	asmDir  = "./tests/asm"

	expectDir = "./expect"
	resultDir = "./result"
	backupDir = "./.backup"

	pikaExt = ".pika"
	asmExt  = ".asm"
	asmoExt = ".asmo"
	txtExt  = ".txt"

	open              = "open"
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
	compiler      = []string{"-ea", "-jar", binDir + "/" + compilerName}
	codegenerator = []string{"-ea", "-jar", binDir + "/" + codegeneratorName}
	optimizer     = []string{"-ea", "-jar", binDir + "/" + optimizerName}
	emulator      = []string{binDir + "/ASMEmu.exe"}

	// TODO reoptimize
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
		buildo: testStorageDir{
			codegenerator:       "",
			compiler:            "./result/buildo/compiler",
			optimizer:           "./result/buildo/optimizer",
			optimizerStandalone: "./result/buildo/optimizer-standalone",
		},
		asmo: testStorageDir{
			codegenerator:       "",
			compiler:            "./result/asmo/compiler",
			optimizer:           "./result/asmo/optimizer",
			optimizerStandalone: "./result/asmo/optimizer-standalone",
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
		buildo: testStorageDir{
			codegenerator:       "",
			compiler:            "./expect/buildo/compiler",
			optimizer:           "./expect/buildo/optimizer",
			optimizerStandalone: "./expect/buildo/optimizer-standalone",
		},
		asmo: testStorageDir{
			codegenerator:       "",
			compiler:            "./expect/asmo/compiler",
			optimizer:           "./expect/asmo/optimizer",
			optimizerStandalone: "./expect/asmo/optimizer-standalone",
		},
	}
)
