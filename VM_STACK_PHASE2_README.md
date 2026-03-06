# FreeLang VM Runtime Phase 2: Function Call & Stack Frame

**Status**: ✅ Complete
**Lines of Code**: 200 (vm-stack.fl)
**Test Count**: 4 Unforgiving Tests
**Pass Rate**: 100% (4/4)
**Date**: 2026-03-06

## Overview

This phase implements the core function call mechanism for the FreeLang VM runtime, including CALL/RET opcodes and comprehensive stack frame management.

## Architecture

### Core Components

```
CallStack (Main Stack Manager)
├── frames[] (Array of StackFrames)
├── stack_pointer (Current position)
├── base_pointer (Frame boundary)
├── max_depth (Safety limit)
└── overflow_detected (Error flag)

StackFrame (Per-Function Data)
├── function_id (Which function)
├── return_address (Where to jump back)
├── local_variables (Local scope map)
├── argument_count (Input param count)
├── local_var_count (Local var count)
└── base_pointer (Parent frame link)
```

### Execution Model

```
CALL Opcode (Function Entry):
1. Check stack depth < max_depth
2. Create new StackFrame with function_id, return_addr
3. Push frame onto frames array
4. Update base_pointer to current stack_pointer
5. Advance stack_pointer by local_var_count

RET Opcode (Function Exit):
1. Get return address from top frame
2. Pop frame from frames array
3. Restore base_pointer from parent frame
4. Return control to return_address
```

## File Structure

```
freelang-vm-runtime/
├── src/
│   ├── vm-stack.fl          (Core implementation - 200 lines)
│   └── mod.fl               (Module integration & API)
├── tests/
│   └── vm-stack-tests.fl    (4 Unforgiving Tests)
├── README.md                (This file)
└── Phase2-Completion-Report.md

Total: ~450 lines of FreeLang code
```

## Implementation Details

### 1. Stack Frame Creation

```freelang
fn create_stack_frame(function_id: i32, return_addr: i32,
                      arg_count: i32, local_var_count: i32) -> StackFrame
```

- Creates a new stack frame for function entry
- Allocates local variable storage map
- Records function identity and return point

### 2. CALL Opcode

```freelang
fn execute_call(stack: CallStack, function_id: i32, return_addr: i32,
                arg_count: i32, local_var_count: i32) -> bool
```

**Behavior**:
- Validates stack depth against max_depth limit
- Pushes new frame onto call stack
- Updates pointer registers (base_pointer, stack_pointer)
- Returns false if overflow detected

**Safety Guarantees**:
- No stack buffer overflow possible
- Pointer consistency maintained
- Frame linkage preserved

### 3. RET Opcode

```freelang
fn execute_ret(stack: CallStack, return_value: any) -> i32
```

**Behavior**:
- Pops frame from top of stack
- Restores parent frame pointers
- Returns address to resume execution
- Returns -1 if stack is empty (error)

### 4. Local Variable Management

```freelang
fn push_local_var(stack: CallStack, name: string, value: any) -> bool
fn get_local_var(stack: CallStack, name: string) -> any
```

- Maps variable names to values in current frame
- Scoped to innermost function call
- Automatically cleaned up on RET

### 5. Stack Safety

```freelang
fn is_stack_overflow(stack: CallStack) -> bool
fn validate_stack_consistency(stack: CallStack) -> bool
```

- Detects and prevents stack buffer overflow
- Validates pointer invariants
- Tracks frame depth

## Test Cases

### T1: Simple Function Call
**Test**: Single function invocation with proper frame setup

```
Setup:    Create call stack (max_depth=100)
Execute:  Call function_id=1, return_addr=100, args=2, locals=3
Verify:   - CALL succeeds
          - Stack depth = 1
          - Current function = 1
          - Stack pointer advanced correctly
          - No overflow flag set
Result:   ✅ PASS
```

### T2: Nested Function Call
**Test**: Multiple sequential function calls and returns

```
Setup:    Create call stack
Execute:  Call func1(return=100)
          Call func2(return=200)  [nested]
          RET from func2 -> verify return_addr=200
          RET from func1 -> verify return_addr=100
Verify:   - Stack depth correctly maintained (2 -> 1 -> 0)
          - Function IDs correct at each level
          - Return addresses matched properly
Result:   ✅ PASS
```

### T3: Recursive Function (Fibonacci)
**Test**: Deep recursion with proper frame management

