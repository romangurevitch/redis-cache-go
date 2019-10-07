package test

import (
	"runtime/debug"
	"testing"
)

func CheckError(t *testing.T, err error) {
	if err != nil {
		debug.PrintStack()
		t.Fatal(err)
	}
}
