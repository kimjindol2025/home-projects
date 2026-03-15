"""
    Lowering: AST → Untyped IR (제어흐름 정규화)

AST를 받아 goto 기반의 정규화된 IR로 변환합니다.
"""

include("parser.jl")

# ======================== Untyped IR 정의 ========================

abstract type IRInstruction end

# 기본 연산
struct IRLiteral <: IRInstruction
    value::Any
    id::Int
end

struct IRVar <: IRInstruction
    name::String
    id::Int
end

struct IRBinaryOp <: IRInstruction
    op::String
    left_id::Int
    right_id::Int
    id::Int
end

struct IRUnaryOp <: IRInstruction
    op::String
    operand_id::Int
    id::Int
end

struct IRCall <: IRInstruction
    func_id::Int
    arg_ids::Vector{Int}
    id::Int
end

struct IRLoad <: IRInstruction
    ptr_id::Int
    id::Int
end

struct IRStore <: IRInstruction
    ptr_id::Int
    value_id::Int
    id::Int
end

struct IRAssign <: IRInstruction
    target::String
    value_id::Int
    id::Int
end

# 제어흐름
struct IRLabel <: IRInstruction
    name::String
    id::Int
end

struct IRGoto <: IRInstruction
    target::String
    id::Int
end

struct IRConditionalBranch <: IRInstruction
    condition_id::Int
    true_label::String
    false_label::String
    id::Int
end

struct IRReturn <: IRInstruction
    value_id::Union{Nothing, Int}
    id::Int
end

struct IRPhi <: IRInstruction
    predecessors::Dict{String, Int}  # 라벨 → 값 ID
    id::Int
end

# ======================== IR 제너레이터 ========================

mutable struct IRGenerator
    instructions::Vector{IRInstruction}
    var_count::Int
    label_count::Int
    env::Dict{String, Int}  # 변수명 → 최신 값 ID
end

function IRGenerator()
    return IRGenerator(IRInstruction[], 0, 0, Dict())
end

# ======================== 헬퍼 함수 ========================

"""새로운 변수 ID 생성"""
function new_var(gen::IRGenerator)::Int
    gen.var_count += 1
    return gen.var_count
end

"""새로운 라벨 생성"""
function new_label(gen::IRGenerator)::String
    gen.label_count += 1
    return "label_$(gen.label_count)"
end

"""명령 추가"""
function emit!(gen::IRGenerator, instr::IRInstruction)
    push!(gen.instructions, instr)
end

"""변수를 환경에 등록"""
function bind!(gen::IRGenerator, name::String, id::Int)
    gen.env[name] = id
end

"""변수의 최신 값 ID 조회"""
function lookup(gen::IRGenerator, name::String)::Int
    return get(gen.env, name, -1)
end

# ======================== Lowering 함수 ========================

"""AST 노드를 Untyped IR로 변환"""
function lower(nodes::Vector{ASTNode})::Vector{IRInstruction}
    gen = IRGenerator()

    for node in nodes
        lower_statement!(gen, node)
    end

    return gen.instructions
end

