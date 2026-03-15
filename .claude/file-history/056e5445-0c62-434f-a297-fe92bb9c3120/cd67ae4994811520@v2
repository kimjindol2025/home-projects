"""
    Type Inference: Untyped IR → Typed SSA IR

모든 변수의 타입을 추론하고 Typed SSA IR로 변환합니다.
"""

include("lowering.jl")

# ======================== 타입 정의 ========================

abstract type JuliaType end

struct IntType <: JuliaType end
struct FloatType <: JuliaType end
struct StringType <: JuliaType end
struct BoolType <: JuliaType end
struct NothingType <: JuliaType end
struct FunctionType <: JuliaType
    arg_types::Vector{JuliaType}
    return_type::JuliaType
end
struct ArrayType <: JuliaType
    element_type::JuliaType
end
struct UnknownType <: JuliaType end

# ======================== Typed IR 정의 ========================

abstract type TypedIRInstruction end

struct TypedIRLiteral <: TypedIRInstruction
    value::Any
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
    arg_types::Vector{JuliaType}
    return_type::JuliaType
    id::Int
end

struct TypedIRAssign <: TypedIRInstruction
    target::String
    value_id::Int
    value_type::JuliaType
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

struct TypedIRPhi <: TypedIRInstruction
    predecessors::Dict{String, Int}
    phi_type::JuliaType
    id::Int
end

# ======================== 타입 추론 엔진 ========================

mutable struct TypeInferencer
    typed_instructions::Vector{TypedIRInstruction}
    type_map::Dict{Int, JuliaType}  # 변수 ID → 타입
    var_types::Dict{String, JuliaType}  # 변수명 → 타입
    built_in_functions::Dict{String, FunctionType}
    counter::Int
end

function TypeInferencer()
    # 내장 함수들
    builtins = Dict(
        "+" => FunctionType([IntType(), IntType()], IntType()),
        "-" => FunctionType([IntType(), IntType()], IntType()),
        "*" => FunctionType([IntType(), IntType()], IntType()),
        "/" => FunctionType([IntType(), IntType()], FloatType()),
        "==" => FunctionType([IntType(), IntType()], BoolType()),
        "!=" => FunctionType([IntType(), IntType()], BoolType()),
        "<" => FunctionType([IntType(), IntType()], BoolType()),
        ">" => FunctionType([IntType(), IntType()], BoolType()),
        "println" => FunctionType([UnknownType()], NothingType()),
    )
    return TypeInferencer(TypedIRInstruction[], Dict(), Dict(), builtins, 0)
end

# ======================== 헬퍼 함수 ========================

"""타입 기록"""
function set_type!(inf::TypeInferencer, id::Int, typ::JuliaType)
    inf.type_map[id] = typ
end

"""타입 조회"""
function get_type(inf::TypeInferencer, id::Int)::JuliaType
    return get(inf.type_map, id, UnknownType())
end

"""Typed 명령 추가"""
function emit_typed!(inf::TypeInferencer, instr::TypedIRInstruction)
    push!(inf.typed_instructions, instr)
end

"""리터럴 값에서 타입 추론"""
function infer_literal_type(value::Any)::JuliaType
    if isa(value, Int) || isa(value, Int64)
        return IntType()
    elseif isa(value, Float64)
        return FloatType()
    elseif isa(value, String)
        return StringType()
    elseif isa(value, Bool)
        return BoolType()
    elseif value === nothing
        return NothingType()
    elseif isa(value, Vector)
        if isempty(value)
            return ArrayType(UnknownType())
        else
            elem_type = infer_literal_type(value[1])
            return ArrayType(elem_type)
        end
    else
        return UnknownType()
    end
end

"""두 타입의 연산 결과 타입 추론"""
function infer_binary_result_type(op::String, left_type::JuliaType, right_type::JuliaType)::JuliaType
    if op in ("+", "-", "*")
        if isa(left_type, IntType) && isa(right_type, IntType)
            return IntType()
        elseif isa(left_type, FloatType) || isa(right_type, FloatType)
            return FloatType()
        end
    elseif op == "/"
        return FloatType()
    elseif op in ("%", "//")
        return IntType()
    elseif op in ("==", "!=", "<", "<=", ">", ">=")
        return BoolType()
    elseif op in ("&&", "||")
        return BoolType()
    elseif op == "index"
        # 배열 인덱싱
        if isa(left_type, ArrayType)
            return left_type.element_type
        end
    end
    return UnknownType()
end

"""단항 연산 결과 타입 추론"""
function infer_unary_result_type(op::String, operand_type::JuliaType)::JuliaType
    if op == "-"
        if isa(operand_type, IntType) || isa(operand_type, FloatType)
            return operand_type
        end
    elseif op == "!"
        return BoolType()
    end
    return UnknownType()
end

# ======================== 타입 추론 함수 ========================

"""Untyped IR을 Typed SSA IR로 변환"""
function infer_types(instructions::Vector{IRInstruction})::Vector{TypedIRInstruction}
    inf = TypeInferencer()

    for instr in instructions
        infer_instruction!(inf, instr)
    end

    return inf.typed_instructions
end

