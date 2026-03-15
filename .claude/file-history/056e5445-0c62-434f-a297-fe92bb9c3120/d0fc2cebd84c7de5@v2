"""
    Parser: AST (추상 구문 트리) 생성

토큰 스트림을 받아 추상 구문 트리로 변환합니다.
"""

include("lexer.jl")

# ======================== AST 노드 타입 ========================

abstract type ASTNode end
abstract type Expression <: ASTNode end
abstract type Statement <: ASTNode end

# 리터럴
struct IntegerLiteral <: Expression
    value::Int64
end

struct FloatLiteral <: Expression
    value::Float64
end

struct StringLiteral <: Expression
    value::String
end

struct SymbolLiteral <: Expression
    value::String
end

struct BooleanLiteral <: Expression
    value::Bool
end

struct NothingLiteral <: Expression end

# 식별자
struct Identifier <: Expression
    name::String
end

# 이항 연산
struct BinaryOp <: Expression
    op::String
    left::Expression
    right::Expression
end

# 단항 연산
struct UnaryOp <: Expression
    op::String
    operand::Expression
end

# 함수 호출
struct Call <: Expression
    func::Expression
    args::Vector{Expression}
end

# 배열 인덱싱
struct IndexAccess <: Expression
    object::Expression
    indices::Vector{Expression}
end

# 멤버 접근
struct MemberAccess <: Expression
    object::Expression
    field::String
end

# 타입 어노테이션
struct TypeAnnotation <: Expression
    value::Expression
    type::Expression
end

# 할당
struct Assignment <: Statement
    target::String
    value::Expression
end

# 함수 선언
struct FunctionDef <: Statement
    name::String
    params::Vector{String}
    return_type::Union{Nothing, String}
    body::Vector{ASTNode}
end

# If 문
struct IfStmt <: Statement
    condition::Expression
    then_body::Vector{ASTNode}
    else_body::Union{Nothing, Vector{ASTNode}}
end

# While 루프
struct WhileStmt <: Statement
    condition::Expression
    body::Vector{ASTNode}
end

# For 루프
struct ForStmt <: Statement
    variable::String
    iterator::Expression
    body::Vector{ASTNode}
end

# Return 문
struct ReturnStmt <: Statement
    value::Union{Nothing, Expression}
end

# 블록
struct Block <: Statement
    statements::Vector{ASTNode}
end

# 배열 리터럴
struct ArrayLiteral <: Expression
    elements::Vector{Expression}
end

# 튜플 리터럴
struct TupleLiteral <: Expression
    elements::Vector{Expression}
end

# ======================== Parser 구조체 ========================

mutable struct Parser
    tokens::Vector{Token}
    current::Int  # 현재 토큰 인덱스
end

function Parser(tokens::Vector{Token})
    return Parser(tokens, 1)
end

# ======================== 헬퍼 함수 ========================

"""현재 토큰 반환"""
function current_token(parser::Parser)::Token
    if parser.current > length(parser.tokens)
        return parser.tokens[end]  # EOF
    end
    return parser.tokens[parser.current]
end

"""다음 토큰 반환"""
function peek_token(parser::Parser, offset=1)::Token
    pos = parser.current + offset
    if pos > length(parser.tokens)
        return parser.tokens[end]  # EOF
    end
    return parser.tokens[pos]
end

"""한 토큰 전진"""
function advance!(parser::Parser)
    if parser.current <= length(parser.tokens)
        parser.current += 1
    end
end

"""특정 타입 토큰 기대"""
function expect!(parser::Parser, type::TokenType)::Token
    token = current_token(parser)
    if token.type != type
        error("Expected $(type), got $(token.type)")
    end
    advance!(parser)
    return token
end

"""공백/줄바꿈 건너뛰기"""
function skip_newlines!(parser::Parser)
    while current_token(parser).type == NEWLINE
        advance!(parser)
    end
end

# ======================== 파싱 함수 ========================

