---
name: C Compiler Learning - Phase 2 (Codegen & x86 Encoding)
description: AST→x86-64 code generation, instruction encoding, fixups, symbol/relocation tracking for ELF
type: project
---

# 🎓 C Compiler Learning - Phase 2: Codegen & x86-64 Encoding

**Status**: Phase 2/5 Complete (2026-03-16)
**Key Insight**: Direct AST→x86-64 (no separate IR stage), using encoder abstraction + fixup mechanism

---

## 1️⃣ **CODE GENERATOR (codegen.c/h - 2,293 lines)**

### Architecture: Direct AST → x86-64

```
AST (from parser/sema)
       ↓
codegen_gen(cg, program)
       ├→ Walk each function def
       ├→ Emit prologue (rbp setup, stack alloc)
       ├→ codegen_stmt(node) for each statement
       ├→ codegen_expr(node) for expressions
       ├→ Emit epilogue (rbp restore, return)
       └→ Track relocations & symbols
       ↓
enc->buf (contains .text machine code)
data_buf / rodata_buf / bss_size
syms / relocs (for ELF writer)
```

**Key difference**: No explicit IR intermediate stage. Instruction selection happens directly during tree-walk codegen.

### CodeGen State Management

```c
typedef struct CodeGen {
    Encoder    *enc;              // x86 instruction encoder

    // Sections
    Buf         data_buf;         // .data initialized global variables
    Buf         rodata_buf;       // .rodata string literals, const arrays
    size_t      bss_size;         // .bss uninitialized globals

    // Tracking
    SymEntry   *syms;             // function/global symbol table
    Reloc      *relocs;           // GOT/PLT relocations
    DataEntry  *data_entries;     // data section layout

    // Function-local state
    char       *cur_func_name;
    int         cur_stack_size;   // bytes allocated for local vars
    char       *cur_return_label;

    // Control flow stacks
    LoopLabel  *loop_stack;       // for break/continue
    SwitchCtx  *switch_stack;     // for switch cases

    int         label_count;      // for unique label generation
} CodeGen;
```

### Stack Frame Layout (System V AMD64 ABI)

```
   High Address
   ┌─────────────────┐
   │ argc, argv      │  (passed to main)
   ├─────────────────┤
   │ return address  │  (pushed by call)
   ├─────────────────┤ ← rbp (after prologue)
   │ old rbp         │
   ├─────────────────┤
   │ local var 1     │  [rbp - 8]
   │ local var 2     │  [rbp - 16]
   │ ...             │
   ├─────────────────┤
   │ [spill space]   │  (for callee-saved regs)
   └─────────────────┘ ← rsp
   Low Address
```

**Key**: Stack grows downward. Locals allocated with negative offsets from rbp.

### Data Section Management

```c
static size_t data_add(CodeGen *cg, const char *label,
                       const uint8_t *bytes, int size, int align,
                       bool is_bss, bool is_rodata) {
    // 1. Calculate alignment padding
    if (align > 1) {
        size_t pad = (align - (off % align)) % align;
        // Add padding zeros
    }

    // 2. Write data bytes to appropriate section
    for (int i = 0; i < size; i++) {
        buf_write8(is_rodata ? &cg->rodata_buf : &cg->data_buf, bytes[i]);
    }

    // 3. Record DataEntry for ELF writer
    DataEntry *de = calloc(1, sizeof(DataEntry));
    de->label = label;        // symbol name
    de->size = size;
    de->align = align;
    de->is_bss = is_bss;      // if true: .bss not .data

    return off;  // for relocations
}
```

**Insight**: Alignment handled per-entry (not per-section), allowing variable alignment.

### Symbol Registration Pattern

```c
static void sym_add(CodeGen *cg, const char *name, CGSymKind kind,
                    size_t offset, int size, bool is_global) {
    SymEntry *s = calloc(1, sizeof(SymEntry));
    s->name = strdup(name);
    s->kind = kind;           // CGSYM_FUNC, CGSYM_OBJECT, CGSYM_BSS, etc.
    s->offset = offset;       // in .text or .data
    s->size = size;
    s->is_global = is_global;
    s->next = cg->syms;
    cg->syms = s;             // linked list LIFO
}
```

**Order**: Symbols added in reverse order (prepend). ELF writer reverses to get forward order.

### Relocation Registration

```c
static void reloc_add(CodeGen *cg, size_t offset, const char *sym,
                      int type, int64_t addend) {
    Reloc *r = calloc(1, sizeof(Reloc));
    r->offset = offset;        // in .text where patch needed
    r->sym_name = strdup(sym); // external symbol name
    r->type = type;            // R_X86_64_PC32, R_X86_64_GLOB_DAT, etc.
    r->addend = addend;        // constant addend
    r->next = cg->relocs;
    cg->relocs = r;
}
```

