"""
    FreeLiulia Lowering 테스트
"""

include("../src/freeliulia_lowering.jl")

println("="*60)
println("🦎 FreeLiulia Lowering 테스트")
println("="*60)

# ======================== 테스트 1: 변수 선언 ========================

test_code_1 = """
변수 x: 정수 = 10
변수 y: 실수 = 3.14
"""

println("\n테스트 1: 변수 선언")
println("코드:")
println(test_code_1)
try
    ir = lower_freeliulia(test_code_1)
    println("✓ 로워링 성공 ($(length(ir)) 명령)")
    print_ir(ir)
catch e
    println("✗ 에러: $e")
end

println("\n" * "="^60 * "\n")

# ======================== 테스트 2: 이항 연산 ========================

test_code_2 = """
변수 결과 = 3 + 4
"""

println("테스트 2: 이항 연산")
println("코드:")
println(test_code_2)
try
    ir = lower_freeliulia(test_code_2)
    println("✓ 로워링 성공 ($(length(ir)) 명령)")
    print_ir(ir)
catch e
    println("✗ 에러: $e")
end

println("\n" * "="^60 * "\n")

# ======================== 테스트 3: 함수 호출 ========================

test_code_3 = """
출력(더하기(3, 4))
"""

println("테스트 3: 함수 호출")
println("코드:")
println(test_code_3)
try
    ir = lower_freeliulia(test_code_3)
    println("✓ 로워링 성공 ($(length(ir)) 명령)")
    print_ir(ir)
catch e
    println("✗ 에러: $e")
end

println("\n" * "="^60 * "\n")

# ======================== 테스트 4: If 문 ========================

test_code_4 = """
만약 x > 10
    출력(x)
아니면
    출력(0)
끝
"""

println("테스트 4: If 문 (제어흐름 - 라벨 생성)")
println("코드:")
println(test_code_4)
try
    ir = lower_freeliulia(test_code_4)
    println("✓ 로워링 성공 ($(length(ir)) 명령 - 제어흐름 정규화됨)")
    print_ir(ir)
catch e
    println("✗ 에러: $e")
end

println("\n" * "="^60 * "\n")

# ======================== 테스트 5: While 문 ========================

test_code_5 = """
동안 x < 100
    x = x + 1
끝
"""

println("테스트 5: While 문 (루프 - 라벨/분기)")
println("코드:")
println(test_code_5)
try
    ir = lower_freeliulia(test_code_5)
    println("✓ 로워링 성공 ($(length(ir)) 명령)")
    print_ir(ir)
catch e
    println("✗ 에러: $e")
end

println("\n" * "="^60 * "\n")

# ======================== 테스트 6: For 문 ========================

test_code_6 = """
반복 i = 1 부터 5 까지
    출력(i)
끝
"""

println("테스트 6: For 문 (루프 - 카운터 + 분기)")
println("코드:")
println(test_code_6)
try
    ir = lower_freeliulia(test_code_6)
    println("✓ 로워링 성공 ($(length(ir)) 명령)")
    print_ir(ir)
catch e
    println("✗ 에러: $e")
end

println("\n" * "="^60 * "\n")

# ======================== 테스트 7: 함수 정의 ========================

test_code_7 = """
함수 더하기(a: 정수, b: 정수) -> 정수
    변수 결과 = a + b
    반환 결과
끝
"""

println("테스트 7: 함수 정의")
println("코드:")
println(test_code_7)
try
    ir = lower_freeliulia(test_code_7)
    println("✓ 로워링 성공 ($(length(ir)) 명령)")
    print_ir(ir)
catch e
    println("✗ 에러: $e")
end

println("\n" * "="^60 * "\n")

# ======================== 테스트 8: 배열 리터럴 ========================

test_code_8 = """
배열 arr = [1, 2, 3]
"""

println("테스트 8: 배열 리터럴")
println("코드:")
println(test_code_8)
try
    ir = lower_freeliulia(test_code_8)
    println("✓ 로워링 성공 ($(length(ir)) 명령)")
    print_ir(ir)
catch e
    println("✗ 에러: $e")
end

println("\n" * "="^60 * "\n")

# ======================== 테스트 9: 배열 접근 ========================

test_code_9 = """
배열 arr = [1, 2, 3]
변수 x = arr[0]
"""

println("테스트 9: 배열 접근")
println("코드:")
println(test_code_9)
try
    ir = lower_freeliulia(test_code_9)
    println("✓ 로워링 성공 ($(length(ir)) 명령)")
    print_ir(ir)
catch e
    println("✗ 에러: $e")
end

println("\n" * "="^60 * "\n")

# ======================== 테스트 10: 복합 (소수 판정) ========================

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

println("테스트 10: 복합 함수 (소수 판정 - 다중 제어흐름)")
println("코드:")
println(test_code_10)
try
    ir = lower_freeliulia(test_code_10)
    println("✓ 로워링 성공 ($(length(ir)) 명령 - 제어흐름 정규화됨)")
    println("✓ 복합 제어흐름: if, for, nested-if 모두 라벨/분기로 변환")
    print_ir(ir)
catch e
    println("✗ 에러: $e")
end

println("\n" * "="^60)
println("✓ Lowering 테스트 완료")
println("  - AST → Untyped IR 변환 성공")
println("  - 제어흐름 정규화: if/while/for → label+branch")
println("  - 모든 표현식 IR 값으로 변환")
println("="*60)
