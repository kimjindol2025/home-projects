# FreeLang-LLC: ManagedPointer Specification

## 📌 Overview

ManagedPointer is FreeLang's answer to **safe low-level memory access**. It provides C-pointer-level power while maintaining FreeLang's integrity guarantees.

**Philosophy**: Direct hardware control without undefined behavior.

---

## 🎯 Design Goals

| Goal | Mechanism | Benefit |
|------|-----------|---------|
| **Safety** | Bounds checking + lifetime tracking | No buffer overflows |
| **Performance** | Zero-copy MMIO + volatile access | Direct hardware control |
| **Debuggability** | Lifetime IDs + access tracking | Clear error messages |
| **Hardware Access** | MMIO registers + cache flushing | True bare-metal |

---

## 🏗️ ManagedPointer Structure

```rust
pub struct ManagedPointer {
    pub address: u64,           // Raw memory address
    pub size_bytes: u64,        // Valid memory region size
    pub is_volatile: bool,      // MMIO flag
    pub element_type: String,   // Type name ("u32", "u64", etc)
    pub lifetime_id: u64,       // Unique allocation ID
    pub created_at_ns: u64,     // Creation timestamp
}
```

### Fields Explained

**address**: Raw 64-bit memory address
- Can point anywhere (RAM, MMIO, device registers)
- No hidden offsets or indirection

**size_bytes**: Valid memory region size
- Used for bounds checking
- Prevents accidental reads beyond region
- Can be arbitrary (1 byte to GB range)

**is_volatile**: Volatile MMIO flag
- True: volatile reads/writes with memory barriers
- False: normal cached reads/writes
- MMIO devices require `is_volatile=true`

**element_type**: Type information
- "u8", "u16", "u32", "u64", "u128" supported
- Used to calculate element size for stride()
- Size determining alignment requirements

**lifetime_id**: Lifetime tracking
- Unique ID for each allocation
- Prevents use-after-free via pointer comparison
- Prevents mixing pointers from different allocations

**created_at_ns**: Creation timestamp
- Enables lifetime expiration validation
- Detects pointer reuse bugs
- Age-based validation for safety

---

## 🔒 Safety Rules

### Rule 1: Bounds Checking
Every pointer dereference is bounds-checked:
```
Required: offset + element_size ≤ size_bytes
```

**Example**:
```freelang
let ptr = ManagedPointer::new(0x1000, 256, "u32".to_string());
// Valid: offset 0, element_size 4, total ≤ 256 ✓
safe_read(&ptr, 0)  // OK

// Invalid: offset 252, element_size 4, total = 256 = size_bytes ✓ (edge case)
safe_read(&ptr, 252)  // OK (inclusive upper bound)

// Invalid: offset 256, element_size 4, total > 256 ✗
safe_read(&ptr, 256)  // ERROR: out of bounds
```

### Rule 2: Lifetime Isolation
Pointers from different allocations cannot be mixed:
```
Can only subtract pointers with lifetime_id₁ == lifetime_id₂
```

**Example**:
```freelang
let ptr1 = ManagedPointer::new(0x1000, 256, "u32".to_string());
let ptr2 = ManagedPointer::new(0x2000, 256, "u32".to_string());

// Invalid: different lifetime_ids
difference(&ptr2, &ptr1)  // ERROR: different allocations
```

### Rule 3: Alignment
Non-byte pointers must have aligned offsets:
```
offset % 4 == 0  (for u32, u64, etc)
offset % 1 == 0  (for u8, allowed unaligned)
```

**Example**:
```freelang
let ptr = ManagedPointer::new(0x1000, 256, "u32".to_string());

safe_read(&ptr, 0)   // OK: 0 % 4 == 0 ✓
safe_read(&ptr, 4)   // OK: 4 % 4 == 0 ✓
safe_read(&ptr, 1)   // ERROR: 1 % 4 != 0 ✗
safe_read(&ptr, 2)   // ERROR: 2 % 4 != 0 ✗
```

