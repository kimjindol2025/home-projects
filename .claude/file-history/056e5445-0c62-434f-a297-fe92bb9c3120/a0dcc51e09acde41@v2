"""
    FreeLiulia Parser: 한글 프로그래밍 언어 파서

토큰 스트림을 받아 추상 구문 트리(AST)로 변환합니다.
"""

include("freeliulia_lexer.jl")

# ======================== AST 노드 정의 ========================

abstract type ASTNode end
abstract type Expression <: ASTNode end
abstract type Statement <: ASTNode end

# ─── 리터럴 ────────────────────────────────────────────

struct IntLiteral <: Expression
    value::Int64
end

struct FloatLiteral <: Expression
    value::Float64
end

struct StringLiteral <: Expression
    value::String
end

struct BoolLiteral <: Expression
    value::Bool
end

struct NothingLiteral <: Expression end

struct ArrayLiteral <: Expression
    elements::Vector{Expression}
end

# ─── 식별자 및 기본 식 ──────────────────────────────

struct Identifier <: Expression
    name::String
end

struct BinaryOp <: Expression
    op::String
    left::Expression
    right::Expression
end

struct UnaryOp <: Expression
    op::String
    operand::Expression
end

struct Call <: Expression
    func::Expression
    args::Vector{Expression}
end

struct ArrayAccess <: Expression
    array::Expression
    index::Expression
end

struct MemberAccess <: Expression
    object::Expression
    member::String
end

struct TypeAnnotation <: Expression
    expr::Expression
    type_name::String
end

# ─── 문장 (Statement) ────────────────────────────────

struct VarDecl <: Statement
    name::String
    type_name::Union{Nothing, String}
    value::Union{Nothing, Expression}
end

struct ConstDecl <: Statement
    name::String
    type_name::Union{Nothing, String}
    value::Expression
end

struct Assignment <: Statement
    target::String
    value::Expression
end

struct FunctionDef <: Statement
    name::String
    params::Vector{Tuple{String, String}}  # (name, type)
    return_type::Union{Nothing, String}
    body::Vector{ASTNode}
end

struct IfStmt <: Statement
    condition::Expression
    then_body::Vector{ASTNode}
    elseif_parts::Vector{Tuple{Expression, Vector{ASTNode}}}
    else_body::Union{Nothing, Vector{ASTNode}}
end

struct WhileStmt <: Statement
    condition::Expression
    body::Vector{ASTNode}
end

struct ForStmt <: Statement
    var::String
    from::Expression
    to::Expression
    body::Vector{ASTNode}
end

struct ReturnStmt <: Statement
    value::Union{Nothing, Expression}
end

struct StructDef <: Statement
    name::String
    fields::Vector{Tuple{String, String}}  # (name, type)
end

struct ExpressionStmt <: Statement
    expr::Expression
end

struct Block <: ASTNode
    statements::Vector{ASTNode}
end

# ======================== 프로그램 ========================

struct Program <: ASTNode
    statements::Vector{ASTNode}
end

# ======================== 파서 구조체 ========================

mutable struct FreeLiuliaParser
    tokens::Vector{Token}
    pos::Int
end

function FreeLiuliaParser(tokens::Vector{Token})
    return FreeLiuliaParser(tokens, 1)
end

# ======================== 헬퍼 함수 ========================

"""현재 토큰 반환"""
function current_token(parser::FreeLiuliaParser)::Token
    if parser.pos > length(parser.tokens)
        return parser.tokens[end]  # EOF 토큰
    end
    return parser.tokens[parser.pos]
end

"""다음 토큰 반환"""
function peek_token(parser::FreeLiuliaParser, offset=1)::Token
    pos = parser.pos + offset
    if pos > length(parser.tokens)
        return parser.tokens[end]  # EOF 토큰
    end
    return parser.tokens[pos]
end

"""토큰 타입 확인"""
function check(parser::FreeLiuliaParser, token_type::TokenType)::Bool
    return current_token(parser).type == token_type
end

"""특정 타입 토큰 소비"""
function consume!(parser::FreeLiuliaParser, token_type::TokenType)::Token
    token = current_token(parser)
    if token.type != token_type
        error("예상 토큰: $token_type, 실제: $(token.type)")
    end
    parser.pos += 1
    return token
end

"""토큰 한 개 전진"""
function advance!(parser::FreeLiuliaParser)::Token
    token = current_token(parser)
    parser.pos += 1
    return token
end

"""개행/세미콜론 건너뛰기"""
function skip_newlines!(parser::FreeLiuliaParser)
    while check(parser, NEWLINE) || check(parser, SEMICOLON)
        advance!(parser)
    end
