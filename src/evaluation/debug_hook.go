package evaluation

type CallFrame struct {
	Name string // function name or "<script>" for top-level
	Line int
}

type DebugHook interface {
	// Called before every statement, blocks when a breakpoint (or step) fires
	OnStatement(line int, frames []CallFrame, env *Environment) error
}
