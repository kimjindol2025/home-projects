package analyzer

// Type represents a FreeLang type
type Type interface {
	TypeName() string
	String() string
}

// PrimitiveType represents primitive types
type PrimitiveType struct {
	name string
}

func (t *PrimitiveType) TypeName() string {
	return t.name
}

func (t *PrimitiveType) String() string {
	return t.name
}

// Predefined primitive types
var (
	TypeI64    = &PrimitiveType{name: "i64"}
	TypeF64    = &PrimitiveType{name: "f64"}
	TypeString = &PrimitiveType{name: "string"}
	TypeBool   = &PrimitiveType{name: "bool"}
	TypeVoid   = &PrimitiveType{name: "void"}
	TypeAny    = &PrimitiveType{name: "any"}
)

// ArrayType represents array types
type ArrayType struct {
	ElementType Type
}

func (t *ArrayType) TypeName() string {
	return "array"
}

func (t *ArrayType) String() string {
	return "[" + t.ElementType.String() + "]"
}

// FunctionType represents function types
type FunctionType struct {
	ParamTypes []Type
	ReturnType Type
}

func (t *FunctionType) TypeName() string {
	return "function"
}

func (t *FunctionType) String() string {
	params := ""
	for i, pt := range t.ParamTypes {
		if i > 0 {
			params += ", "
		}
		params += pt.String()
	}
	return "(" + params + ") -> " + t.ReturnType.String()
}

// StructType represents struct types
type StructType struct {
	Name   string
	Fields map[string]Type
}

func (t *StructType) TypeName() string {
	return "struct"
}

func (t *StructType) String() string {
	return t.Name
}

// EnumType represents enum types
type EnumType struct {
	Name   string
	Values []string
}

func (t *EnumType) TypeName() string {
	return "enum"
}

func (t *EnumType) String() string {
	return t.Name
}

// UnionType represents union types (future)
type UnionType struct {
	Types []Type
}

func (t *UnionType) TypeName() string {
	return "union"
}

func (t *UnionType) String() string {
	result := ""
	for i, ty := range t.Types {
		if i > 0 {
			result += " | "
		}
		result += ty.String()
	}
	return result
}

// GenericType represents generic types (future)
type GenericType struct {
	Name       string
	TypeParams []string
	BaseType   Type
}

func (t *GenericType) TypeName() string {
	return t.Name
}

func (t *GenericType) String() string {
	return t.Name + "<...>"
}

// Helper function to create types from strings
func ParseType(typeStr string) Type {
	switch typeStr {
	case "i64":
		return TypeI64
	case "f64":
		return TypeF64
	case "string":
		return TypeString
	case "bool":
		return TypeBool
	case "void":
		return TypeVoid
	case "any", "":
		return TypeAny
	default:
		// Could be a custom type (struct, enum, etc.)
		return &PrimitiveType{name: typeStr}
	}
}

// IsAssignableTo checks if type from can be assigned to type to
func IsAssignableTo(from, to Type) bool {
	if from == nil || to == nil {
		return true // Unknown types are assignable
	}

	// Same type is always assignable
	if from.TypeName() == to.TypeName() {
		return true
	}

	// Numeric type coercion
	if (from.TypeName() == "i64" || from.TypeName() == "f64") &&
		(to.TypeName() == "i64" || to.TypeName() == "f64") {
		return true
	}

	// Any type is assignable to anything and vice versa
	if from.TypeName() == "any" || to.TypeName() == "any" {
		return true
	}

	return false
}
