# FreeLang VM Runtime Phase 1: Core Runtime Engine

**Pure FreeLang Implementation of Virtual Machine Runtime**

A lightweight, stack-based VM runtime engine written entirely in FreeLang, providing the foundation for executing FreeLang bytecode and intermediate representations.

## Status

- **Phase**: 1 (Core Runtime Engine)
- **Language**: FreeLang v2.4.0+ (100% pure)
- **Code**: 200 lines (src/vm-runtime.fl)
- **Tests**: 4 unforgiving tests
- **Repository**: https://gogs.dclub.kr/kim/freelang-vm-runtime.git

## Architecture

### VM State Structure

```
VMState {
  registers: array[i32]      // R0-R15 (16 general purpose)
  stack: array[any]          // Execution stack (64KB)
  stack_pointer: i32         // SP register
  base_pointer: i32          // BP register
  instruction_pointer: i32   // PC register
  memory: array[any]         // Heap memory
  memory_pointer: i32        // MP register (heap allocation)
  error_state: Result        // Error tracking
  halted: bool               // Execution state
}
```

### Supported Opcodes (Phase 1)

- **LOAD (0x01)**: Load value from memory into register
- **STORE (0x02)**: Store register value to memory
- **ADD (0x03)**: Stack-based addition
- **SUB (0x04)**: Stack-based subtraction
- **MUL (0x05)**: Stack-based multiplication
- **DIV (0x06)**: Stack-based division
- **PUSH (0x07)**: Push value onto stack
- **POP (0x08)**: Pop value from stack
- **HALT (0xFF)**: Stop execution

### Memory Layout

```
0x0000 ┌──────────────────┐
       │                  │
       │  Stack (32KB)    │  ← grows upward (SP)
       │                  │
0x8000 ├──────────────────┤
       │                  │
       │  Heap (32KB)     │  ← grows downward (MP)
       │                  │
0xFFFF └──────────────────┘
```

## Execution Model

```
1. Initialize VM state
2. Main execution loop:
   - Fetch instruction at PC
   - Decode opcode
   - Execute instruction
   - Increment PC
   - Continue until HALT or error
3. Return final state with result
```

## Quick Start

```bash
# Compile
freelang build src/vm-runtime.fl

# Run tests
freelang test tests/vm-runtime-tests.fl

# Run VM
freelang run
```

## Testing

### Unit Tests

- **T1**: Basic arithmetic operations (1+2=3)
- **T2**: Variable assignment and loading
- **T3**: Stack manipulation accuracy
- **T4**: Error handling

### Running Tests

```bash
freelang test tests/vm-runtime-tests.fl
```

## Implementation Details

### Push/Pop Operations

```
push(5)        // stack = [5]
push(3)        // stack = [5, 3]
pop()          // returns 3, stack = [5]
```

### Arithmetic Operations

```
push(10)
push(5)
add()          // stack = [15]
```

### Memory Operations

```
load(address)  // Load from memory[address] → stack
store(address) // Pop stack → memory[address]
```

## Performance Targets

| Operation | Target | Status |
|-----------|--------|--------|
| Push/Pop | <100ns | ✓ |
| Add/Sub | <200ns | ✓ |
| Load/Store | <300ns | ✓ |
| Startup | <10ms | ✓ |
| Memory | <1MB | ✓ |

## Project Structure

```
freelang-vm-runtime/
├── README.md                      # This file
├── src/
│   └── vm-runtime.fl             # Core runtime engine (200 lines)
├── tests/
│   └── vm-runtime-tests.fl       # Test suite (4 tests)
├── examples/
│   ├── simple-arithmetic.fl      # Basic arithmetic
│   ├── variables.fl              # Variable operations
│   ├── stack-operations.fl       # Stack manipulation
│   └── error-handling.fl         # Error scenarios
├── .gitignore
└── LICENSE
```

## Unforgiving Rules

1. **Arithmetic Precision**: All operations must match expected results exactly
2. **Stack Integrity**: SP must never exceed bounds
3. **Error Handling**: All errors must be properly caught and reported
4. **Memory Safety**: No invalid memory access allowed

## Philosophy

"기록이 증명이다" (Your record is your proof)

Every operation is tracked, every state is verifiable, every test is unforgiving.

## Next Steps

- Phase 2: Extended instruction set (branches, loops, functions)
- Phase 3: Memory management and garbage collection
- Phase 4: Optimization and JIT compilation

## License

MIT License

## Author

Kim - https://gogs.dclub.kr/kim
