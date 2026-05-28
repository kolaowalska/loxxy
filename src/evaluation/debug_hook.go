package evaluation

type CallFrame struct {
	Name string
	Line int
}

type DebugHook interface {
	// OnStatement called before every statement, blocks when a breakpoint (or step) fires
	OnStatement(line int, frames []CallFrame, env *Environment) error
}
