"""
    Codegen: Typed SSA IR → C 코드

Typed SSA IR을 받아 C 코드를 생성합니다.
"""

include("type_inference.jl")

# ======================== C 코드 생성기 ========================

mutable struct CodeGenerator
    code::Vector{String}
    var_map::Dict{Int, String}  # 변수 ID → C 변수명
    var_counter::Int
    indent_level::Int
end

function CodeGenerator()
    return CodeGenerator(String[], Dict(), 0, 0)
end

# ======================== 헬퍼 함수 ========================

"""들여쓰기 추가"""
function indent(gen::CodeGenerator)::String
    return repeat("    ", gen.indent_level)
end

"""코드 라인 추가"""
function emit_line!(gen::CodeGenerator, line::String)
    push!(gen.code, indent(gen) * line)
end

"""C 변수명 생성 (v0, v1, ...)"""
function get_var_name(gen::CodeGenerator, id::Int)::String
    if !haskey(gen.var_map, id)
        gen.var_counter += 1
        gen.var_map[id] = "v$(gen.var_counter)"
    end
    return gen.var_map[id]
end

"""타입을 C 타입으로 변환"""
function julia_type_to_c_type(typ::JuliaType)::String
    if isa(typ, IntType)
        return "int64_t"
    elseif isa(typ, FloatType)
        return "double"
    elseif isa(typ, StringType)
        return "const char*"
    elseif isa(typ, BoolType)
        return "bool"
    elseif isa(typ, NothingType)
        return "void"
    elseif isa(typ, ArrayType)
        elem_type = julia_type_to_c_type(typ.element_type)
        return "$(elem_type)*"  # 포인터로 표현
    else
        return "void*"
    end
end

"""값을 C 코드로 변환"""
function value_to_c_literal(value::Any)::String
    if isa(value, Int) || isa(value, Int64)
        return "$(value)LL"
    elseif isa(value, Float64)
        return "$(value)"
    elseif isa(value, String)
        # 문자열 이스케이프
        escaped = replace(value, "\"" => "\\\"", "\n" => "\\n")
        return "\"$(escaped)\""
    elseif isa(value, Bool)
        return value ? "true" : "false"
    elseif value === nothing
        return "NULL"
    else
        return "NULL"
    end
end

# ======================== 코드 생성 함수 ========================

"""Typed IR을 C 코드로 변환"""
function generate_code(instructions::Vector{TypedIRInstruction})::String
    gen = CodeGenerator()

    # C 헤더
    push!(gen.code, "#include <stdio.h>")
    push!(gen.code, "#include <stdint.h>")
    push!(gen.code, "#include <stdbool.h>")
    push!(gen.code, "")

    # 변수 선언
    var_declarations = Dict{String, String}()
    for instr in instructions
        if isa(instr, TypedIRLiteral) || isa(instr, TypedIRBinaryOp) ||
           isa(instr, TypedIRUnaryOp) || isa(instr, TypedIRCall)
            var_name = get_var_name(gen, instr.id)
            c_type = julia_type_to_c_type(instr.result_type)
            var_declarations[var_name] = c_type
        elseif isa(instr, TypedIRAssign)
            var_name = instr.target
            c_type = julia_type_to_c_type(instr.value_type)
            var_declarations[var_name] = c_type
        end
    end

    # 메인 함수
    push!(gen.code, "int main() {")
    gen.indent_level = 1

    # 변수 선언들
    for (var_name, c_type) in var_declarations
        emit_line!(gen, "$(c_type) $(var_name);")
    end
    if !isempty(var_declarations)
        push!(gen.code, "")
    end

    # 명령들 생성
    for (idx, instr) in enumerate(instructions)
        generate_instruction!(gen, instr)
    end

    gen.indent_level = 0
    emit_line!(gen, "return 0;")
    push!(gen.code, "}")

    return join(gen.code, "\n")
end

