package representation

import (
	"fmt"
	"strings"
)

func Print(expr Expr) string {
	switch e := expr.(type) {

	case *Binary:
		return parenthesize(e.Operator.Lexeme, e.Left, e.Right)

	case *Grouping:
		return parenthesize("group", e.Expression)

	case *Literal:
		if e.Value == nil {
			return "nil"
		}
		return fmt.Sprintf("%v", e.Value)

	case *Unary:
		return parenthesize(e.Operator.Lexeme, e.Right)

	default:
		return "UNKOWN_EXPR"
	}
}

func parenthesize(name string, exprs ...Expr) string {
	var builder strings.Builder

	builder.WriteString("(" + name)
	for _, expr := range exprs {
		builder.WriteString(" ")
		builder.WriteString(Print(expr))
	}
	builder.WriteString(")")

	return builder.String()
}
