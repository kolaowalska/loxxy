package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/kolaowalska/loxxy/src/evaluation"
	parser "github.com/kolaowalska/loxxy/src/parsing"
	scanner "github.com/kolaowalska/loxxy/src/scanning"
)

var hadError = false
var hadRuntimeError = false

// LoxReporter - Concrete implementation of scanner.ErrorReporter
type LoxReporter struct{}

func (r LoxReporter) Error(line int, message string) {
	report(line, "", message)
}

func report(line int, where string, message string) {
	log.Print("[line: " + strconv.Itoa(line) + "] error" + where + ": " + message)
	hadError = true
}

func (r LoxReporter) TokenError(t scanner.Token, message string) {
	if t.TokenType == scanner.EOF {
		report(t.Line, " at end", message)
	} else {
		report(t.Line, " at '"+t.Lexeme+"'", message)
	}

func reportRuntimeError(err *evaluation.RuntimeError) {
	msg := fmt.Sprintf("%s\n[line %d]", err.Message, err.Token.Line)
	log.Print(msg)
	hadRuntimeError = true
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
	if hadRuntimeError {
		os.Exit(70)
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
	tokens := newScanner.ScanTokens()

	newParser := parser.NewParser(tokens, reporter)
	expression := newParser.Parse()

	if hadError {
		return
	}
	result, err := evaluation.Evaluate(expression)
	if err != nil {
		if runtimeErr, ok := err.(*evaluation.RuntimeError); ok {
			reportRuntimeError(runtimeErr)
		} else {
			log.Printf("Error: %v", err)
		}
		return
	}
	fmt.Println(stringify(result))
}
func stringify(obj any) string {
	if obj == nil {
		return "nil"
	}
	if val, ok := obj.(float64); ok {
		text := fmt.Sprintf("%v", val)
		return strings.TrimSuffix(text, ".0")
	}

	return fmt.Sprintf("%v", obj)
}