**Example**:
- Call instruction at .text offset 0x1000 → `call main` becomes `call 0x00000000` with relocation type=PC32 for "main"
- ELF writer applies relocation: `addr = symbol_address + addend - reloc_offset`

### Control Flow Context Stacks

```c
// Break/Continue targets
typedef struct LoopLabel {
    char *break_label;      // goto label on break;
    char *cont_label;       // goto label on continue;
    struct LoopLabel *next;
} LoopLabel;

// Switch state
typedef struct SwitchCtx {
    char *end_label;
    char *default_label;
    char **case_labels;     // pre-allocated labels
    int ncases;
    int case_idx;           // incremented as cases emitted
} SwitchCtx;
```

**Key insight**: Labels pre-allocated before switch body (two-pass dispatch). Allows fallthrough within single pass.

---

## 2️⃣ **x86-64 INSTRUCTION ENCODER (x86_encode.c/h - 34K)**

### Encoder Architecture

```c
typedef struct Encoder {
    Buf      buf;              // output machine-code buffer
    Fixup   *fixups;           // unresolved relative jumps
    LabelDef *labels;          // label definitions (name → offset)
    int      label_counter;    // for auto-generated labels
} Encoder;

// Fixup: a relative 32-bit displacement needing patch
typedef struct Fixup {
    size_t  patch_off;         // offset where disp goes
    char   *label;             // target label name
    size_t  inst_end;          // for PC-relative calc
} Fixup;
```

### Instruction Emission Pattern

```c
// Naming convention: enc_<mnemonic>_<operand-form>
// Forms: rr (reg,reg), ri (reg,imm), rm (reg,[mem]), mr ([mem],reg), etc.

void enc_add_rr(Encoder *e, Reg dst, Reg src, OpSize sz) {
    // Emit x86-64 ADD instruction
    // 1. Emit REX prefix if needed (for 64-bit, r8-r15 regs)
    // 2. Emit opcode byte(s)
    // 3. Emit ModRM byte (dst reg, src reg)
}

void enc_lea_rm(Encoder *e, Reg dst, Reg base, int disp, OpSize sz) {
    // LEA dst, [base + disp]
    // Useful for: address calculation, pointer arithmetic
    // NO memory access — just computes address
}

void enc_jmp_label(Encoder *e, const char *label) {
    // JMP rel32 label
    // Creates fixup if label not yet defined
    // Will be patched by enc_resolve_fixups()
}
```

### Relative Jump Resolution

```c
void enc_resolve_fixups(Encoder *e) {
    // Two-pass:
    // Pass 1: Build LabelDef table (name → offset)
    // Pass 2: Patch each fixup

    for (Fixup *f = e->fixups; f; f = f->next) {
        LabelDef *ld = find_label(e, f->label);

        // Calculate PC-relative displacement
        int32_t disp = (int32_t)(ld->off - f->inst_end);

        // Patch the 4-byte immediate
        buf_writei32_at(&e->buf, f->patch_off, disp);
    }
}
```

**Key insight**: Emit as 00000000, then patch after all labels known. PC-relative = target_offset - instruction_end.

### Register Usage Conventions

```c
// Function argument passing (System V AMD64 ABI)
static const Reg ARG_REGS[] = {
    REG_RDI,  // arg 0
    REG_RSI,  // arg 1
    REG_RDX,  // arg 2
    REG_RCX,  // arg 3
    REG_R8,   // arg 4
    REG_R9,   // arg 5
    // arg 6+ go on stack (right-to-left)
};

// Return value
// RAX: integer (or low 64 bits of 128-bit)
// RDX: high 64 bits (for 128-bit values)

// Caller-saved (can be clobbered):
// rax, rcx, rdx, rsi, rdi, r8, r9, r10, r11

// Callee-saved (must preserve):
// rbx, rbp, r12, r13, r14, r15
```

### Opcode Encoding (ModRM Byte)

```
ModRM Byte:  [MOD(2)][REG(3)][RM(3)]

Examples:
  ADD rax, rbx    → opcode=01, ModRM: 11 000 011 (rax=000, rbx=011, MOD=11=reg)
  ADD [rax], rbx  → opcode=01, ModRM: 00 011 000 (rax=000, rbx=011, MOD=00=mem)

REX prefix (for 64-bit, high registers):
  Format: 0x40 | W(1) | R(1) | X(1) | B(1)
  W=1: 64-bit operand
  R: high bit of REG field (r8-r15)
  X: high bit of index register
  B: high bit of RM field (r8-r15)
```

---

## 3️⃣ **ELF BINARY WRITING (elf_writer.c - 16K)**

### ELF Structure