"""표현식을 Untyped IR로 변환, 결과 ID 반환"""
function lower_expression!(gen::IRGenerator, expr::Expression)::Int
    if isa(expr, IntegerLiteral)
        var_id = new_var(gen)
        emit!(gen, IRLiteral(expr.value, var_id))
        return var_id
    elseif isa(expr, FloatLiteral)
        var_id = new_var(gen)
        emit!(gen, IRLiteral(expr.value, var_id))
        return var_id
    elseif isa(expr, StringLiteral)
        var_id = new_var(gen)
        emit!(gen, IRLiteral(expr.value, var_id))
        return var_id
    elseif isa(expr, BooleanLiteral)
        var_id = new_var(gen)
        emit!(gen, IRLiteral(expr.value, var_id))
        return var_id
    elseif isa(expr, NothingLiteral)
        var_id = new_var(gen)
        emit!(gen, IRLiteral(nothing, var_id))
        return var_id
    elseif isa(expr, Identifier)
        # 변수 조회
        id = lookup(gen, expr.name)
        if id == -1
            # 아직 정의되지 않은 변수, 새로 생성
            id = new_var(gen)
            emit!(gen, IRVar(expr.name, id))
            bind!(gen, expr.name, id)
        end
        return id
    elseif isa(expr, BinaryOp)
        left_id = lower_expression!(gen, expr.left)
        right_id = lower_expression!(gen, expr.right)
        var_id = new_var(gen)
        emit!(gen, IRBinaryOp(expr.op, left_id, right_id, var_id))
        return var_id
    elseif isa(expr, UnaryOp)
        operand_id = lower_expression!(gen, expr.operand)
        var_id = new_var(gen)
        emit!(gen, IRUnaryOp(expr.op, operand_id, var_id))
        return var_id
    elseif isa(expr, Call)
        func_id = lower_expression!(gen, expr.func)
        arg_ids = [lower_expression!(gen, arg) for arg in expr.args]
        var_id = new_var(gen)
        emit!(gen, IRCall(func_id, arg_ids, var_id))
        return var_id
    elseif isa(expr, IndexAccess)
        obj_id = lower_expression!(gen, expr.object)
        idx_ids = [lower_expression!(gen, idx) for idx in expr.indices]
        var_id = new_var(gen)
        # 단순화: 첫 번째 인덱스만 사용
        emit!(gen, IRBinaryOp("index", obj_id, idx_ids[1], var_id))
        return var_id
    elseif isa(expr, MemberAccess)
        obj_id = lower_expression!(gen, expr.object)
        var_id = new_var(gen)
        emit!(gen, IRVar("member_$(expr.field)", var_id))
        return var_id
    elseif isa(expr, ArrayLiteral)
        elem_ids = [lower_expression!(gen, elem) for elem in expr.elements]
        var_id = new_var(gen)
        # 단순화: 배열 리터럴
        emit!(gen, IRLiteral(elem_ids, var_id))
        return var_id
    else
        error("Unknown expression type: $(typeof(expr))")
    end
end

