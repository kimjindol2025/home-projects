#!/usr/bin/env julia

"""
    FreeLiulia 통합 컴파일러

    사용법:
        julia compile.jl <input_file> [output_file]

    예제:
        julia compile.jl examples/example1.fl examples/example1.c
"""

# 모듈 include
include("src/freeliulia_codegen.jl")

function main()
    if length(ARGS) < 1
        println("사용법: julia compile.jl <input_file> [output_file]")
        exit(1)
    end

    input_file = ARGS[1]
    output_file = length(ARGS) > 1 ? ARGS[2] : replace(input_file, ".fl" => ".c")

    if !isfile(input_file)
        println("에러: 파일을 찾을 수 없습니다: $input_file")
        exit(1)
    end

    println("🦎 FreeLiulia 컴파일러")
    println("=" * 60)
    println("입력:  $input_file")
    println("출력:  $output_file")
    println("=" * 60)

    try
        # 1. 소스 코드 읽기
        println("\n[1/5] 소스 코드 읽기...")
        source = read(input_file, String)
        println("✓ 읽음: $(length(source)) 바이트")

        # 2. 토크나이징
        println("\n[2/5] 토크나이징 (한글 토크나이저)...")
        tokens = tokenize_freeliulia(source)
        println("✓ 토큰: $(length(tokens))개")

        # 3. 파싱
        println("\n[3/5] 파싱 (AST 생성)...")
        ast = parse_freeliulia(source)
        println("✓ AST 생성 완료")

        # 4. 로워링 + 타입 추론
        println("\n[4/5] 로워링 및 타입 추론...")
        ir = lower_program(ast, LoweringContext())
        ctx = TypeInferenceContext()
        typed_ir = infer_types(ir, ctx)
        println("✓ IR 생성: $(length(typed_ir))개 명령")

        # 5. 코드 생성
        println("\n[5/5] C 코드 생성...")
        codegen_ctx = CodegenContext()
        c_code = generate_c_code(typed_ir, codegen_ctx)
        println("✓ C 코드 생성 완료")

        # 6. 파일에 저장
        println("\n[6/5] 파일에 저장...")
        write(output_file, c_code)
        println("✓ 저장 완료: $output_file")

        # 통계
        c_lines = length(split(c_code, "\n"))
        println("\n" * "=" * 60)
        println("✅ 컴파일 성공!")
        println("  한글 코드: $(length(source)) 바이트")
        println("  C 코드: $c_lines 줄")
        println("=" * 60)

    catch e
        println("\n❌ 컴파일 에러: $e")
        import Base.showerror
        showerror(stdout, e, catch_backtrace())
        exit(1)
    end
end

main()
