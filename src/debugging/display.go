package debugging

import (
	"fmt"
	"sort"
	"strings"

	"github.com/kolaowalska/loxxy/src/evaluation"
)

func (d *Debugger) renderSourceContext(line int) {
	start := line - 3
	if start < 1 {
		start = 1
	}

	end := line + 2
	if end > len(d.sourceLines) {
		end = len(d.sourceLines)
	}

	for i := start; i <= end; i++ {
		src := ""
		if i-1 < len(d.sourceLines) {
			src = d.sourceLines[i-1]
		}
		marker := "  "
		if i == line {
			marker = "→ "
		}
		_, err := fmt.Fprintf(d.out, "%s%3d | %s\n", marker, i, src)
		if err != nil {
			return
		}
	}
}

func (d *Debugger) renderLocals(env *evaluation.Environment) {
	snap := env.Snapshot()
	keys := make([]string, 0, len(snap))
	for k := range snap {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		_, err := fmt.Fprintf(d.out, "  %-12s = %v\n", k, snap[k])
		if err != nil {
			return
		}
	}
}

func (d *Debugger) renderCallStack(frames []evaluation.CallFrame) {
	lastIndex := len(frames) - 1
	for i, f := range frames {
		marker := "  "
		if i == lastIndex {
			marker = "→ "
		}
		_, err := fmt.Fprintf(d.out, "  %s%-16s line %d\n", marker, f.Name, f.Line)
		if err != nil {
			return
		}
	}
}

func (d *Debugger) renderPause(
	line int,
	frames []evaluation.CallFrame,
	env *evaluation.Environment,
) {
	sep := strings.Repeat("—", 50)
	_, err := fmt.Fprintf(d.out, "\n—— PAUSED — line %d %s\n\n", line, sep)
	if err != nil {
		return
	}
	d.renderSourceContext(line)
	_, err = fmt.Fprintf(d.out, "\n—— LOCALS ———— %s\n", sep)
	if err != nil {
		return
	}
	d.renderLocals(env)
	_, err = fmt.Fprintf(d.out, "\n—— CALL STACK %s\n", sep)
	if err != nil {
		return
	}
	d.renderCallStack(frames)
	_, err = fmt.Fprintf(d.out, "\n%s\n", sep)
	if err != nil {
		return
	}
	_, err = fmt.Fprintf(d.out, "[enter] step  [c] continue  [b N] breakpoint  [d N] delete  [q] quit\n")
	if err != nil {
		return
	}
}
