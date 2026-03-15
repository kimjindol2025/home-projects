"""
    FreeLiulia Lowering: AST → Untyped IR (제어흐름 정규화)

한글 AST를 받아 goto 기반의 정규화된 Untyped IR로 변환합니다.
"""

include("freeliulia_parser.jl")

# ======================== Untyped IR 정의 ========================

abstract type IRInstruction end

# ─── 기본 연산 ──────────────────────────────────────

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

struct IRAssign <: IRInstruction
    target::String
    value_id::Int
    id::Int
end

struct IRIndex <: IRInstruction
    array_id::Int
    index_id::Int
    id::Int
end

struct IRMember <: IRInstruction
    object_id::Int
    member::String
    id::Int
end

# ─── 제어흐름 ──────────────────────────────────────

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

# ─── 메타 명령 ──────────────────────────────────────

struct IRFunctionBegin <: IRInstruction
    name::String
    params::Vector{String}
    id::Int
end

struct IRFunctionEnd <: IRInstruction
    id::Int
end

struct IRVariableDecl <: IRInstruction
    name::String
    type_name::Union{Nothing, String}
    id::Int
end

# ======================== 로워링 컨텍스트 ========================

mutable struct LoweringContext
    instructions::Vector{IRInstruction}
    next_id::Int
    label_counter::Int
    var_counter::Int
    value_stack::Dict{Int, Any}
end

function LoweringContext()
    return LoweringContext(IRInstruction[], 0, 0, 0, Dict())
end

# ======================== 헬퍼 함수 ========================

"""다음 IR 아이디 반환"""
function next_id!(ctx::LoweringContext)::Int
    ctx.next_id += 1
    return ctx.next_id
end

"""새 라벨 생성"""
function new_label!(ctx::LoweringContext)::String
    ctx.label_counter += 1
    return "label_$(ctx.label_counter)"
end

"""IR 명령 추가"""
function emit!(ctx::LoweringContext, instr::IRInstruction)
    push!(ctx.instructions, instr)
    return instr.id
end

"""리터럴 IR 생성"""
function emit_literal!(ctx::LoweringContext, value::Any)::Int
    id = next_id!(ctx)
    emit!(ctx, IRLiteral(value, id))
    return id
end

"""변수 참조 IR 생성"""
function emit_var!(ctx::LoweringContext, name::String)::Int
    id = next_id!(ctx)
    emit!(ctx, IRVar(name, id))
    return id
end

"""라벨 추가"""
function emit_label!(ctx::LoweringContext, label::String)::Int
    id = next_id!(ctx)
    emit!(ctx, IRLabel(label, id))
    return id
end

"""조건부 분기 생성"""
function emit_cond_branch!(ctx::LoweringContext, cond_id::Int, true_label::String, false_label::String)::Int
    id = next_id!(ctx)
    emit!(ctx, IRConditionalBranch(cond_id, true_label, false_label, id))
    return id
end

"""무조건 분기"""
function emit_goto!(ctx::LoweringContext, label::String)::Int
    id = next_id!(ctx)
    emit!(ctx, IRGoto(label, id))
    return id
end

"""반환 생성"""
function emit_return!(ctx::LoweringContext, value_id::Union{Nothing, Int})::Int
    id = next_id!(ctx)
    emit!(ctx, IRReturn(value_id, id))
    return id
end

# ======================== 로워링 함수 ========================

"""프로그램 로워링"""
function lower_program(prog::Program, ctx::LoweringContext)
    for stmt in prog.statements
        lower_statement(stmt, ctx)
    end
    return ctx.instructions
end

"""문장 로워링"""
function lower_statement(stmt::ASTNode, ctx::LoweringContext)
    if stmt isa VarDecl
        lower_var_decl(stmt, ctx)
    elseif stmt isa ConstDecl
        lower_const_decl(stmt, ctx)
    elseif stmt isa FunctionDef
        lower_function_def(stmt, ctx)
    elseif stmt isa IfStmt
        lower_if_stmt(stmt, ctx)
    elseif stmt isa WhileStmt
        lower_while_stmt(stmt, ctx)
    elseif stmt isa ForStmt
        lower_for_stmt(stmt, ctx)
    elseif stmt isa ReturnStmt
        lower_return_stmt(stmt, ctx)
    elseif stmt isa ExpressionStmt
        lower_expression(stmt.expr, ctx)
    end
