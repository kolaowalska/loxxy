package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"

	scanner "github.com/kolaowalska/loxxy/src/scanning"
)

var hadError = false

// LoxReporter - Concrete implementation of scanner.ErrorReporter
type LoxReporter struct{}

func (r LoxReporter) Error(line int, message string) {
	report(line, "", message)
}

func report(line int, where string, message string) {
	log.Printf("[line: " + strconv.Itoa(line) + "] error" + where + ": " + message)
	hadError = true
}

func init() {
	log.SetFlags(0)
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
		log.Printf("error: %v\n", err)
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
		_, err := fmt.Fprint(out, "> ")
		if err != nil {
			log.Printf("error writing prompt: %v\n", err)
			break
		}
		line, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				// exit on end of file (ctrl+d)
				break
			}
			log.Printf("error reading prompt: %v\n", err)
			break
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
