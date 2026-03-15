package analyzer

// Symbol represents a named entity (variable, function, struct, etc.)
type Symbol struct {
	Name       string
	Type       string // "i64", "f64", "string", "bool", "array", "function", "struct", "enum", etc.
	IsMutable  bool
	IsUsed     bool
	IsGlobal   bool
	InitialValue interface{}
}

// Scope represents a lexical scope (global, function, block)
type Scope struct {
	parent  *Scope
	name    string
	symbols map[string]*Symbol
}

// NewScope creates a new scope
func NewScope(parent *Scope, name string) *Scope {
	return &Scope{
		parent:  parent,
		name:    name,
		symbols: make(map[string]*Symbol),
	}
}

// define defines a symbol in this scope
func (s *Scope) define(name string, symbol *Symbol) {
	s.symbols[name] = symbol
}

// isDefined checks if a symbol is defined in this scope (not parent scopes)
func (s *Scope) isDefined(name string) bool {
	_, ok := s.symbols[name]
	return ok
}

// resolve looks up a symbol in this scope or parent scopes
func (s *Scope) resolve(name string) *Symbol {
	if symbol, ok := s.symbols[name]; ok {
		return symbol
	}
	if s.parent != nil {
		return s.parent.resolve(name)
	}
	return nil
}

// All returns all symbols in this scope
func (s *Scope) All() map[string]*Symbol {
	return s.symbols
}

// Parent returns the parent scope
func (s *Scope) Parent() *Scope {
	return s.parent
}

// Name returns the scope name
func (s *Scope) Name() string {
	return s.name
}