"""프로그램 파싱"""
function parse_program(parser::Parser)::Vector{ASTNode}
    statements = ASTNode[]

    skip_newlines!(parser)
    while current_token(parser).type != EOF
        stmt = parse_statement!(parser)
        if stmt !== nothing
            push!(statements, stmt)
        end
        skip_newlines!(parser)
    end

    return statements
end

"""문 파싱"""
function parse_statement!(parser::Parser)::Union{Nothing, ASTNode}
    skip_newlines!(parser)

    token = current_token(parser)

    if token.type == KW_FUNCTION
        return parse_function_def!(parser)
    elseif token.type == KW_IF
        return parse_if_stmt!(parser)
    elseif token.type == KW_WHILE
        return parse_while_stmt!(parser)
    elseif token.type == KW_FOR
        return parse_for_stmt!(parser)
    elseif token.type == KW_RETURN
        return parse_return_stmt!(parser)
    elseif token.type == IDENTIFIER
        # 할당 확인
        if peek_token(parser).type == EQ
            return parse_assignment!(parser)
        else
            return parse_expression_stmt!(parser)
        end
    else
        return parse_expression_stmt!(parser)
    end
end

"""함수 정의 파싱"""
function parse_function_def!(parser::Parser)::FunctionDef
    expect!(parser, KW_FUNCTION)

    # 함수명
    name = expect!(parser, IDENTIFIER).value

    # 매개변수
    expect!(parser, LPAREN)
    params = String[]
    while current_token(parser).type != RPAREN
        param = expect!(parser, IDENTIFIER).value
        push!(params, param)
        if current_token(parser).type == COMMA
            advance!(parser)
        end
    end
    expect!(parser, RPAREN)

    # 본체
    skip_newlines!(parser)
    body = ASTNode[]
    while current_token(parser).type != KW_END
        stmt = parse_statement!(parser)
        if stmt !== nothing
            push!(body, stmt)
        end
        skip_newlines!(parser)
    end
    expect!(parser, KW_END)

    return FunctionDef(name, params, nothing, body)
end

"""If 문 파싱"""
function parse_if_stmt!(parser::Parser)::IfStmt
    expect!(parser, KW_IF)

    # 조건
    condition = parse_expression!(parser)

    # 본체
    skip_newlines!(parser)
    then_body = ASTNode[]
    while current_token(parser).type ∉ (KW_ELSE, KW_ELSEIF, KW_END)
        stmt = parse_statement!(parser)
        if stmt !== nothing
            push!(then_body, stmt)
        end
        skip_newlines!(parser)
    end

    # Else 부분
    else_body = nothing
    if current_token(parser).type == KW_ELSE
        advance!(parser)
        skip_newlines!(parser)
        else_body = ASTNode[]
        while current_token(parser).type != KW_END
            stmt = parse_statement!(parser)
            if stmt !== nothing
                push!(else_body, stmt)
            end
            skip_newlines!(parser)
        end
    end

    expect!(parser, KW_END)
    return IfStmt(condition, then_body, else_body)
end

"""While 루프 파싱"""
function parse_while_stmt!(parser::Parser)::WhileStmt
    expect!(parser, KW_WHILE)

    condition = parse_expression!(parser)

    skip_newlines!(parser)
    body = ASTNode[]
    while current_token(parser).type != KW_END
        stmt = parse_statement!(parser)
        if stmt !== nothing
            push!(body, stmt)
        end
        skip_newlines!(parser)
    end
    expect!(parser, KW_END)

    return WhileStmt(condition, body)
end

"""For 루프 파싱"""
function parse_for_stmt!(parser::Parser)::ForStmt
    expect!(parser, KW_FOR)

    variable = expect!(parser, IDENTIFIER).value
    expect!(parser, KW_IN)
    iterator = parse_expression!(parser)

    skip_newlines!(parser)
    body = ASTNode[]
    while current_token(parser).type != KW_END
        stmt = parse_statement!(parser)
        if stmt !== nothing
            push!(body, stmt)
        end
        skip_newlines!(parser)
    end
    expect!(parser, KW_END)

    return ForStmt(variable, iterator, body)
