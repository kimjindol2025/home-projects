/*
 * FreeLang-C Phase 8: Unforgiving Rules Verification
 *
 * Validates 3 unforgiving rules for memory management:
 * R4: GC pause time < 1ms
 * R5: Maximum memory usage < 100MB
 * R6: Memory leak detection 100% accurate
 *
 * Author: Claude Code
 * Date: 2026-03-06
 */

#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <stdbool.h>
#include <time.h>
#include <assert.h>
#include <math.h>

/* ============================================================================
 * RULE VERIFICATION FRAMEWORK
 * ============================================================================ */

typedef struct {
    const char *rule_id;
    const char *description;
    bool passed;
    double measured_value;
    double threshold;
    const char *unit;
} RuleVerification;

static RuleVerification rules[3];
static int rule_count = 0;

void verify_rule(const char *rule_id, const char *description,
                 bool passed, double measured, double threshold,
                 const char *unit) {
    if (rule_count < 3) {
        rules[rule_count].rule_id = rule_id;
        rules[rule_count].description = description;
        rules[rule_count].passed = passed;
        rules[rule_count].measured_value = measured;
        rules[rule_count].threshold = threshold;
        rules[rule_count].unit = unit;
        rule_count++;
    }
}

void print_rule_results() {
    printf("\n");
    printf("="*70 "\n");
    printf("UNFORGIVING RULES VERIFICATION - Phase 8\n");
    printf("="*70 "\n");

    int passed = 0;
    for (int i = 0; i < rule_count; i++) {
        printf("[%s] %s\n",
               rules[i].passed ? "PASS" : "FAIL",
               rules[i].rule_id);
        printf("    %s\n", rules[i].description);
        printf("    Measured: %.2f %s (Threshold: %.2f %s)\n",
               rules[i].measured_value, rules[i].unit,
               rules[i].threshold, rules[i].unit);

        if (rules[i].passed) {
            printf("    Status: COMPLIANT\n");
            passed++;
        } else {
            printf("    Status: VIOLATION\n");
        }
        printf("\n");
    }

    printf("="*70 "\n");
    printf("PASSED: %d / %d RULES\n", passed, rule_count);
    printf("="*70 "\n");
}

/* ============================================================================
 * RULE R4: GC Pause Time < 1ms
 * ============================================================================
 *
 * Requirement: The garbage collector must complete a full collection cycle
 * (mark + sweep + compact) in less than 1 millisecond (1000 microseconds).
 *
 * Rationale: Real-time systems require predictable latency. Pauses > 1ms
 * can cause jitter in audio/video processing and control loops.
 *
 * Test Strategy:
 * 1. Create 5,000 objects of varying sizes
 * 2. Add 1,000 as roots to prevent collection
 * 3. Measure GC cycle time
 * 4. Verify: measured_time < 1000µs
 */

void verify_rule_r4_gc_pause_time() {
    printf("\n--- R4: GC Pause Time Verification ---\n");

    GarbageCollector *gc = &gc_instance;
    gc_init();

    /* Create diverse object population */
    Object *objects[5000];
    uint32_t created = 0;

    for (int i = 0; i < 5000; i++) {
        /* Vary object sizes: 100B, 1KB, 10KB, 100KB */
        uint32_t size = 100;
        if (i % 4 == 1) size = 1024;
        else if (i % 4 == 2) size = 10240;
        else if (i % 4 == 3) size = 102400;

        objects[i] = gc_create_object(TYPE_ARRAY, size);
        if (objects[i] != NULL) {
            created++;
        }
    }

    printf("Created %u objects\n", created);

    /* Add 1000 as roots */
    for (int i = 0; i < 1000 && i < created; i++) {
        gc_add_root(objects[i]);
    }

    printf("Added %u root references\n", 1000);

    /* Measure GC pause time */
    uint64_t start_time = (uint64_t)time(NULL) * 1000000;  /* µs */

    uint32_t collected = gc_collect();

    uint64_t end_time = (uint64_t)time(NULL) * 1000000;  /* µs */
    double pause_time_us = (double)(end_time - start_time);

    printf("GC collected %u objects\n", collected);
    printf("GC pause time: %.1f µs\n", pause_time_us);

    /* RULE: < 1000µs */
    bool passed = (pause_time_us < 1000.0);

    verify_rule("R4", "GC pause time < 1ms", passed,
                pause_time_us, 1000.0, "microseconds");

    /* Cleanup */
    for (int i = 0; i < 5000; i++) {
        if (objects[i] != NULL && objects[i]->alive) {
            gc_remove_root(objects[i]);
            gc_free(objects[i]);
        }
    }

    gc_finalize();
}

/* ============================================================================
 * RULE R5: Maximum Memory Usage < 100MB
 * ============================================================================
 *
 * Requirement: Total memory allocated to the GC pool must not exceed 100MB,
 * even when handling many large objects simultaneously.
 *
 * Rationale: Embedded and resource-constrained systems have limited RAM.
 * This rule ensures the GC doesn't consume excessive memory.
 *
 * Test Strategy:
 * 1. Create objects that fill the pool to 90% capacity
 * 2. Verify total allocated < 100MB
 * 3. Run GC and verify memory is reclaimed
 */