end

"""변수 선언 로워링"""
function lower_var_decl(var_decl::VarDecl, ctx::LoweringContext)
    id = next_id!(ctx)
    emit!(ctx, IRVariableDecl(var_decl.name, var_decl.type_name, id))

    if var_decl.value !== nothing
        value_id = lower_expression(var_decl.value, ctx)
        id = next_id!(ctx)
        emit!(ctx, IRAssign(var_decl.name, value_id, id))
    end
end

"""상수 선언 로워링"""
function lower_const_decl(const_decl::ConstDecl, ctx::LoweringContext)
    id = next_id!(ctx)
    emit!(ctx, IRVariableDecl(const_decl.name, const_decl.type_name, id))

    value_id = lower_expression(const_decl.value, ctx)
    id = next_id!(ctx)
    emit!(ctx, IRAssign(const_decl.name, value_id, id))
end

"""함수 정의 로워링"""
function lower_function_def(func_def::FunctionDef, ctx::LoweringContext)
    params = [name for (name, _) in func_def.params]

    id = next_id!(ctx)
    emit!(ctx, IRFunctionBegin(func_def.name, params, id))

    for stmt in func_def.body
        lower_statement(stmt, ctx)
    end

    # 명시적 반환이 없으면 없음 반환
    has_return = any(stmt isa ReturnStmt for stmt in func_def.body)
    if !has_return
        emit_return!(ctx, nothing)
    end

    id = next_id!(ctx)
    emit!(ctx, IRFunctionEnd(id))
end

"""If 문 로워링"""
function lower_if_stmt(if_stmt::IfStmt, ctx::LoweringContext)
    cond_id = lower_expression(if_stmt.condition, ctx)

    then_label = new_label!(ctx)
    else_label = if_stmt.else_body !== nothing ? new_label!(ctx) : nothing
    elseif_labels = [new_label!(ctx) for _ in if_stmt.elseif_parts]
    end_label = new_label!(ctx)

    # 첫 번째 elseif가 있으면 그것으로, 없으면 else 또는 end
    next_label = if !isempty(if_stmt.elseif_parts)
        elseif_labels[1]
    elseif else_label !== nothing
        else_label
    else
        end_label
    end

    emit_cond_branch!(ctx, cond_id, then_label, next_label)

    # Then 블록
    emit_label!(ctx, then_label)
    for stmt in if_stmt.then_body
        lower_statement(stmt, ctx)
    end
    emit_goto!(ctx, end_label)

    # Elseif 블록들
    for (i, (elseif_cond, elseif_body)) in enumerate(if_stmt.elseif_parts)
        emit_label!(ctx, elseif_labels[i])

        elseif_cond_id = lower_expression(elseif_cond, ctx)

        next_elseif_label = i < length(elseif_labels) ? elseif_labels[i+1] :
                           else_label !== nothing ? else_label : end_label

        emit_cond_branch!(ctx, elseif_cond_id, new_label!(ctx), next_elseif_label)

        body_label = new_label!(ctx)
        emit_label!(ctx, body_label)
        for stmt in elseif_body
            lower_statement(stmt, ctx)
        end
        emit_goto!(ctx, end_label)
    end

    # Else 블록
    if else_label !== nothing && if_stmt.else_body !== nothing
        emit_label!(ctx, else_label)
        for stmt in if_stmt.else_body
            lower_statement(stmt, ctx)
        end
    end

    emit_label!(ctx, end_label)
end

"""While 문 로워링"""
function lower_while_stmt(while_stmt::WhileStmt, ctx::LoweringContext)
    loop_label = new_label!(ctx)
    body_label = new_label!(ctx)
    end_label = new_label!(ctx)

    emit_label!(ctx, loop_label)

    cond_id = lower_expression(while_stmt.condition, ctx)
    emit_cond_branch!(ctx, cond_id, body_label, end_label)

    emit_label!(ctx, body_label)
    for stmt in while_stmt.body
        lower_statement(stmt, ctx)
    end
    emit_goto!(ctx, loop_label)

    emit_label!(ctx, end_label)
