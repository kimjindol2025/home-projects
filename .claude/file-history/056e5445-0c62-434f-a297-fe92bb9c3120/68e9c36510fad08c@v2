"""
    Julia Self-Hosting Compiler

완전한 Julia 컴파일러 파이프라인
Tokenize → Parse → Lowering → Type Inference → Codegen
"""

include("codegen.jl")

# ======================== 통합 컴파일러 ========================

struct CompilationResult
    success::Bool
    c_code::Union{Nothing, String}
    ir::Union{Nothing, Vector{IRInstruction}}
    typed_ir::Union{Nothing, Vector{TypedIRInstruction}}
    ast::Union{Nothing, Vector{ASTNode}}
    error_message::String
end

"""주요 컴파일 함수"""
function compile(source::String)::CompilationResult
    try
        # Phase 1: Tokenize
        tokens = tokenize(source)

        # Phase 2: Parse
        ast = parse_program(Parser(tokens))

        # Phase 3: Lowering
        ir = lower(ast)

        # Phase 4: Type Inference
        typed_ir = infer_types(ir)

        # Phase 5: Optimization
        optimized_ir = optimize_ir(typed_ir)

        # Phase 6: Codegen
        c_code = generate_code(optimized_ir)

        return CompilationResult(true, c_code, ir, typed_ir, ast, "")

    catch e
        error_msg = "Compilation Error: $(e)"
        return CompilationResult(false, nothing, nothing, nothing, nothing, error_msg)
    end
end

"""AST만 생성"""
function parse_only(source::String)::CompilationResult
    try
        tokens = tokenize(source)
        ast = parse_program(Parser(tokens))
        return CompilationResult(true, nothing, nothing, nothing, ast, "")
    catch e
        return CompilationResult(false, nothing, nothing, nothing, nothing, "Parse Error: $(e)")
    end
end

"""IR까지 생성"""
function lower_only(source::String)::CompilationResult
    try
        tokens = tokenize(source)
        ast = parse_program(Parser(tokens))
        ir = lower(ast)
        return CompilationResult(true, nothing, ir, nothing, ast, "")
    catch e
        return CompilationResult(false, nothing, nothing, nothing, nothing, "Lowering Error: $(e)")
    end
end

"""Typed IR까지 생성"""
function infer_only(source::String)::CompilationResult
    try
        tokens = tokenize(source)
        ast = parse_program(Parser(tokens))
        ir = lower(ast)
        typed_ir = infer_types(ir)
        return CompilationResult(true, nothing, ir, typed_ir, ast, "")
    catch e
        return CompilationResult(false, nothing, nothing, nothing, nothing, "Type Inference Error: $(e)")
    end
end

# ======================== 결과 출력 ========================

"""컴파일 결과 출력"""
function print_result(result::CompilationResult)
    if !result.success
        println("❌ Compilation Failed")
        println("Error: $(result.error_message)")
        return
    end

    println("✅ Compilation Successful")
    println()

    if result.ast !== nothing
        println("=== AST ===")
        for (idx, node) in enumerate(result.ast)
            println("$idx: $(ast_to_string(node))")
        end
        println()
    end

    if result.ir !== nothing
        println("=== Untyped IR ===")
        print_ir(result.ir)
        println()
    end

    if result.typed_ir !== nothing
        println("=== Typed IR ===")
        print_typed_ir(result.typed_ir)
        println()
    end

    if result.c_code !== nothing
        println("=== Generated C Code ===")
        println(result.c_code)
    end
end

# ======================== 테스트 프로그램 ========================

"""간단한 덧셈"""
function test_simple_addition()
    println("\n" * "=" * 60)
    println("TEST 1: Simple Addition")
    println("=" * 60)

    source = "x = 3\ny = 4\nz = x + y"
    result = compile(source)
    print_result(result)
end

"""함수 정의"""
function test_function()
    println("\n" * "=" * 60)
    println("TEST 2: Function Definition")
    println("=" * 60)

    source = """
    function add(a, b)
        return a + b
    end
    result = add(5, 6)
    """
    result = compile(source)
    print_result(result)
