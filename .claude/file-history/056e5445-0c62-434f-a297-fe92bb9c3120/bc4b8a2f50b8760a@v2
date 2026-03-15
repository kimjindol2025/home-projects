"""
    FreeLiulia Codegen: Typed SSA IR → C 코드 생성

Typed SSA IR을 받아 C 코드를 생성합니다.
"""

include("freeliulia_type_inference.jl")

# ======================== C 코드 생성 컨텍스트 ========================

mutable struct CodegenContext
    code::Vector{String}
    indentation::Int
    value_names::Dict{Int, String}
    next_var::Int
end

function CodegenContext()
    return CodegenContext(String[], 0, Dict(), 0)
end

# ======================== 헬퍼 함수 ========================

"""다음 임시 변수 이름"""
function next_var_name!(ctx::CodegenContext)::String
    ctx.next_var += 1
    return "t$(ctx.next_var)"
end

"""코드 추가"""
function emit_code!(ctx::CodegenContext, line::String)
    indent = repeat("  ", ctx.indentation)
    push!(ctx.code, indent * line)
end

"""C 타입 문자열"""
function c_type_string(t::JuliaType)::String
    if t isa IntType
        return "int64_t"
    elseif t isa FloatType
        return "double"
    elseif t isa BoolType
        return "bool"
    elseif t isa StringType
        return "char*"
    elseif t isa NothingType
        return "void"
    elseif t isa ArrayType
        return "void*"  # 일반 배열 포인터
    else
        return "void*"
    end
end

"""값 이름 가져오기"""
function get_value_name(ctx::CodegenContext, value_id::Int)::String
    if haskey(ctx.value_names, value_id)
        return ctx.value_names[value_id]
    else
        var_name = next_var_name!(ctx)
        ctx.value_names[value_id] = var_name
        return var_name
    end
end

# ======================== 코드 생성 함수 ========================

"""Typed IR에서 C 코드 생성"""
function generate_c_code(typed_instructions::Vector{TypedIRInstruction}, ctx::CodegenContext)
    # 헤더
    push!(ctx.code, "#include <stdio.h>")
    push!(ctx.code, "#include <stdbool.h>")
    push!(ctx.code, "#include <stdint.h>")
    push!(ctx.code, "")

    # 함수 선언들
    for instr in typed_instructions
        if instr isa TypedIRFunctionBegin
            emit_function_declaration(instr, typed_instructions, ctx)
        end
    end

    push!(ctx.code, "")

    # 함수 정의들
    i = 1
    while i <= length(typed_instructions)
        instr = typed_instructions[i]
        if instr isa TypedIRFunctionBegin
            i = emit_function_definition(instr, typed_instructions, i, ctx)
        else
            i += 1
        end
    end

    # Main 함수 (테스트)
    emit_code!(ctx, "int main() {")
    ctx.indentation += 1

    # 메인의 명령들
    for instr in typed_instructions
        if !(instr isa TypedIRFunctionBegin || instr isa TypedIRFunctionEnd)
            if !(instr isa TypedIRLabel || instr isa TypedIRGoto)
                emit_typed_ir_instruction(instr, ctx)
            end
        end
    end

    emit_code!(ctx, "return 0;")
    ctx.indentation -= 1
    emit_code!(ctx, "}")

    return join(ctx.code, "\n")
end

"""함수 선언 생성"""
function emit_function_declaration(func_begin::TypedIRFunctionBegin, all_instr::Vector{TypedIRInstruction}, ctx::CodegenContext)
    param_strs = String[]
    for (param_name, param_type) in zip(func_begin.params, func_begin.param_types)
        push!(param_strs, c_type_string(param_type) * " " * param_name)
    end

    param_str = join(param_strs, ", ")
    return_type = c_type_string(func_begin.return_type)

    emit_code!(ctx, "$return_type $(func_begin.name)($param_str);")
end

"""함수 정의 생성"""
function emit_function_definition(func_begin::TypedIRFunctionBegin, all_instr::Vector{TypedIRInstruction}, start_idx::Int, ctx::CodegenContext)
    param_strs = String[]
    for (param_name, param_type) in zip(func_begin.params, func_begin.param_types)
        push!(param_strs, c_type_string(param_type) * " " * param_name)
    end

    param_str = join(param_strs, ", ")
    return_type = c_type_string(func_begin.return_type)

    emit_code!(ctx, "$return_type $(func_begin.name)($param_str) {")
    ctx.indentation += 1

    # 함수 본체 명령들
    i = start_idx + 1
    while i <= length(all_instr)
        instr = all_instr[i]

        if instr isa TypedIRFunctionEnd
            break
        end

        emit_typed_ir_instruction(instr, ctx)
        i += 1
    end

    ctx.indentation -= 1
    emit_code!(ctx, "}")
    emit_code!(ctx, "")

    return i + 1