end

"""Return 문 파싱"""
function parse_return_stmt!(parser::Parser)::ReturnStmt
    expect!(parser, KW_RETURN)

    value = nothing
    if current_token(parser).type ∉ (NEWLINE, EOF, KW_END)
        value = parse_expression!(parser)
    end

    return ReturnStmt(value)
end

"""할당 파싱"""
function parse_assignment!(parser::Parser)::Assignment
    target = expect!(parser, IDENTIFIER).value
    expect!(parser, EQ)
    value = parse_expression!(parser)
    return Assignment(target, value)
end

"""표현식 문 파싱"""
function parse_expression_stmt!(parser::Parser)::Union{Nothing, Expression}
    expr = parse_expression!(parser)
    return expr
end

"""표현식 파싱"""
function parse_expression!(parser::Parser)::Expression
    return parse_additive!(parser)
end

"""덧셈/뺄셈 파싱"""
function parse_additive!(parser::Parser)::Expression
    left = parse_multiplicative!(parser)

    while current_token(parser).type ∈ (PLUS, MINUS)
        op = current_token(parser).value
        advance!(parser)
        right = parse_multiplicative!(parser)
        left = BinaryOp(op, left, right)
    end

    return left
end

"""곱셈/나눗셈 파싱"""
function parse_multiplicative!(parser::Parser)::Expression
    left = parse_unary!(parser)

    while current_token(parser).type ∈ (STAR, SLASH, PERCENT)
        op = current_token(parser).value
        advance!(parser)
        right = parse_unary!(parser)
        left = BinaryOp(op, left, right)
    end

    return left
end

"""단항 연산 파싱"""
function parse_unary!(parser::Parser)::Expression
    if current_token(parser).type ∈ (MINUS, NOT)
        op = current_token(parser).value
        advance!(parser)
        operand = parse_unary!(parser)
        return UnaryOp(op, operand)
    end

    return parse_postfix!(parser)
end

"""Postfix 연산 (함수 호출, 배열 인덱싱) 파싱"""
function parse_postfix!(parser::Parser)::Expression
    expr = parse_primary!(parser)

    while true
        if current_token(parser).type == LPAREN
            # 함수 호출
            advance!(parser)
            args = Expression[]
            while current_token(parser).type != RPAREN
                push!(args, parse_expression!(parser))
                if current_token(parser).type == COMMA
                    advance!(parser)
                end
            end
            expect!(parser, RPAREN)
            expr = Call(expr, args)
        elseif current_token(parser).type == LBRACKET
            # 배열 인덱싱
            advance!(parser)
            indices = Expression[]
            while current_token(parser).type != RBRACKET
                push!(indices, parse_expression!(parser))
                if current_token(parser).type == COMMA
                    advance!(parser)
                end
            end
            expect!(parser, RBRACKET)
            expr = IndexAccess(expr, indices)
        elseif current_token(parser).type == DOT
            # 멤버 접근
            advance!(parser)
            field = expect!(parser, IDENTIFIER).value
            expr = MemberAccess(expr, field)
        else
            break
        end
    end

    return expr
end

"""기본 식 파싱"""
function parse_primary!(parser::Parser)::Expression
    token = current_token(parser)

    if token.type == INTEGER
        advance!(parser)
        return IntegerLiteral(parse(Int64, token.value))
    elseif token.type == FLOAT
        advance!(parser)
        return FloatLiteral(parse(Float64, token.value))
    elseif token.type == STRING
        advance!(parser)
        return StringLiteral(token.value)
    elseif token.type == SYMBOL
        advance!(parser)
        return SymbolLiteral(token.value)
    elseif token.type == KW_TRUE
        advance!(parser)
        return BooleanLiteral(true)
    elseif token.type == KW_FALSE
        advance!(parser)
        return BooleanLiteral(false)
    elseif token.type == KW_NOTHING
        advance!(parser)
        return NothingLiteral()
    elseif token.type == IDENTIFIER
        advance!(parser)
        return Identifier(token.value)
    elseif token.type == LPAREN
        advance!(parser)
        expr = parse_expression!(parser)
        expect!(parser, RPAREN)
        return expr
    elseif token.type == LBRACKET
        advance!(parser)
        elements = Expression[]
        while current_token(parser).type != RBRACKET
            push!(elements, parse_expression!(parser))
            if current_token(parser).type == COMMA
                advance!(parser)
            end
        end
        expect!(parser, RBRACKET)
        return ArrayLiteral(elements)
    else
        error("Unexpected token: $(token.type) = '$(token.value)'")
    end
