package codegen

import (
	"fmt"
	"strings"

	"github.com/freelang-ai/gofree/internal/compiler"
)

// Codegen generates native code from IR
type Codegen struct {
	program *compiler.IRProgram
	output  strings.Builder
	indent  int
	errors  []CodegenError
}

// CodegenError represents a codegen error
type CodegenError struct {
	Message string
}

// NewCodegen creates a new code generator
func NewCodegen(program *compiler.IRProgram) *Codegen {
	return &Codegen{
		program: program,
		indent:  0,
		errors:  []CodegenError{},
	}
}

// GenerateC generates C code from IR
func (cg *Codegen) GenerateC() (string, error) {
	if cg.program == nil {
		return "", fmt.Errorf("program is nil")
	}

	cg.emitHeader()
	cg.emitMain()
	cg.emitFunctions()
	cg.emitFooter()

	if len(cg.errors) > 0 {
		return "", fmt.Errorf("%d codegen errors", len(cg.errors))
	}

	return cg.output.String(), nil
}

// emitHeader emits C header includes
func (cg *Codegen) emitHeader() {
	cg.emit("#include <stdio.h>")
	cg.emit("#include <stdlib.h>")
	cg.emit("#include <string.h>")
	cg.emit("#include <math.h>")
	cg.emit("")
	cg.emit("// FreeLang compiled to C")
	cg.emit("")
	cg.emit("typedef union {")
	cg.indent++
	cg.emit("double number;")
	cg.emit("char* string;")
	cg.emit("int boolean;")
	cg.emit("void* object;")
	cg.indent--
	cg.emit("} Value;")
	cg.emit("")
}

// emitMain emits main function
func (cg *Codegen) emitMain() {
	cg.emit("int main() {")
	cg.indent++

	// Emit main code
	for _, ins := range cg.program.Instructions {
		cg.emitInstruction(ins)
	}

	cg.emit("return 0;")
	cg.indent--
	cg.emit("}")
	cg.emit("")
}

// emitFunctions emits function definitions
func (cg *Codegen) emitFunctions() {
	for _, fn := range cg.program.Functions {
		cg.emitFunction(fn)
	}
}

// emitFunction emits a single function
func (cg *Codegen) emitFunction(fn *compiler.IRFunction) {
	cg.emit(fmt.Sprintf("Value %s(", fn.Name))
	cg.indent++
	for i := 0; i < fn.ParamCount; i++ {
		if i > 0 {
			cg.output.WriteString(", ")
		}
		cg.output.WriteString(fmt.Sprintf("Value param%d", i))
	}
	cg.indent--
	cg.emit(") {")

	cg.indent++
	for _, ins := range fn.Instructions {
		cg.emitInstruction(ins)
	}
	cg.emit("Value result;")
	cg.emit("result.number = 0;")
	cg.emit("return result;")
	cg.indent--

	cg.emit("}")
	cg.emit("")
}

// emitInstruction emits a single instruction
func (cg *Codegen) emitInstruction(ins *compiler.Instruction) {
	switch ins.Opcode {
	case compiler.OpLoadConst:
		if len(ins.Args) > 0 {
			cg.emit(fmt.Sprintf("// LOAD_CONST %v", ins.Args[0]))
		}

	case compiler.OpAdd:
		cg.emit("// ADD")

	case compiler.OpSub:
		cg.emit("// SUB")

	case compiler.OpMul:
		cg.emit("// MUL")

	case compiler.OpDiv:
		cg.emit("// DIV")

	case compiler.OpCall:
		if len(ins.Args) >= 1 {
			fn := ins.Args[0].(string)
			cg.emit(fmt.Sprintf("%s();", fn))
		}

	case compiler.OpReturn:
		cg.emit("return result;")

	case compiler.OpPop:
		cg.emit("// POP")

	default:
		cg.emit(fmt.Sprintf("// %s", ins.Opcode.String()))
	}
}

// emitFooter emits C footer
func (cg *Codegen) emitFooter() {
	// Footer if needed
}

// Helper methods

func (cg *Codegen) emit(code string) {
	if code == "" {
		cg.output.WriteString("\n")
		return
	}

	// Add indentation
	for i := 0; i < cg.indent; i++ {
		cg.output.WriteString("    ")
	}

	cg.output.WriteString(code)
	cg.output.WriteString("\n")
}

// GenerateWASM generates WASM (stub)
func (cg *Codegen) GenerateWASM() ([]byte, error) {
	// TODO: Implement WASM generation
	return nil, fmt.Errorf("WASM generation not implemented")
}

// GenerateJavaScript generates JavaScript code from IR
func (cg *Codegen) GenerateJavaScript() (string, error) {
	var js strings.Builder

	js.WriteString("// Generated JavaScript from FreeLang IR\n\n")

	// Generate variable declarations
	for name := range cg.program.GlobalVars {
		js.WriteString(fmt.Sprintf("let %s = 0;\n", name))
	}

	js.WriteString("\n")

	// Generate main code
	for _, ins := range cg.program.Instructions {
		cg.emitJSInstruction(ins, &js)
	}

	// Generate functions
	for _, fn := range cg.program.Functions {
		cg.emitJSFunction(fn, &js)
	}

	return js.String(), nil
}

// emitJSInstruction emits a JavaScript instruction
func (cg *Codegen) emitJSInstruction(ins *compiler.Instruction, js *strings.Builder) {
	switch ins.Opcode {
	case compiler.OpLoadConst:
		if len(ins.Args) > 0 {
			js.WriteString(fmt.Sprintf("// %v\n", ins.Args[0]))
		}

	case compiler.OpCall:
		if len(ins.Args) >= 1 {
			fn := ins.Args[0].(string)
			js.WriteString(fmt.Sprintf("%s();\n", fn))
		}

	default:
		js.WriteString(fmt.Sprintf("// %s\n", ins.Opcode.String()))
	}
}

// emitJSFunction emits a JavaScript function
func (cg *Codegen) emitJSFunction(fn *compiler.IRFunction, js *strings.Builder) {
	js.WriteString(fmt.Sprintf("function %s(", fn.Name))
	for i := 0; i < fn.ParamCount; i++ {
		if i > 0 {
			js.WriteString(", ")
		}
		js.WriteString(fmt.Sprintf("param%d", i))
	}
	js.WriteString(") {\n")

	for _, ins := range fn.Instructions {
		switch ins.Opcode {
		case compiler.OpReturn:
			js.WriteString("    return;\n")
		default:
			js.WriteString("    // Instruction\n")
		}
	}

	js.WriteString("}\n\n")
}

// Errors returns codegen errors
func (cg *Codegen) Errors() []CodegenError {
	return cg.errors
}

func (cg *Codegen) addError(message string) {
	cg.errors = append(cg.errors, CodegenError{Message: message})
}
