package app

import (
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	// setup before tests
	_ = os.Remove("traces.txt")

	initTrace()

	// run tests
	exitcode := m.Run()

	// cleanup

	// exit
	os.Exit(exitcode)
}
