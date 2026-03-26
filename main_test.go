package main

import (
	"testing"

	"github.com/kolaowalska/loxxy/src/reports"
)

func TestRun_ValidSource(t *testing.T) {
	defer reporter.Clear()
	reporter.Clear()

	source := "var a = 1; print a;"

	run(source)

	if reporter.HadError {
		t.Errorf("run() wybombilo a kod jest niby dobry")
	}
}

func TestRun_InvalidSource(t *testing.T) {
	defer reporter.Clear()
	reporter.Clear()

	source := "var a = @ gugugaga;"

	run(source)

	if !reporter.HadError {
		t.Errorf("run() zadzialalo a powinno sie wybombic")
	}
}
