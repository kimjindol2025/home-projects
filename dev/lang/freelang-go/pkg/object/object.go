package object

import "fmt"

// Object is the base interface for all runtime objects
type Object interface {
	Inspect() string
}

// Integer represents an integer value
type Integer struct {
	Value int64
}

func (i *Integer) Inspect() string {
	return fmt.Sprintf("%d", i.Value)
}

// Float represents a floating-point value
type Float struct {
	Value float64
}

func (f *Float) Inspect() string {
	return fmt.Sprintf("%f", f.Value)
}

// String represents a string value
type String struct {
	Value string
}

func (s *String) Inspect() string {
	return s.Value
}

// Boolean represents a boolean value
type Boolean struct {
	Value bool
}

func (b *Boolean) Inspect() string {
	return fmt.Sprintf("%v", b.Value)
}

// Null represents a null/nil value
type Null struct{}

func (n *Null) Inspect() string {
	return "null"
}

// Array represents an array value
type Array struct {
	Elements []Object
}

func (a *Array) Inspect() string {
	result := "["
	for i, elem := range a.Elements {
		if i > 0 {
			result += ", "
		}
		result += elem.Inspect()
	}
	result += "]"
	return result
}

// Hash represents a hash/map value
type Hash struct {
	Pairs map[string]Object
}

func (h *Hash) Inspect() string {
	result := "{"
	first := true
	for key, val := range h.Pairs {
		if !first {
			result += ", "
		}
		first = false
		result += fmt.Sprintf("%s: %s", key, val.Inspect())
	}
	result += "}"
	return result
}

// Function represents a function value
type Function struct {
	Name string
}

func (f *Function) Inspect() string {
	return fmt.Sprintf("fn:%s", f.Name)
}
