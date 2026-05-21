package debugging

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

func (d *Debugger) commandLoop() {
	for {
		_, _ = fmt.Fprint(d.out, "> ")

		line, _ := d.in.ReadString('\n')
		line = strings.TrimSpace(line)

		switch {

		case line == "" || line == "s":
			// Enter or s → step
			d.StepMode = true
			return

		case line == "c":
			// continue
			return

		case strings.HasPrefix(line, "b "):
			n, err := strconv.Atoi(strings.TrimPrefix(line, "b "))
			if err == nil {
				d.SetBreakpoint(n)
			}

		case strings.HasPrefix(line, "d "):
			n, err := strconv.Atoi(strings.TrimPrefix(line, "d "))
			if err == nil {
				d.DeleteBreakpoint(n)
			}

		case line == "l":
			for _, bp := range d.ListBreakpoints() {
				_, _ = fmt.Fprintf(d.out, "  breakpoint at line %d\n", bp)
			}

		case line == "stack":
			// re-render just the stack section (cached from last pause)

		case line == "q":
			os.Exit(0)

		default:
			_, _ = fmt.Fprintf(d.out, "unknown command %q\n", line)
		}
	}
}
