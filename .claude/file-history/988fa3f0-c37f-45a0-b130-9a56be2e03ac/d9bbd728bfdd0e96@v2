package memory

import "fmt"

// GC represents a garbage collector
type GC struct {
	objects        map[int]*Object
	roots          []*Object
	marked         map[int]bool
	collectionCount int
	nextID         int
}

// Object represents a managed object
type Object struct {
	ID        int
	Value     interface{}
	Marked    bool
	RefCount  int
	Size      int // Size in bytes
	Type      string
}

// NewGC creates a new garbage collector
func NewGC() *GC {
	return &GC{
		objects: make(map[int]*Object),
		roots:   []*Object{},
		marked:  make(map[int]bool),
		nextID:  0,
	}
}

// Allocate allocates a new object
func (gc *GC) Allocate(value interface{}, size int, objType string) *Object {
	obj := &Object{
		ID:       gc.nextID,
		Value:    value,
		Marked:   false,
		RefCount: 1,
		Size:     size,
		Type:     objType,
	}

	gc.objects[gc.nextID] = obj
	gc.nextID++

	// Trigger GC if needed (threshold-based)
	if len(gc.objects) > 1000 {
		gc.Collect()
	}

	return obj
}

// AddRoot adds an object as a GC root
func (gc *GC) AddRoot(obj *Object) {
	if obj != nil {
		gc.roots = append(gc.roots, obj)
		obj.RefCount++
	}
}

// RemoveRoot removes a GC root
func (gc *GC) RemoveRoot(obj *Object) {
	if obj != nil {
		for i, root := range gc.roots {
			if root.ID == obj.ID {
				gc.roots = append(gc.roots[:i], gc.roots[i+1:]...)
				obj.RefCount--
				break
			}
		}
	}
}

// Collect performs garbage collection (Mark-Sweep)
func (gc *GC) Collect() {
	// Mark phase
	gc.marked = make(map[int]bool)
	for _, root := range gc.roots {
		gc.mark(root)
	}

	// Sweep phase
	for id := range gc.objects {
		if !gc.marked[id] {
			delete(gc.objects, id)
		}
	}

	gc.collectionCount++
}

// mark marks an object and its references
func (gc *GC) mark(obj *Object) {
	if obj == nil || gc.marked[obj.ID] {
		return
	}

	gc.marked[obj.ID] = true
	obj.Marked = true

	// Mark referenced objects (simplified)
	switch v := obj.Value.(type) {
	case []interface{}:
		for _, elem := range v {
			if childObj, ok := elem.(*Object); ok {
				gc.mark(childObj)
			}
		}
	case map[string]interface{}:
		for _, elem := range v {
			if childObj, ok := elem.(*Object); ok {
				gc.mark(childObj)
			}
		}
	}
}

// Free frees an object (manual deallocation)
func (gc *GC) Free(obj *Object) {
	if obj != nil {
		obj.RefCount--
		if obj.RefCount <= 0 {
			delete(gc.objects, obj.ID)
		}
	}
}

// Stats returns GC statistics
type Stats struct {
	TotalObjects   int
	TotalSize      int
	CollectionCount int
	MarkedObjects  int
}

// GetStats returns garbage collection statistics
func (gc *GC) GetStats() Stats {
	totalSize := 0
	for _, obj := range gc.objects {
		totalSize += obj.Size
	}

	return Stats{
		TotalObjects:    len(gc.objects),
		TotalSize:       totalSize,
		CollectionCount: gc.collectionCount,
		MarkedObjects:   len(gc.marked),
	}
}

// String returns GC statistics as string
func (gc *GC) String() string {
	stats := gc.GetStats()
	return fmt.Sprintf("GC Stats: %d objects, %d bytes, %d collections",
		stats.TotalObjects, stats.TotalSize, stats.CollectionCount)
}

// MemoryPool represents a memory pool for fast allocation
type MemoryPool struct {
	pools map[int][]*Object
	size  int
}

// NewMemoryPool creates a new memory pool
func NewMemoryPool() *MemoryPool {
	return &MemoryPool{
		pools: make(map[int][]*Object),
		size:  0,
	}
}

// Allocate allocates from pool
func (mp *MemoryPool) Allocate(size int) *Object {
	// Try to reuse pooled object
	if objects, exists := mp.pools[size]; exists && len(objects) > 0 {
		obj := objects[0]
		mp.pools[size] = objects[1:]
		return obj
	}

	// Allocate new
	return &Object{
		Size: size,
	}
}

// Free returns object to pool
func (mp *MemoryPool) Free(obj *Object) {
	if obj == nil {
		return
	}

	size := obj.Size
	if _, exists := mp.pools[size]; !exists {
		mp.pools[size] = []*Object{}
	}

	// Limit pool size
	if len(mp.pools[size]) < 10 {
		mp.pools[size] = append(mp.pools[size], obj)
	}
}