```
ELF Header (64 bytes)
  ├─ Magic: 0x7F 'E' 'L' 'F'
  ├─ Bitness: 64-bit
  ├─ Endianness: little-endian
  ├─ Entry point: typically _start
  └─ e_shoff: offset to section header table

Program Header Table (Phdrs)
  ├─ PT_LOAD: loadable segment (.text, .data, .bss)
  └─ PT_DYNAMIC: dynamic linking info

Section Header Table (Shdrs)
  ├─ .text (executable code)
  ├─ .rodata (read-only data)
  ├─ .data (initialized globals)
  ├─ .bss (uninitialized globals)
  ├─ .symtab (symbol table)
  ├─ .strtab (string table)
  ├─ .rel.text (relocations for .text)
  └─ ...

Symbol Table (.symtab)
  ├─ STT_FUNC: function symbols
  ├─ STT_OBJECT: global variable symbols
  └─ STT_NOTYPE: unresolved external symbols

Relocation Table (.rel.text)
  ├─ R_X86_64_PC32: PC-relative 32-bit (for calls, direct jumps)
  ├─ R_X86_64_GLOB_DAT: global data address
  └─ R_X86_64_PLT32: PLT indirect call
```

### C Compiler's ELF Pipeline

```c
void codegen_emit_elf(CodeGen *cg, const char *outfile) {
    ElfWriter *ew = elf_writer_new();

    // 1. Add symbol entries
    for (SymEntry *s = cg->syms; s; s = s->next) {
        elf_add_symbol(ew, s->name, s->kind, s->offset, s->size, s->is_global);
    }

    // 2. Add relocation entries
    for (Reloc *r = cg->relocs; r; r = r->next) {
        elf_add_relocation(ew, r->offset, r->sym_name, r->type, r->addend);
    }

    // 3. Write sections
    elf_write_section(ew, ".text", cg->enc->buf.data, cg->enc->buf.size);
    elf_write_section(ew, ".rodata", cg->rodata_buf.data, cg->rodata_buf.size);
    elf_write_section(ew, ".data", cg->data_buf.data, cg->data_buf.size);
    elf_write_section_bss(ew, ".bss", cg->bss_size);

    // 4. Write ELF file
    elf_write_file(ew, outfile);
}
```

---

## 📊 **Design Patterns Applied**

### ✅ Pattern 1: State Stacks for Nested Contexts
```c
// Push when entering loop/switch
LoopLabel *new_loop = malloc(sizeof(LoopLabel));
new_loop->next = cg->loop_stack;
cg->loop_stack = new_loop;

// Pop when exiting
cg->loop_stack = cg->loop_stack->next;
free(new_loop);
```

### ✅ Pattern 2: Two-Pass Compilation
```c
// Pass 1: Emit all instructions with placeholder targets
enc_jmp_label(e, ".L_unknown");  // Creates fixup

// Pass 2: Resolve all labels
enc_label(e, ".L_unknown");

// Pass 3: Patch all fixups
enc_resolve_fixups(e);
```

### ✅ Pattern 3: Union-based Data Sections
```c
// Instead of 3 separate buffers, use tagged union
enum { SECT_TEXT, SECT_DATA, SECT_RODATA, SECT_BSS };
struct {
    Buf bufs[4];      // indexed by section type
    // or: smartly switch on flag
}
```

### ✅ Pattern 4: Symbol Reversal
```c
// Symbols added LIFO (prepend)
// ELF writer reverses before writing
// Result: correct order in output without extra work
```

---

## 🔗 **Integration with freelang-to-c**

### Key Takeaways for C Code Generation

1. **Direct tree-walk codegen** works well for simple languages
   - No need for explicit IR if instruction selection is straightforward
   - Best for: statically typed, simple control flow

2. **Encoder abstraction** separates:
   - High-level instruction selection (codegen.c)
   - Low-level machine encoding (x86_encode.c)
   - ELF binary writing (elf_writer.c)
   - **For freelang-to-c**: Can emit C code at codegen stage, then call C compiler

3. **Fixup mechanism** enables:
   - Single-pass code generation with forward references
   - Lazily-resolved labels
   - **For freelang-to-c**: Labels in generated C code become fixups in C compiler

4. **Section separation** (.text, .data, .rodata, .bss):
   - Static data separate from code
   - Read-only data protected
   - Uninitialized data in BSS (no disk space)
   - **For freelang-to-c**: Map FreeLang data sections → C static/const/global

---

## 🎯 **Next Steps (Phase 3-4)**

### Phase 3: Semantic Analysis (sema.c - 50K)
- Type checking and inference
- Symbol resolution
- Array/struct field calculation
- Function signature validation

### Phase 4: IR Optimization (ir.c - 38K)
- Constant folding
- Dead code elimination
- Copy propagation
- Common subexpression elimination

---

**Status**: ✅ Phase 2 Complete | Phase 3-4 Pending