end

"""Typed IR 명령 코드 생성"""
function emit_typed_ir_instruction(instr::TypedIRInstruction, ctx::CodegenContext)
    if instr isa TypedIRLiteral
        var_name = get_value_name(ctx, instr.id)
        c_type = c_type_string(instr.type)

        if instr.type isa StringType
            emit_code!(ctx, "$c_type $var_name = \"$(instr.value)\";")
        else
            emit_code!(ctx, "$c_type $var_name = $(instr.value);")
        end

    elseif instr isa TypedIRVar
        # 변수 참조는 직접 사용

    elseif instr isa TypedIRBinaryOp
        left_var = get_value_name(ctx, instr.left_id)
        right_var = get_value_name(ctx, instr.right_id)
        result_var = get_value_name(ctx, instr.id)
        c_type = c_type_string(instr.result_type)

        c_op = translate_operator(instr.op)
        emit_code!(ctx, "$c_type $result_var = $left_var $c_op $right_var;")

    elseif instr isa TypedIRUnaryOp
        operand_var = get_value_name(ctx, instr.operand_id)
        result_var = get_value_name(ctx, instr.id)
        c_type = c_type_string(instr.result_type)

        if instr.op == "-"
            emit_code!(ctx, "$c_type $result_var = -($operand_var);")
        elseif instr.op == "!"
            emit_code!(ctx, "$c_type $result_var = !($operand_var);")
        end

    elseif instr isa TypedIRCall
        func_var = get_value_name(ctx, instr.func_id)
        arg_vars = [get_value_name(ctx, arg_id) for arg_id in instr.arg_ids]
        args_str = join(arg_vars, ", ")

        result_var = get_value_name(ctx, instr.id)
        c_type = c_type_string(instr.return_type)

        emit_code!(ctx, "$c_type $result_var = $func_var($args_str);")

    elseif instr isa TypedIRAssign
        value_var = get_value_name(ctx, instr.value_id)
        emit_code!(ctx, "$(instr.target) = $value_var;")

    elseif instr isa TypedIRLabel
        emit_code!(ctx, "$(instr.name):")

    elseif instr isa TypedIRGoto
        emit_code!(ctx, "goto $(instr.target);")

    elseif instr isa TypedIRConditionalBranch
        cond_var = get_value_name(ctx, instr.condition_id)
        emit_code!(ctx, "if ($cond_var) { goto $(instr.true_label); } else { goto $(instr.false_label); }")

    elseif instr isa TypedIRReturn
        if instr.value_id !== nothing
            value_var = get_value_name(ctx, instr.value_id)
            emit_code!(ctx, "return $value_var;")
        else
            emit_code!(ctx, "return;")
        end

    elseif instr isa TypedIRVariableDecl
        c_type = c_type_string(instr.type)
        emit_code!(ctx, "$c_type $(instr.name);")
    end
end

"""연산자 번역 (한글/Julia → C)"""
function translate_operator(op::String)::String
    dict = Dict(
        "+" => "+",
        "-" => "-",
        "*" => "*",
        "/" => "/",
        "%" => "%",
        "==" => "==",
        "!=" => "!=",
        "<" => "<",
        "<=" => "<=",
        ">" => ">",
        ">=" => ">=",
        "&&" => "&&",
        "||" => "||",
        "++" => "strcat"  # 특수: 문자열 연결
    )
    return get(dict, op, op)
end

# ======================== 공개 인터페이스 ========================

"""FreeLiulia 코드 → C 코드 생성"""
function generate_c_from_freeliulia(source::String)::String
    ast = parse_freeliulia(source)
    ir = lower_program(ast, LoweringContext())
    ctx = TypeInferenceContext()
    typed_ir = infer_types(ir, ctx)

    codegen_ctx = CodegenContext()
    return generate_c_code(typed_ir, codegen_ctx)
end

# Export
export CodegenContext, generate_c_from_freeliulia