end

# ======================== 파서 함수 ========================

"""프로그램 파싱"""
function parse_program(parser::FreeLiuliaParser)::Program
    statements = ASTNode[]

    skip_newlines!(parser)
    while !check(parser, EOF)
        stmt = parse_statement(parser)
        if stmt !== nothing
            push!(statements, stmt)
        end
        skip_newlines!(parser)
    end

    return Program(statements)
end

"""문장 파싱"""
function parse_statement(parser::FreeLiuliaParser)::Union{Nothing, ASTNode}
    skip_newlines!(parser)

    if check(parser, EOF)
        return nothing
    end

    # 변수 선언: 변수 x: 정수 = 10
    if check(parser, KW_변수)
        return parse_var_decl(parser)
    end

    # 상수 선언: 상수 PI: 실수 = 3.14
    if check(parser, KW_상수)
        return parse_const_decl(parser)
    end

    # 함수 정의: 함수 더하기(x: 정수, y: 정수) -> 정수 ... 끝
    if check(parser, KW_함수)
        return parse_function_def(parser)
    end

    # 구조체: 구조체 점 ... 끝
    if check(parser, KW_구조체)
        return parse_struct_def(parser)
    end

    # If 문: 만약 ... 끝
    if check(parser, KW_만약)
        return parse_if_stmt(parser)
    end

    # While 문: 동안 ... 끝
    if check(parser, KW_동안)
        return parse_while_stmt(parser)
    end

    # For 문: 반복 i = 1 부터 5 까지 ... 끝
    if check(parser, KW_반복)
        return parse_for_stmt(parser)
    end

    # Return 문: 반환 값
    if check(parser, KW_반환)
        return parse_return_stmt(parser)
    end

    # 표현식 문: 식만 있는 문장
    expr = parse_expression(parser)
    return ExpressionStmt(expr)
end

"""변수 선언 파싱"""
function parse_var_decl(parser::FreeLiuliaParser)::VarDecl
    consume!(parser, KW_변수)

    name_token = consume!(parser, IDENTIFIER)
    name = name_token.value

    type_name = nothing

    # 타입 어노테이션: : 정수
    if check(parser, COLON)
        advance!(parser)
        type_token = consume!(parser, IDENTIFIER)
        type_name = type_token.value
    end

    value = nothing

    # 초기값: = 10
    if check(parser, EQ)
        advance!(parser)
        value = parse_expression(parser)
    end

    skip_newlines!(parser)
    return VarDecl(name, type_name, value)
end

"""상수 선언 파싱"""
function parse_const_decl(parser::FreeLiuliaParser)::ConstDecl
    consume!(parser, KW_상수)

    name_token = consume!(parser, IDENTIFIER)
    name = name_token.value

    type_name = nothing

    # 타입: : 실수
    if check(parser, COLON)
        advance!(parser)
        type_token = consume!(parser, IDENTIFIER)
        type_name = type_token.value
    end

    # 값 (필수): = 3.14
    consume!(parser, EQ)
    value = parse_expression(parser)

    skip_newlines!(parser)
    return ConstDecl(name, type_name, value)
end

"""함수 정의 파싱"""
function parse_function_def(parser::FreeLiuliaParser)::FunctionDef
    consume!(parser, KW_함수)

    name_token = consume!(parser, IDENTIFIER)
    name = name_token.value

    # 매개변수: (x: 정수, y: 정수)
    consume!(parser, LPAREN)

    params = Tuple{String, String}[]

    if !check(parser, RPAREN)
        while true
            param_name_token = consume!(parser, IDENTIFIER)
            param_name = param_name_token.value

            consume!(parser, COLON)

            param_type_token = consume!(parser, IDENTIFIER)
            param_type = param_type_token.value

            push!(params, (param_name, param_type))

            if !check(parser, COMMA)
                break
            end
            advance!(parser)  # consume comma
        end
    end

    consume!(parser, RPAREN)

    # 반환 타입: -> 정수
    return_type = nothing
    if check(parser, ARROW)
        advance!(parser)
        return_type_token = consume!(parser, IDENTIFIER)
        return_type = return_type_token.value
    end

    skip_newlines!(parser)

    # 함수 본체
    body = parse_block_until(parser, KW_끝)

    consume!(parser, KW_끝)
    skip_newlines!(parser)

    return FunctionDef(name, params, return_type, body)
end