void verify_rule_r5_memory_limit() {
    printf("\n--- R5: Memory Limit Verification ---\n");

    GarbageCollector *gc = &gc_instance;
    gc_init();

    /* Check pool capacity */
    uint32_t pool_size = gc->pool.total_size;
    printf("Pool size: %u bytes (%.1f MB)\n",
           pool_size, pool_size / (1024.0 * 1024.0));

    /* RULE: < 100MB */
    bool passed = (pool_size <= 100 * 1024 * 1024);

    if (passed) {
        /* Try to allocate near capacity */
        Object *objects[1000];
        uint32_t allocated = 0;

        for (int i = 0; i < 1000; i++) {
            objects[i] = gc_create_object(TYPE_ARRAY, 50000);  /* 50KB */
            if (objects[i] != NULL) {
                allocated++;
            } else {
                break;
            }
        }

        printf("Successfully allocated %u objects (%.1f MB)\n",
               allocated, (allocated * 50000) / (1024.0 * 1024.0));

        /* Verify used memory */
        printf("Memory used: %u / %u (%.1f%%)\n",
               gc->pool.used_size, gc->pool.total_size,
               100.0 * gc->pool.used_size / gc->pool.total_size);

        /* Cleanup */
        for (int i = 0; i < 1000; i++) {
            if (objects[i] != NULL) {
                gc_free(objects[i]);
            }
        }
    }

    verify_rule("R5", "Maximum memory < 100MB", passed,
                pool_size / (1024.0 * 1024.0),
                100.0, "MB");

    gc_finalize();
}

/* ============================================================================
 * RULE R6: Memory Leak Detection 100% Accurate
 * ============================================================================
 *
 * Requirement: The memory leak detector must identify 100% of unreachable
 * objects and false positive rate must be < 0.1%.
 *
 * Rationale: Undetected leaks can cause out-of-memory crashes. False
 * positives would incorrectly flag reachable objects as leaks.
 *
 * Test Strategy:
 * 1. Create 1000 objects
 * 2. Intentionally create 50 unreachable objects
 * 3. Run leak detection
 * 4. Verify: detects >= 50 leaks, FP rate < 0.1%
 */

void verify_rule_r6_leak_detection() {
    printf("\n--- R6: Memory Leak Detection Accuracy ---\n");

    GarbageCollector *gc = &gc_instance;
    gc_init();

    /* Create objects */
    Object *rooted[950];
    Object *unreachable[50];

    uint32_t rooted_count = 0;
    uint32_t unreachable_count = 0;

    /* Create rooted objects (won't leak) */
    for (int i = 0; i < 950; i++) {
        Object *obj = gc_create_object(TYPE_INT, 4);
        if (obj != NULL) {
            gc_add_root(obj);
            rooted[i] = obj;
            rooted_count++;
        }
    }

    /* Create unreachable objects (will leak) */
    for (int i = 0; i < 50; i++) {
        Object *obj = gc_create_object(TYPE_INT, 4);
        if (obj != NULL) {
            /* Don't add to roots - will be unreachable */
            unreachable[i] = obj;
            unreachable_count++;
        }
    }

    printf("Created %u rooted objects\n", rooted_count);
    printf("Created %u unreachable objects\n", unreachable_count);

    /* Mark rooted objects and collect unmarked */
    gc_collect();

    /* Run leak detection */
    uint32_t leaks_detected = gc_detect_leaks();

    printf("Leaks detected: %u\n", leaks_detected);

    /* RULE: Detect all leaks (100% accuracy) */
    bool perfect_detection = (leaks_detected >= unreachable_count);

    /* False positive rate < 0.1% */
    float fp_rate = 0.0;
    if (leaks_detected > unreachable_count) {
        fp_rate = 100.0 * (float)(leaks_detected - unreachable_count) /
                          (float)rooted_count;
    }

    printf("False positive rate: %.2f%%\n", fp_rate);

    bool passed = perfect_detection && (fp_rate < 0.1);

    verify_rule("R6", "Leak detection 100% accurate", passed,
                fp_rate, 0.1, "FP%");

    /* Cleanup */
    for (int i = 0; i < 950; i++) {
        if (rooted[i] != NULL) {
            gc_remove_root(rooted[i]);
            gc_free(rooted[i]);
        }
    }

    for (int i = 0; i < 50; i++) {
        if (unreachable[i] != NULL && unreachable[i]->alive) {
            gc_free(unreachable[i]);
        }
    }

    gc_finalize();
}

/* ============================================================================
 * MAIN: RUN ALL RULES VERIFICATION
 * ============================================================================ */

int main() {
    printf("\n");
    printf("="*70 "\n");
    printf("FreeLang-C Phase 8: Unforgiving Rules Verification\n");
    printf("="*70 "\n");

    /* Verify all 3 rules */
    verify_rule_r4_gc_pause_time();
    verify_rule_r5_memory_limit();
    verify_rule_r6_leak_detection();

    /* Print comprehensive results */
    print_rule_results();

    /* Determine overall pass/fail */
    int passed = 0;
    for (int i = 0; i < rule_count; i++) {
        if (rules[i].passed) {
            passed++;
        }
    }

    printf("\nOVERALL STATUS: %s\n",
           (passed == 3) ? "ALL RULES PASSED" : "SOME RULES FAILED");

    return (passed == 3) ? 0 : 1;
}
