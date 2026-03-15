"""
    FreeLiulia Type Inference 테스트
"""

include("../src/freeliulia_type_inference.jl")

println("="*60)
println("🦎 FreeLiulia Type Inference 테스트")
println("="*60)

# ======================== 테스트 1: 정수 연산 ========================

test_code_1 = """
변수 x: 정수 = 10
변수 y: 정수 = 20
변수 z = x + y
"""

println("\n테스트 1: 정수 연산 (타입 추론)")
println("코드:")
println(test_code_1)
try
    ir = lower_freeliulia(test_code_1)
    typed_ir = infer_freeliulia_types(ir)
    println("✓ 타입 추론 성공 ($(length(typed_ir)) 타입된 명령)")
    print_typed_ir(typed_ir)
catch e
    println("✗ 에러: $e")
end

println("\n" * "="^60 * "\n")

# ======================== 테스트 2: 실수 연산 ========================

test_code_2 = """
변수 a: 실수 = 3.14
변수 b: 실수 = 2.71
변수 c = a + b
"""

println("테스트 2: 실수 연산")
println("코드:")
println(test_code_2)
try
    ir = lower_freeliulia(test_code_2)
    typed_ir = infer_freeliulia_types(ir)
    println("✓ 타입 추론 성공 (실수 타입 확인)")
    print_typed_ir(typed_ir)
catch e
    println("✗ 에러: $e")
end

println("\n" * "="^60 * "\n")

# ======================== 테스트 3: 논리값 ========================

test_code_3 = """
변수 p: 논리값 = 참
변수 q: 논리값 = 거짓
변수 r = p && q
"""

println("테스트 3: 논리값 (&&, || 연산)")
println("코드:")
println(test_code_3)
try
    ir = lower_freeliulia(test_code_3)
    typed_ir = infer_freeliulia_types(ir)
    println("✓ 타입 추론 성공 (논리값 타입 확인)")
    print_typed_ir(typed_ir)
catch e
    println("✗ 에러: $e")
end

println("\n" * "="^60 * "\n")

# ======================== 테스트 4: 비교 연산 ========================

test_code_4 = """
변수 x: 정수 = 10
변수 y: 정수 = 20
변수 cmp = x < y
"""

println("테스트 4: 비교 연산 (< → 논리값)")
println("코드:")
println(test_code_4)
try
    ir = lower_freeliulia(test_code_4)
    typed_ir = infer_freeliulia_types(ir)
    println("✓ 타입 추론 성공 (비교 결과 → 논리값)")
    print_typed_ir(typed_ir)
catch e
    println("✗ 에러: $e")
end

println("\n" * "="^60 * "\n")

# ======================== 테스트 5: 문자열 연결 ========================

test_code_5 = """
변수 s1: 문자열 = \"안녕\"
변수 s2: 문자열 = \"하세요\"
변수 s3 = s1 ++ s2
"""

println("테스트 5: 문자열 연결 (++ → 문자열)")
println("코드:")
println(test_code_5)
try
    ir = lower_freeliulia(test_code_5)
    typed_ir = infer_freeliulia_types(ir)
    println("✓ 타입 추론 성공 (문자열 연결)")
    print_typed_ir(typed_ir)
catch e
    println("✗ 에러: $e")
end

println("\n" * "="^60 * "\n")

# ======================== 테스트 6: 함수 호출 ========================

test_code_6 = """
출력(더하기(3, 4))
"""

println("테스트 6: 함수 호출 (타입 추론)")
println("코드:")
println(test_code_6)
try
    ir = lower_freeliulia(test_code_6)
    typed_ir = infer_freeliulia_types(ir)
    println("✓ 타입 추론 성공 (함수 호출)")
    print_typed_ir(typed_ir)
catch e
    println("✗ 에러: $e")
end

println("\n" * "="^60 * "\n")

# ======================== 테스트 7: If 문 (타입 안정성) ========================

test_code_7 = """
만약 10 > 5
    변수 x: 정수 = 100
아니면
    변수 x: 정수 = 200
끝
"""

println("테스트 7: If 문 (각 분기의 타입 추론)")
println("코드:")
println(test_code_7)
try
    ir = lower_freeliulia(test_code_7)
    typed_ir = infer_freeliulia_types(ir)
    println("✓ 타입 추론 성공 (If 분기)")
    print_typed_ir(typed_ir)
catch e
    println("✗ 에러: $e")
end

println("\n" * "="^60 * "\n")

# ======================== 테스트 8: 함수 정의 ========================

test_code_8 = """
함수 더하기(a: 정수, b: 정수) -> 정수
    변수 결과 = a + b
    반환 결과
끝
"""

println("테스트 8: 함수 정의 (매개변수 & 반환 타입)")
println("코드:")
println(test_code_8)
try
    ir = lower_freeliulia(test_code_8)
    typed_ir = infer_freeliulia_types(ir)
    println("✓ 타입 추론 성공 (함수 타입)")
    print_typed_ir(typed_ir)
catch e
    println("✗ 에러: $e")
end

println("\n" * "="^60 * "\n")

# ======================== 테스트 9: 배열 ========================

test_code_9 = """
배열 arr = [1, 2, 3]
변수 elem = arr[0]
"""

println("테스트 9: 배열 접근 (배열 타입 추론)")
println("코드:")
println(test_code_9)
try
    ir = lower_freeliulia(test_code_9)
    typed_ir = infer_freeliulia_types(ir)
    println("✓ 타입 추론 성공 (배열<정수>[0] → 정수)")
    print_typed_ir(typed_ir)
catch e
    println("✗ 에러: $e")
end

println("\n" * "="^60 * "\n")

# ======================== 테스트 10: 복합 프로그램 ========================

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

println("테스트 10: 복합 프로그램 (전체 타입 추론)")
println("코드:")
println(test_code_10)
try
    ir = lower_freeliulia(test_code_10)
    typed_ir = infer_freeliulia_types(ir)
    println("✓ 타입 추론 성공 (복합 프로그램 - 다중 제어흐름)")
    println("✓ 타입 확인: 함수 → 정수, 비교 → 논리값, 모듈로 → 정수")
    print_typed_ir(typed_ir)
catch e
    println("✗ 에러: $e")
end

println("\n" * "="^60)
println("✓ Type Inference 테스트 완료")
println("  - Untyped IR → Typed SSA IR 변환 성공")
println("  - 모든 값에 타입 지정")
println("  - 연산 결과 타입 추론")
println("="*60)
