"""
    FreeLiulia Advanced Optimizer

Typed IR에 대한 3가지 최적화 패스:
1. Dead Code Elimination (DCE) - 사용되지 않는 지시어 제거
2. Function Inlining - 작은 함수를 호출처에 직접 삽입
3. Loop Unrolling - 고정 반복 횟수 루프 펼치기
"""

include("freeliulia_type_inference.jl")

# ======================== 최적화 통합 ========================

"""
    run_all_optimizations(ir::Vector)

모든 최적화 패스를 순서대로 실행
"""
function run_all_optimizations(ir::Vector)
    # Pass 1: Dead Code Elimination
    ir = dead_code_elimination(ir)

    # Pass 2: Function Inlining (소형 함수)
    ir = inline_small_functions(ir)

    # Pass 3: Loop Unrolling
    ir = loop_unrolling(ir)

    return ir
end

# ======================== Dead Code Elimination ========================

"""
    dead_code_elimination(ir::Vector)

사용되지 않는 변수 정의와 지시어를 제거합니다.
참조 카운트가 0인 것은 제거 대상입니다.
"""
function dead_code_elimination(ir::Vector)
    if isempty(ir)
        return ir
    end

    # Step 1: 모든 값의 참조 카운트 계산
    ref_count = Dict{Int, Int}()
    defined_vars = Set{Int}()

    for instr in ir
        # 각 지시어에서 정의된 변수
        if haskey(instr, :id)
            defined_vars = push!(defined_vars, instr[:id])
            ref_count[instr[:id]] = 0
        end
    end

    # Step 2: 참조 카운트 증가
    for instr in ir
        if haskey(instr, :left_id)
            ref_count[instr[:left_id]] = get(ref_count, instr[:left_id], 0) + 1
        end
        if haskey(instr, :right_id)
            ref_count[instr[:right_id]] = get(ref_count, instr[:right_id], 0) + 1
        end
        if haskey(instr, :operand_id)
            ref_count[instr[:operand_id]] = get(ref_count, instr[:operand_id], 0) + 1
        end
        if haskey(instr, :func_id)
            ref_count[instr[:func_id]] = get(ref_count, instr[:func_id], 0) + 1
        end
        if haskey(instr, :args)
            for arg_id in instr[:args]
                ref_count[arg_id] = get(ref_count, arg_id, 0) + 1
            end
        end
    end

    # Step 3: 참조 카운트가 0인 변수 정의 제거
    # (단, 함수 호출, 반환, 점프 등은 제거 대상 아님)
    optimized_ir = []

    for instr in ir
        should_keep = true

        # 변수 정의이고 참조되지 않으면 제거
        if haskey(instr, :type) &&
           (instr[:type] == :literal || instr[:type] == :var) &&
           haskey(instr, :id) &&
           ref_count[instr[:id]] == 0
            should_keep = false
        end

        # 부작용이 없는 단항/이항 연산이 참조되지 않으면 제거
        if haskey(instr, :type) &&
           (instr[:type] == :binary_op || instr[:type] == :unary_op) &&
           haskey(instr, :id) &&
           ref_count[instr[:id]] == 0
            should_keep = false
        end

        # 함수 호출, 반환, 점프, I/O는 항상 유지
        if haskey(instr, :type) &&
           (instr[:type] == :call || instr[:type] == :return ||
            instr[:type] == :jump || instr[:type] == :cond_jump ||
            instr[:type] == :print || instr[:type] == :input)
            should_keep = true
        end

        if should_keep
            push!(optimized_ir, instr)
        end
    end

    return optimized_ir
end

# ======================== Function Inlining ========================

"""
    inline_small_functions(ir::Vector)

3줄 이하의 작은 함수를 호출처에 직접 삽입합니다.
"""
function inline_small_functions(ir::Vector)
    if length(ir) < 2
        return ir
    end

    # Step 1: 함수 정의 찾기 및 크기 측정
    func_defs = Dict{String, Vector{Dict}}()
    func_sizes = Dict{String, Int}()

    i = 1
    while i <= length(ir)
        instr = ir[i]
        if haskey(instr, :type) && instr[:type] == :func_def
            func_name = instr[:name]
            body = []

            # 함수 본체 수집 (다음 return까지)
            i += 1
            while i <= length(ir)
                body_instr = ir[i]
                if haskey(body_instr, :type) && body_instr[:type] == :return
                    push!(body, body_instr)
                    break
                end
                push!(body, body_instr)
                i += 1
            end

            func_defs[func_name] = body
            func_sizes[func_name] = length(body)
        end
        i += 1
    end

    # Step 2: 호출 가능한 함수 필터링 (3줄 이하)
    inlineable = Set{String}()
    for (fname, size) in func_sizes
        if size <= 3
            push!(inlineable, fname)
        end
    end

    if isempty(inlineable)
        return ir
    end

    # Step 3: 호출을 인라인화로 치환
    # (현재 구현: 단순화를 위해 함수 호출만 표시하고 실제 인라인은 제한적)
    # 실제로는 변수 이름 충돌 해결(renaming)이 필요하지만,
    # 여기서는 마킹만 수행

    optimized_ir = []
    for instr in ir
        if haskey(instr, :type) && instr[:type] == :call
            func_name = instr[:func]
            if func_name in inlineable
                # 마킹: inline 최적화 적용 가능
                instr = copy(instr)
                instr[:inlining_candidate] = true
            end
        end
        push!(optimized_ir, instr)
    end

    return optimized_ir
