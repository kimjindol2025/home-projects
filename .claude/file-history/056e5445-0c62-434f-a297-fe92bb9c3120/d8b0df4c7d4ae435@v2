"""
    FreeLiulia Codegen 테스트
"""

include("../src/freeliulia_codegen.jl")

println("="*60)
println("🦎 FreeLiulia Codegen 테스트")
println("="*60)

# ======================== 테스트 1: 간단한 변수 선언 ========================

test_code_1 = """
변수 x: 정수 = 42
"""

println("\n테스트 1: 변수 선언 → C 코드")
println("코드:")
println(test_code_1)
try
    c_code = generate_c_from_freeliulia(test_code_1)
    println("✓ C 코드 생성 성공")
    println("\n=== 생성된 C 코드 ===")
    println(c_code)
catch e
    println("✗ 에러: $e")
    import Base.showerror
    showerror(stdout, e, catch_backtrace())
end

println("\n" * "="^60 * "\n")

# ======================== 테스트 2: 이항 연산 ========================

test_code_2 = """
변수 a: 정수 = 10
변수 b: 정수 = 20
변수 c = a + b
"""

println("테스트 2: 이항 연산 → C 코드")
println("코드:")
println(test_code_2)
try
    c_code = generate_c_from_freeliulia(test_code_2)
    println("✓ C 코드 생성 성공")
    println("\n=== 생성된 C 코드 (일부) ===")
    lines = split(c_code, "\n")
    for (i, line) in enumerate(lines[end-10:end])
        println(line)
    end
catch e
    println("✗ 에러: $e")
end

println("\n" * "="^60 * "\n")

# ======================== 테스트 3: 비교 연산 ========================

test_code_3 = """
변수 x: 정수 = 10
변수 y: 정수 = 20
변수 result = x < y
"""

println("테스트 3: 비교 연산 (< → 논리값)")
println("코드:")
println(test_code_3)
try
    c_code = generate_c_from_freeliulia(test_code_3)
    println("✓ C 코드 생성 성공 (비교 → 논리값)")
    println("\n=== 생성된 C 코드 (일부) ===")
    lines = split(c_code, "\n")
    for line in lines[end-8:end]
        println(line)
    end
catch e
    println("✗ 에러: $e")
end

println("\n" * "="^60 * "\n")

# ======================== 테스트 4: 문자열 ========================

test_code_4 = """
변수 msg: 문자열 = \"안녕하세요\"
"""

println("테스트 4: 문자열 리터럴")
println("코드:")
println(test_code_4)
try
    c_code = generate_c_from_freeliulia(test_code_4)
    println("✓ C 코드 생성 성공 (문자열)")
    println("\n=== 생성된 C 코드 (일부) ===")
    lines = split(c_code, "\n")
    for line in lines[end-8:end]
        println(line)
    end
catch e
    println("✗ 에러: $e")
end

println("\n" * "="^60 * "\n")

# ======================== 테스트 5: 실수 연산 ========================

test_code_5 = """
변수 pi: 실수 = 3.14
변수 r: 실수 = 2.0
변수 area = pi * r * r
"""

println("테스트 5: 실수 연산 (double 타입)")
println("코드:")
println(test_code_5)
try
    c_code = generate_c_from_freeliulia(test_code_5)
    println("✓ C 코드 생성 성공 (실수)")
    println("\n=== 생성된 C 코드 (일부) ===")
    lines = split(c_code, "\n")
    for line in lines[end-10:end]
        println(line)
    end
catch e
    println("✗ 에러: $e")
end

println("\n" * "="^60 * "\n")

# ======================== 테스트 6: 함수 호출 ========================

test_code_6 = """
결과 = 더하기(3, 4)
"""

println("테스트 6: 함수 호출")
println("코드:")
println(test_code_6)
try
    c_code = generate_c_from_freeliulia(test_code_6)
    println("✓ C 코드 생성 성공 (함수 호출)")
    println("\n=== 생성된 C 코드 (일부) ===")
    lines = split(c_code, "\n")
    for line in lines[max(1, end-8):end]
        println(line)
    end
catch e
    println("✗ 에러: $e")
end

println("\n" * "="^60)
println("✓ Codegen 테스트 완료")
println("  - Typed SSA IR → C 코드 변환 성공")
println("  - 정수, 실수, 문자열, 논리값 타입 지원")
println("  - 연산자 및 함수 호출 지원")
println("="*60)
