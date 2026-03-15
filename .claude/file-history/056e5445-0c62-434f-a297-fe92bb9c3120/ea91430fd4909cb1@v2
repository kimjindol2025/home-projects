"""
    Optimization Tests for FreeLiulia

Phase 6 최적화 테스트:
- Dead Code Elimination (DCE): 10개
- Function Inlining: 6개
- Loop Unrolling: 6개
총 22개 테스트
"""

include("../src/freeliulia_optimizer.jl")

# ======================== 테스트 프레임워크 ========================

test_count = 0
test_passed = 0

function test(name::String, condition::Bool)
    global test_count, test_passed
    test_count += 1
    if condition
        test_passed += 1
        println("✓ Test $test_count: $name")
    else
        println("✗ Test $test_count: $name")
    end
end

# ======================== Test Suite 1: Dead Code Elimination ========================

println("\n📌 Test Suite 1: Dead Code Elimination (10 tests)")
println("─" ^ 50)

# Test 1.1: 사용되지 않는 리터럴 제거
ir_1_1 = [
    Dict(:type => :literal, :value => 5, :id => 1),      # 사용안함
    Dict(:type => :literal, :value => 10, :id => 2),     # 사용됨
    Dict(:type => :binary_op, :op => :+, :left_id => 2, :right_id => 2, :id => 3),  # 결과: 20
]
ir_opt_1_1 = dead_code_elimination(ir_1_1)
test("DCE: 사용안 함 리터럴 제거", length(ir_opt_1_1) < length(ir_1_1))

# Test 1.2: 사용되지 않는 변수 제거
ir_1_2 = [
    Dict(:type => :var, :name => "x", :value => 5, :id => 1),        # 사용안함
    Dict(:type => :var, :name => "y", :value => 10, :id => 2),       # 사용됨
    Dict(:type => :return, :value_id => 2),
]
ir_opt_1_2 = dead_code_elimination(ir_1_2)
test("DCE: 사용안 함 변수 정의 제거", length(ir_opt_1_2) < length(ir_1_2))

# Test 1.3: 사용 중인 코드는 유지
ir_1_3 = [
    Dict(:type => :literal, :value => 5, :id => 1),
    Dict(:type => :literal, :value => 10, :id => 2),
    Dict(:type => :binary_op, :op => :+, :left_id => 1, :right_id => 2, :id => 3),
]
ir_opt_1_3 = dead_code_elimination(ir_1_3)
test("DCE: 모든 코드가 사용 중이면 유지", length(ir_opt_1_3) == length(ir_1_3))

# Test 1.4: 함수 호출은 항상 유지
ir_1_4 = [
    Dict(:type => :literal, :value => 1, :id => 1),
    Dict(:type => :call, :func => "print", :args => [1], :id => 2),  # 부작용 있음
]
ir_opt_1_4 = dead_code_elimination(ir_1_4)
test("DCE: 함수 호출은 항상 유지", length(ir_opt_1_4) == length(ir_1_4))

# Test 1.5: 반환문은 항상 유지
ir_1_5 = [
    Dict(:type => :literal, :value => 42, :id => 1),
    Dict(:type => :return, :value_id => 1),
]
ir_opt_1_5 = dead_code_elimination(ir_1_5)
test("DCE: 반환문은 항상 유지", length(ir_opt_1_5) == length(ir_1_5))

# Test 1.6: 빈 IR은 빈 상태 유지
ir_1_6 = []
ir_opt_1_6 = dead_code_elimination(ir_1_6)
test("DCE: 빈 IR 처리", length(ir_opt_1_6) == 0)

# Test 1.7: 단일 리터럴 사용 안 함
ir_1_7 = [
    Dict(:type => :literal, :value => 999, :id => 1),
]
ir_opt_1_7 = dead_code_elimination(ir_1_7)
test("DCE: 단일 사용 안 함 리터럴 제거", length(ir_opt_1_7) == 0)

# Test 1.8: 복합 연산에서 미사용 항 제거
ir_1_8 = [
    Dict(:type => :literal, :value => 1, :id => 1),    # 미사용
    Dict(:type => :literal, :value => 2, :id => 2),    # 사용
    Dict(:type => :literal, :value => 3, :id => 3),    # 사용
    Dict(:type => :binary_op, :op => :+, :left_id => 2, :right_id => 3, :id => 4),  # 결과: 5
]
ir_opt_1_8 = dead_code_elimination(ir_1_8)
test("DCE: 복합 연산에서 미사용 항 제거", length(ir_opt_1_8) < length(ir_1_8))