"""개별 명령 생성"""
function generate_instruction!(gen::CodeGenerator, instr::TypedIRInstruction)
    if isa(instr, TypedIRLiteral)
        var_name = get_var_name(gen, instr.id)
        lit_value = value_to_c_literal(instr.value)
        emit_line!(gen, "$(var_name) = $(lit_value);")

    elseif isa(instr, TypedIRBinaryOp)
        var_name = get_var_name(gen, instr.id)
        left_name = get_var_name(gen, instr.left_id)
        right_name = get_var_name(gen, instr.right_id)

        # C 연산자로 변환
        c_op = instr.op
        if c_op == "=="
            c_op = "=="
        elseif c_op == "!="
            c_op = "!="
        elseif c_op == "<"
            c_op = "<"
        elseif c_op == ">"
            c_op = ">"
        elseif c_op == "<="
            c_op = "<="
        elseif c_op == ">="
            c_op = ">="
        elseif c_op == "index"
            emit_line!(gen, "$(var_name) = $(left_name)[(int)$(right_name)];")
            return
        end

        emit_line!(gen, "$(var_name) = $(left_name) $(c_op) $(right_name);")

    elseif isa(instr, TypedIRUnaryOp)
        var_name = get_var_name(gen, instr.id)
        operand_name = get_var_name(gen, instr.operand_id)

        c_op = instr.op
        if c_op == "!"
            c_op = "!"
        end

        emit_line!(gen, "$(var_name) = $(c_op)$(operand_name);")

    elseif isa(instr, TypedIRCall)
        var_name = get_var_name(gen, instr.id)
        func_name = get_var_name(gen, instr.func_id)
        args = join([get_var_name(gen, id) for id in instr.arg_ids], ", ")

        emit_line!(gen, "$(var_name) = $(func_name)($(args));")

    elseif isa(instr, TypedIRAssign)
        var_name = instr.target
        value_name = get_var_name(gen, instr.value_id)
        emit_line!(gen, "$(var_name) = $(value_name);")

    elseif isa(instr, TypedIRLabel)
        push!(gen.code, "$(instr.name):")

    elseif isa(instr, TypedIRGoto)
        emit_line!(gen, "goto $(instr.target);")

    elseif isa(instr, TypedIRConditionalBranch)
        cond_name = get_var_name(gen, instr.condition_id)
        emit_line!(gen, "if ($(cond_name)) { goto $(instr.true_label); } else { goto $(instr.false_label); }")

    elseif isa(instr, TypedIRReturn)
        if instr.value_id !== nothing
            value_name = get_var_name(gen, instr.value_id)
            emit_line!(gen, "return $(value_name);")
        else
            emit_line!(gen, "return 0;")
        end
    end
end

# ======================== 최적화 ========================

"""Typed IR 최적화"""
function optimize_ir(instructions::Vector{TypedIRInstruction})::Vector{TypedIRInstruction}
    optimized = TypedIRInstruction[]

    # 상수 전파 (Constant Propagation)
    constants = Dict{Int, Any}()

    for instr in instructions
        if isa(instr, TypedIRLiteral)
            constants[instr.id] = instr.value
            push!(optimized, instr)
        elseif isa(instr, TypedIRBinaryOp)
            # 상수 접기 (Constant Folding)
            if haskey(constants, instr.left_id) && haskey(constants, instr.right_id)
                left_val = constants[instr.left_id]
                right_val = constants[instr.right_id]

                try
                    if instr.op == "+"
                        result = left_val + right_val
                    elseif instr.op == "-"
                        result = left_val - right_val
                    elseif instr.op == "*"
                        result = left_val * right_val
                    elseif instr.op == "/"
                        result = left_val / right_val
                    else
                        result = nothing
                    end

                    if result !== nothing
                        constants[instr.id] = result
                        push!(optimized, TypedIRLiteral(result, instr.result_type, instr.id))
                        continue
                    end
                catch
                    # 계산 실패, 원래 명령 유지
                end
            end
            push!(optimized, instr)
        else
            push!(optimized, instr)
        end
    end

    return optimized
end

# ======================== 컴파일 파이프라인 ========================

"""완전한 컴파일 파이프라인"""
function compile_to_c(source::String)::String
    # Step 1: Parse
    ast = parse(source)

    # Step 2: Lowering
    ir = lower(ast)

    # Step 3: Type Inference
    typed_ir = infer_types(ir)

    # Step 4: Optimization
    optimized_ir = optimize_ir(typed_ir)

    # Step 5: Codegen
    c_code = generate_code(optimized_ir)

    return c_code
end

# Export
export CodeGenerator, generate_code, optimize_ir, compile_to_c
export julia_type_to_c_type, value_to_c_literal
