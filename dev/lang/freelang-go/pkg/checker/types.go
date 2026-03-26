package checker

import (
	"fmt"
	"strings"
)

// Type represents a FreeLang type
type Type interface {
	String() string
	TypeName() string
}

// BasicType represents primitive types
type BasicType struct {
	name string
}

func (bt *BasicType) String() string {
	return bt.name
}

func (bt *BasicType) TypeName() string {
	return bt.name
}

var (
	// Primitive types
	TypeInt    = &BasicType{"int"}
	TypeFloat  = &BasicType{"float"}
	TypeString = &BasicType{"string"}
	TypeBool   = &BasicType{"bool"}
	TypeNull   = &BasicType{"null"}
	TypeVoid   = &BasicType{"void"}
	TypeAny    = &BasicType{"any"}
)

// ArrayType represents an array type
type ArrayType struct {
	ElementType Type
}

func (at *ArrayType) String() string {
	return fmt.Sprintf("[%s]", at.ElementType.String())
}

func (at *ArrayType) TypeName() string {
	return "array"
}

func NewArrayType(elementType Type) *ArrayType {
	return &ArrayType{ElementType: elementType}
}

// HashType represents a hash/object type
type HashType struct {
	KeyType   Type
	ValueType Type
}

func (ht *HashType) String() string {
	return fmt.Sprintf("{%s: %s}", ht.KeyType.String(), ht.ValueType.String())
}

func (ht *HashType) TypeName() string {
	return "hash"
}

func NewHashType(keyType, valueType Type) *HashType {
	return &HashType{
		KeyType:   keyType,
		ValueType: valueType,
	}
}

// FunctionType represents a function type
type FunctionType struct {
	ParameterTypes []Type
	ReturnType     Type
}

func (ft *FunctionType) String() string {
	params := make([]string, len(ft.ParameterTypes))
	for i, t := range ft.ParameterTypes {
		params[i] = t.String()
	}
	return fmt.Sprintf("func(%s) %s", strings.Join(params, ", "), ft.ReturnType.String())
}

func (ft *FunctionType) TypeName() string {
	return "function"
}

func NewFunctionType(paramTypes []Type, returnType Type) *FunctionType {
	return &FunctionType{
		ParameterTypes: paramTypes,
		ReturnType:     returnType,
	}
}

// UnionType represents a union type (type1 | type2)
type UnionType struct {
	Types []Type
}

func (ut *UnionType) String() string {
	typeStrs := make([]string, len(ut.Types))
	for i, t := range ut.Types {
		typeStrs[i] = t.String()
	}
	return strings.Join(typeStrs, " | ")
}

func (ut *UnionType) TypeName() string {
	return "union"
}

func NewUnionType(types ...Type) *UnionType {
	return &UnionType{Types: types}
}

// TypesEqual checks if two types are equal
func TypesEqual(a, b Type) bool {
	// Handle nil cases
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}

	// Handle basic types (same instance)
	if aBasic, ok := a.(*BasicType); ok {
		if bBasic, ok := b.(*BasicType); ok {
			return aBasic.name == bBasic.name
		}
		return false
	}

	// Handle array types
	if aArray, ok := a.(*ArrayType); ok {
		if bArray, ok := b.(*ArrayType); ok {
			return TypesEqual(aArray.ElementType, bArray.ElementType)
		}
		return false
	}

	// Handle hash types
	if aHash, ok := a.(*HashType); ok {
		if bHash, ok := b.(*HashType); ok {
			return TypesEqual(aHash.KeyType, bHash.KeyType) &&
				TypesEqual(aHash.ValueType, bHash.ValueType)
		}
		return false
	}

	// Handle function types
	if aFunc, ok := a.(*FunctionType); ok {
		if bFunc, ok := b.(*FunctionType); ok {
			if len(aFunc.ParameterTypes) != len(bFunc.ParameterTypes) {
				return false
			}
			for i, at := range aFunc.ParameterTypes {
				if !TypesEqual(at, bFunc.ParameterTypes[i]) {
					return false
				}
			}
			return TypesEqual(aFunc.ReturnType, bFunc.ReturnType)
		}
		return false
	}

	return false
}

// IsCompatible checks if type a is compatible with type b
func IsCompatible(a, b Type) bool {
	if TypesEqual(a, b) {
		return true
	}

	// Any type is compatible with anything
	if aBasic, ok := a.(*BasicType); ok && aBasic.name == "any" {
		return true
	}
	if bBasic, ok := b.(*BasicType); ok && bBasic.name == "any" {
		return true
	}

	// Handle union types
	if aUnion, ok := a.(*UnionType); ok {
		for _, t := range aUnion.Types {
			if IsCompatible(t, b) {
				return true
			}
		}
		return false
	}

	return false
}