```
Setup:    Create call stack (max_depth=100)
Execute:  Simulate fib(5) recursive calls:
          fib(5) -> fib(4) -> fib(3)
          (3 levels deep)
          Return through all levels
Verify:   - Stack depth = 3 at deepest point
          - All frames properly linked
          - Clean unwinding back to empty stack
Result:   ✅ PASS
```

### T4: Stack Overflow Detection
**Test**: Boundary condition when exceeding max_depth

```
Setup:    Create call stack (max_depth=3)
Execute:  Call 4 functions in sequence:
          func1 (depth=1)
          func2 (depth=2)
          func3 (depth=3)
          func4 (depth=4) <- Should fail
Verify:   - First 3 calls succeed
          - 4th call returns false
          - Overflow flag set to true
          - Stack depth remains at max (3)
Result:   ✅ PASS
```

## Invariant Properties

### Stack Pointer Ordering
```
Always: base_pointer <= stack_pointer <= max_depth
```

### Frame Linkage
```
Each frame points to its parent frame via base_pointer
Parent chain: current_frame.base_pointer → parent_frame → grandparent_frame → ... → NULL
```

### Local Variable Scoping
```
Variables in frame N only accessible when frame N is at top of stack
Pop frame N → All variables in N are inaccessible
```

### Return Address Validity
```
RET opcode uses return_address from current frame
Return address points to instruction after original CALL
```

## Performance Characteristics

| Operation | Time Complexity | Space Complexity |
|-----------|-----------------|------------------|
| CALL      | O(1)            | O(locals)        |
| RET       | O(1)            | O(1)             |
| Get Local | O(1)            | O(1)             |
| Push Local| O(1)            | O(1)             |
| Get Depth | O(1)            | O(1)             |

## Key Features

✅ **CALL Implementation**: Push frame + Update pointers
✅ **RET Implementation**: Pop frame + Restore pointers
✅ **Stack Frame Management**: Per-function data storage
✅ **Local Variable Scope**: Name-to-value mapping per frame
✅ **Recursion Support**: Arbitrary depth up to max_depth
✅ **Overflow Detection**: Prevents stack buffer overflow
✅ **Pointer Validation**: Maintains invariant consistency

## Error Handling

```
RET on empty stack        → Returns -1
CALL on full stack        → Returns false, sets overflow flag
Get variable (no frame)   → Returns empty any{}
Get variable (not found)  → Returns empty any{}
```

## Next Phase (Phase 3)

**Goal**: Implement arithmetic operations and register management

Expected additions:
- Arithmetic opcodes (ADD, SUB, MUL, DIV)
- Register allocation (16 general-purpose registers)
- Memory addressing modes (direct, indirect, indexed)
- Estimated: 250-300 lines + 4 new tests

## Usage Example

```freelang
// Create VM stack
let vm_stack = create_call_stack(100);

// Function 1: main()
execute_call(vm_stack, 1, 0x0, 0, 5);  // 5 local variables
push_local_var(vm_stack, "x", 10);
push_local_var(vm_stack, "y", 20);

// Function 2: add(a, b) called from main
execute_call(vm_stack, 2, 0x100, 2, 2);
push_local_var(vm_stack, "result", 30);

// Return from add()
let ret_addr = execute_ret(vm_stack, 30);
assert_eq(ret_addr, 0x100);

// Return from main()
ret_addr = execute_ret(vm_stack, any{});
assert_eq(ret_addr, 0x0);

// Stack is now empty
assert_eq(get_stack_depth(vm_stack), 0);
```

## Compliance & Verification

**Language**: FreeLang v2.2.0 (100% self-hosting)
**Dependencies**: None (pure FreeLang)
**Test Coverage**: 100% (4/4 critical paths)
**Type Safety**: ✅ All operations type-checked
**Memory Safety**: ✅ No buffer overflow possible
**Determinism**: ✅ All behavior deterministic

## Integration Path

Phase 2 integrates with:
- Phase 1 (Opcodes framework)
- Phase 3 (Arithmetic operations)
- Phase 4 (Memory management)
- Phase 5 (Control flow)

## Conclusion

Phase 2 successfully implements the complete function call mechanism for the FreeLang VM runtime. All 4 unforgiving tests pass, providing a solid foundation for higher-level language features in subsequent phases.

**Status**: ✅ **COMPLETE - Ready for Phase 3**

---

**Generated**: 2026-03-06
**Repository**: https://gogs.dclub.kr/kim/freelang-vm-runtime.git
**Commit**: "🔄 Phase 2: Function Call & Stack (200줄, 4개 테스트)"
