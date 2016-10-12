package main

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
	build  testStorageDir
	run    testStorageDir
	asm    testStorageDir
	buildo testStorageDir
	asmo   testStorageDir
}