end

# ======================== Loop Unrolling ========================

"""
    loop_unrolling(ir::Vector)

고정 반복 횟수 루프를 펼칩니다.
범위가 상수이고 4회 이하일 때만 언롤링합니다.
"""
function loop_unrolling(ir::Vector)
    if length(ir) < 2
        return ir
    end

    optimized_ir = []
    i = 1

    while i <= length(ir)
        instr = ir[i]

        # for_range 루프 찾기
        if haskey(instr, :type) && instr[:type] == :for_range
            start_id = instr[:start_id]
            end_id = instr[:end_id]
            loop_body = []
            loop_var = instr[:var]

            # 루프 본체 수집
            i += 1
            depth = 1
            while i <= length(ir) && depth > 0
                loop_body_instr = ir[i]

                if haskey(loop_body_instr, :type) &&
                   loop_body_instr[:type] == :for_range
                    depth += 1
                elseif haskey(loop_body_instr, :type) &&
                       loop_body_instr[:type] == :end_for
                    depth -= 1
                    if depth == 0
                        break
                    end
                end

                push!(loop_body, loop_body_instr)
                i += 1
            end

            # 상수 범위 확인
            start_val = nothing
            end_val = nothing

            # start_id의 값 찾기
            for prev_instr in optimized_ir
                if haskey(prev_instr, :id) && prev_instr[:id] == start_id &&
                   haskey(prev_instr, :type) && prev_instr[:type] == :literal
                    start_val = prev_instr[:value]
                end
            end

            # end_id의 값 찾기
            for prev_instr in optimized_ir
                if haskey(prev_instr, :id) && prev_instr[:id] == end_id &&
                   haskey(prev_instr, :type) && prev_instr[:type] == :literal
                    end_val = prev_instr[:value]
                end
            end

            # 루프 언롤링 가능 확인 (상수 + 4회 이하)
            if start_val !== nothing && end_val !== nothing &&
               isa(start_val, Int) && isa(end_val, Int)

                iterations = end_val - start_val + 1
                if iterations > 0 && iterations <= 4
                    # 루프 언롤링 수행
                    push!(optimized_ir, Dict(
                        :type => :comment,
                        :text => "# Loop unrolled from $(start_val) to $(end_val)"
                    ))

                    for iter in start_val:end_val
                        # 각 반복에 대해 본체 복사
                        for body_instr in loop_body
                            body_copy = copy(body_instr)
                            push!(optimized_ir, body_copy)
                        end
                    end

                    push!(optimized_ir, Dict(
                        :type => :comment,
                        :text => "# Loop unrolling complete"
                    ))

                    i += 1
                    continue
                end
            end

            # 언롤링 불가능: 원본 루프 유지
            push!(optimized_ir, instr)
            for body_instr in loop_body
                push!(optimized_ir, body_instr)
            end
            push!(optimized_ir, Dict(:type => :end_for))
            i += 1
            continue
        end

        push!(optimized_ir, instr)
        i += 1
    end

    return optimized_ir
end

# ======================== 최적화 통계 ========================

"""
    optimization_stats(ir_before::Vector, ir_after::Vector)

최적화 결과의 통계를 반환합니다.
"""
function optimization_stats(ir_before::Vector, ir_after::Vector)
    instructions_removed = length(ir_before) - length(ir_after)
    reduction_percent = if length(ir_before) > 0
        round(100 * instructions_removed / length(ir_before); digits=2)
    else
        0
    end

    return Dict(
        :original_size => length(ir_before),
        :optimized_size => length(ir_after),
        :instructions_removed => instructions_removed,
        :reduction_percent => reduction_percent
    )
end

end  # module
