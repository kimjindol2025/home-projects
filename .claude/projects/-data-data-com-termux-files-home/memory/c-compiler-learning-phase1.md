---
name: C Compiler Learning - Phase 1 (Lexer & Parser)
description: Bottom-up analysis of c-compiler-from-scratch: lexer tokenization, parser AST construction, bump allocation pattern
type: project
---

# 🎓 C Compiler Learning - Phase 1: Lexer & Parser

**Status**: Phase 1/5 Complete (2026-03-16)
**Approach**: Bottom-up (Lexer → Parser → Codegen → x86_encode → ELF)

---

## 1️⃣ **LEXER (lexer.c - 33K)**

### Architecture Pattern
```
Source Code (string)
       ↓
skip_whitespace_and_comments()
       ↓
scan_one() [dispatcher by char type]
       ├→ scan_ident() → keyword lookup or TK_IDENT
       ├→ scan_number() → decimal/hex/octal/binary, suffixes
       ├→ scan_char() → escape sequences, codepoint
       ├→ scan_string() → escape sequences, heap allocation
       └→ scan_operator() → multi-char operators (<<, ->, etc.)
       ↓
Token (struct with union: ival/fval/cval/sval)
```

### Key Design Decisions

#### 1. **Token Union Design**
```c
typedef struct {
    TokenKind kind;      // discriminator
    union {
        unsigned long long ival;  // TK_INT
        double fval;              // TK_FLOAT
        long long cval;           // TK_CHAR
        char *sval;               // TK_STRING, TK_IDENT
    };
    IntSuffix isuf;       // 1-2 bits
    FloatSuffix fsuf;     // 1-2 bits
} Token;
```
**Insight**: Discriminator `kind` determines which union field is valid. This saves memory vs having separate token types for each literal variant.

#### 2. **String Buffer Pattern**
```c
typedef struct {
    char *data;      // heap allocation
    size_t len;      // current length
    size_t cap;      // capacity
} StrBuf;

static void sbuf_push(StrBuf *b, char c) {
    if (b->len + 1 >= b->cap) {
        b->cap *= 2;                    // exponential growth
        b->data = realloc(b->data, b->cap);
    }
    b->data[b->len++] = c;
}
```
**Insight**: Doubling capacity amortizes reallocation cost to O(1) per character over time.

#### 3. **Line/Column Tracking**
```c
static void advance(Lexer *l) {
    if (l->src[l->pos] == '\n') {
        l->line++;
        l->col = 1;
    } else {
        l->col++;
    }
    l->pos++;
}
```
**Insight**: Tracks position separately from pointer for better error messages.

#### 4. **1-Token Lookahead**
```c
struct Lexer {
    bool peeked;        // is peek valid?
    Token peek_tok;     // cached token
};

Token lexer_peek(Lexer *l);   // non-consuming
Token lexer_next(Lexer *l);   // consuming
```
**Insight**: Simple caching strategy for parser's lookahead needs (needed for disambiguating constructs).

#### 5. **Escape Sequence Decoding**
```c
static long long decode_escape(Lexer *l) {
    // Handles: \n, \t, \xHH, \uHHHH, \UHHHHHHHH, \ddd (octal)
    // Returns: byte value or -1 on error
    // Advances: lexer past escape
}
```
**Insight**: Unified decoder handles all escape forms (octal, hex, unicode) — called from both char and string scanners.

### Token Kinds Coverage
- **Literals**: INT, FLOAT, CHAR, STRING
- **Keywords**: 48 (traditional C + C11 _Keywords)
- **Operators**: 40+ (binary, unary, assignment, compound-assign)
- **Punctuation**: (), {}, [], =>, ., etc.
- **Special**: # (preprocessor), ## (token paste), EOF, ERROR

### Error Handling
```c
static void lex_error(const Lexer *l, int line, int col, const char *msg) {
    fprintf(stderr, "%s:%d:%d: lexer error: %s\n", l->filename, line, col, msg);
    // No recovery — caller may continue or abort
}
```
**Insight**: Best-effort recovery (returns partial token, increments error_count) allows continued parsing for multi-error reporting.

---

## 2️⃣ **PARSER (parser.c - 60K)**

### AST Design Pattern

#### Dual-Layer Node Structure
```c
struct Node {
    NodeKind kind;
    Type *type;           // resolved by sema

    // Layer 1: Flat fields (for sema/codegen)
    Node *lhs, *rhs;      // binary operands
    Node *body, *cond;    // control flow
    Node *then, *els;     // branches

    // Layer 2: Named union (for parser/printer)
    union {
        struct { Node *left; Node *right; } binary;      // ND_ADD, ND_MUL, etc.
        struct { Node *cond; Node *then; Node *else_; } if_;
        struct { Node *init; Node *cond; Node *step; Node *body; } for_;
        struct { const char *name; Type *ret_type; Node **params; Node *body; } func;
        // ... more variants
    };
};
```