### Rule 4: Volatile-Safe Access
Volatile access only allowed on volatile pointers:
```
volatile_read/volatile_write requires is_volatile == true
```

**Example**:
```freelang
let ptr = ManagedPointer::new(0x1000, 256, "u32".to_string());
volatile_read(&ptr, 0)  // ERROR: not volatile

let mmio = ManagedPointer::new_volatile(0xDEAD0000, 4096, "u32".to_string());
volatile_read(&mmio, 0)  // OK: volatile flag set ✓
```

---

## 📊 Element Size Table

| Type | Size | Alignment | Use Case |
|------|------|-----------|----------|
| **u8** | 1 byte | No align req | Byte access, strings |
| **u16** | 2 bytes | 2-byte align | 16-bit registers |
| **u32** | 4 bytes | 4-byte align | 32-bit registers, ints |
| **u64** | 8 bytes | 4-byte align | 64-bit registers, pointers |
| **u128** | 16 bytes | 4-byte align | XMM registers |

---

## 🔄 Memory Operations

### Safe Read
```freelang
fn safe_read(ptr: &ManagedPointer, offset: u64) -> Result<u64, String>
```
- ✓ Bounds checking
- ✓ Lifetime validation
- ✓ No memory barriers (normal cached access)
- Returns value read from (address + offset)

**Example**:
```freelang
let ptr = ManagedPointer::new(0x1000, 256, "u32".to_string());
match safe_read(&ptr, 8) {
    Ok(value) => println!("Read: {}", value),
    Err(e) => println!("Error: {}", e),
}
```

### Safe Write
```freelang
fn safe_write(ptr: &ManagedPointer, offset: u64, value: u64) -> Result<(), String>
```
- ✓ Bounds checking
- ✓ Lifetime validation
- ✓ No memory barriers
- Writes value to (address + offset)

**Example**:
```freelang
let ptr = ManagedPointer::new(0x1000, 256, "u32".to_string());
safe_write(&ptr, 8, 0x12345678)?;
```

### Volatile Read (MMIO)
```freelang
fn volatile_read(ptr: &ManagedPointer, offset: u64) -> Result<u64, String>
```
- ✓ Requires `is_volatile=true`
- ✓ Memory barrier before read
- ✓ Memory barrier after read
- ✓ Prevents compiler optimization/reordering
- Use for hardware registers, device memory

**Example**:
```freelang
let status_reg = ManagedPointer::new_volatile(0xDEAD0000, 8, "u32".to_string());
let status = volatile_read(&status_reg, 0)?;  // Read device status
```

### Volatile Write (MMIO)
```freelang
fn volatile_write(ptr: &ManagedPointer, offset: u64, value: u64) -> Result<(), String>
```
- ✓ Requires `is_volatile=true`
- ✓ Memory barrier before write
- ✓ Memory barrier after write
- ✓ Guarantees write is visible to hardware
- Use for device control registers

**Example**:
```freelang
let control_reg = ManagedPointer::new_volatile(0xDEAD0004, 8, "u32".to_string());
volatile_write(&control_reg, 0, 0x00000001)?;  // Trigger device
```

---

## ➕ Pointer Arithmetic

### offset()
```freelang
fn offset(ptr: &ManagedPointer, byte_offset: u64) -> Result<ManagedPointer, String>
```
Returns new pointer with adjusted address:
```
new_address = ptr.address + byte_offset
new_size = ptr.size_bytes - byte_offset
```

**Example**:
```freelang
let ptr = ManagedPointer::new(0x1000, 256, "u32".to_string());
let ptr2 = offset(&ptr, 32)?;  // Point to 0x1020
assert_eq!(ptr2.address, 0x1020);
assert_eq!(ptr2.size_bytes, 224);  // 256 - 32
```

### stride()
```freelang
fn stride(ptr: &ManagedPointer, element_index: u64) -> Result<ManagedPointer, String>
```
Advances by N elements (automatic size calculation):
```
byte_offset = element_index * element_size
```