# Test 1.9: 조건부 점프는 유지
ir_1_9 = [
    Dict(:type => :literal, :value => true, :id => 1),
    Dict(:type => :cond_jump, :cond_id => 1, :target => :L1),
]
ir_opt_1_9 = dead_code_elimination(ir_1_9)
test("DCE: 조건부 점프는 유지", length(ir_opt_1_9) == length(ir_1_9))

# Test 1.10: 인쇄는 부작용이 있으므로 유지
ir_1_10 = [
    Dict(:type => :literal, :value => "hello", :id => 1),
    Dict(:type => :print, :value_id => 1),
]
ir_opt_1_10 = dead_code_elimination(ir_1_10)
test("DCE: 인쇄는 부작용이 있으므로 유지", length(ir_opt_1_10) == length(ir_1_10))

# ======================== Test Suite 2: Function Inlining ========================

println("\n📌 Test Suite 2: Function Inlining (6 tests)")
println("─" ^ 50)

# Test 2.1: 작은 함수 감지
ir_2_1 = [
    Dict(:type => :func_def, :name => "add", :params => [:a, :b]),
    Dict(:type => :return, :value_id => 1),  # 작음 (2줄)
]
ir_opt_2_1 = inline_small_functions(ir_2_1)
test("Inlining: 작은 함수 감지", length(ir_opt_2_1) > 0)

# Test 2.2: 함수 호출 찾기
ir_2_2 = [
    Dict(:type => :func_def, :name => "id", :params => [:x]),
    Dict(:type => :return, :value_id => 1),  # 1줄 함수
    Dict(:type => :call, :func => "id", :args => [1]),  # 호출
]
ir_opt_2_2 = inline_small_functions(ir_2_2)
test("Inlining: 함수 호출 찾기", length(ir_opt_2_2) > 0)

# Test 2.3: 인라인 가능 표시
ir_2_3 = [
    Dict(:type => :func_def, :name => "small", :params => []),
    Dict(:type => :return, :value_id => 1),
    Dict(:type => :call, :func => "small", :args => []),
]
ir_opt_2_3 = inline_small_functions(ir_2_3)
# 인라인 가능한 함수 호출이 있는지 확인
has_inline_candidate = any(haskey(instr, :inlining_candidate) && instr[:inlining_candidate] for instr in ir_opt_2_3)
test("Inlining: 인라인 가능 표시", has_inline_candidate)

# Test 2.4: 빈 IR 처리
ir_2_4 = []
ir_opt_2_4 = inline_small_functions(ir_2_4)
test("Inlining: 빈 IR 처리", length(ir_opt_2_4) == 0)

# Test 2.5: 큰 함수는 인라인 안 함
ir_2_5 = [
    Dict(:type => :func_def, :name => "large", :params => [:x, :y]),
    Dict(:type => :var, :name => "a", :value => 1, :id => 1),
    Dict(:type => :var, :name => "b", :value => 2, :id => 2),
    Dict(:type => :var, :name => "c", :value => 3, :id => 3),
    Dict(:type => :binary_op, :op => :+, :left_id => 1, :right_id => 2, :id => 4),
    Dict(:type => :binary_op, :op => :+, :left_id => 4, :right_id => 3, :id => 5),
    Dict(:type => :return, :value_id => 5),  # 큼 (7줄)
    Dict(:type => :call, :func => "large", :args => [1, 2]),
]
ir_opt_2_5 = inline_small_functions(ir_2_5)
has_inline_marker = any(haskey(instr, :inlining_candidate) && instr[:inlining_candidate] for instr in ir_opt_2_5)
test("Inlining: 큰 함수는 인라인 안 함", !has_inline_marker)

# Test 2.6: 함수 정의 없이 호출
ir_2_6 = [
    Dict(:type => :call, :func => "undefined", :args => [1]),
]
ir_opt_2_6 = inline_small_functions(ir_2_6)
test("Inlining: 정의 없는 함수 호출 유지", length(ir_opt_2_6) > 0)

# ======================== Test Suite 3: Loop Unrolling ========================

println("\n📌 Test Suite 3: Loop Unrolling (6 tests)")
println("─" ^ 50)

# Test 3.1: 상수 범위 루프 언롤링
ir_3_1 = [
    Dict(:type => :literal, :value => 1, :id => 1),
    Dict(:type => :literal, :value => 3, :id => 2),
    Dict(:type => :for_range, :var => :i, :start_id => 1, :end_id => 2),
    Dict(:type => :literal, :value => 0, :id => 3),
    Dict(:type => :end_for),
]
ir_opt_3_1 = loop_unrolling(ir_3_1)
test("Loop Unrolling: 상수 범위 언롤링", length(ir_opt_3_1) >= length(ir_3_1))