end

# ======================== 메인 파서 함수 ========================

"""Julia 코드를 AST로 파싱"""
function parse(source::String)::Vector{ASTNode}
    tokens = tokenize(source)
    parser = Parser(tokens)
    return parse_program(parser)
end

# ======================== AST 출력 ========================

"""AST를 문자열로 변환"""
function ast_to_string(node::ASTNode, indent=0)::String
    ind = repeat(" ", indent)

    if isa(node, IntegerLiteral)
        return "$(ind)Integer($(node.value))"
    elseif isa(node, FloatLiteral)
        return "$(ind)Float($(node.value))"
    elseif isa(node, StringLiteral)
        return "$(ind)String(\"$(node.value)\")"
    elseif isa(node, SymbolLiteral)
        return "$(ind)Symbol(:$(node.value))"
    elseif isa(node, BooleanLiteral)
        return "$(ind)Boolean($(node.value))"
    elseif isa(node, NothingLiteral)
        return "$(ind)Nothing"
    elseif isa(node, Identifier)
        return "$(ind)Identifier($(node.name))"
    elseif isa(node, BinaryOp)
        return "$(ind)BinaryOp($(node.op))\n" *
               ast_to_string(node.left, indent+2) * "\n" *
               ast_to_string(node.right, indent+2)
    elseif isa(node, UnaryOp)
        return "$(ind)UnaryOp($(node.op))\n" *
               ast_to_string(node.operand, indent+2)
    elseif isa(node, Call)
        args_str = join([ast_to_string(arg, indent+2) for arg in node.args], "\n")
        return "$(ind)Call\n" *
               ast_to_string(node.func, indent+2) * "\n" *
               "$(ind)  Args:\n" * args_str
    elseif isa(node, Assignment)
        return "$(ind)Assignment($(node.target))\n" *
               ast_to_string(node.value, indent+2)
    elseif isa(node, FunctionDef)
        params = join(node.params, ", ")
        body_str = join([ast_to_string(stmt, indent+2) for stmt in node.body], "\n")
        return "$(ind)FunctionDef($(node.name))($(params))\n$(body_str)"
    elseif isa(node, ReturnStmt)
        if node.value === nothing
            return "$(ind)Return"
        else
            return "$(ind)Return\n" * ast_to_string(node.value, indent+2)
        end
    elseif isa(node, IfStmt)
        then_str = join([ast_to_string(stmt, indent+2) for stmt in node.then_body], "\n")
        result = "$(ind)If\n$(ind)  Condition:\n" *
                 ast_to_string(node.condition, indent+4) * "\n" *
                 "$(ind)  Then:\n" * then_str
        if node.else_body !== nothing
            else_str = join([ast_to_string(stmt, indent+2) for stmt in node.else_body], "\n")
            result *= "\n$(ind)  Else:\n" * else_str
        end
        return result
    else
        return "$(ind)$(typeof(node))"
    end
end

# Export
export Parser, ASTNode, Expression, Statement
export IntegerLiteral, FloatLiteral, StringLiteral, BooleanLiteral, Identifier
export BinaryOp, UnaryOp, Call, IndexAccess, MemberAccess
export Assignment, FunctionDef, IfStmt, WhileStmt, ForStmt, ReturnStmt
export parse, ast_to_string
