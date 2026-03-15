"""
    FreeLiulia Lexer: 한글 프로그래밍 언어 토크나이저

FreeLang + Julia 혼합 언어의 한글 토큰화
"""

# ======================== 토큰 타입 정의 ========================

@enum TokenType begin
    # 기본 토큰
    EOF
    NEWLINE
    COMMENT

    # 한글 키워드 (25개)
    KW_함수          # function
    KW_변수          # variable/let
    KW_상수          # const
    KW_만약          # if
    KW_아니면        # else
    KW_아니면만약    # elseif
    KW_반복          # for
    KW_동안          # while
    KW_반환          # return
    KW_끝            # end
    KW_없음          # nothing
    KW_참            # true
    KW_거짓          # false
    KW_구조체        # struct
    KW_배열          # array
    KW_딕셔너리      # dict
    KW_부터          # from
    KW_까지          # to
    KW_까지_1        # to (alternative)
    KW_하기          # do
    KW_시작          # begin
    KW_사용          # using/import
    KW_포함          # include
    KW_전역          # global
    KW_지역          # local

    # 리터럴
    INTEGER
    FLOAT
    STRING
    SYMBOL
    IDENTIFIER

    # 연산자
    PLUS            # +
    MINUS           # -
    STAR            # *
    SLASH           # /
    PERCENT         # %
    CARET           # ^
    DOUBLE_SLASH    # //

    # 비교
    EQ_EQ           # ==
    NOT_EQ          # !=
    LT              # <
    LE              # <=
    GT              # >
    GE              # >=

    # 논리
    AND_AND         # &&
    OR_OR           # ||
    NOT             # !

    # 할당
    EQ              # =
    PLUS_EQ         # +=
    MINUS_EQ        # -=
    STAR_EQ         # *=
    SLASH_EQ        # /=

    # 문자열 연결
    CONCAT          # ++

    # 구분자
    LPAREN          # (
    RPAREN          # )
    LBRACKET        # [
    RBRACKET        # ]
    LBRACE          # {
    RBRACE          # }
    COMMA           # ,
    SEMICOLON       # ;
    COLON           # :
    DOUBLE_COLON    # ::
    DOT             # .
    ARROW           # ->
    TYPE_ARROW      # =>

    # 특수
    AT              # @
    QUESTION        # ?
    BACKSLASH       # \
end

# ======================== 토큰 구조체 ========================

struct Position
    line::Int
    column::Int
    start::Int
    length::Int
end

struct Token
    type::TokenType
    value::String
    pos::Position
end

# ======================== 한글 키워드 매핑 ========================

const KOREAN_KEYWORDS = Dict(
    "함수" => KW_함수,
    "변수" => KW_변수,
    "상수" => KW_상수,
    "만약" => KW_만약,
    "아니면" => KW_아니면,
    "아니면만약" => KW_아니면만약,
    "반복" => KW_반복,
    "동안" => KW_동안,
    "반환" => KW_반환,
    "끝" => KW_끝,
    "없음" => KW_없음,
    "참" => KW_참,
    "거짓" => KW_거짓,
    "구조체" => KW_구조체,
    "배열" => KW_배열,
    "딕셔너리" => KW_딕셔너리,
    "부터" => KW_부터,
    "까지" => KW_까지,
    "하기" => KW_하기,
    "시작" => KW_시작,
    "사용" => KW_사용,
    "포함" => KW_포함,
    "전역" => KW_전역,
    "지역" => KW_지역,
)

# ======================== Lexer 구조체 ========================

mutable struct FreeLiuliaLexer
    source::String
    pos::Int
    line::Int
    column::Int
    tokens::Vector{Token}
end

function FreeLiuliaLexer(source::String)
    return FreeLiuliaLexer(source, 1, 1, 1, Token[])
end

# ======================== 한글 문자 판별 ========================

"""한글 음절 범위 확인"""
function is_korean(c::Char)::Bool
    code = UInt32(c)
    # 한글: U+AC00 ~ U+D7AF
    return 0xAC00 ≤ code ≤ 0xD7AF
end

"""한글 또는 ASCII 알파벳"""
function is_korean_or_alpha(c::Char)::Bool
    return is_korean(c) || isalpha(c) || c == '_'
end

"""한글 또는 ASCII 알파벳 또는 숫자"""
function is_korean_or_alnum(c::Char)::Bool
    return is_korean(c) || isalnum(c) || c == '_'
end

# ======================== 헬퍼 함수 ========================

"""현재 문자 반환"""
function current_char(lexer::FreeLiuliaLexer)::Char
    if lexer.pos > length(lexer.source)
        return Char(0)
    end
    return lexer.source[lexer.pos]
end

"""다음 문자 반환"""
function peek_char(lexer::FreeLiuliaLexer, offset=1)::Char
    pos = lexer.pos + offset
    if pos > length(lexer.source)
        return Char(0)
    end
    return lexer.source[pos]
end

"""한 문자 전진 (UTF-8 인식)"""
function advance!(lexer::FreeLiuliaLexer)
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
function skip_whitespace!(lexer::FreeLiuliaLexer)
    while current_char(lexer) in (' ', '\t', '\r')
        advance!(lexer)
    end