end

"""For 문 로워링 (반복 i = from 부터 to 까지)"""
function lower_for_stmt(for_stmt::ForStmt, ctx::LoweringContext)
    # 초기값: i = from
    from_id = lower_expression(for_stmt.from, ctx)
    id = next_id!(ctx)
    emit!(ctx, IRAssign(for_stmt.var, from_id, id))

    loop_label = new_label!(ctx)
    body_label = new_label!(ctx)
    update_label = new_label!(ctx)
    end_label = new_label!(ctx)

    # 루프 조건 검사: i <= to
    emit_label!(ctx, loop_label)
    to_id = lower_expression(for_stmt.to, ctx)

    var_id = emit_var!(ctx, for_stmt.var)
    id = next_id!(ctx)
    cond_id = next_id!(ctx)
    emit!(ctx, IRBinaryOp("<=", var_id, to_id, cond_id))
    emit_cond_branch!(ctx, cond_id, body_label, end_label)

    # 루프 본체
    emit_label!(ctx, body_label)
    for stmt in for_stmt.body
        lower_statement(stmt, ctx)
    end
    emit_goto!(ctx, update_label)

    # 증가 (i = i + 1)
    emit_label!(ctx, update_label)
    var_id = emit_var!(ctx, for_stmt.var)
    one_id = emit_literal!(ctx, 1)
    id = next_id!(ctx)
    increment_id = next_id!(ctx)
    emit!(ctx, IRBinaryOp("+", var_id, one_id, increment_id))
    emit!(ctx, IRAssign(for_stmt.var, increment_id, id))
    emit_goto!(ctx, loop_label)

    emit_label!(ctx, end_label)
end

"""반환 문 로워링"""
function lower_return_stmt(return_stmt::ReturnStmt, ctx::LoweringContext)
    value_id = nothing
    if return_stmt.value !== nothing
        value_id = lower_expression(return_stmt.value, ctx)
    end
    emit_return!(ctx, value_id)
end

# ======================== 표현식 로워링 ========================

"""표현식 로워링 - 모든 표현식을 IR 값으로 변환"""
function lower_expression(expr::Expression, ctx::LoweringContext)::Int
    if expr isa IntLiteral
        return emit_literal!(ctx, expr.value)
    elseif expr isa FloatLiteral
        return emit_literal!(ctx, expr.value)
    elseif expr isa StringLiteral
        return emit_literal!(ctx, expr.value)
    elseif expr isa BoolLiteral
        return emit_literal!(ctx, expr.value)
    elseif expr isa NothingLiteral
        return emit_literal!(ctx, nothing)
    elseif expr isa Identifier
        return emit_var!(ctx, expr.name)
    elseif expr isa BinaryOp
        return lower_binary_op(expr, ctx)
    elseif expr isa UnaryOp
        return lower_unary_op(expr, ctx)
    elseif expr isa Call
        return lower_call(expr, ctx)
    elseif expr isa ArrayAccess
        return lower_array_access(expr, ctx)
    elseif expr isa MemberAccess
        return lower_member_access(expr, ctx)
    elseif expr isa ArrayLiteral
        return lower_array_literal(expr, ctx)
    else
        error("미지원 표현식: $(typeof(expr))")
    end
end

"""이항 연산 로워링"""
function lower_binary_op(bin_op::BinaryOp, ctx::LoweringContext)::Int
    left_id = lower_expression(bin_op.left, ctx)
    right_id = lower_expression(bin_op.right, ctx)

    id = next_id!(ctx)
    emit!(ctx, IRBinaryOp(bin_op.op, left_id, right_id, id))
    return id
end

"""단항 연산 로워링"""
function lower_unary_op(unary_op::UnaryOp, ctx::LoweringContext)::Int
    operand_id = lower_expression(unary_op.operand, ctx)

    id = next_id!(ctx)
    emit!(ctx, IRUnaryOp(unary_op.op, operand_id, id))
    return id
end

"""함수 호출 로워링"""
function lower_call(call::Call, ctx::LoweringContext)::Int
    func_id = lower_expression(call.func, ctx)

    arg_ids = Int[]
    for arg in call.args
        push!(arg_ids, lower_expression(arg, ctx))
    end

    id = next_id!(ctx)
    emit!(ctx, IRCall(func_id, arg_ids, id))
    return id
