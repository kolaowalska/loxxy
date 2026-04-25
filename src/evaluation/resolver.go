package evaluation

type Resolver struct {
	interpreter *Interpreter
	scopes      []map[string]bool
}
