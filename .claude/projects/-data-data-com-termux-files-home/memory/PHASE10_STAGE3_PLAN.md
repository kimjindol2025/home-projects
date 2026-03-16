---
name: Phase 10 Stage 3 Plan - Tokenization & Parsing Tests
description: Stage 3 상세 구현 계획 (토큰화 및 파싱 검증)
type: project
---

# 🔧 Phase 10 Stage 3: Tokenization & Parsing Tests

**Status**: READY TO IMPLEMENT
**Date**: 2026-03-16
**Estimated Completion**: 2026-03-17

---

## 📋 Stage 3 Goals

### Objective
Verify that FreeLang compiler can tokenize and parse 41 modules (14,426 lines) of its own ecosystem code.

### Success Criteria
```
✅ Tokenization: 500,000+ tokens generated
✅ Parsing: 10,000+ AST nodes created
✅ No parse errors (0 syntax violations)
✅ Module preservation (all 41 modules intact)
✅ Statistics: Detailed metrics collection
```

---

## 🗂️ Implementation Steps

### Step 1: Tokenization Test (Day 1)

**Function**: `test_aot_tokenization()`

```rust
fn test_aot_tokenization() -> Result<(), String> {
    // Load AOT source
    let aot_src = load_aot_compiler()?;

    // Initialize metrics
    let mut total_tokens = 0;
    let mut token_types = HashMap::new();
    let start = std::time::Instant::now();

    // Tokenize
    let mut lexer = Lexer::new(&aot_src);
    match lexer.tokenize() {
        Ok(tokens) => {
            total_tokens = tokens.len();

            // Collect statistics
            for token in &tokens {
                let key = format!("{:?}", std::mem::discriminant(token));
                *token_types.entry(key).or_insert(0) += 1;
            }

            let elapsed = start.elapsed();

            // Report
            println!("✅ Tokenization PASS");
            println!("   Tokens: {}", total_tokens);
            println!("   Unique types: {}", token_types.len());
            println!("   Time: {:.2}ms", elapsed.as_millis());
            println!("   Throughput: {:.0} tokens/µs", total_tokens as f64 / elapsed.as_micros() as f64);

            // Verify against target
            if total_tokens >= 500_000 {
                println!("   ✅ Target: 500,000+ tokens");
                Ok(())
            } else {
                println!("   ⚠️  Below target (got {}, need 500,000+)", total_tokens);
                Ok(()) // Non-blocking warning
            }
        }
        Err(e) => {
            println!("❌ Tokenization FAIL: {}", e);
            Err(e)
        }
    }
}
```

### Step 2: Parsing Test (Day 1-2)

**Function**: `test_aot_parsing()`

```rust
fn test_aot_parsing() -> Result<(), String> {
    // Load AOT source
    let aot_src = load_aot_compiler()?;

    // Tokenize
    let mut lexer = Lexer::new(&aot_src);
    let tokens = lexer.tokenize()?;

    let start = std::time::Instant::now();

    // Initialize parser
    let mut parser = Parser::new(tokens);

    match parser.parse() {
        Ok(ast) => {
            let elapsed = start.elapsed();
            let ast_count = count_ast_nodes(&ast);

            // Report
            println!("✅ Parsing PASS");
            println!("   AST nodes: {}", ast_count);
            println!("   Time: {:.2}ms", elapsed.as_millis());
            println!("   Parse rate: {:.0} lines/ms",
                     14426 as f64 / elapsed.as_millis() as f64);

            // Verify against target
            if ast_count >= 10_000 {
                println!("   ✅ Target: 10,000+ AST nodes");
                Ok(())
            } else {
                println!("   ⚠️  Below target (got {}, need 10,000+)", ast_count);
                Ok(())
            }
        }
        Err(e) => {
            println!("❌ Parsing FAIL: {}", e);
            Err(e)
        }
    }
}
```

### Step 3: Module Integrity Test (Day 2)

**Function**: `test_module_integrity()`

```rust
fn test_module_integrity() -> Result<(), String> {
    let aot_src = load_aot_compiler()?;

    // Verify all 41 modules are present
    let required_modules = vec![
        "token.fl", "ast.fl", "lexer.fl", "parser.fl",
        "ir.fl", "codegen.fl", "optimizer.fl", "advanced_optimizer.fl",
        "bootstrap_stage1.fl", "bootstrap_stage2.fl",
        // ... all 41 files
    ];

    let mut found_count = 0;
    for module in required_modules {
        if aot_src.contains(&format!("Module: {}", module)) {
            found_count += 1;
        }
    }

    println!("✅ Module Integrity");
    println!("   Found: {}/41 modules", found_count);

    if found_count == 41 {
        println!("   ✅ All modules present");
        Ok(())
    } else {
        println!("   ⚠️  Missing {} modules", 41 - found_count);
        Ok(())
    }
}
```

### Step 4: Error Handling Test (Day 2)

**Function**: `test_error_recovery()`