end

"""배열 접근 로워링"""
function lower_array_access(array_access::ArrayAccess, ctx::LoweringContext)::Int
    array_id = lower_expression(array_access.array, ctx)
    index_id = lower_expression(array_access.index, ctx)

    id = next_id!(ctx)
    emit!(ctx, IRIndex(array_id, index_id, id))
    return id
end

"""멤버 접근 로워링"""
function lower_member_access(member_access::MemberAccess, ctx::LoweringContext)::Int
    object_id = lower_expression(member_access.object, ctx)

    id = next_id!(ctx)
    emit!(ctx, IRMember(object_id, member_access.member, id))
    return id
end

"""배열 리터럴 로워링"""
function lower_array_literal(array_literal::ArrayLiteral, ctx::LoweringContext)::Int
    elem_ids = Int[]
    for elem in array_literal.elements
        push!(elem_ids, lower_expression(elem, ctx))
    end

    # 배열을 특수한 Call로 표현
    id = next_id!(ctx)
    array_func_id = emit_var!(ctx, "배열")
    emit!(ctx, IRCall(array_func_id, elem_ids, id))
    return id
end

# ======================== 공개 인터페이스 ========================

"""한글 코드 로워링"""
function lower_freeliulia(source::String)::Vector{IRInstruction}
    ast = parse_freeliulia(source)
    ctx = LoweringContext()
    return lower_program(ast, ctx)
end

"""IR 출력 (디버깅용)"""
function print_ir(instructions::Vector{IRInstruction})
    println("=== Untyped IR ===\n")

    for (i, instr) in enumerate(instructions)
        print_ir_instruction(instr, i)
    end
end

"""단일 IR 명령 출력"""
function print_ir_instruction(instr::IRInstruction, index::Int)
    if instr isa IRLiteral
        println("[$index] Literal: $(instr.value) → id=$(instr.id)")
    elseif instr isa IRVar
        println("[$index] Var: $(instr.name) → id=$(instr.id)")
    elseif instr isa IRBinaryOp
        println("[$index] BinOp: $(instr.left_id) $(instr.op) $(instr.right_id) → id=$(instr.id)")
    elseif instr isa IRUnaryOp
        println("[$index] UnaryOp: $(instr.op) $(instr.operand_id) → id=$(instr.id)")
    elseif instr isa IRCall
        args_str = join(instr.arg_ids, ", ")
        println("[$index] Call: func=$(instr.func_id) args=[$args_str] → id=$(instr.id)")
    elseif instr isa IRAssign
        println("[$index] Assign: $(instr.target) = $(instr.value_id) → id=$(instr.id)")
    elseif instr isa IRIndex
        println("[$index] Index: $(instr.array_id)[$(instr.index_id)] → id=$(instr.id)")
    elseif instr isa IRMember
        println("[$index] Member: $(instr.object_id).$(instr.member) → id=$(instr.id)")
    elseif instr isa IRLabel
        println("[$index] Label: $(instr.name)")
    elseif instr isa IRGoto
        println("[$index] Goto: $(instr.target)")
    elseif instr isa IRConditionalBranch
        println("[$index] Branch: if $(instr.condition_id) then $(instr.true_label) else $(instr.false_label)")
    elseif instr isa IRReturn
        ret_str = instr.value_id !== nothing ? "return $(instr.value_id)" : "return"
        println("[$index] $ret_str")
    elseif instr isa IRFunctionBegin
        params_str = join(instr.params, ", ")
        println("[$index] FuncBegin: $(instr.name)([$params_str])")
    elseif instr isa IRFunctionEnd
        println("[$index] FuncEnd")
    elseif instr isa IRVariableDecl
        type_str = instr.type_name !== nothing ? ": $(instr.type_name)" : ""
        println("[$index] VarDecl: $(instr.name)$type_str")
    else
        println("[$index] ??? $(typeof(instr))")
    end
end

# Export
export LoweringContext, lower_freeliulia, print_ir
export IRInstruction, IRLiteral, IRVar, IRBinaryOp, IRCall, IRAssign
export IRLabel, IRGoto, IRConditionalBranch, IRReturn
