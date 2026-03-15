"""
    FreeLiulia Type Inference: Untyped IR → Typed SSA IR

Untyped IR을 받아 타입 추론을 수행하고 Typed SSA IR로 변환합니다.
"""

include("freeliulia_lowering.jl")

# ======================== 타입 정의 ========================

abstract type JuliaType end

struct IntType <: JuliaType end
struct FloatType <: JuliaType end
struct BoolType <: JuliaType end
struct StringType <: JuliaType end
struct NothingType <: JuliaType end
struct ArrayType <: JuliaType
    elem_type::JuliaType
end
struct FunctionType <: JuliaType
    param_types::Vector{JuliaType}
    return_type::JuliaType
end
struct UnknownType <: JuliaType end

# ======================== Typed IR 정의 ========================

abstract type TypedIRInstruction end

struct TypedIRLiteral <: TypedIRInstruction
    value::Any
    type::JuliaType
    id::Int
end

struct TypedIRVar <: TypedIRInstruction
    name::String
    type::JuliaType
    id::Int
end

struct TypedIRBinaryOp <: TypedIRInstruction
    op::String
    left_id::Int
    right_id::Int
    result_type::JuliaType
    id::Int
end

struct TypedIRUnaryOp <: TypedIRInstruction
    op::String
    operand_id::Int
    result_type::JuliaType
    id::Int
end

struct TypedIRCall <: TypedIRInstruction
    func_id::Int
    arg_ids::Vector{Int}
    return_type::JuliaType
    id::Int
end

struct TypedIRAssign <: TypedIRInstruction
    target::String
    value_id::Int
    type::JuliaType
    id::Int
end

struct TypedIRIndex <: TypedIRInstruction
    array_id::Int
    index_id::Int
    result_type::JuliaType
    id::Int
end

struct TypedIRMember <: TypedIRInstruction
    object_id::Int
    member::String
    result_type::JuliaType
    id::Int
end

struct TypedIRPhi <: TypedIRInstruction
    incoming_ids::Vector{Int}
    incoming_types::Vector{JuliaType}
    result_type::JuliaType
    id::Int
end

struct TypedIRLabel <: TypedIRInstruction
    name::String
    id::Int
end

struct TypedIRGoto <: TypedIRInstruction
    target::String
    id::Int
end

struct TypedIRConditionalBranch <: TypedIRInstruction
    condition_id::Int
    true_label::String
    false_label::String
    id::Int
end

struct TypedIRReturn <: TypedIRInstruction
    value_id::Union{Nothing, Int}
    return_type::JuliaType
    id::Int
end

struct TypedIRFunctionBegin <: TypedIRInstruction
    name::String
    params::Vector{String}
    param_types::Vector{JuliaType}
    return_type::JuliaType
    id::Int
end

struct TypedIRFunctionEnd <: TypedIRInstruction
    id::Int
end

struct TypedIRVariableDecl <: TypedIRInstruction
    name::String
    type::JuliaType
    id::Int
end

# ======================== 타입 추론 컨텍스트 ========================

mutable struct TypeInferenceContext
    typed_instructions::Vector{TypedIRInstruction}
    value_types::Dict{Int, JuliaType}
    variable_types::Dict{String, JuliaType}
    next_id::Int
    label_counter::Int
end

function TypeInferenceContext()
    return TypeInferenceContext(TypedIRInstruction[], Dict(), Dict(), 0, 0)
end

# ======================== 헬퍼 함수 ========================

"""다음 ID"""
function next_id!(ctx::TypeInferenceContext)::Int
    ctx.next_id += 1
    return ctx.next_id
end

"""타입 출력"""
function type_string(t::JuliaType)::String
    if t isa IntType
        return "정수"
    elseif t isa FloatType
        return "실수"
    elseif t isa BoolType
        return "논리값"
    elseif t isa StringType
        return "문자열"
    elseif t isa NothingType
        return "없음"
    elseif t isa ArrayType
        return "배열<$(type_string(t.elem_type))>"
    elseif t isa FunctionType
        param_str = join(type_string.(t.param_types), ", ")
        return "함수($param_str) -> $(type_string(t.return_type))"
    elseif t isa UnknownType
        return "?"
    else
        return "알수없음"
    end
end

"""리터럴 값의 타입 추론"""
function infer_literal_type(value::Any)::JuliaType
    if value isa Int || value isa Int64
        return IntType()
    elseif value isa Float64
        return FloatType()
    elseif value isa Bool
        return BoolType()
    elseif value isa String
        return StringType()
    elseif value === nothing
        return NothingType()
    else
        return UnknownType()
    end
end