**Design Insight**:
- **Parser uses**: `node->binary.left`, `node->binary.right` (named, readable)
- **Sema/Codegen use**: `node->lhs`, `node->rhs` (flat, uniform)
- **Both storage present**: Struct is larger but no shared memory issues

#### Bump Allocation for Nodes
```c
struct Parser {
    Arena *arena;  // bump allocator
};

Node *parser_new_node(Parser *p, NodeKind kind, int line, int col) {
    Node *n = arena_alloc(p->arena, sizeof(Node));
    n->kind = kind;
    n->line = line;
    n->col = col;
    return n;
}
```

**Advantage**: Zero-copy deallocation, cache-friendly allocation pattern.

### Precedence Climbing Parser
```c
// Typical pattern in parser.c:
Node *parse_assignment(Parser *p) {
    Node *expr = parse_logical_or(p);

    if (match(p, TK_EQ)) {
        Node *rhs = parse_assignment(p);  // right-associative
        expr = parser_new_node(p, ND_ASSIGN, ...);
        expr->lhs = expr;
        expr->rhs = rhs;
    }
    return expr;
}

Node *parse_logical_or(Parser *p) {
    Node *left = parse_logical_and(p);

    while (match(p, TK_PIPE_PIPE)) {
        Node *right = parse_logical_and(p);
        Node *expr = parser_new_node(p, ND_LOGIC_OR, ...);
        expr->lhs = left;
        expr->rhs = right;
        left = expr;
    }
    return left;
}
```

**Pattern**:
- Higher precedence → deeper recursion
- Left-associative: loop with `while (match)`
- Right-associative: recurse into same level

### Type System in Parser
```c
enum TypeKind {
    TY_INT, TY_CHAR, TY_FLOAT, TY_DOUBLE,  // scalar
    TY_PTR, TY_ARRAY, TY_FUNC,             // derived
    TY_STRUCT, TY_UNION, TY_ENUM,          // aggregate
};

struct Type {
    TypeKind kind;
    int size, align;
    Type *base;           // TY_PTR → pointee, TY_ARRAY → element
    Type **params;        // TY_FUNC → param types
    int param_count;
    bool is_variadic;
    bool is_const, is_volatile, is_restrict, is_atomic;
};
```

**Key insight**: Types are immutable after construction; sema fills in size/align.

---

## 📊 **Key Patterns to Apply to freelang-to-c**

### ✅ Pattern 1: Token-Based Dispatch
Use discriminator (kind) to determine semantic meaning.
```c
// Instead of 100 token types, use union with kind discriminator
```

### ✅ Pattern 2: Bump Allocation
```c
Arena *codegen_arena = arena_new();
// ... allocate IR nodes, never free individually
arena_free(codegen_arena);  // free all at once
```
**For freelang-to-c**: Allocate IR, AST, symbol table from single arena.

### ✅ Pattern 3: Dual-Layer Structure
```c
// IRNode has both:
// 1. Flat fields for IR operations (lhs, rhs, type, etc.)
// 2. Named union members for readability in codegen (binary.left, etc.)
```

### ✅ Pattern 4: Escape Sequence Decoder
FreeLang may need string handling — reuse lexer's escape decoder pattern.

### ✅ Pattern 5: 1-Token Lookahead
For FreeLang parser, peek buffer helps disambiguate:
- `fn foo()` vs `fn (x)` → need lookahead

---

## 🔍 **Integration Points for C Codegen**

### Where Lexer Output Needed
- **Symbol resolution**: Identifiers → symbol table lookup
- **Type inference**: Number literals → integer/float type

### Where Parser Output Needed
- **Semantic checking**: Type nodes → sema analysis
- **Code generation**: Statement/expression AST → IR

### Where Type System Used
- **Memory layout**: Calculate struct offsets, array sizes
- **Type checking**: Verify assignment compatibility
- **Register allocation**: Size determines register width

---

## 📝 **Next Steps (Phase 2)**

1. Study **codegen.c** (82K) - AST → IR conversion
   - How IR is structured (3-address code)
   - How functions/statements/expressions become IR
   - Optimization passes (const folding, DCE)

2. Study **x86_encode.c** (34K) - IR → machine code
   - Register allocation strategy
   - How IR ops map to x86-64 instructions
   - Calling convention implementation

3. Study **elf_writer.c** (16K) - machine code → ELF binary
   - ELF header/sections/symbols
   - Relocation handling
   - How globals/functions are placed

---

**Status**: ✅ Phase 1 Complete | Phase 2 Starting Next

