"""
    Lexer: Julia 코드를 토큰으로 변환

문자 스트림을 읽고 토큰 스트림으로 변환하는 첫 단계입니다.
"""

# ======================== 토큰 타입 정의 ========================

@enum TokenType begin
    # 기본 토큰
    EOF
    NEWLINE
    COMMENT

    # 키워드 (35개)
    KW_FUNCTION
    KW_RETURN
    KW_IF
    KW_ELSE
    KW_ELSEIF
    KW_END
    KW_FOR
    KW_WHILE
    KW_BREAK
    KW_CONTINUE
    KW_BEGIN
    KW_LET
    KW_DO
    KW_CONST
    KW_GLOBAL
    KW_LOCAL
    KW_TRY
    KW_CATCH
    KW_FINALLY
    KW_IMPORT
    KW_USING
    KW_MODULE
    KW_EXPORT
    KW_INCLUDE
    KW_ABSTRACT
    KW_STRUCT
    KW_MUTABLE
    KW_MACRO
    KW_QUOTE
    KW_TRUE
    KW_FALSE
    KW_NOTHING
    KW_IN
    KW_END_KW  # 마지막 키워드 마커

    # 리터럴
    INTEGER
    FLOAT
    STRING
    SYMBOL
    IDENTIFIER

    # 연산자 & 구분자 (30+개)
    # 산술 연산자
    PLUS
    MINUS
    STAR
    SLASH
    PERCENT
    DOUBLE_SLASH
    CARET
    AT

    # 비교 연산자
    EQ_EQ
    NOT_EQ
    LT
    LE
    GT
    GE

    # 논리 연산자
    AND_AND
    OR_OR
    NOT

    # 비트 연산자
    AND
    OR
    TILDE
    LSHIFT
    RSHIFT

    # 할당 연산자
    EQ
    PLUS_EQ
    MINUS_EQ
    STAR_EQ
    SLASH_EQ
    PERCENT_EQ
    AND_EQ
    OR_EQ
    CARET_EQ

    # 구분자
    LPAREN
    RPAREN
    LBRACKET
    RBRACKET
    LBRACE
    RBRACE
    COMMA
    SEMICOLON
    COLON
    DOUBLE_COLON
    DOT
    ARROW
    DOUBLE_ARROW

    # 특수
    QUESTION
    DOLLAR
    BACKSLASH
end

# ======================== 토큰 구조체 ========================

struct Position
    line::Int
    column::Int
    start::Int  # 파일 내 위치
    length::Int
end

struct Token
    type::TokenType
    value::String
    pos::Position
end

# ======================== 키워드 매핑 ========================

const KEYWORDS = Dict(
    "function" => KW_FUNCTION,
    "return" => KW_RETURN,
    "if" => KW_IF,
    "else" => KW_ELSE,
    "elseif" => KW_ELSEIF,
    "end" => KW_END,
    "for" => KW_FOR,
    "while" => KW_WHILE,
    "break" => KW_BREAK,
    "continue" => KW_CONTINUE,
    "begin" => KW_BEGIN,
    "let" => KW_LET,
    "do" => KW_DO,
    "const" => KW_CONST,
    "global" => KW_GLOBAL,
    "local" => KW_LOCAL,
    "try" => KW_TRY,
    "catch" => KW_CATCH,
    "finally" => KW_FINALLY,
    "import" => KW_IMPORT,
    "using" => KW_USING,
    "module" => KW_MODULE,
    "export" => KW_EXPORT,
    "include" => KW_INCLUDE,
    "abstract" => KW_ABSTRACT,
    "struct" => KW_STRUCT,
    "mutable" => KW_MUTABLE,
    "macro" => KW_MACRO,
    "quote" => KW_QUOTE,
    "true" => KW_TRUE,
    "false" => KW_FALSE,
    "nothing" => KW_NOTHING,
    "in" => KW_IN,
)

# ======================== Lexer 구조체 ========================

