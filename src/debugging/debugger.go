package debugging

import (
	"bufio"
	"io"
	"sort"
	"strings"

	"github.com/kolaowalska/loxxy/src/evaluation"
)

type PauseSnapshot struct {
	Line   int
	Frames []evaluation.CallFrame
	Locals map[string]any
}

type Debugger struct {
	breakpoints map[int]bool
	stepMode    bool
	sourceLines []string

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
		stepMode:    true,
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
) {
	if !d.stepMode && !d.breakpoints[line] {
		return
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

	d.stepMode = false
	d.renderPause(line, frames, env)
	d.commandLoop()
}
