package main

import (
	"bufio"
	"fmt"
	"strconv"

	scanner "github.com/kolaowalska/loxxy/src/scanning"

	"io"
	"os"
)

var hadError = false

// LoxReporter - Concrete implementation of scanner.ErrorReporter
type LoxReporter struct{}

func (r LoxReporter) Error(line int, message string) {
	report(line, "", message)
}

func report(line int, where string, message string) {
	fmt.Fprintln(os.Stderr, "[line: "+strconv.Itoa(line)+"] error"+where+": "+message)
	hadError = true
}

func main() {
	args := os.Args[1:]
	if len(args) > 1 {
		fmt.Println("usage: loxxy [script]") //TODO: later think about changing to loxxy
		os.Exit(64)
	} else if len(args) == 1 {
		runFile(args[0])
	} else {
		runPrompt(os.Stdin, os.Stdout)
	}
}

func runFile(path string) {
	bytes, err := os.ReadFile(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %s\n", err)
		os.Exit(66)
	}
	run(string(bytes))

	if hadError {
		os.Exit(65)
	}
}

func runPrompt(in io.Reader, out io.Writer) {
	reader := bufio.NewReader(in)
	for {
		fmt.Fprint(out, "> ")
		line, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				// exit on end of file (ctrl+d)
				break
			}
		}
		run(line)
		hadError = false
	}
}

func run(source string) {
	reporter := LoxReporter{}

	newScanner := scanner.NewScanner(source, reporter)
	newScanner.ScanTokens()

	// for debugging tests
	/*
		tokens := newScanner.ScanTokens()
		for index := range tokens {
			fmt.Println(tokens[index])
		}
	*/
}