"""이항 연산 결과 타입"""
function infer_binary_op_type(op::String, left_type::JuliaType, right_type::JuliaType)::JuliaType
    if op in ["+", "-", "*", "/", "%"]
        if left_type isa IntType && right_type isa IntType
            return op == "/" ? FloatType() : IntType()
        elseif left_type isa FloatType || right_type isa FloatType
            return FloatType()
        end
    elseif op in ["==", "!=", "<", "<=", ">", ">="]
        return BoolType()
    elseif op in ["&&", "||"]
        return BoolType()
    elseif op == "++"
        return StringType()
    end
    return UnknownType()
end

"""단항 연산 결과 타입"""
function infer_unary_op_type(op::String, operand_type::JuliaType)::JuliaType
    if op in ["-"]
        return operand_type
    elseif op == "!"
        return BoolType()
    end
    return UnknownType()
end

"""Typed IR 명령 추가"""
function emit_typed!(ctx::TypeInferenceContext, instr::TypedIRInstruction)
    push!(ctx.typed_instructions, instr)
    if hasproperty(instr, :id) && hasproperty(instr, :result_type)
        ctx.value_types[instr.id] = instr.result_type
    elseif hasproperty(instr, :id) && hasproperty(instr, :type)
        ctx.value_types[instr.id] = instr.type
    end
end

# ======================== 타입 추론 함수 ========================

"""프로그램 타입 추론"""
function infer_types(instructions::Vector{IRInstruction}, ctx::TypeInferenceContext)
    for instr in instructions
        infer_instruction_type(instr, ctx)
    end
    return ctx.typed_instructions
end

"""IR 명령 타입 추론"""
function infer_instruction_type(instr::IRInstruction, ctx::TypeInferenceContext)
    if instr isa IRLiteral
        literal_type = infer_literal_type(instr.value)
        emit_typed!(ctx, TypedIRLiteral(instr.value, literal_type, instr.id))
        ctx.value_types[instr.id] = literal_type

    elseif instr isa IRVar
        var_type = get(ctx.variable_types, instr.name, UnknownType())
        emit_typed!(ctx, TypedIRVar(instr.name, var_type, instr.id))
        ctx.value_types[instr.id] = var_type

    elseif instr isa IRBinaryOp
        left_type = get(ctx.value_types, instr.left_id, UnknownType())
        right_type = get(ctx.value_types, instr.right_id, UnknownType())
        result_type = infer_binary_op_type(instr.op, left_type, right_type)

        emit_typed!(ctx, TypedIRBinaryOp(instr.op, instr.left_id, instr.right_id, result_type, instr.id))
        ctx.value_types[instr.id] = result_type

    elseif instr isa IRUnaryOp
        operand_type = get(ctx.value_types, instr.operand_id, UnknownType())
        result_type = infer_unary_op_type(instr.op, operand_type)

        emit_typed!(ctx, TypedIRUnaryOp(instr.op, instr.operand_id, result_type, instr.id))
        ctx.value_types[instr.id] = result_type

    elseif instr isa IRCall
        # 함수 호출 반환 타입 (기본: Unknown)
        return_type = UnknownType()
        emit_typed!(ctx, TypedIRCall(instr.func_id, instr.arg_ids, return_type, instr.id))
        ctx.value_types[instr.id] = return_type

    elseif instr isa IRAssign
        value_type = get(ctx.value_types, instr.value_id, UnknownType())
        ctx.variable_types[instr.target] = value_type

        emit_typed!(ctx, TypedIRAssign(instr.target, instr.value_id, value_type, instr.id))

    elseif instr isa IRIndex
        array_type = get(ctx.value_types, instr.array_id, UnknownType())
        result_type = if array_type isa ArrayType
            array_type.elem_type
        else
            UnknownType()
        end

        emit_typed!(ctx, TypedIRIndex(instr.array_id, instr.index_id, result_type, instr.id))
        ctx.value_types[instr.id] = result_type

    elseif instr isa IRMember
        # 멤버 접근 타입 (기본: Unknown)
        result_type = UnknownType()
        emit_typed!(ctx, TypedIRMember(instr.object_id, instr.member, result_type, instr.id))
        ctx.value_types[instr.id] = result_type

    elseif instr isa IRLabel
        emit_typed!(ctx, TypedIRLabel(instr.name, instr.id))

    elseif instr isa IRGoto
        emit_typed!(ctx, TypedIRGoto(instr.target, instr.id))

    elseif instr isa IRConditionalBranch
        emit_typed!(ctx, TypedIRConditionalBranch(instr.condition_id, instr.true_label, instr.false_label, instr.id))

    elseif instr isa IRReturn
        return_type = if instr.value_id !== nothing
            get(ctx.value_types, instr.value_id, UnknownType())
        else
            NothingType()
        end

        emit_typed!(ctx, TypedIRReturn(instr.value_id, return_type, instr.id))

    elseif instr isa IRFunctionBegin
        # 함수 매개변수 타입 (기본: Unknown)
        param_types = fill(UnknownType(), length(instr.params))
        return_type = UnknownType()

        for param in instr.params
            ctx.variable_types[param] = UnknownType()
        end

        emit_typed!(ctx, TypedIRFunctionBegin(instr.name, instr.params, param_types, return_type, instr.id))

    elseif instr isa IRFunctionEnd
        emit_typed!(ctx, TypedIRFunctionEnd(instr.id))

    elseif instr isa IRVariableDecl
        var_type = UnknownType()
        if instr.type_name !== nothing
            # 타입 이름에서 타입 추론
            var_type = case_type_name(instr.type_name)
        end

        ctx.variable_types[instr.name] = var_type
        emit_typed!(ctx, TypedIRVariableDecl(instr.name, var_type, instr.id))
    end