end

"""주석 건너뛰기 (#로 시작)"""
function skip_comment!(lexer::FreeLiuliaLexer)
    if current_char(lexer) == '#'
        while current_char(lexer) != '\n' && current_char(lexer) != Char(0)
            advance!(lexer)
        end
    end
end

"""토큰 추가"""
function add_token!(lexer::FreeLiuliaLexer, type::TokenType, value::String)
    token = Token(type, value, Position(lexer.line, lexer.column, lexer.pos - length(value), length(value)))
    push!(lexer.tokens, token)
end

# ======================== 토큰화 함수 ========================

"""정수 또는 실수 토큰화"""
function tokenize_number!(lexer::FreeLiuliaLexer)
    start_pos = lexer.pos

    while isdigit(current_char(lexer))
        advance!(lexer)
    end

    if current_char(lexer) == '.' && isdigit(peek_char(lexer))
        advance!(lexer)
        while isdigit(current_char(lexer))
            advance!(lexer)
        end

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
function tokenize_string!(lexer::FreeLiuliaLexer)
    quote_char = current_char(lexer)
    advance!(lexer)

    value = ""
    while current_char(lexer) != quote_char && current_char(lexer) != Char(0)
        if current_char(lexer) == '\\'
            advance!(lexer)
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
        advance!(lexer)
    end

    add_token!(lexer, STRING, value)
end

"""심볼 토큰화"""
function tokenize_symbol!(lexer::FreeLiuliaLexer)
    advance!(lexer)  # ':' 건너뛰기

    start_pos = lexer.pos
    while is_korean_or_alnum(current_char(lexer)) || current_char(lexer) in ('?', '!')
        advance!(lexer)
    end

    value = lexer.source[start_pos:lexer.pos-1]
    add_token!(lexer, SYMBOL, value)
end

"""한글 또는 영문 식별자/키워드 토큰화"""
function tokenize_identifier!(lexer::FreeLiuliaLexer)
    start_pos = lexer.pos

    while is_korean_or_alnum(current_char(lexer))
        advance!(lexer)
    end

    value = lexer.source[start_pos:lexer.pos-1]

    # 한글 키워드 확인
    if haskey(KOREAN_KEYWORDS, value)
        add_token!(lexer, KOREAN_KEYWORDS[value], value)
    else
        add_token!(lexer, IDENTIFIER, value)
    end
end

# ======================== 메인 토크나이저 ========================

"""FreeLiulia 소스 코드 토크나이저"""
function tokenize_freeliulia(source::String)::Vector{Token}
    lexer = FreeLiuliaLexer(source)

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
        # 한글 또는 영문 식별자/키워드
        elseif is_korean_or_alpha(char)
            tokenize_identifier!(lexer)
        # 연산자 & 구분자
        elseif char == '+'
            advance!(lexer)
            if current_char(lexer) == '+'
                advance!(lexer)
                add_token!(lexer, CONCAT, "++")
            elseif current_char(lexer) == '='
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
            add_token!(lexer, PERCENT, "%")
        elseif char == '^'
            advance!(lexer)
            add_token!(lexer, CARET, "^")
        elseif char == '='
            advance!(lexer)
            if current_char(lexer) == '='
                advance!(lexer)
                add_token!(lexer, EQ_EQ, "==")
            elseif current_char(lexer) == '>'
                advance!(lexer)
                add_token!(lexer, TYPE_ARROW, "=>")
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
            else
                add_token!(lexer, LT, "<")
            end
        elseif char == '>'
            advance!(lexer)
            if current_char(lexer) == '='
                advance!(lexer)
                add_token!(lexer, GE, ">=")
            else
                add_token!(lexer, GT, ">")
            end
        elseif char == '&'
            advance!(lexer)
            if current_char(lexer) == '&'
                advance!(lexer)
                add_token!(lexer, AND_AND, "&&")
            else
                add_token!(lexer, NOT, "&")
            end
        elseif char == '|'
            advance!(lexer)
            if current_char(lexer) == '|'
                advance!(lexer)
                add_token!(lexer, OR_OR, "||")
            else
                add_token!(lexer, NOT, "|")
            end
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
        elseif char == '?'
            advance!(lexer)
            add_token!(lexer, QUESTION, "?")
        elseif char == '\\'
            advance!(lexer)
            add_token!(lexer, BACKSLASH, "\\")
        else
            advance!(lexer)
        end
    end

    add_token!(lexer, EOF, "")
    return lexer.tokens
end

# ======================== 유틸리티 ========================

"""토큰 출력"""
function print_tokens_ko(tokens::Vector{Token})
    for token in tokens
        println("$(token.type): '$(token.value)' @ $(token.pos.line):$(token.pos.column)")
    end
end

# Export
export FreeLiuliaLexer, Token, Position, TokenType
export tokenize_freeliulia, print_tokens_ko
export is_korean, is_korean_or_alpha, is_korean_or_alnum
export KOREAN_KEYWORDS