**Example**:
```freelang
let ptr = ManagedPointer::new(0x1000, 256, "u32".to_string());
let ptr2 = stride(&ptr, 5)?;  // 5 × 4 bytes = +20 bytes
assert_eq!(ptr2.address, 0x1014);  // 0x1000 + 20
```

### difference()
```freelang
fn difference(ptr1: &ManagedPointer, ptr2: &ManagedPointer) -> Result<i64, String>
```
Returns byte distance between pointers (same allocation only):
```
difference = ptr1.address - ptr2.address
```

**Example**:
```freelang
let ptr1 = ManagedPointer::new(0x1000, 256, "u32".to_string());
let ptr2 = offset(&ptr1, 40)?;
let diff = difference(&ptr2, &ptr1)?;
assert_eq!(diff, 40);
```

---

## 🔐 Lifetime Management

### Lifetime Tracking
Each allocation gets unique ID:
```freelang
let ptr1 = ManagedPointer::new(0x1000, 256, "u32".to_string());
let ptr2 = ManagedPointer::new(0x2000, 256, "u32".to_string());

// ptr1 and ptr2 have different lifetime_id values
assert_ne!(ptr1.lifetime_id, ptr2.lifetime_id);
```

### Lifetime Validation
```freelang
fn validate_lifetime(ptr: &ManagedPointer, max_age_ns: u64) -> bool
```
Check if pointer is still "alive":
```
age = current_ns() - created_at_ns
valid = age ≤ max_age_ns
```

**Example**:
```freelang
let ptr = ManagedPointer::new(0x1000, 256, "u32".to_string());
assert!(validate_lifetime(&ptr, 1_000_000));  // Valid if <1 second old
```

---

## 🎓 Common Patterns

### MMIO Device Driver
```freelang
// Map device memory region
let device_regs = ManagedPointer::new_volatile(0xDEAD0000, 4096, "u32".to_string());

// Read status register (offset 0)
let status = volatile_read(&device_regs, 0)?;

// Read data register (offset 4)
let data = volatile_read(&device_regs, 4)?;

// Write control register (offset 8)
volatile_write(&device_regs, 8, 0x00000001)?;
```

### Array Access
```freelang
// Allocate array in memory
let base = ManagedPointer::new(0x2000, 1024, "u32".to_string());

// Access element [0]
safe_read(&base, 0)?;

// Access element [5] via stride
let elem5 = stride(&base, 5)?;
safe_read(&elem5, 0)?;

// Or direct offset
safe_read(&base, 5 * 4)?;  // 5 × 4 bytes
```

### Memory Copy
```freelang
let src = ManagedPointer::new(0x1000, 256, "u8".to_string());
let dst = ManagedPointer::new(0x2000, 256, "u8".to_string());

for i in 0..256 {
    let byte = safe_read(&src, i)?;
    safe_write(&dst, i, byte)?;
}
```

---

## 🧪 Unforgiving Rules (Validation)

| Rule | Requirement | Test |
|------|-------------|------|
| **R1** | Bounds checking 100% | Out-of-bounds read/write fails ✓ |
| **R2** | MMIO latency <500ns | volatile_read/write ≤500ns ✓ |
| **R3** | Volatile ordering | Barriers prevent reordering ✓ |
| **R4** | No UB in safe code | Safe functions never panic ✓ |
| **R5** | Lifetime isolation | Different allocs can't be subtracted ✓ |

---

## 📝 Summary

| Feature | Status | Use Case |
|---------|--------|----------|
| **Safe Read/Write** | ✅ | General memory access |
| **Volatile MMIO** | ✅ | Hardware registers |
| **Pointer Arithmetic** | ✅ | Array/struct navigation |
| **Lifetime Tracking** | ✅ | Use-after-free prevention |
| **Cache Control** | ✅ (via intrinsics.fl) | Performance tuning |
| **Atomic Ops** | ✅ (via intrinsics.fl) | Concurrency |

---

**Version**: 1.0
**Status**: Phase 1 Complete
**Philosophy**: "기록이 증명이다" — Safety through design, not restriction.
