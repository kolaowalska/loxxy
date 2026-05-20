package tests

import (
	"bytes"
	"testing"

	"github.com/kolaowalska/loxxy/src/debugging"
	"github.com/kolaowalska/loxxy/src/evaluation"
	parser "github.com/kolaowalska/loxxy/src/parsing"
	resolving "github.com/kolaowalska/loxxy/src/resolving"
	scanner "github.com/kolaowalska/loxxy/src/scanning"
	"github.com/kolaowalska/loxxy/src/testutils"
)

func runDebuggerProgram(
	t *testing.T,
	source string,
	input string,
	setup func(*debugging.Debugger),
) (*debugging.Debugger, string, string) {
	t.Helper()

	reporter := &testutils.TestReporter{}

	s := scanner.NewScanner(source, reporter)
	tokens := s.ScanTokens()

	p := parser.NewParser(tokens, reporter)

	statements, err := p.Parse()
	if err != nil {
		t.Fatal(err)
	}

	var progOut bytes.Buffer
	var dbgOut bytes.Buffer

	interp := evaluation.NewInterpreter()
	interp.Stdout = &progOut

	in := bytes.NewBufferString(input)

	dbg := debugging.NewDebugger(in, &dbgOut)

	if setup != nil {
		setup(dbg)
	}

	interp.Hook = dbg

	resolver := resolving.NewResolver(interp, reporter)

	if err := resolver.ResolveStatements(statements); err != nil {
		t.Fatal(err)
	}

	err = interp.Interpret(statements)
	if err != nil {
		t.Fatal(err)
	}

	return dbg, progOut.String(), dbgOut.String()
}

func TestDebuggerBreakpoint(t *testing.T) {
	source := `var a = 10;
var b = 20;
print a + b;`

	dbg, progOut, _ := runDebuggerProgram(
		t,
		source,
		"c\nc\n",
		func(d *debugging.Debugger) {
			d.SetBreakpoint(3)
		},
	)

	if dbg.LastPausedLine != 3 {
		t.Fatalf(
			"expected breakpoint at line 3, got %d",
			dbg.LastPausedLine,
		)
	}

	if dbg.LastLocals["a"] != 10.0 {
		t.Fatalf("expected local a=10")
	}

	if dbg.LastLocals["b"] != 20.0 {
		t.Fatalf("expected local b=20")
	}

	if progOut != "30\n" {
		t.Fatalf("unexpected program output: %q", progOut)
	}
}

func TestDebuggerFunctionFrame(t *testing.T) {
	source := `fun add(x, y) {
  var z = x + y;
  print z;
}

add(3, 4);`

	dbg, progOut, _ := runDebuggerProgram(
		t,
		source,
		"c\nc\n",
		func(d *debugging.Debugger) {
			d.SetBreakpoint(2)
		},
	)

	found := false

	for _, pause := range dbg.PauseHistory {
		if pause.Line != 2 {
			continue
		}

		if pause.Locals["x"] == 3.0 &&
			pause.Locals["y"] == 4.0 {

			found = true
			break
		}
	}

	if !found {
		t.Fatalf("expected function locals x=3 y=4")
	}
	if progOut != "7\n" {
		t.Fatalf("unexpected program output: %q", progOut)
	}
}

func TestDebuggerLoopBreakpoint(t *testing.T) {
	source := `for (var i = 0; i < 3; i = i + 1) {
  print i;
}`

	dbg, progOut, _ := runDebuggerProgram(
		t,
		source,
		"c\nc\nc\nc\n",
		func(d *debugging.Debugger) {
			d.SetBreakpoint(2)
		},
	)

	if dbg.LastPausedLine != 2 {
		t.Fatalf(
			"expected pause at line 2, got %d",
			dbg.LastPausedLine,
		)
	}

	if progOut != "0\n1\n2\n" {
		t.Fatalf("unexpected program output: %q", progOut)
	}
}

func TestDebuggerNestedCalls(t *testing.T) {
	source := `fun c() {
  var x = 123;
  print x;
}

fun b() {
  c();
}

fun a() {
  b();
}

a();`

	dbg, progOut, _ := runDebuggerProgram(
		t,
		source,
		"c\nc\nc\n",
		func(d *debugging.Debugger) {
			d.SetBreakpoint(2)
		},
	)

	found := false

	for _, pause := range dbg.PauseHistory {
		if len(pause.Frames) == 0 {
			continue
		}

		top := pause.Frames[len(pause.Frames)-1]

		if top.Name == "c" {
			found = true
			break
		}
	}

	if !found {
		t.Fatalf("expected pause inside function c")
	}
	if progOut != "123\n" {
		t.Fatalf("unexpected program output: %q", progOut)
	}
}

func TestDebuggerRecursiveFunction(t *testing.T) {
	source := `fun fact(n) {
  if (n <= 1) return 1;
  return n * fact(n - 1);
}

print fact(3);`

	dbg, progOut, _ := runDebuggerProgram(
		t,
		source,
		"c\nc\nc\nc\nc\n",
		func(d *debugging.Debugger) {
			d.SetBreakpoint(2)
		},
	)

	foundRecursive := false

	for _, pause := range dbg.PauseHistory {
		if len(pause.Frames) >= 2 {
			foundRecursive = true
			break
		}
	}

	if !foundRecursive {
		t.Fatalf("expected recursive call stack")
	}

	if progOut != "6\n" {
		t.Fatalf("unexpected program output: %q", progOut)
	}
}

func TestDebuggerClosureEnvironment(t *testing.T) {
	source := `fun outer() {
  var x = 123;

  fun inner() {
    print x;
  }

  inner();
}

outer();`

	dbg, progOut, _ := runDebuggerProgram(
		t,
		source,
		"c\nc\n",
		func(d *debugging.Debugger) {
			d.SetBreakpoint(5)
		},
	)

	if dbg.LastLocals["x"] != 123.0 {
		t.Fatalf("expected closure variable x=123")
	}

	if progOut != "123\n" {
		t.Fatalf("unexpected program output: %q", progOut)
	}
}
