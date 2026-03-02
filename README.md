# FreeLang Rust Runtime

**Phase B - Independent Rust-based Execution Engine**

![Status](https://img.shields.io/badge/Status-Development-yellow)
![Rust](https://img.shields.io/badge/Rust-2021-orange)
![License](https://img.shields.io/badge/License-MIT-green)

---

## 🚀 **Overview**

FreeLang Rust Runtime is a complete, independent execution engine written entirely in Rust. It provides:

- **Complete Type System**: 8 core value types (null, bool, number, string, array, object, function, error)
- **Memory Management**: Reference counting + Mark & Sweep GC
- **Standard Library**: 50+ built-in functions (I/O, String, Array, Math, System, Crypto, JSON)
- **High Performance**: <100ms startup, <10ms per function execution
- **Zero Dependencies**: Minimal external dependencies (parking_lot only)

---

## 📐 **Architecture**

### 4-Layer Design

```
┌─────────────────────────────────────┐
│ Layer 4: Standard Library (50+ funcs)│
├─────────────────────────────────────┤
│ Layer 3: Runtime Engine             │
├─────────────────────────────────────┤
│ Layer 2: Memory Management          │
├─────────────────────────────────────┤
│ Layer 1: Core Types & Traits        │
└─────────────────────────────────────┘
```

### Module Structure

```
freelang-runtime/
├── src/
│   ├── core/          # Type system
│   │   ├── value.rs
│   │   ├── object.rs
│   │   └── traits.rs
│   ├── memory/        # Memory management
│   │   ├── allocator.rs
│   │   ├── gc.rs
│   │   └── refcount.rs
│   ├── runtime/       # Execution engine
│   │   ├── vm.rs
│   │   ├── executor.rs
│   │   └── context.rs
│   └── stdlib/        # Standard library
│       ├── io.rs
│       ├── string.rs
│       ├── array.rs
│       ├── math.rs
│       ├── system.rs
│       ├── crypto.rs
│       └── json.rs
```

---

## 📦 **Value Types**

```rust
pub enum Value {
    Null,                          // null
    Bool(bool),                    // true/false
    Number(f64),                   // Integer & Float
    String(String),                // UTF-8 String
    Array(Rc<RefCell<Vec>>),      // Mutable Array
    Object(Rc<RefCell<HashMap>>), // Mutable Object
    Function(Box<dyn Callable>),  // Function
    Error(String),                 // Error
}
```

---

## 📚 **Standard Library**

### I/O Functions (15)
```
print()          println()      input()
read_file()      write_file()   append_file()
file_exists()    delete_file()  list_files()
mkdir()          rmdir()        get_cwd()
cd()             env()          system()
```

### String Functions (18)
```
length()         upper()        lower()
trim()           split()        join()
substring()      contains()     startswith()
endswith()       replace()      find()
reverse()        repeat()       pad_left()
pad_right()      to_upper_case() to_lower_case()
```

### Array Functions (12)
```
length()         push()         pop()
shift()          unshift()      slice()
concat()         reverse()      sort()
unique()         contains()     index_of()
```

### Math Functions (15)
```
abs()            round()        floor()
ceil()           sqrt()         pow()
sin()            cos()          tan()
log()            log10()        exp()
random()         min()          max()
```

### System Functions (8)
```
now()            sleep()        exit()
typeof()         is_null()      is_array()
is_object()      memory_usage()
```

### Crypto Functions (8)
```
hash_md5()       hash_sha256()  random_bytes()
base64_encode()  base64_decode() hex_encode()
hex_decode()     md5()
```

### JSON Functions (5)
```
json_parse()     json_stringify() json_pretty()
json_validate()  json_minify()
```

---

## ⚡ **Performance Goals**

| Metric | Goal | Status |
|--------|------|--------|
| **Startup Time** | <100ms | 🟡 In Progress |
| **Function Call** | <10ms | 🟡 In Progress |
| **Memory Usage** | <10MB base | 🟡 In Progress |
| **Throughput** | 10,000 calls/sec | 🟡 In Progress |

---

## 🔧 **Building**

```bash
# Clone repository
git clone https://gogs.dclub.kr/kim/freelang-runtime.git
cd freelang-runtime

# Build
cargo build --release

# Run tests
cargo test

# Run benchmarks
cargo bench
```

---

## 📋 **Development Timeline**

### Week 1: Core + Memory (3/2-3/8)
- [ ] Value & Object types
- [ ] MemoryAllocator
- [ ] Reference Counting
- [ ] Mark & Sweep GC
- [ ] 50+ Unit Tests

### Week 2: Runtime Engine (3/9-3/15)
- [ ] ExecutionContext
- [ ] CallStack
- [ ] IR Interpreter
- [ ] Function Calls
- [ ] 30+ Integration Tests

### Week 3: Standard Library (3/16-3/22)
- [ ] I/O functions (15)
- [ ] String functions (18)
- [ ] Array functions (12)
- [ ] Math functions (15)
- [ ] System functions (8)
- [ ] 68+ Function Tests

### Week 4: Crypto + Optimization (3/23-3/30)
- [ ] Crypto functions (8)
- [ ] JSON functions (5)
- [ ] Performance Optimization
- [ ] Documentation
- [ ] 20+ Integration Tests

---

## 📊 **Progress**

```
Phase B Status: Week 1 Preparation
├── Project Structure       ✅ Complete
├── Core Types             🟡 In Progress
├── Memory Management      🟡 In Progress
├── Runtime Engine         🟡 In Progress
├── Standard Library       🟡 In Progress
└── Testing & Benchmarks   🟡 In Progress

Overall: 20% (Foundation Phase)
```

---

## 📝 **Code Quality Standards**

- **Coverage**: >90%
- **Lint**: Zero warnings (Clippy)
- **Tests**: 100+ unit + 50+ integration
- **Documentation**: Complete API docs

---

## 🤝 **Contributing**

FreeLang is an open-source project. Contributions welcome!

- Report issues
- Submit PRs
- Improve documentation
- Write tests

---

## 📄 **License**

MIT License - See LICENSE file

---

## 🔗 **Related Projects**

- [FreeLang Compiler](https://gogs.dclub.kr/kim/v2-freelang-ai) - Phase A
- [FreeLang Bank System](https://gogs.dclub.kr/kim/freelang-bank-system) - Example
- [FreeLang Specification](https://gogs.dclub.kr/kim/freelang-v4-integrated) - Language Spec

---

**Phase B - Rust Runtime Development**
*Making FreeLang completely independent* 🚀

Last Updated: 2026-03-02