"""개별 명령의 타입 추론"""
function infer_instruction!(inf::TypeInferencer, instr::IRInstruction)
    if isa(instr, IRLiteral)
        typ = infer_literal_type(instr.value)
        set_type!(inf, instr.id, typ)
        emit_typed!(inf, TypedIRLiteral(instr.value, typ, instr.id))

    elseif isa(instr, IRVar)
        if haskey(inf.var_types, instr.name)
            typ = inf.var_types[instr.name]
        else
            typ = UnknownType()
        end
        set_type!(inf, instr.id, typ)
        emit_typed!(inf, TypedIRLiteral(nothing, typ, instr.id))

    elseif isa(instr, IRBinaryOp)
        left_type = get_type(inf, instr.left_id)
        right_type = get_type(inf, instr.right_id)
        result_type = infer_binary_result_type(instr.op, left_type, right_type)
        set_type!(inf, instr.id, result_type)
        emit_typed!(inf, TypedIRBinaryOp(instr.op, instr.left_id, instr.right_id, result_type, instr.id))

    elseif isa(instr, IRUnaryOp)
        operand_type = get_type(inf, instr.operand_id)
        result_type = infer_unary_result_type(instr.op, operand_type)
        set_type!(inf, instr.id, result_type)
        emit_typed!(inf, TypedIRUnaryOp(instr.op, instr.operand_id, result_type, instr.id))

    elseif isa(instr, IRCall)
        # 함수 타입 추론 (단순화)
        arg_types = [get_type(inf, id) for id in instr.arg_ids]
        return_type = NothingType()

        # 내장 함수 확인
        if haskey(inf.built_in_functions, "func")
            func_type = inf.built_in_functions["func"]
            return_type = func_type.return_type
        end

        set_type!(inf, instr.id, return_type)
        emit_typed!(inf, TypedIRCall(instr.func_id, instr.arg_ids, arg_types, return_type, instr.id))

    elseif isa(instr, IRAssign)
        value_type = get_type(inf, instr.value_id)
        inf.var_types[instr.target] = value_type
        set_type!(inf, instr.id, NothingType())
        emit_typed!(inf, TypedIRAssign(instr.target, instr.value_id, value_type, instr.id))

    elseif isa(instr, IRLabel)
        emit_typed!(inf, TypedIRLabel(instr.name, instr.id))

    elseif isa(instr, IRGoto)
        emit_typed!(inf, TypedIRGoto(instr.target, instr.id))

    elseif isa(instr, IRConditionalBranch)
        emit_typed!(inf, TypedIRConditionalBranch(instr.condition_id, instr.true_label, instr.false_label, instr.id))

    elseif isa(instr, IRReturn)
        if instr.value_id !== nothing
            return_type = get_type(inf, instr.value_id)
        else
            return_type = NothingType()
        end
        emit_typed!(inf, TypedIRReturn(instr.value_id, return_type, instr.id))

    else
        # 미처리 명령
    end
end

# ======================== 타입 출력 ========================

"""타입을 문자열로 변환"""
function type_to_string(typ::JuliaType)::String
    if isa(typ, IntType)
        return "Int"
    elseif isa(typ, FloatType)
        return "Float64"
    elseif isa(typ, StringType)
        return "String"
    elseif isa(typ, BoolType)
        return "Bool"
    elseif isa(typ, NothingType)
        return "Nothing"
    elseif isa(typ, FunctionType)
        args = join([type_to_string(t) for t in typ.arg_types], ", ")
        return "$(args) → $(type_to_string(typ.return_type))"
    elseif isa(typ, ArrayType)
        return "Array{$(type_to_string(typ.element_type))}"
    else
        return "Unknown"
    end
end

"""Typed IR 명령을 문자열로 변환"""
function typed_ir_to_string(instr::TypedIRInstruction)::String
    if isa(instr, TypedIRLiteral)
        return "v$(instr.id)::$(type_to_string(instr.type)) = literal($(instr.value))"
    elseif isa(instr, TypedIRBinaryOp)
        return "v$(instr.id)::$(type_to_string(instr.result_type)) = v$(instr.left_id) $(instr.op) v$(instr.right_id)"
    elseif isa(instr, TypedIRUnaryOp)
        return "v$(instr.id)::$(type_to_string(instr.result_type)) = $(instr.op) v$(instr.operand_id)"
    elseif isa(instr, TypedIRCall)
        args = join(["v$(id)" for id in instr.arg_ids], ", ")
        return "v$(instr.id)::$(type_to_string(instr.return_type)) = call v$(instr.func_id)($args)"
    elseif isa(instr, TypedIRAssign)
        return "$(instr.target)::$(type_to_string(instr.value_type)) = v$(instr.value_id)"
    elseif isa(instr, TypedIRLabel)
        return "$(instr.name):"
    elseif isa(instr, TypedIRGoto)
        return "goto $(instr.target)"
    elseif isa(instr, TypedIRConditionalBranch)
        return "if v$(instr.condition_id) then $(instr.true_label) else $(instr.false_label)"
    elseif isa(instr, TypedIRReturn)
        if instr.value_id !== nothing
            return "return v$(instr.value_id)::$(type_to_string(instr.return_type))"
        else
            return "return::$(type_to_string(instr.return_type))"
        end
    else
        return "$(typeof(instr))"
    end
end

"""모든 Typed IR 명령 출력"""
function print_typed_ir(instructions::Vector{TypedIRInstruction})
    for (idx, instr) in enumerate(instructions)
        println("$idx: $(typed_ir_to_string(instr))")
    end
end

# Export
export JuliaType, IntType, FloatType, StringType, BoolType, NothingType
export FunctionType, ArrayType, UnknownType
export TypedIRInstruction, TypedIRLiteral, TypedIRBinaryOp, TypedIRUnaryOp
export TypedIRCall, TypedIRAssign, TypedIRLabel, TypedIRGoto
export TypedIRConditionalBranch, TypedIRReturn, TypedIRPhi
export infer_types, print_typed_ir, type_to_string, typed_ir_to_string
