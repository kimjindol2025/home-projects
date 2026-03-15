"""
    FreeLiulia Complete Compiler with Advanced Optimizations

6단계 파이프라인:
1. Lexer       → 한글 토크나이저
2. Parser      → AST 생성
3. Lowering    → Untyped IR 생성
4. Type Infer  → Typed SSA IR 생성
5. Optimize    → 고급 최적화 (DCE, Inlining, Loop Unrolling)
6. Codegen     → C 코드 생성

Phase 6: 고급 최적화 통합
"""

include("freeliulia_lexer.jl")
include("freeliulia_parser.jl")
include("freeliulia_lowering.jl")
include("freeliulia_type_inference.jl")
include("freeliulia_optimizer.jl")
include("freeliulia_codegen.jl")

# ======================== 컴파일 결과 구조 ========================

struct FreeLiuliaCompilationResult
    success::Bool
    c_code::Union{Nothing, String}
    ir::Union{Nothing, Vector{Dict}}
    typed_ir::Union{Nothing, Vector{Dict}}
    optimized_ir::Union{Nothing, Vector{Dict}}
    ast::Union{Nothing, Vector{Dict}}
    error_message::String
    optimization_stats::Union{Nothing, Dict}
end

# ======================== 통합 컴파일러 (6단계) ========================

"""
    compile_with_optimization(source::String)

최적화를 포함한 완전한 6단계 컴파일 파이프라인
"""
function compile_with_optimization(source::String)::FreeLiuliaCompilationResult
    try
        # Phase 1: Lexing
        tokens = tokenize(source)

        # Phase 2: Parsing
        ast = parse_program(tokens)

        # Phase 3: Lowering
        ir = lower_to_ir(ast)

        # Phase 4: Type Inference
        typed_ir = type_inference(ir)

        # Phase 5: Optimization (새로 추가)
        optimized_ir = run_all_optimizations(typed_ir)

        optimization_stats = optimization_stats(typed_ir, optimized_ir)

        # Phase 6: Codegen
        c_code = codegen(optimized_ir)

        return FreeLiuliaCompilationResult(
            true, c_code, ir, typed_ir, optimized_ir, ast,
            "", optimization_stats
        )

    catch e
        error_msg = "Compilation Error: $(e)"
        return FreeLiuliaCompilationResult(
            false, nothing, nothing, nothing, nothing, nothing,
            error_msg, nothing
        )
    end
end

"""
    compile_without_optimization(source::String)

최적화 없이 5단계로 컴파일 (비교용)
"""
function compile_without_optimization(source::String)::FreeLiuliaCompilationResult
    try
        tokens = tokenize(source)
        ast = parse_program(tokens)
        ir = lower_to_ir(ast)
        typed_ir = type_inference(ir)
        c_code = codegen(typed_ir)

        return FreeLiuliaCompilationResult(
            true, c_code, ir, typed_ir, nothing, ast,
            "", nothing
        )

    catch e
        error_msg = "Compilation Error (no optimization): $(e)"
        return FreeLiuliaCompilationResult(
            false, nothing, nothing, nothing, nothing, nothing,
            error_msg, nothing
        )
    end
end

# ======================== 단계별 컴파일 함수 ========================

"""
    compile_until_lexing(source::String)

렉싱까지만 수행
"""
function compile_until_lexing(source::String)
    try
        tokens = tokenize(source)
        return (success=true, tokens=tokens, error="")
    catch e
        return (success=false, tokens=nothing, error="Lexing Error: $(e)")
    end
end

"""
    compile_until_parsing(source::String)

파싱까지만 수행
"""
function compile_until_parsing(source::String)
    try
        tokens = tokenize(source)
        ast = parse_program(tokens)
        return (success=true, ast=ast, error="")
    catch e
        return (success=false, ast=nothing, error="Parsing Error: $(e)")
    end
end

"""
    compile_until_lowering(source::String)

로워링까지만 수행
"""
function compile_until_lowering(source::String)
    try
        tokens = tokenize(source)
        ast = parse_program(tokens)
        ir = lower_to_ir(ast)
        return (success=true, ir=ir, error="")
    catch e
        return (success=false, ir=nothing, error="Lowering Error: $(e)")
    end
end

"""
    compile_until_type_inference(source::String)

타입 추론까지만 수행
"""
function compile_until_type_inference(source::String)
    try
        tokens = tokenize(source)
        ast = parse_program(tokens)
        ir = lower_to_ir(ast)
        typed_ir = type_inference(ir)
        return (success=true, typed_ir=typed_ir, error="")
    catch e
        return (success=false, typed_ir=nothing, error="Type Inference Error: $(e)")
    end
end