"""구조체 정의 파싱"""
function parse_struct_def(parser::FreeLiuliaParser)::StructDef
    consume!(parser, KW_구조체)

    name_token = consume!(parser, IDENTIFIER)
    name = name_token.value

    skip_newlines!(parser)

    fields = Tuple{String, String}[]

    while !check(parser, KW_끝)
        field_name_token = consume!(parser, IDENTIFIER)
        field_name = field_name_token.value

        consume!(parser, COLON)

        field_type_token = consume!(parser, IDENTIFIER)
        field_type = field_type_token.value

        push!(fields, (field_name, field_type))

        skip_newlines!(parser)
    end

    consume!(parser, KW_끝)
    skip_newlines!(parser)

    return StructDef(name, fields)
end

"""If 문 파싱"""
function parse_if_stmt(parser::FreeLiuliaParser)::IfStmt
    consume!(parser, KW_만약)

    condition = parse_expression(parser)
    skip_newlines!(parser)

    then_body = parse_block_until(parser, Union{KW_아니면, KW_아니면만약, KW_끝})

    elseif_parts = Tuple{Expression, Vector{ASTNode}}[]

    # 아니면만약
    while check(parser, KW_아니면만약)
        advance!(parser)
        elseif_cond = parse_expression(parser)
        skip_newlines!(parser)
        elseif_body = parse_block_until(parser, Union{KW_아니면, KW_아니면만약, KW_끝})
        push!(elseif_parts, (elseif_cond, elseif_body))
    end

    else_body = nothing

    # 아니면
    if check(parser, KW_아니면)
        advance!(parser)
        skip_newlines!(parser)
        else_body = parse_block_until(parser, KW_끝)
    end

    consume!(parser, KW_끝)
    skip_newlines!(parser)

    return IfStmt(condition, then_body, elseif_parts, else_body)
end

"""While 문 파싱"""
function parse_while_stmt(parser::FreeLiuliaParser)::WhileStmt
    consume!(parser, KW_동안)

    condition = parse_expression(parser)
    skip_newlines!(parser)

    body = parse_block_until(parser, KW_끝)

    consume!(parser, KW_끝)
    skip_newlines!(parser)

    return WhileStmt(condition, body)
end

"""For 문 파싱: 반복 i = 1 부터 5 까지"""
function parse_for_stmt(parser::FreeLiuliaParser)::ForStmt
    consume!(parser, KW_반복)

    var_token = consume!(parser, IDENTIFIER)
    var = var_token.value

    consume!(parser, EQ)
    from = parse_expression(parser)

    consume!(parser, KW_부터)
    to = parse_expression(parser)

    consume!(parser, KW_까지)
    skip_newlines!(parser)

    body = parse_block_until(parser, KW_끝)

    consume!(parser, KW_끝)
    skip_newlines!(parser)

    return ForStmt(var, from, to, body)
end

"""Return 문 파싱"""
function parse_return_stmt(parser::FreeLiuliaParser)::ReturnStmt
    consume!(parser, KW_반환)

    value = nothing

    if !check(parser, NEWLINE) && !check(parser, EOF) && !check(parser, KW_끝)
        value = parse_expression(parser)
    end

    skip_newlines!(parser)
    return ReturnStmt(value)
end

"""블록 파싱: 특정 끝 토큰까지"""
function parse_block_until(parser::FreeLiuliaParser, end_token::Union{TokenType, Type})::Vector{ASTNode}
    statements = ASTNode[]

    skip_newlines!(parser)

    while !check_any(parser, end_token)
        if check(parser, EOF)
            error("예상치 못한 EOF")
        end
        stmt = parse_statement(parser)
        if stmt !== nothing
            push!(statements, stmt)
        end
        skip_newlines!(parser)
    end

    return statements
end

"""여러 토큰 타입 확인"""
function check_any(parser::FreeLiuliaParser, types::Union{TokenType, Type})
    if types isa TokenType
        return check(parser, types)
    elseif types isa Union
        for t in types.types
            if check(parser, t)
                return true
            end
        end
        return false
    else
        return check(parser, types)
    end
end

# ======================== 표현식 파싱 ========================

"""표현식 파싱 (최소 우선순위)"""
function parse_expression(parser::FreeLiuliaParser)::Expression
    return parse_assignment(parser)
end

"""할당 파싱"""
function parse_assignment(parser::FreeLiuliaParser)::Expression
    expr = parse_logical_or(parser)

    if check(parser, EQ) && expr isa Identifier
        advance!(parser)
        value = parse_assignment(parser)
        # 할당을 표현식으로 취급 (문장으로 변환됨)
        return expr
    end

    return expr
end

