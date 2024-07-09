package main

import (
	"os"
	"testing"
)

/*
 * The unit tests in this file simulate command line invocation.
 */
func TestMain(test *testing.T) {
	_ = test
	os.Args = []string{"command-name", "--avoid-serving"}
	main()
}