mutable struct Lexer
    source::String
    pos::Int  # 현재 위치
    line::Int
    column::Int
    tokens::Vector{Token}
end

function Lexer(source::String)
    return Lexer(source, 1, 1, 1, Token[])
end

# ======================== 헬퍼 함수 ========================

"""현재 문자 반환 (범위 밖이면 Char(0))"""
function current_char(lexer::Lexer)::Char
    if lexer.pos > length(lexer.source)
        return Char(0)
    end
    return lexer.source[lexer.pos]
end

"""다음 문자 반환 (범위 밖이면 Char(0))"""
function peek_char(lexer::Lexer, offset=1)::Char
    pos = lexer.pos + offset
    if pos > length(lexer.source)
        return Char(0)
    end
    return lexer.source[pos]
end

"""한 문자 전진"""
function advance!(lexer::Lexer)
    if lexer.pos <= length(lexer.source)
        if lexer.source[lexer.pos] == '\n'
            lexer.line += 1
            lexer.column = 1
        else
            lexer.column += 1
        end
        lexer.pos += 1
    end
end

"""공백 건너뛰기"""
function skip_whitespace!(lexer::Lexer)
    while current_char(lexer) in (' ', '\t', '\r')
        advance!(lexer)
    end
end

"""주석 건너뛰기"""
function skip_comment!(lexer::Lexer)
    if current_char(lexer) == '#'
        while current_char(lexer) != '\n' && current_char(lexer) != Char(0)
            advance!(lexer)
        end
    end
end

"""토큰 추가"""
function add_token!(lexer::Lexer, type::TokenType, value::String)
    token = Token(type, value, Position(lexer.line, lexer.column, lexer.pos - length(value), length(value)))
    push!(lexer.tokens, token)
end

# ======================== 토큰화 함수 ========================

"""정수 또는 실수 토큰화"""
function tokenize_number!(lexer::Lexer)
    start_pos = lexer.pos
    start_col = lexer.column

    # 정수 부분 읽기
    while isdigit(current_char(lexer))
        advance!(lexer)
    end

    # 실수 부분 확인
    if current_char(lexer) == '.' && isdigit(peek_char(lexer))
        advance!(lexer)  # '.' 건너뛰기
        while isdigit(current_char(lexer))
            advance!(lexer)
        end

        # 과학 기수법 확인 (e.g., 1.5e-3)
        if current_char(lexer) in ('e', 'E')
            advance!(lexer)
            if current_char(lexer) in ('+', '-')
                advance!(lexer)
            end
            while isdigit(current_char(lexer))
                advance!(lexer)
            end
        end

        value = lexer.source[start_pos:lexer.pos-1]
        add_token!(lexer, FLOAT, value)
    else
        value = lexer.source[start_pos:lexer.pos-1]
        add_token!(lexer, INTEGER, value)
    end
end

"""문자열 토큰화"""
function tokenize_string!(lexer::Lexer)
    quote_char = current_char(lexer)
    advance!(lexer)  # 따옴표 건너뛰기

    start_pos = lexer.pos
    value = ""

    while current_char(lexer) != quote_char && current_char(lexer) != Char(0)
        if current_char(lexer) == '\\'
            advance!(lexer)
            # 이스케이프 시퀀스 처리
            char = current_char(lexer)
            if char == 'n'
                value *= '\n'
            elseif char == 't'
                value *= '\t'
            elseif char == 'r'
                value *= '\r'
            elseif char == '\\'
                value *= '\\'
            elseif char == quote_char
                value *= quote_char
            else
                value *= char
            end
            advance!(lexer)
        else
            value *= current_char(lexer)
            advance!(lexer)
        end
    end

    if current_char(lexer) == quote_char
        advance!(lexer)  # 닫는 따옴표 건너뛰기
    end

    add_token!(lexer, STRING, value)
end

