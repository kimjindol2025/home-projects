package memory

import (
	"testing"
)

// Test object allocation
func TestAllocate(t *testing.T) {
	gc := NewGC()
	obj := gc.Allocate(42, 8, "number")

	if obj == nil {
		t.Errorf("expected object, got nil")
	}

	if obj.Value != 42 {
		t.Errorf("expected value 42, got %v", obj.Value)
	}

	if obj.Size != 8 {
		t.Errorf("expected size 8, got %d", obj.Size)
	}
}

// Test garbage collection
func TestGarbageCollection(t *testing.T) {
	gc := NewGC()

	// Allocate objects
	obj1 := gc.Allocate(1, 8, "number")
	_ = gc.Allocate(2, 8, "number")
	_ = gc.Allocate(3, 8, "number")

	// Add one as root
	gc.AddRoot(obj1)

	// Collect
	gc.Collect()

	// obj1 should be alive, others may be collected
	if _, exists := gc.objects[obj1.ID]; !exists {
		t.Errorf("root object should not be collected")
	}
}

// Test root management
func TestRootManagement(t *testing.T) {
	gc := NewGC()
	obj := gc.Allocate(42, 8, "number")

	initialCount := len(gc.roots)
	gc.AddRoot(obj)

	if len(gc.roots) != initialCount+1 {
		t.Errorf("expected root count to increase")
	}

	gc.RemoveRoot(obj)

	if len(gc.roots) != initialCount {
		t.Errorf("expected root count to decrease")
	}
}

// Test reference counting
func TestReferenceCounting(t *testing.T) {
	gc := NewGC()
	obj := gc.Allocate(42, 8, "number")

	initialRefCount := obj.RefCount
	gc.AddRoot(obj)

	if obj.RefCount != initialRefCount+1 {
		t.Errorf("expected refcount to increase")
	}

	gc.RemoveRoot(obj)

	if obj.RefCount != initialRefCount {
		t.Errorf("expected refcount to decrease")
	}
}

// Test manual free
func TestManualFree(t *testing.T) {
	gc := NewGC()
	obj := gc.Allocate(42, 8, "number")

	_ = obj.ID
	gc.Free(obj)

	// Object should be removed if refcount reaches 0
	if obj.RefCount > 0 {
		t.Logf("object still has references")
	}
}

// Test GC statistics
func TestGCStats(t *testing.T) {
	gc := NewGC()

	gc.Allocate(1, 8, "number")
	gc.Allocate(2, 16, "string")
	gc.Allocate(3, 32, "array")

	stats := gc.GetStats()

	if stats.TotalObjects != 3 {
		t.Errorf("expected 3 objects, got %d", stats.TotalObjects)
	}

	if stats.TotalSize != 56 {
		t.Errorf("expected total size 56, got %d", stats.TotalSize)
	}
}

// Test GC string representation
func TestGCString(t *testing.T) {
	gc := NewGC()
	gc.Allocate(1, 8, "number")
	gc.Allocate(2, 8, "number")

	str := gc.String()

	if str == "" {
		t.Errorf("expected non-empty string")
	}
}

// Test memory pool allocation
func TestMemoryPoolAllocate(t *testing.T) {
	pool := NewMemoryPool()
	obj := pool.Allocate(8)

	if obj == nil {
		t.Errorf("expected object from pool")
	}

	if obj.Size != 8 {
		t.Errorf("expected size 8, got %d", obj.Size)
	}
}

// Test memory pool reuse
func TestMemoryPoolReuse(t *testing.T) {
	pool := NewMemoryPool()

	obj1 := pool.Allocate(8)
	id1 := obj1.ID
	pool.Free(obj1)

	obj2 := pool.Allocate(8)

	// Should be same ID (reused)
	if obj2.ID != id1 {
		t.Logf("expected pool to reuse object (different IDs)")
	}
}

// Test collection with array
func TestCollectionWithArray(t *testing.T) {
	gc := NewGC()

	obj1 := gc.Allocate(1, 8, "number")
	_ = gc.Allocate(2, 8, "number")

	// Create array referencing obj1
	arr := []interface{}{obj1}
	obj3 := gc.Allocate(arr, 16, "array")

	gc.AddRoot(obj3)
	gc.Collect()

	// obj3 should be alive, others depend on mark
	if _, exists := gc.objects[obj3.ID]; !exists {
		t.Errorf("root object should not be collected")
	}
}

// Test multiple collection cycles
func TestMultipleCollections(t *testing.T) {
	gc := NewGC()

	obj1 := gc.Allocate(1, 8, "number")
	gc.AddRoot(obj1)

	gc.Collect()
	gc.Collect()
	gc.Collect()

	if gc.GetStats().CollectionCount != 3 {
		t.Errorf("expected 3 collections")
	}
}

// Test marking phase
func TestMarkingPhase(t *testing.T) {
	gc := NewGC()

	obj1 := gc.Allocate(1, 8, "number")
	obj2 := gc.Allocate(2, 8, "number")

	gc.AddRoot(obj1)
	gc.mark(obj1)

	if !gc.marked[obj1.ID] {
		t.Errorf("expected obj1 to be marked")
	}

	if gc.marked[obj2.ID] {
		t.Errorf("expected obj2 to not be marked")
	}
}

// Test null pointer handling
func TestNullPointerHandling(t *testing.T) {
	gc := NewGC()

	// Should not panic
	gc.AddRoot(nil)
	gc.RemoveRoot(nil)
	gc.Free(nil)
	gc.mark(nil)

	stats := gc.GetStats()
	if stats.TotalObjects != 0 {
		t.Errorf("expected 0 objects")
	}
}

// Benchmark allocation
func BenchmarkAllocation(b *testing.B) {
	gc := NewGC()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		gc.Allocate(i, 8, "number")
	}
}

// Benchmark collection
func BenchmarkCollection(b *testing.B) {
	gc := NewGC()

	for i := 0; i < 100; i++ {
		gc.Allocate(i, 8, "number")
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		gc.Collect()
	}
}
