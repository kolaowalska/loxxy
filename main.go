package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"

	"github.com/kolaowalska/loxxy/src/evaluation"
	parser "github.com/kolaowalska/loxxy/src/parsing"
	"github.com/kolaowalska/loxxy/src/resolving"
	scanner "github.com/kolaowalska/loxxy/src/scanning"
)

var hadError = false
var hadRuntimeError = false

var interpreter = evaluation.NewInterpreter()

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
		hadRuntimeError = false
	}
}

func run(source string) {
	reporter := LoxReporter{}

	s := scanner.NewScanner(source, reporter)
	tokens := s.ScanTokens()

	p := parser.NewParser(tokens, reporter)
	statements, _ := p.Parse()

	if hadError {
		return
	}

	interpreter.ClearLocals()

	resolver := resolving.NewResolver(interpreter, reporter)
	_ = resolver.ResolveStatements(statements)
	if hadError {
		return
	}

	err := interpreter.Interpret(statements)
	if err != nil {
		if rterr, ok := errors.AsType[*evaluation.RuntimeError](err); ok {
			reportRuntimeError(rterr)
		}
	}

}