"""심볼 토큰화 (:symbol)"""
function tokenize_symbol!(lexer::Lexer)
    advance!(lexer)  # ':' 건너뛰기

    start_pos = lexer.pos
    while isalnum(current_char(lexer)) || current_char(lexer) in ('_', '?', '!')
        advance!(lexer)
    end

    value = lexer.source[start_pos:lexer.pos-1]
    add_token!(lexer, SYMBOL, value)
end

"""식별자 또는 키워드 토큰화"""
function tokenize_identifier!(lexer::Lexer)
    start_pos = lexer.pos

    while isalnum(current_char(lexer)) || current_char(lexer) in ('_', '?', '!')
        advance!(lexer)
    end

    value = lexer.source[start_pos:lexer.pos-1]

    # 키워드 확인
    if haskey(KEYWORDS, value)
        add_token!(lexer, KEYWORDS[value], value)
    else
        add_token!(lexer, IDENTIFIER, value)
    end
end

# ======================== 메인 토크나이저 ========================

"""Julia 소스 코드를 토큰화"""
function tokenize(source::String)::Vector{Token}
    lexer = Lexer(source)

    while lexer.pos <= length(lexer.source)
        skip_whitespace!(lexer)

        if lexer.pos > length(lexer.source)
            break
        end

        char = current_char(lexer)

        # 주석
        if char == '#'
            skip_comment!(lexer)
        # 개행
        elseif char == '\n'
            add_token!(lexer, NEWLINE, "\n")
            advance!(lexer)
        # 숫자
        elseif isdigit(char)
            tokenize_number!(lexer)
        # 문자열
        elseif char in ('"', '\'')
            tokenize_string!(lexer)
        # 심볼
        elseif char == ':' && peek_char(lexer) != ':'
            tokenize_symbol!(lexer)
        # 식별자/키워드
        elseif isalpha(char) || char == '_'
            tokenize_identifier!(lexer)
        # 연산자 & 구분자
        elseif char == '+'
            advance!(lexer)
            if current_char(lexer) == '='
                advance!(lexer)
                add_token!(lexer, PLUS_EQ, "+=")
            else
                add_token!(lexer, PLUS, "+")
            end
        elseif char == '-'
            advance!(lexer)
            if current_char(lexer) == '='
                advance!(lexer)
                add_token!(lexer, MINUS_EQ, "-=")
            elseif current_char(lexer) == '>'
                advance!(lexer)
                add_token!(lexer, ARROW, "->")
            else
                add_token!(lexer, MINUS, "-")
            end
        elseif char == '*'
            advance!(lexer)
            if current_char(lexer) == '='
                advance!(lexer)
                add_token!(lexer, STAR_EQ, "*=")
            else
                add_token!(lexer, STAR, "*")
            end
        elseif char == '/'
            advance!(lexer)
            if current_char(lexer) == '='
                advance!(lexer)
                add_token!(lexer, SLASH_EQ, "/=")
            elseif current_char(lexer) == '/'
                advance!(lexer)
                add_token!(lexer, DOUBLE_SLASH, "//")
            else
                add_token!(lexer, SLASH, "/")
            end
        elseif char == '%'
            advance!(lexer)
            if current_char(lexer) == '='
                advance!(lexer)
                add_token!(lexer, PERCENT_EQ, "%=")
            else
                add_token!(lexer, PERCENT, "%")
            end
        elseif char == '^'
            advance!(lexer)
            if current_char(lexer) == '='
                advance!(lexer)
                add_token!(lexer, CARET_EQ, "^=")
            else
                add_token!(lexer, CARET, "^")
            end
        elseif char == '='
            advance!(lexer)
            if current_char(lexer) == '='
                advance!(lexer)
                add_token!(lexer, EQ_EQ, "==")
            elseif current_char(lexer) == '>'
                advance!(lexer)
                add_token!(lexer, DOUBLE_ARROW, "=>")
            else
                add_token!(lexer, EQ, "=")
            end
        elseif char == '!'
            advance!(lexer)
            if current_char(lexer) == '='
                advance!(lexer)
                add_token!(lexer, NOT_EQ, "!=")
            else
                add_token!(lexer, NOT, "!")
            end
        elseif char == '<'
            advance!(lexer)
            if current_char(lexer) == '='
                advance!(lexer)
                add_token!(lexer, LE, "<=")
            elseif current_char(lexer) == '<'
                advance!(lexer)
                add_token!(lexer, LSHIFT, "<<")
            else
                add_token!(lexer, LT, "<")
            end
        elseif char == '>'
            advance!(lexer)
            if current_char(lexer) == '='
                advance!(lexer)
                add_token!(lexer, GE, ">=")
            elseif current_char(lexer) == '>'
                advance!(lexer)
                add_token!(lexer, RSHIFT, ">>")
            else
                add_token!(lexer, GT, ">")
            end
        elseif char == '&'
            advance!(lexer)
            if current_char(lexer) == '&'
                advance!(lexer)
                add_token!(lexer, AND_AND, "&&")
            elseif current_char(lexer) == '='
                advance!(lexer)
                add_token!(lexer, AND_EQ, "&=")
            else
                add_token!(lexer, AND, "&")
            end
        elseif char == '|'
            advance!(lexer)
            if current_char(lexer) == '|'
                advance!(lexer)
                add_token!(lexer, OR_OR, "||")
            elseif current_char(lexer) == '='
                advance!(lexer)
                add_token!(lexer, OR_EQ, "|=")
            else
                add_token!(lexer, OR, "|")
            end
        elseif char == '~'
            advance!(lexer)
            add_token!(lexer, TILDE, "~")
        elseif char == '('
            advance!(lexer)
            add_token!(lexer, LPAREN, "(")
        elseif char == ')'
            advance!(lexer)
            add_token!(lexer, RPAREN, ")")
        elseif char == '['
            advance!(lexer)
            add_token!(lexer, LBRACKET, "[")
        elseif char == ']'
            advance!(lexer)
            add_token!(lexer, RBRACKET, "]")
        elseif char == '{'
            advance!(lexer)
            add_token!(lexer, LBRACE, "{")
        elseif char == '}'
            advance!(lexer)
            add_token!(lexer, RBRACE, "}")
        elseif char == ','
            advance!(lexer)
            add_token!(lexer, COMMA, ",")
        elseif char == ';'
            advance!(lexer)
            add_token!(lexer, SEMICOLON, ";")
        elseif char == ':'
            advance!(lexer)
            if current_char(lexer) == ':'
                advance!(lexer)
                add_token!(lexer, DOUBLE_COLON, "::")
            else
                add_token!(lexer, COLON, ":")
            end
        elseif char == '.'
            advance!(lexer)
            add_token!(lexer, DOT, ".")
        elseif char == '@'
            advance!(lexer)
            add_token!(lexer, AT, "@")
        elseif char == '$'
            advance!(lexer)
            add_token!(lexer, DOLLAR, "$")
        elseif char == '?'
            advance!(lexer)
            add_token!(lexer, QUESTION, "?")
        else
            advance!(lexer)  # 알 수 없는 문자 건너뛰기
        end
    end

    add_token!(lexer, EOF, "")
    return lexer.tokens
end

# ======================== 유틸리티 함수 ========================

"""토큰 출력"""
function print_tokens(tokens::Vector{Token})
    for token in tokens
        println("$(token.type): '$(token.value)' @ Line $(token.pos.line):$(token.pos.column)")
    end
end

# Export
export Lexer, Token, Position, TokenType, tokenize, print_tokens
export KEYWORDS, FLOAT, INTEGER, STRING, SYMBOL, IDENTIFIER, EOF
export KW_FUNCTION, KW_RETURN, KW_IF, KW_ELSE, KW_END, KW_FOR, KW_WHILE