"""논리 OR 파싱"""
function parse_logical_or(parser::FreeLiuliaParser)::Expression
    left = parse_logical_and(parser)

    while check(parser, OR_OR)
        op_token = advance!(parser)
        right = parse_logical_and(parser)
        left = BinaryOp(op_token.value, left, right)
    end

    return left
end

"""논리 AND 파싱"""
function parse_logical_and(parser::FreeLiuliaParser)::Expression
    left = parse_equality(parser)

    while check(parser, AND_AND)
        op_token = advance!(parser)
        right = parse_equality(parser)
        left = BinaryOp(op_token.value, left, right)
    end

    return left
end

"""등호 파싱"""
function parse_equality(parser::FreeLiuliaParser)::Expression
    left = parse_comparison(parser)

    while check(parser, EQ_EQ) || check(parser, NOT_EQ)
        op_token = advance!(parser)
        right = parse_comparison(parser)
        left = BinaryOp(op_token.value, left, right)
    end

    return left
end

"""비교 파싱"""
function parse_comparison(parser::FreeLiuliaParser)::Expression
    left = parse_additive(parser)

    while check(parser, LT) || check(parser, LE) || check(parser, GT) || check(parser, GE)
        op_token = advance!(parser)
        right = parse_additive(parser)
        left = BinaryOp(op_token.value, left, right)
    end

    return left
end

"""덧셈/뺄셈 파싱"""
function parse_additive(parser::FreeLiuliaParser)::Expression
    left = parse_multiplicative(parser)

    while check(parser, PLUS) || check(parser, MINUS) || check(parser, CONCAT)
        op_token = advance!(parser)
        right = parse_multiplicative(parser)
        left = BinaryOp(op_token.value, left, right)
    end

    return left
end

"""곱셈/나눗셈 파싱"""
function parse_multiplicative(parser::FreeLiuliaParser)::Expression
    left = parse_unary(parser)

    while check(parser, STAR) || check(parser, SLASH) || check(parser, PERCENT)
        op_token = advance!(parser)
        right = parse_unary(parser)
        left = BinaryOp(op_token.value, left, right)
    end

    return left
end

"""단항 연산 파싱"""
function parse_unary(parser::FreeLiuliaParser)::Expression
    if check(parser, NOT) || check(parser, MINUS)
        op_token = advance!(parser)
        operand = parse_unary(parser)
        return UnaryOp(op_token.value, operand)
    end

    return parse_postfix(parser)
end

"""후위 연산 파싱 (함수 호출, 배열 접근, 멤버 접근)"""
function parse_postfix(parser::FreeLiuliaParser)::Expression
    expr = parse_primary(parser)

    while true
        if check(parser, LPAREN)
            # 함수 호출
            advance!(parser)
            args = parse_arguments(parser)
            consume!(parser, RPAREN)
            expr = Call(expr, args)
        elseif check(parser, LBRACKET)
            # 배열 접근
            advance!(parser)
            index = parse_expression(parser)
            consume!(parser, RBRACKET)
            expr = ArrayAccess(expr, index)
        elseif check(parser, DOT)
            # 멤버 접근
            advance!(parser)
            member_token = consume!(parser, IDENTIFIER)
            expr = MemberAccess(expr, member_token.value)
        else
            break
        end
    end

    return expr
end

"""기본 식 파싱"""
function parse_primary(parser::FreeLiuliaParser)::Expression
    # 정수
    if check(parser, INTEGER)
        token = advance!(parser)
        value = parse(Int64, token.value)
        return IntLiteral(value)
    end

    # 실수
    if check(parser, FLOAT)
        token = advance!(parser)
        value = parse(Float64, token.value)
        return FloatLiteral(value)
    end

    # 문자열
    if check(parser, STRING)
        token = advance!(parser)
        return StringLiteral(token.value)
    end

    # 참/거짓
    if check(parser, KW_참)
        advance!(parser)
        return BoolLiteral(true)
    end

    if check(parser, KW_거짓)
        advance!(parser)
        return BoolLiteral(false)
    end

    # 없음
    if check(parser, KW_없음)
        advance!(parser)
        return NothingLiteral()
    end

    # 배열 리터럴: [1, 2, 3]
    if check(parser, LBRACKET)
        advance!(parser)
        elements = Expression[]

        if !check(parser, RBRACKET)
            while true
                push!(elements, parse_expression(parser))
                if !check(parser, COMMA)
                    break
                end
                advance!(parser)
            end
        end

        consume!(parser, RBRACKET)
        return ArrayLiteral(elements)
    end

    # 괄호로 묶인 식: (1 + 2)
    if check(parser, LPAREN)
        advance!(parser)
        expr = parse_expression(parser)
        consume!(parser, RPAREN)
        return expr
    end

    # 식별자
    if check(parser, IDENTIFIER)
        token = advance!(parser)
        return Identifier(token.value)
    end

    error("예상치 못한 토큰: $(current_token(parser))")
