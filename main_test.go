package main

import "testing"

// func TestRun_ValidSource(t *testing.T) {
// 	hadError = false
// 	source := "var a = 1; print a;"
// 	run(source)
//
// 	if hadError {
// 		t.Errorf("run() wybombilo a kod jest niby dobry")
// 	}
// 	hadError = false
// }

func TestRun_InvalidSource(t *testing.T) {
	hadError = false
	source := "var a = @ gugugaga;"
	run(source)

	if !hadError {
		t.Errorf("run() zadzialalo a powinno sie wybombic")
	}
	hadError = false
}
