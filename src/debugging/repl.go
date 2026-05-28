package debugging

import (
	"fmt"
	"strconv"
	"strings"
)

func (d *Debugger) commandLoop() error {
	for {
		_, _ = fmt.Fprint(d.out, "> ")
		line, _ := d.in.ReadString('\n')
		line = strings.TrimSpace(line)

		switch {
		case line == "" || line == "s":
			d.StepMode = true
			return nil

		case line == "c":
			return nil

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
			d.renderCallStack(d.LastFrames)

		case line == "q":
			return ErrQuit

		default:
			_, _ = fmt.Fprintf(d.out, "unknown command %q\n", line)
		}
	}
}