# ======================== 최적화 전후 비교 ========================

"""
    compare_optimizations(source::String)

최적화 전후 결과를 비교합니다.
"""
function compare_optimizations(source::String)
    try
        # 최적화 없이 컴파일
        result_no_opt = compile_without_optimization(source)
        if !result_no_opt.success
            return (success=false, error=result_no_opt.error_message)
        end

        # 최적화 포함 컴파일
        result_with_opt = compile_with_optimization(source)
        if !result_with_opt.success
            return (success=false, error=result_with_opt.error_message)
        end

        # 비교
        no_opt_lines = isnothing(result_no_opt.c_code) ? 0 : length(split(result_no_opt.c_code, "\n"))
        with_opt_lines = isnothing(result_with_opt.c_code) ? 0 : length(split(result_with_opt.c_code, "\n"))

        return (
            success=true,
            no_optimization=(
                c_code_lines=no_opt_lines,
                ir_size=length(result_no_opt.ir),
                typed_ir_size=length(result_no_opt.typed_ir)
            ),
            with_optimization=(
                c_code_lines=with_opt_lines,
                ir_size=length(result_with_opt.optimized_ir),
                typed_ir_size=length(result_with_opt.typed_ir)
            ),
            stats=result_with_opt.optimization_stats,
            reduction_percent=if no_opt_lines > 0
                round(100 * (no_opt_lines - with_opt_lines) / no_opt_lines; digits=2)
            else
                0
            end
        )

    catch e
        return (success=false, error="Comparison Error: $(e)")
    end
end

# ======================== 컴파일 결과 출력 ========================

"""
    print_compilation_result(result::FreeLiuliaCompilationResult)

컴파일 결과를 보기 좋게 출력합니다.
"""
function print_compilation_result(result::FreeLiuliaCompilationResult)
    if !result.success
        println("❌ Compilation Failed")
        println("Error: $(result.error_message)")
        return
    end

    println("✅ Compilation Successful (6-Stage Pipeline)")
    println()

    println("📊 Pipeline Statistics:")
    println("─" ^ 50)

    if result.ast !== nothing
        println("  AST Size: $(length(result.ast)) nodes")
    end

    if result.ir !== nothing
        println("  Untyped IR Size: $(length(result.ir)) instructions")
    end

    if result.typed_ir !== nothing
        println("  Typed SSA IR Size: $(length(result.typed_ir)) instructions")
    end

    if result.optimized_ir !== nothing
        println("  Optimized IR Size: $(length(result.optimized_ir)) instructions")

        if result.optimization_stats !== nothing
            stats = result.optimization_stats
            println("\n🔧 Optimization Details:")
            println("  ├─ Dead Code Eliminated: $(stats[:instructions_removed]) instructions")
            println("  ├─ Size Reduction: $(stats[:reduction_percent])%")
            println("  └─ Optimized Size: $(stats[:optimized_size]) instructions")
        end
    end

    if result.c_code !== nothing
        c_lines = length(split(result.c_code, "\n"))
        println("\n💻 Generated C Code:")
        println("  Generated $(c_lines) lines of C code")
        println("\n  First 20 lines:")
        println("  " * "─" ^ 46)
        for line in split(result.c_code, "\n")[1:min(20, length(split(result.c_code, "\n")))]
            println("  $line")
        end
    end

    println("\n" * "=" ^ 50)
end

"""
    print_comparison_result(result::Dict)

최적화 비교 결과를 출력합니다.
"""
function print_comparison_result(result::Dict)
    if !result[:success]
        println("❌ Comparison Failed")
        println("Error: $(result[:error])")
        return
    end

    println("📊 Optimization Comparison Report")
    println("=" ^ 50)

    no_opt = result[:no_optimization]
    with_opt = result[:with_optimization]
    stats = result[:stats]

    println("\n📈 C Code Size:")
    println("  Without Optimization: $(no_opt[:c_code_lines]) lines")
    println("  With Optimization:    $(with_opt[:c_code_lines]) lines")
    println("  Reduction:            $(result[:reduction_percent])%")

    println("\n📈 IR Size (Typed):")
    println("  Before Optimization: $(no_opt[:typed_ir_size]) instructions")
    println("  After Optimization:  $(with_opt[:ir_size]) instructions")
    println("  Removed:             $(stats[:instructions_removed]) instructions")

    println("\n🎯 Optimization Impact:")
    if stats[:reduction_percent] > 10
        println("  ✨ Excellent optimization achieved!")
    elseif stats[:reduction_percent] > 5
        println("  ✓ Good optimization results")
    else
        println("  - Minimal optimization opportunities")
    end

    println("\n" * "=" ^ 50)
end

end  # module
