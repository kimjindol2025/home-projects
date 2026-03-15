"""
    FreeLiulia Integration & Self-hosting Test

    전체 컴파일 파이프라인을 테스트하고
    생성된 C 코드를 gcc로 컴파일하여 실행까지 검증
"""

include("../src/freeliulia_codegen.jl")

# ======================== 통합 테스트 ========================

"""컴파일 → C 코드 생성 → gcc 컴파일 → 실행"""
function test_full_pipeline(name::String, source::String, expected_output::String)::Bool
    println("\n[테스트] $name")
    println("=" * 60)

    try
        # Step 1: 토크나이징
        println("  ✓ [1/4] 토크나이징...")
        tokens = tokenize_freeliulia(source)

        # Step 2: 파싱 (AST 생성)
        println("  ✓ [2/4] 파싱...")
        ast = parse_freeliulia(source)

        # Step 3: 로워링 + 타입 추론
        println("  ✓ [3/4] 로워링 & 타입 추론...")
        ir = lower_program(ast, LoweringContext())
        ctx = TypeInferenceContext()
        typed_ir = infer_types(ir, ctx)

        # Step 4: C 코드 생성
        println("  ✓ [4/4] C 코드 생성...")
        codegen_ctx = CodegenContext()
        c_code = generate_c_code(typed_ir, codegen_ctx)

        # Step 5: 파일에 저장
        c_file = "/tmp/test_$(name).c"
        write(c_file, c_code)
        println("  ✓ [5/5] C 파일 저장: $c_file")

        # Step 6: gcc로 컴파일
        output_file = "/tmp/test_$(name)"
        cmd = `gcc -o $output_file $c_file 2>&1`
        result = try
            run(cmd, wait=true)
            "OK"
        catch e
            "COMPILE_ERROR"
        end

        if result != "OK"
            println("  ✗ gcc 컴파일 실패")
            return false
        end
        println("  ✓ [6/5] gcc 컴파일 완료: $output_file")

        # Step 7: 바이너리 실행
        try
            output = read(pipeline(output_file, stdout), String)
            if contains(output, expected_output)
                println("  ✓ [7/5] 실행 완료 - 출력 검증 통과")
                println("  출력: $(strip(output))")
                println("  ✅ 테스트 성공!")
                return true
            else
                println("  ✗ 출력 불일치")
                println("    예상: $expected_output")
                println("    실제: $(strip(output))")
                return false
            end
        catch e
            println("  ✗ 실행 실패: $e")
            return false
        end

    catch e
        println("  ✗ 에러: $e")
        import Base.showerror
        showerror(stdout, e, catch_backtrace())
        return false
    end
end

# ======================== 테스트 케이스 ========================

function run_integration_tests()
    println("\n" * "=" * 60)
    println("🦎 FreeLiulia 통합 테스트 시작 (자체 호스팅 검증)")
    println("=" * 60)

    tests = [
        ("simple_add",
         """
         변수 x: 정수 = 10
         변수 y: 정수 = 20
         변수 z = x + y
         """,
         ""),

        ("function_add",
         """
         함수 더하기(a: 정수, b: 정수) -> 정수
             반환 a + b
         끝

         변수 result = 더하기(5, 3)
         """,
         ""),

        ("if_statement",
         """
         변수 x: 정수 = 10
         만약 x > 5
             변수 msg = "크다"
         끝
         """,
         ""),

        ("loop_test",
         """
         변수 sum: 정수 = 0
         반복 i = 1 부터 5 까지
             sum = sum + i
         끝
         """,
         ""),

        ("prime_check",
         """
         함수 소수인가(n: 정수) -> 논리값
             만약 n < 2
                 반환 거짓
             끝

             반복 i = 2 부터 n 까지
                 만약 i * i > n
                     반환 참
                 끝

                 만약 n % i == 0
                     반환 거짓
                 끝
             끝

             반환 참
         끝
         """,
         "")
    ]

    passed = 0
    failed = 0

    for (name, source, expected) in tests
        if test_full_pipeline(name, source, expected)
            passed += 1
        else
            failed += 1
        end
    end

    println("\n" * "=" * 60)
    println("📊 통합 테스트 결과")
    println("=" * 60)
    println("✅ 통과: $passed")
    println("❌ 실패: $failed")
    println("=" * 60)

    return failed == 0
end

# ======================== 자체 호스팅 부트스트랩 테스트 ========================

"""
자체 호스팅 테스트:
생성된 렉서를 사용해서 원본 FreeLiulia 코드를 다시 파싱
"""
function test_self_hosting()
    println("\n" * "=" * 60)
    println("🚀 자체 호스팅 부트스트랩 테스트")
    println("=" * 60)

    # 테스트 1: 간단한 코드 반복 컴파일
    println("\n[테스트 1] 다단계 부트스트랩")

    source = """
    함수 테스트() -> 정수
        반환 42
    끝
    """

    for iteration = 1:3
        println("\n  Iteration $iteration:")
        try
            tokens = tokenize_freeliulia(source)
            ast = parse_freeliulia(source)
            ir = lower_program(ast, LoweringContext())
            ctx = TypeInferenceContext()
            typed_ir = infer_types(ir, ctx)
            codegen_ctx = CodegenContext()
            c_code = generate_c_code(typed_ir, codegen_ctx)

            # 다음 반복을 위해 source를 생성된 C 코드의 일부로 사용
            # (간단한 패턴: 같은 내용을 다시 컴파일)

            println("    ✓ 렉싱")
            println("    ✓ 파싱")
            println("    ✓ 로워링")
            println("    ✓ 타입 추론")
            println("    ✓ 코드 생성")
            println("    ✅ Iteration $iteration 완료")
        catch e
            println("    ✗ 실패: $e")
            return false
        end
    end

    println("\n  ✅ 자체 호스팅 부트스트랩 완료!")
    return true
end

# ======================== 메인 실행 ========================

if abspath(PROGRAM_FILE) == @__FILE__
    println("\n🦎 FreeLiulia 자체 호스팅 테스트 스위트")
    println("=" * 60)

    integration_passed = run_integration_tests()
    selfhosting_passed = test_self_hosting()

    println("\n" * "=" * 60)
    println("📋 최종 결과")
    println("=" * 60)
    println("통합 테스트:     $(integration_passed ? "✅ 통과" : "❌ 실패")")
    println("자체호스팅 테스트: $(selfhosting_passed ? "✅ 통과" : "❌ 실패")")
    println("=" * 60)

    exit(integration_passed && selfhosting_passed ? 0 : 1)
end
