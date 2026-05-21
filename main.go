/*
package main provides the loxxy command-line interface

usage:

	loxxy [script]
*/
package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"

	debugging "github.com/kolaowalska/loxxy/src/debugging"
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

func runFileDebug(path string, initialBreak int) {
	src, err := os.ReadFile(path)
	if err != nil {
		log.Printf("error: %v\n", err)
		os.Exit(66)
	}

	dbg := debugging.NewDebugger(os.Stdin, os.Stderr)
	dbg.LoadSource(string(src))
	if initialBreak > 0 {
		dbg.SetBreakpoint(initialBreak)
		dbg.StepMode = false
	}

	interpreter.Hook = dbg

	run(string(src))

	if hadError {
		os.Exit(65)
	}
	if hadRuntimeError {
		os.Exit(70)
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

func main() {
	debugFlag := flag.Bool("debug", false, "run with debugger")
	breakFlag := flag.Int("break", 0, "set an initial breakpoint at line N (requires -debug)")

	flag.Parse()
	args := flag.Args()

	if *debugFlag {
		if len(args) != 1 {
			fmt.Println("usage: loxxy -debug [-break N] <script>")
			os.Exit(64)
		}

		runFileDebug(args[0], *breakFlag)
		return
	}

	if len(args) > 1 {
		fmt.Println("usage: loxxy [script]")
		os.Exit(64)
	}

	if len(args) == 1 {
		runFile(args[0])
	} else {
		runPrompt(os.Stdin, os.Stdout)
	}
}