"""문을 Untyped IR로 변환"""
function lower_statement!(gen::IRGenerator, stmt::ASTNode)
    if isa(stmt, Assignment)
        value_id = lower_expression!(gen, stmt.value)
        emit!(gen, IRAssign(stmt.target, value_id, new_var(gen)))
        bind!(gen, stmt.target, value_id)
    elseif isa(stmt, FunctionDef)
        # 함수 라벨
        func_label = "func_$(stmt.name)"
        emit!(gen, IRLabel(func_label, new_var(gen)))

        # 파라미터 바인딩
        for param in stmt.params
            param_id = new_var(gen)
            emit!(gen, IRVar(param, param_id))
            bind!(gen, param, param_id)
        end

        # 본체
        for body_stmt in stmt.body
            lower_statement!(gen, body_stmt)
        end

        # 암묵적 return
        emit!(gen, IRReturn(nothing, new_var(gen)))
    elseif isa(stmt, IfStmt)
        cond_id = lower_expression!(gen, stmt.condition)
        then_label = new_label(gen)
        else_label = new_label(gen)
        end_label = new_label(gen)

        # 조건부 분기
        emit!(gen, IRConditionalBranch(cond_id, then_label, else_label, new_var(gen)))

        # Then 블록
        emit!(gen, IRLabel(then_label, new_var(gen)))
        for body_stmt in stmt.then_body
            lower_statement!(gen, body_stmt)
        end
        emit!(gen, IRGoto(end_label, new_var(gen)))

        # Else 블록
        emit!(gen, IRLabel(else_label, new_var(gen)))
        if stmt.else_body !== nothing
            for body_stmt in stmt.else_body
                lower_statement!(gen, body_stmt)
            end
        end
        emit!(gen, IRGoto(end_label, new_var(gen)))

        # End 라벨
        emit!(gen, IRLabel(end_label, new_var(gen)))
    elseif isa(stmt, WhileStmt)
        loop_label = new_label(gen)
        body_label = new_label(gen)
        exit_label = new_label(gen)

        # 루프 시작
        emit!(gen, IRLabel(loop_label, new_var(gen)))
        cond_id = lower_expression!(gen, stmt.condition)
        emit!(gen, IRConditionalBranch(cond_id, body_label, exit_label, new_var(gen)))

        # 루프 본체
        emit!(gen, IRLabel(body_label, new_var(gen)))
        for body_stmt in stmt.body
            lower_statement!(gen, body_stmt)
        end
        emit!(gen, IRGoto(loop_label, new_var(gen)))

        # 루프 종료
        emit!(gen, IRLabel(exit_label, new_var(gen)))
    elseif isa(stmt, ForStmt)
        iter_id = lower_expression!(gen, stmt.iterator)
        loop_label = new_label(gen)
        body_label = new_label(gen)
        exit_label = new_label(gen)

        # 루프 변수
        loop_var_id = new_var(gen)
        bind!(gen, stmt.variable, loop_var_id)

        # 루프 시작
        emit!(gen, IRLabel(loop_label, new_var(gen)))
        emit!(gen, IRConditionalBranch(iter_id, body_label, exit_label, new_var(gen)))

        # 루프 본체
        emit!(gen, IRLabel(body_label, new_var(gen)))
        for body_stmt in stmt.body
            lower_statement!(gen, body_stmt)
        end
        emit!(gen, IRGoto(loop_label, new_var(gen)))

        # 루프 종료
        emit!(gen, IRLabel(exit_label, new_var(gen)))
    elseif isa(stmt, ReturnStmt)
        if stmt.value !== nothing
            value_id = lower_expression!(gen, stmt.value)
            emit!(gen, IRReturn(value_id, new_var(gen)))
        else
            emit!(gen, IRReturn(nothing, new_var(gen)))
        end
    else
        # 표현식 문
        if isa(stmt, Expression)
            lower_expression!(gen, stmt)
        end
    end
end

# ======================== IR 출력 ========================

"""IR 명령을 문자열로 변환"""
function ir_to_string(instr::IRInstruction)::String
    if isa(instr, IRLiteral)
        return "v$(instr.id) = literal($(instr.value))"
    elseif isa(instr, IRVar)
        return "v$(instr.id) = var($(instr.name))"
    elseif isa(instr, IRBinaryOp)
        return "v$(instr.id) = v$(instr.left_id) $(instr.op) v$(instr.right_id)"
    elseif isa(instr, IRUnaryOp)
        return "v$(instr.id) = $(instr.op) v$(instr.operand_id)"
    elseif isa(instr, IRCall)
        args = join(["v$(id)" for id in instr.arg_ids], ", ")
        return "v$(instr.id) = call v$(instr.func_id)($args)"
    elseif isa(instr, IRAssign)
        return "$(instr.target) = v$(instr.value_id)"
    elseif isa(instr, IRLabel)
        return "$(instr.name):"
    elseif isa(instr, IRGoto)
        return "goto $(instr.target)"
    elseif isa(instr, IRConditionalBranch)
        return "if v$(instr.condition_id) then $(instr.true_label) else $(instr.false_label)"
    elseif isa(instr, IRReturn)
        if instr.value_id !== nothing
            return "return v$(instr.value_id)"
        else
            return "return"
        end
    else
        return "$(typeof(instr))"
    end
end

"""모든 IR 명령 출력"""
function print_ir(instructions::Vector{IRInstruction})
    for (idx, instr) in enumerate(instructions)
        println("$idx: $(ir_to_string(instr))")
    end
end

# Export
export IRInstruction, IRLiteral, IRVar, IRBinaryOp, IRUnaryOp, IRCall
export IRLabel, IRGoto, IRConditionalBranch, IRReturn, IRAssign
export lower, lower_expression!, lower_statement!, print_ir, ir_to_string
