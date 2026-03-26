package checker

// Symbol represents a symbol in the symbol table
type Symbol struct {
	Name       string
	Type       Type
	Kind       SymbolKind
	Scope      *Scope
	Mutable    bool
	Initialized bool
}

// SymbolKind represents the kind of symbol
type SymbolKind int

const (
	SymbolVariable SymbolKind = iota
	SymbolFunction
	SymbolParameter
	SymbolBuiltin
)

func (sk SymbolKind) String() string {
	switch sk {
	case SymbolVariable:
		return "variable"
	case SymbolFunction:
		return "function"
	case SymbolParameter:
		return "parameter"
	case SymbolBuiltin:
		return "builtin"
	default:
		return "unknown"
	}
}

// Scope represents a scope level
type Scope struct {
	parent   *Scope
	symbols  map[string]*Symbol
	level    int
}

// NewScope creates a new scope
func NewScope(parent *Scope) *Scope {
	level := 0
	if parent != nil {
		level = parent.level + 1
	}

	return &Scope{
		parent:  parent,
		symbols: make(map[string]*Symbol),
		level:   level,
	}
}

// Define defines a symbol in the current scope
func (s *Scope) Define(name string, symbol *Symbol) {
	symbol.Scope = s
	s.symbols[name] = symbol
}

// Resolve resolves a symbol by name, looking up the scope chain
func (s *Scope) Resolve(name string) (*Symbol, bool) {
	symbol, ok := s.symbols[name]
	if ok {
		return symbol, true
	}

	if s.parent != nil {
		return s.parent.Resolve(name)
	}

	return nil, false
}

// ResolveLocal resolves a symbol only in the current scope
func (s *Scope) ResolveLocal(name string) (*Symbol, bool) {
	symbol, ok := s.symbols[name]
	return symbol, ok
}

// Parent returns the parent scope
func (s *Scope) Parent() *Scope {
	return s.parent
}

// Level returns the scope level (0 = global)
func (s *Scope) Level() int {
	return s.level
}

// Symbols returns all symbols in this scope
func (s *Scope) Symbols() map[string]*Symbol {
	return s.symbols
}

// SymbolTable manages all scopes
type SymbolTable struct {
	globals   *Scope
	current   *Scope
	builtins  map[string]*Symbol
}

// NewSymbolTable creates a new symbol table
func NewSymbolTable() *SymbolTable {
	st := &SymbolTable{
		globals:  NewScope(nil),
		builtins: make(map[string]*Symbol),
	}

	st.current = st.globals
	st.defineBuiltins()

	return st
}

// defineBuiltins defines built-in functions and types
func (st *SymbolTable) defineBuiltins() {
	// Built-in functions
	builtins := map[string]*Symbol{
		"len": {
			Name: "len",
			Type: NewFunctionType(
				[]Type{TypeAny},
				TypeInt,
			),
			Kind:     SymbolBuiltin,
			Mutable:  false,
		},
		"print": {
			Name: "print",
			Type: NewFunctionType(
				[]Type{TypeAny},
				TypeVoid,
			),
			Kind:    SymbolBuiltin,
			Mutable: false,
		},
		"push": {
			Name: "push",
			Type: NewFunctionType(
				[]Type{TypeAny, TypeAny},
				TypeVoid,
			),
			Kind:    SymbolBuiltin,
			Mutable: false,
		},
		"pop": {
			Name: "pop",
			Type: NewFunctionType(
				[]Type{TypeAny},
				TypeAny,
			),
			Kind:    SymbolBuiltin,
			Mutable: false,
		},
		"type": {
			Name: "type",
			Type: NewFunctionType(
				[]Type{TypeAny},
				TypeString,
			),
			Kind:    SymbolBuiltin,
			Mutable: false,
		},
	}

	for name, symbol := range builtins {
		st.builtins[name] = symbol
		st.globals.Define(name, symbol)
	}
}

// Define defines a symbol in the current scope
func (st *SymbolTable) Define(name string, kind SymbolKind, typ Type, mutable bool) *Symbol {
	symbol := &Symbol{
		Name:    name,
		Type:    typ,
		Kind:    kind,
		Mutable: mutable,
	}

	st.current.Define(name, symbol)
	return symbol
}

// Resolve resolves a symbol from the current scope
func (st *SymbolTable) Resolve(name string) (*Symbol, bool) {
	return st.current.Resolve(name)
}

// ResolveLocal resolves a symbol only in the current scope
func (st *SymbolTable) ResolveLocal(name string) (*Symbol, bool) {
	return st.current.ResolveLocal(name)
}

// PushScope pushes a new scope
func (st *SymbolTable) PushScope() {
	st.current = NewScope(st.current)
}

// PopScope pops the current scope
func (st *SymbolTable) PopScope() {
	if st.current.parent != nil {
		st.current = st.current.parent
	}
}

// CurrentScope returns the current scope
func (st *SymbolTable) CurrentScope() *Scope {
	return st.current
}

// GlobalScope returns the global scope
func (st *SymbolTable) GlobalScope() *Scope {
	return st.globals
}

// ScopeLevel returns the current scope level
func (st *SymbolTable) ScopeLevel() int {
	return st.current.Level()
}