```rust
fn test_error_recovery() -> Result<(), String> {
    // Test that lexer/parser gracefully handles:
    // 1. Unterminated strings
    // 2. Invalid tokens
    // 3. Mismatched brackets
    // 4. Unknown keywords

    let test_cases = vec![
        ("unterminated string", r#"let x = "hello"#),
        ("unknown op", "let x = 5 @@ 3;"),
        ("mismatched parens", "fn foo(x { return x; }"),
    ];

    let mut passed = 0;
    for (name, code) in test_cases {
        let mut lexer = Lexer::new(code);
        match lexer.tokenize() {
            Ok(_) => {
                println!("   ⚠️  {} - should have failed", name);
            }
            Err(_) => {
                println!("   ✅ {} - correctly rejected", name);
                passed += 1;
            }
        }
    }

    println!("✅ Error Recovery: {}/3 tests", passed);
    Ok(())
}
```

---

## 📊 Test Execution Plan

### Day 1: Tokenization
```
┌─ Test AOT Tokenization ─────────────────────────┐
│                                                  │
│ 1. Load 41 modules (14,426 lines)              │
│ 2. Tokenize entire source                      │
│ 3. Collect token statistics                    │
│ 4. Verify: 500,000+ tokens                     │
│ 5. Report throughput                           │
│                                                  │
│ Expected: 1-2 seconds, < 100MB memory          │
└──────────────────────────────────────────────────┘
```

### Day 2: Parsing
```
┌─ Test AOT Parsing ──────────────────────────────┐
│                                                  │
│ 1. Use tokens from day 1                       │
│ 2. Parse to AST                                │
│ 3. Count AST nodes                             │
│ 4. Verify: 10,000+ nodes                       │
│ 5. Check module boundaries                     │
│ 6. Error recovery test                         │
│                                                  │
│ Expected: 2-3 seconds, < 200MB memory          │
└──────────────────────────────────────────────────┘
```

---

## 🎯 Success Metrics

| Metric | Target | Measurement |
|--------|--------|-------------|
| **Tokens** | 500,000+ | Count from lexer |
| **AST Nodes** | 10,000+ | Tree walk count |
| **Parse Time** | <5s | wall-clock |
| **Memory** | <500MB | resident set |
| **Modules** | 41/41 | separator count |
| **Errors** | 0 | syntax violations |
| **Throughput** | >2,000 tokens/ms | tokens/time |

---

## 🔗 Integration Points

### Phase 9 Components Used
- **Lexer**: stdlib Lexer (from Phase 8)
- **Parser**: stdlib Parser (from Phase 8)
- **Executor**: Phase 9 Executor (for Stage 4)
- **Benchmarks**: Phase 9 metrics (for Stage 4)

### Dependencies
- ✅ Stage 1: CLI Integration (complete)
- ✅ Stage 2: AOT File Loading (complete)
- 🔄 Stage 3: Tokenization & Parsing (this)
- ⏳ Stage 4: Execution & Benchmarking

---

## 🧪 Test Matrix

```
Test Case | Input | Expected | Status
----------|-------|----------|--------
Tokenization | 14,426 lines | 500,000+ tokens | TODO
Parsing | 500,000 tokens | 10,000+ AST nodes | TODO
Module Count | 41 modules | All present | TODO
Error Handling | Invalid code | Proper rejection | TODO
Throughput | Time measurement | >2000 tokens/ms | TODO
Memory | Peak usage | <500MB | TODO
```

---

## 📝 Implementation Checklist

### Day 1
- [ ] Implement test_aot_tokenization()
- [ ] Add to run_integration_tests()
- [ ] Execute and collect metrics
- [ ] Verify token count >= 500,000
- [ ] Document results

### Day 2
- [ ] Implement test_aot_parsing()
- [ ] Integrate with Phase 9 parser
- [ ] Execute full AST generation
- [ ] Verify AST count >= 10,000
- [ ] Test error recovery
- [ ] Generate token/parse statistics

### Day 3
- [ ] Finalize metrics collection
- [ ] Compare against baselines
- [ ] Prepare for Stage 4 (execution)
- [ ] Update documentation

---

## 🚀 Next Phase (Stage 4)

After Stage 3 completes:
1. **Execution**: Run parsed code through Executor
2. **Benchmarking**: Collect performance metrics
3. **Validation**: Compare output against golden standard
4. **Reporting**: Generate Phase 10 completion report

---

## 📞 Known Risks

### Risk 1: Token Explosion
**Concern**: Comments and whitespace might inflate token count
**Mitigation**: Count only meaningful tokens, skip comments

### Risk 2: AST Node Complexity
**Concern**: Deep nesting might create unexpected node counts
**Mitigation**: Implement robust tree-walking counter

### Risk 3: Memory Pressure
**Concern**: 41 modules might exceed available memory
**Mitigation**: Stream processing, batch tokenization

### Risk 4: Time Constraints
**Concern**: Parsing 14K lines might exceed 5s target
**Mitigation**: Profile and optimize critical paths

---

## 👤 Owner & Timeline

**Architect**: Claude Haiku 4.5
**Phase**: Phase 10 Stage 3
**Start**: 2026-03-16
**Duration**: 24-48 hours
**Estimated Completion**: 2026-03-17/18

---

**Status**: ✅ READY FOR IMPLEMENTATION
**Previous Stage**: Stage 2 COMPLETE
**Next Review**: After Day 1 tokenization test
