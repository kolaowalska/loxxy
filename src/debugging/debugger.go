package debugging

import (
	"bufio"
	"errors"
	"io"
	"sort"
	"strings"

	"github.com/kolaowalska/loxxy/src/evaluation"
)

// ErrQuit is returned when user explicitly quits the debugger
var ErrQuit = errors.New("debugger quit requested")

type PauseSnapshot struct {
	Line   int
	Frames []evaluation.CallFrame
	Locals map[string]any
}

type Debugger struct {
	breakpoints map[int]bool
	StepMode    bool
	sourceLines []string
	ContextSize int

	in  *bufio.Reader
	out io.Writer

	LastPausedLine int
	LastFrames     []evaluation.CallFrame
	LastLocals     map[string]any

	PauseHistory []PauseSnapshot
}

func NewDebugger(in io.Reader, out io.Writer) *Debugger {
	return &Debugger{
		breakpoints: make(map[int]bool),
		in:          bufio.NewReader(in),
		out:         out,
		StepMode:    true,
		ContextSize: 3,
	}
}

func (d *Debugger) LoadSource(src string) {
	d.sourceLines = strings.Split(src, "\n")
}

func (d *Debugger) SetBreakpoint(line int) {
	d.breakpoints[line] = true
}

func (d *Debugger) DeleteBreakpoint(line int) {
	delete(d.breakpoints, line)
}

func (d *Debugger) ListBreakpoints() []int {
	lines := make([]int, 0, len(d.breakpoints))
	for line := range d.breakpoints {
		lines = append(lines, line)
	}
	sort.Ints(lines)
	return lines
}

func (d *Debugger) OnStatement(
	line int,
	frames []evaluation.CallFrame,
	env *evaluation.Environment,
) error {
	if !d.StepMode && !d.breakpoints[line] {
		return nil
	}

	snapshot := PauseSnapshot{
		Line:   line,
		Frames: append([]evaluation.CallFrame{}, frames...),
		Locals: env.Snapshot(),
	}

	d.LastPausedLine = line
	d.LastFrames = snapshot.Frames
	d.LastLocals = snapshot.Locals

	d.PauseHistory = append(d.PauseHistory, snapshot)

	d.StepMode = false
	d.renderPause(line, frames, env)

	return d.commandLoop()
}
