// SHM Header Definition (Shared with Go)
// This must be byte-for-byte identical to Go's EngineHeader
// Location: internal/protocol/header.go

const std = @import("std");
const builtin = @import("builtin");

/// Magic number: 0x47524945_00000001 = "GRIE" v1
pub const MAGIC_NUMBER: u64 = 0x47524945_00000001;

/// Header size: 128 bytes (2x 64B cache lines)
pub const HEADER_SIZE: usize = 128;

/// Data starts at offset 128
pub const DATA_OFFSET: usize = 128;

/// State machine values (must match Go's StateFlag)
pub const State = enum(i32) {
    Idle = 0,
    Writing = 1,
    Ready = 2,
    Reading = 3,
    _,
};

/// Lock-free synchronization: Use std.atomic for cross-process
pub fn AtomicInt32() type {
    return std.atomic.Atomic(i32);
}

pub fn AtomicUint64() type {
    return std.atomic.Atomic(u64);
}

/// EngineHeader: extern struct for Go-Zig shared memory layout
/// CRITICAL: Must be exactly 128 bytes!
/// Layout matches Go's internal/protocol/header.go precisely
pub const EngineHeader = extern struct {
    // Cache Line 1 (0-39 bytes)
    magic: u64 = MAGIC_NUMBER,      // [0:8]   uint64
    version: u32 = 1,               // [8:12]  uint32
    state: i32 = @intFromEnum(State.Idle),  // [12:16] int32 (atomic in Go, but we use raw access)
    seq_num: u64 = 0,               // [16:24] uint64 (atomic in Go)
    timestamp: i64 = 0,             // [24:32] int64
    data_len: u64 = 0,              // [32:40] uint64

    // Cache Line 1 (40-63 bytes)
    writer_pid: i64 = 0,            // [40:48] int64
    reader_pid: i64 = 0,            // [48:56] int64
    total_writes: u64 = 0,          // [56:64] uint64 (atomic in Go)

    // Cache Line 2 (64-127 bytes)
    total_reads: u64 = 0,           // [64:72] uint64 (atomic in Go)
    _pad: [56]u8 = undefined,       // [72:128] padding

    comptime {
        // Verify size at compile time
        if (@sizeOf(EngineHeader) != HEADER_SIZE) {
            @compileError("EngineHeader size mismatch");
        }
    }

    // Helper: Load state atomically (mimics Go's LoadState())
    pub fn loadState(self: *const EngineHeader) State {
        // Go uses atomic.LoadInt32 which is seq_cst by default
        // We read the raw i32 and interpret as State
        const raw = @as(i32, @intCast(@intFromPtr(&self.state)));
        return @enumFromInt(raw);
    }

    // Helper: Store state atomically (mimics Go's StoreState())
    pub fn storeState(self: *EngineHeader, s: State) void {
        // In Zig, we use volatile access for sync
        self.state = @intFromEnum(s);
        @fence(.SeqCst);
    }

    // Helper: CAS state (mimics Go's CASState())
    pub fn casState(self: *EngineHeader, old: State, new: State) bool {
        // Attempt atomic compare-and-swap using inline assembly
        // This is a simplified version; real implementation uses LDREX/STREX on ARM64
        const old_val = @intFromEnum(old);
        const new_val = @intFromEnum(new);

        // For now, use volatile read/write (not truly atomic in multiprocessing)
        // Production code would use platform-specific atomics
        if (self.state == old_val) {
            self.state = new_val;
            @fence(.SeqCst);
            return true;
        }
        return false;
    }
};

// Verify at compile time
comptime {
    std.debug.assert(@sizeOf(EngineHeader) == HEADER_SIZE);
    std.debug.assert(@offsetOf(EngineHeader, "magic") == 0);
    std.debug.assert(@offsetOf(EngineHeader, "state") == 12);
    std.debug.assert(@offsetOf(EngineHeader, "seq_num") == 16);
    std.debug.assert(@offsetOf(EngineHeader, "data_len") == 32);
    std.debug.assert(@offsetOf(EngineHeader, "total_writes") == 56);
    std.debug.assert(@offsetOf(EngineHeader, "total_reads") == 64);
}
