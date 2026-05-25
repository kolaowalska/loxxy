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

const (
	exitSuccess  = 0
	exitFailure  = 1
	exitUsage    = 64 // incorrect command line usage
	exitDataErr  = 65 // incorrect input data
	exitNoInput  = 66 // input file does not exist or is unreadable
	exitSoftware = 70 // internal software error
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
		os.Exit(exitNoInput)
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
		os.Exit(exitDataErr)
	}
	if hadRuntimeError {
		os.Exit(exitSoftware)
	}
}

func runFile(path string) {
	bytes, err := os.ReadFile(path)
	if err != nil {
		log.Printf("error: %v\n", err)
		os.Exit(exitNoInput)
	}
	run(string(bytes))

	if hadError {
		os.Exit(exitDataErr)
	}
	if hadRuntimeError {
		os.Exit(exitSoftware)
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
		if errors.Is(err, debugging.ErrQuit) {
			os.Exit(exitSuccess)
		}

		if rterr, ok := errors.AsType[*evaluation.RuntimeError](err); ok {
			reportRuntimeError(rterr)
		}
	}
}

func main() {
	debugFlag := flag.Bool("debug", false, "run with debugger")
	breakFlag := flag.Int("break", 0, "set an initial breakpoint at line N")

	flag.Parse()
	args := flag.Args()

	// silently set debug mode
	isDebugMode := *debugFlag || *breakFlag > 0

	if isDebugMode {
		if len(args) != 1 {
			fmt.Println("usage: loxxy -debug [-break N] <script>")
			os.Exit(exitUsage)
		}

		runFileDebug(args[0], *breakFlag)
		return
	}

	if len(args) > 1 {
		fmt.Println("usage: loxxy [script]")
		os.Exit(exitUsage)
	}

	if len(args) == 1 {
		runFile(args[0])
	} else {
		runPrompt(os.Stdin, os.Stdout)
	}
}
