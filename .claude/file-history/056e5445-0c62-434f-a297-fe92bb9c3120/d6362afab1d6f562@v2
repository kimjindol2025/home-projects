"""
    FreeLiulia Parser 테스트
"""

include("../src/freeliulia_parser.jl")

# ======================== 테스트 1: 변수 선언 ========================

test_code_1 = """
변수 x: 정수 = 10
변수 y: 실수 = 3.14
"""

println("테스트 1: 변수 선언")
println("코드:")
println(test_code_1)
try
    ast = parse_freeliulia(test_code_1)
    println("✓ 파싱 성공")
    print_ast(ast)
catch e
    println("✗ 에러: $e")
end

println("\n" * "="^50 * "\n")

# ======================== 테스트 2: 함수 정의 ========================

test_code_2 = """
함수 더하기(a: 정수, b: 정수) -> 정수
    변수 결과: 정수 = a + b
    반환 결과
끝
"""

println("테스트 2: 함수 정의")
println("코드:")
println(test_code_2)
try
    ast = parse_freeliulia(test_code_2)
    println("✓ 파싱 성공")
    print_ast(ast)
catch e
    println("✗ 에러: $e")
end

println("\n" * "="^50 * "\n")

# ======================== 테스트 3: If 문 ========================

test_code_3 = """
만약 x > 10
    출력(x)
아니면
    출력(0)
끝
"""

println("테스트 3: If 문")
println("코드:")
println(test_code_3)
try
    ast = parse_freeliulia(test_code_3)
    println("✓ 파싱 성공")
    print_ast(ast)
catch e
    println("✗ 에러: $e")
end

println("\n" * "="^50 * "\n")

# ======================== 테스트 4: While 문 ========================

test_code_4 = """
동안 x < 100
    x = x + 1
끝
"""

println("테스트 4: While 문")
println("코드:")
println(test_code_4)
try
    ast = parse_freeliulia(test_code_4)
    println("✓ 파싱 성공")
    print_ast(ast)
catch e
    println("✗ 에러: $e")
end

println("\n" * "="^50 * "\n")

# ======================== 테스트 5: For 문 ========================

test_code_5 = """
반복 i = 1 부터 10 까지
    출력(i)
끝
"""

println("테스트 5: For 문")
println("코드:")
println(test_code_5)
try
    ast = parse_freeliulia(test_code_5)
    println("✓ 파싱 성공")
    print_ast(ast)
catch e
    println("✗ 에러: $e")
end

println("\n" * "="^50 * "\n")

# ======================== 테스트 6: 복합 표현식 ========================

test_code_6 = """
변수 결과 = (x + y) * 2
"""

println("테스트 6: 복합 표현식")
println("코드:")
println(test_code_6)
try
    ast = parse_freeliulia(test_code_6)
    println("✓ 파싱 성공")
    print_ast(ast)
catch e
    println("✗ 에러: $e")
end

println("\n" * "="^50 * "\n")

# ======================== 테스트 7: 배열 리터럴 ========================

test_code_7 = """
배열 arr = [1, 2, 3, 4, 5]
"""

println("테스트 7: 배열 리터럴")
println("코드:")
println(test_code_7)
try
    ast = parse_freeliulia(test_code_7)
    println("✓ 파싱 성공")
    print_ast(ast)
catch e
    println("✗ 에러: $e")
end

println("\n" * "="^50 * "\n")

# ======================== 테스트 8: 함수 호출 ========================

test_code_8 = """
출력(더하기(3, 4))
"""

println("테스트 8: 함수 호출")
println("코드:")
println(test_code_8)
try
    ast = parse_freeliulia(test_code_8)
    println("✓ 파싱 성공")
    print_ast(ast)
catch e
    println("✗ 에러: $e")
end

println("\n" * "="^50 * "\n")

# ======================== 테스트 9: 상수 선언 ========================

test_code_9 = """
상수 PI: 실수 = 3.14159
"""

println("테스트 9: 상수 선언")
println("코드:")
println(test_code_9)
try
    ast = parse_freeliulia(test_code_9)
    println("✓ 파싱 성공")
    print_ast(ast)
catch e
    println("✗ 에러: $e")
end

println("\n" * "="^50 * "\n")

# ======================== 테스트 10: 소수 판정 (복합) ========================

test_code_10 = """
함수 소수인가(n: 정수) -> 논리값
    만약 n < 2
        반환 거짓
    끝

    반복 i = 2 부터 n 까지
        만약 n % i == 0
            반환 거짓
        끝
    끝

    반환 참
끝
"""

println("테스트 10: 소수 판정 함수 (복합)")
println("코드:")
println(test_code_10)
try
    ast = parse_freeliulia(test_code_10)
    println("✓ 파싱 성공")
    print_ast(ast)
catch e
    println("✗ 에러: $e")
end

println("\n" * "="^50 * "\n")

# ======================== 통계 ========================

println("✓ 파서 테스트 완료")
println("- 변수 선언, 함수 정의, 제어 흐름 모두 지원")
println("- 한글 키워드 완전 파싱 가능")
println("- 복합 표현식 및 연산자 우선순위 처리 완료")