end

"""함수 인수 파싱"""
function parse_arguments(parser::FreeLiuliaParser)::Vector{Expression}
    args = Expression[]

    if !check(parser, RPAREN)
        while true
            push!(args, parse_expression(parser))
            if !check(parser, COMMA)
                break
            end
            advance!(parser)
        end
    end

    return args
end

# ======================== 공개 인터페이스 ========================

"""FreeLiulia 코드 파싱"""
function parse_freeliulia(source::String)::Program
    tokens = tokenize_freeliulia(source)
    parser = FreeLiuliaParser(tokens)
    return parse_program(parser)
end

"""AST 출력 (디버깅용)"""
function print_ast(node::ASTNode, indent=0)
    prefix = repeat("  ", indent)

    if node isa Program
        println("$prefix프로그램:")
        for stmt in node.statements
            print_ast(stmt, indent + 1)
        end
    elseif node isa VarDecl
        println("$prefix변수: $(node.name) $(node.type_name !== nothing ? ": $(node.type_name)" : "")")
        if node.value !== nothing
            print_ast(node.value, indent + 1)
        end
    elseif node isa ConstDecl
        println("$prefix상수: $(node.name) $(node.type_name !== nothing ? ": $(node.type_name)" : "")")
        print_ast(node.value, indent + 1)
    elseif node isa FunctionDef
        params_str = join(["$name: $type" for (name, type) in node.params], ", ")
        return_type_str = node.return_type !== nothing ? " -> $(node.return_type)" : ""
        println("$prefix함수: $(node.name)($params_str)$return_type_str")
        for stmt in node.body
            print_ast(stmt, indent + 1)
        end
    elseif node isa IfStmt
        println("$prefix만약:")
        print_ast(node.condition, indent + 1)
        println("$prefix  그럼:")
        for stmt in node.then_body
            print_ast(stmt, indent + 2)
        end
        for (cond, body) in node.elseif_parts
            println("$prefix아니면만약:")
            print_ast(cond, indent + 1)
            for stmt in body
                print_ast(stmt, indent + 2)
            end
        end
        if node.else_body !== nothing
            println("$prefix아니면:")
            for stmt in node.else_body
                print_ast(stmt, indent + 2)
            end
        end
    elseif node isa WhileStmt
        println("$prefix동안:")
        print_ast(node.condition, indent + 1)
        for stmt in node.body
            print_ast(stmt, indent + 1)
        end
    elseif node isa ForStmt
        println("$prefix반복: $(node.var)")
        print_ast(node.from, indent + 1)
        print_ast(node.to, indent + 1)
        for stmt in node.body
            print_ast(stmt, indent + 1)
        end
    elseif node isa ReturnStmt
        println("$prefix반환")
        if node.value !== nothing
            print_ast(node.value, indent + 1)
        end
    elseif node isa ExpressionStmt
        print_ast(node.expr, indent)
    elseif node isa BinaryOp
        println("$prefix$(node.op):")
        print_ast(node.left, indent + 1)
        print_ast(node.right, indent + 1)
    elseif node isa UnaryOp
        println("$prefix단항 $(node.op):")
        print_ast(node.operand, indent + 1)
    elseif node isa Call
        println("$prefix호출:")
        print_ast(node.func, indent + 1)
        for arg in node.args
            print_ast(arg, indent + 1)
        end
    elseif node isa Identifier
        println("$prefix식별자: $(node.name)")
    elseif node isa IntLiteral
        println("$prefix정수: $(node.value)")
    elseif node isa FloatLiteral
        println("$prefix실수: $(node.value)")
    elseif node isa StringLiteral
        println("$prefix문자열: \"$(node.value)\"")
    elseif node isa BoolLiteral
        println("$prefix논리값: $(node.value)")
    elseif node isa NothingLiteral
        println("$prefix없음")
    elseif node isa ArrayLiteral
        println("$prefix배열:")
        for elem in node.elements
            print_ast(elem, indent + 1)
        end
    else
        println("$prefix??? $(typeof(node))")
    end
end

# Export
export FreeLiuliaParser, parse_freeliulia, print_ast
export Program, VarDecl, ConstDecl, FunctionDef, IfStmt, WhileStmt, ForStmt
export IntLiteral, FloatLiteral, StringLiteral, BoolLiteral, Identifier, BinaryOp, Call
