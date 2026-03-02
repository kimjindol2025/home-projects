// GRIE 3-A: Zig Build Configuration (v0.14+)
// Generates static library (.a) and shared library (.so) for Go FFI
// Supports multiple targets: aarch64-linux-android, x86_64-linux-gnu, aarch64-linux-gnu
//
// Usage:
//   # Default (native build)
//   zig build -Doptimize=ReleaseFast
//
//   # Cross-compile for ARM64 Android
//   zig build --prefix=/tmp/grie-kernel -Dtarget=aarch64-linux-android -Doptimize=ReleaseFast
//
//   # Cross-compile for x86_64 Linux
//   zig build --prefix=/tmp/grie-kernel -Dtarget=x86_64-linux-gnu -Doptimize=ReleaseFast
//
//   # Run tests
//   zig build test
//
const std = @import("std");
const builtin = @import("builtin");

pub fn build(b: *std.Build) void {
    // Standard options
    const target = b.standardTargetOptions(.{});
    const optimize = b.standardOptimizeOption(.{});

    // Custom option: Build mode (static, shared, both)
    const build_mode = b.option([]const u8, "mode", "Build mode: static, shared, or both (default: both)") orelse "both";

    // Static library (for embedding in Go binary)
    const lib_static = b.addStaticLibrary(.{
        .name = "grie_kernel",
        .root_source_file = b.path("src/main.zig"),
        .target = target,
        .optimize = optimize,
    });

    // Add standard library
    lib_static.linkLibC();

    // Shared library (for dynamic linking)
    const lib_shared = b.addSharedLibrary(.{
        .name = "grie_kernel",
        .root_source_file = b.path("src/main.zig"),
        .target = target,
        .optimize = optimize,
    });

    lib_shared.linkLibC();

    // Install artifacts based on build_mode
    if (std.mem.eql(u8, build_mode, "static") or std.mem.eql(u8, build_mode, "both")) {
        b.installArtifact(lib_static);
    }

    if (std.mem.eql(u8, build_mode, "shared") or std.mem.eql(u8, build_mode, "both")) {
        b.installArtifact(lib_shared);
    }

    // Unit tests
    const unit_tests = b.addTest(.{
        .root_source_file = b.path("src/main.zig"),
        .target = target,
        .optimize = optimize,
    });

    const run_unit_tests = b.addRunArtifact(unit_tests);

    // Test step
    const test_step = b.step("test", "Run unit tests");
    test_step.dependOn(&run_unit_tests.step);

    // Add default test step
    b.default_step.dependOn(test_step);

    // Info step
    const info_step = b.step("info", "Print build information");
    const info = b.addSystemCommand(&.{ "echo", "GRIE Kernel Build Configuration" });
    info_step.dependOn(&info.step);
}