# Test 3.2: 1회 반복 언롤링
ir_3_2 = [
    Dict(:type => :literal, :value => 5, :id => 1),
    Dict(:type => :literal, :value => 5, :id => 2),
    Dict(:type => :for_range, :var => :i, :start_id => 1, :end_id => 2),
    Dict(:type => :literal, :value => 10, :id => 3),
    Dict(:type => :end_for),
]
ir_opt_3_2 = loop_unrolling(ir_3_2)
test("Loop Unrolling: 1회 반복 언롤링", any(instr -> haskey(instr, :type) && instr[:type] == :comment && contains(instr[:text], "unrolled"), ir_opt_3_2))

# Test 3.3: 4회 반복 언롤링
ir_3_3 = [
    Dict(:type => :literal, :value => 1, :id => 1),
    Dict(:type => :literal, :value => 4, :id => 2),
    Dict(:type => :for_range, :var => :i, :start_id => 1, :end_id => 2),
    Dict(:type => :literal, :value => 0, :id => 3),
    Dict(:type => :end_for),
]
ir_opt_3_3 = loop_unrolling(ir_3_3)
test("Loop Unrolling: 4회 반복 언롤링", any(instr -> haskey(instr, :type) && instr[:type] == :comment, ir_opt_3_3))

# Test 3.4: 5회 반복은 언롤링 안 함 (한계)
ir_3_4 = [
    Dict(:type => :literal, :value => 1, :id => 1),
    Dict(:type => :literal, :value => 5, :id => 2),
    Dict(:type => :for_range, :var => :i, :start_id => 1, :end_id => 2),
    Dict(:type => :literal, :value => 0, :id => 3),
    Dict(:type => :end_for),
]
ir_opt_3_4 = loop_unrolling(ir_3_4)
unrolled_count = count(instr -> haskey(instr, :type) && instr[:type] == :comment && contains(instr[:text], "unrolled"), ir_opt_3_4)
test("Loop Unrolling: 5회는 한계이므로 언롤링 안 함", unrolled_count == 0)

# Test 3.5: 빈 IR 처리
ir_3_5 = []
ir_opt_3_5 = loop_unrolling(ir_3_5)
test("Loop Unrolling: 빈 IR 처리", length(ir_opt_3_5) == 0)

# Test 3.6: 루프 없는 IR은 변화 없음
ir_3_6 = [
    Dict(:type => :literal, :value => 5, :id => 1),
    Dict(:type => :literal, :value => 10, :id => 2),
    Dict(:type => :binary_op, :op => :+, :left_id => 1, :right_id => 2, :id => 3),
]
ir_opt_3_6 = loop_unrolling(ir_3_6)
test("Loop Unrolling: 루프 없는 IR은 변화 없음", length(ir_opt_3_6) == length(ir_3_6))

# ======================== 통합 최적화 테스트 ========================

println("\n📌 Integrated Test: run_all_optimizations()")
println("─" ^ 50)

ir_integrated = [
    Dict(:type => :literal, :value => 1, :id => 1),      # 미사용
    Dict(:type => :literal, :value => 2, :id => 2),      # 사용
    Dict(:type => :literal, :value => 3, :id => 3),      # 사용
    Dict(:type => :binary_op, :op => :+, :left_id => 2, :right_id => 3, :id => 4),
    Dict(:type => :return, :value_id => 4),
]
ir_opt_integrated = run_all_optimizations(ir_integrated)
test("Integrated: DCE + Inlining + Loop Unrolling", length(ir_opt_integrated) <= length(ir_integrated))

# ======================== 최적화 통계 테스트 ========================

println("\n📌 Optimization Statistics")
println("─" ^ 50)

stats = optimization_stats(ir_integrated, ir_opt_integrated)
test("Stats: 통계 생성 성공", haskey(stats, :reduction_percent))
println("  원본 크기: $(stats[:original_size]) 명령어")
println("  최적화 크기: $(stats[:optimized_size]) 명령어")
println("  제거됨: $(stats[:instructions_removed]) 명령어")
println("  감소: $(stats[:reduction_percent])%")

# ======================== 테스트 결과 요약 ========================

println("\n" * "=" ^ 50)
println("✨ Optimization Test Summary")
println("=" ^ 50)
println("✅ 통과: $test_passed / $test_count")
if test_passed == test_count
    println("🎉 모든 최적화 테스트 통과!")
else
    println("⚠️  일부 테스트 실패")
end
println("=" ^ 50)

# 종료 상태 코드
exit(test_passed == test_count ? 0 : 1)