end

"""한글 타입 이름을 JuliaType로 변환"""
function case_type_name(type_name::String)::JuliaType
    if type_name == "정수"
        return IntType()
    elseif type_name == "실수"
        return FloatType()
    elseif type_name == "논리값"
        return BoolType()
    elseif type_name == "문자열"
        return StringType()
    elseif type_name == "없음"
        return NothingType()
    elseif startswith(type_name, "배열")
        return ArrayType(UnknownType())
    else
        return UnknownType()
    end
end

# ======================== 공개 인터페이스 ========================

"""IR 타입 추론"""
function infer_freeliulia_types(instructions::Vector{IRInstruction})::Vector{TypedIRInstruction}
    ctx = TypeInferenceContext()
    return infer_types(instructions, ctx)
end

"""Typed IR 출력"""
function print_typed_ir(typed_instructions::Vector{TypedIRInstruction})
    println("=== Typed SSA IR ===\n")

    for (i, instr) in enumerate(typed_instructions)
        print_typed_ir_instruction(instr, i)
    end
end

"""단일 Typed IR 명령 출력"""
function print_typed_ir_instruction(instr::TypedIRInstruction, index::Int)
    if instr isa TypedIRLiteral
        println("[$index] Literal[$(type_string(instr.type))]: $(instr.value) → id=$(instr.id)")
    elseif instr isa TypedIRVar
        println("[$index] Var[$(type_string(instr.type))]: $(instr.name) → id=$(instr.id)")
    elseif instr isa TypedIRBinaryOp
        println("[$index] BinOp[$(type_string(instr.result_type))]: $(instr.left_id) $(instr.op) $(instr.right_id) → id=$(instr.id)")
    elseif instr isa TypedIRUnaryOp
        println("[$index] UnaryOp[$(type_string(instr.result_type))]: $(instr.op) $(instr.operand_id) → id=$(instr.id)")
    elseif instr isa TypedIRCall
        args_str = join(instr.arg_ids, ", ")
        println("[$index] Call[$(type_string(instr.return_type))]: func=$(instr.func_id) args=[$args_str] → id=$(instr.id)")
    elseif instr isa TypedIRAssign
        println("[$index] Assign: target = $(instr.value_id)[$(type_string(instr.type))] → id=$(instr.id)")
    elseif instr isa TypedIRIndex
        println("[$index] Index[$(type_string(instr.result_type))]: $(instr.array_id)[$(instr.index_id)] → id=$(instr.id)")
    elseif instr isa TypedIRMember
        println("[$index] Member[$(type_string(instr.result_type))]: $(instr.object_id).$(instr.member) → id=$(instr.id)")
    elseif instr isa TypedIRLabel
        println("[$index] Label: $(instr.name)")
    elseif instr isa TypedIRGoto
        println("[$index] Goto: $(instr.target)")
    elseif instr isa TypedIRConditionalBranch
        println("[$index] Branch: if $(instr.condition_id) then $(instr.true_label) else $(instr.false_label)")
    elseif instr isa TypedIRReturn
        ret_str = instr.value_id !== nothing ? "return $(instr.value_id)[$(type_string(instr.return_type))]" : "return[$(type_string(instr.return_type))]"
        println("[$index] $ret_str")
    elseif instr isa TypedIRFunctionBegin
        params_str = join(["$(instr.params[i])[$(type_string(instr.param_types[i]))]" for i in 1:length(instr.params)], ", ")
        println("[$index] FuncBegin[$(type_string(instr.return_type))]: $(instr.name)([$params_str])")
    elseif instr isa TypedIRFunctionEnd
        println("[$index] FuncEnd")
    elseif instr isa TypedIRVariableDecl
        println("[$index] VarDecl[$(type_string(instr.type))]: $(instr.name)")
    else
        println("[$index] ??? $(typeof(instr))")
    end
end

# Export
export TypeInferenceContext, infer_freeliulia_types, print_typed_ir
export JuliaType, IntType, FloatType, BoolType, StringType, NothingType, ArrayType
export TypedIRInstruction, TypedIRLiteral, TypedIRVar, TypedIRBinaryOp, TypedIRCall