end

"""조건문"""
function test_if_statement()
    println("\n" * "=" * 60)
    println("TEST 3: If Statement")
    println("=" * 60)

    source = """
    x = 10
    if x > 5
        y = 1
    else
        y = 0
    end
    """
    result = compile(source)
    print_result(result)
end

"""반복문"""
function test_while_loop()
    println("\n" * "=" * 60)
    println("TEST 4: While Loop")
    println("=" * 60)

    source = """
    x = 0
    while x < 5
        x = x + 1
    end
    """
    result = compile(source)
    print_result(result)
end

"""배열"""
function test_array()
    println("\n" * "=" * 60)
    println("TEST 5: Array Literal")
    println("=" * 60)

    source = """
    arr = [1, 2, 3, 4, 5]
    x = arr[0]
    """
    result = compile(source)
    print_result(result)
end

"""복잡한 표현식"""
function test_complex_expression()
    println("\n" * "=" * 60)
    println("TEST 6: Complex Expression")
    println("=" * 60)

    source = """
    a = 10
    b = 20
    c = (a + b) * 2
    d = c / 4
    """
    result = compile(source)
    print_result(result)
end

# ======================== 벤치마크 ========================

"""컴파일 성능 측정"""
function benchmark_compilation()
    println("\n" * "=" * 60)
    println("BENCHMARK: Compilation Performance")
    println("=" * 60)

    test_cases = [
        ("Simple", "x = 1"),
        ("Arithmetic", "x = 1 + 2 * 3 - 4 / 5"),
        ("Function", "function f(x)\n    return x + 1\nend"),
        ("Control Flow", "if x > 0\n    y = 1\nelse\n    y = 0\nend"),
        ("Multiple Statements", "a=1\nb=2\nc=3\nd=a+b+c"),
    ]

    for (name, code) in test_cases
        time_start = time()
        result = compile(code)
        time_elapsed = (time() - time_start) * 1000  # ms

        status = result.success ? "✅" : "❌"
        println("$status $name: $(round(time_elapsed, digits=2))ms")
    end
end

# ======================== 통계 ========================

"""컴파일러 통계"""
function compiler_statistics()
    println("\n" * "=" * 60)
    println("COMPILER STATISTICS")
    println("=" * 60)

    # 토큰 타입 개수
    token_types = length(instances(TokenType))
    println("Token Types: $token_types")

    # AST 노드 타입 개수 (근사)
    println("AST Node Types: 15+")

    # IR 명령 타입
    ir_types = 11
    println("IR Instruction Types: $ir_types")

    # Typed IR 명령 타입
    typed_ir_types = 11
    println("Typed IR Instruction Types: $typed_ir_types")

    println()
    println("Pipeline Stages: 6")
    println("  1. Tokenize (Lexer)")
    println("  2. Parse (Parser)")
    println("  3. Lower (Lowering)")
    println("  4. Type Inference")
    println("  5. Optimize")
    println("  6. Codegen (C)")
end

# ======================== 메인 실행 ========================

"""모든 테스트 실행"""
function run_all_tests()
    println("╔" * "=" * 58 * "╗")
    println("║" * " " * "JULIA SELF-HOSTING COMPILER - ALL PHASES" * " " * "║")
    println("╚" * "=" * 58 * "╝")

    test_simple_addition()
    test_function()
    test_if_statement()
    test_while_loop()
    test_array()
    test_complex_expression()

    compiler_statistics()
    benchmark_compilation()

    println("\n" * "=" * 60)
    println("ALL TESTS COMPLETED")
    println("=" * 60)
end

# Export
export compile, parse_only, lower_only, infer_only
export CompilationResult, print_result
export test_simple_addition, test_function, test_if_statement
export test_while_loop, test_array, test_complex_expression
export benchmark_compilation, compiler_statistics, run_all_tests

# 메인 프로그램 진입점
if abspath(PROGRAM_FILE) == @__FILE__
    run_all_tests()
end
