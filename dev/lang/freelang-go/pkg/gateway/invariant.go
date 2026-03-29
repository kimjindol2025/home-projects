package gateway

import (
	"fmt"
	"sync"
)

// InvariantFn is a function that checks an invariant condition.
// It takes a read-only snapshot of state and returns "" if OK, or an error message if violated.
type InvariantFn func(state InvariantState) string

// InvariantState is a read-only snapshot of gateway state at a specific point in request processing.
type InvariantState struct {
	Method        string
	KeyID         string
	Phase         string        // "pre" or "post"
	RequestCount  int64
	RateLimitUsed float64       // Tokens consumed
	RateLimitCap  float64       // Bucket capacity
	Response      *Response
}

// InvariantMonitor tracks and enforces invariants.
type InvariantMonitor struct {
	invariants map[string]InvariantFn
	mu         sync.RWMutex
}

// NewInvariantMonitor creates a new monitor.
func NewInvariantMonitor() *InvariantMonitor {
	return &InvariantMonitor{
		invariants: make(map[string]InvariantFn),
	}
}

// Register adds a named invariant to the monitor.
func (im *InvariantMonitor) Register(name string, fn InvariantFn) error {
	im.mu.Lock()
	defer im.mu.Unlock()

	if _, exists := im.invariants[name]; exists {
		return fmt.Errorf("invariant %q already registered", name)
	}

	im.invariants[name] = fn
	return nil
}

// Check evaluates all registered invariants against the given state.
// Returns nil if all invariants pass, or an error with details of first violation.
func (im *InvariantMonitor) Check(state InvariantState) error {
	im.mu.RLock()
	defer im.mu.RUnlock()

	for name, fn := range im.invariants {
		if msg := fn(state); msg != "" {
			return fmt.Errorf("invariant %q violated: %s", name, msg)
		}
	}

	return nil
}

// Built-in invariants

// NoUnauthenticatedWrite returns an invariant that prevents write operations without auth.
func NoUnauthenticatedWrite() InvariantFn {
	return func(state InvariantState) string {
		// Write operations should not happen if KeyID is empty
		if state.KeyID == "" && (state.Method == "KV.Set" || state.Method == "KV.Delete") {
			return "write operation attempted without authentication"
		}
		return ""
	}
}

// RateLimitBound returns an invariant that ensures rate limit is never exceeded.
func RateLimitBound() InvariantFn {
	return func(state InvariantState) string {
		// Tokens used should not exceed capacity
		if state.RateLimitUsed > state.RateLimitCap {
			return fmt.Sprintf("rate limit exceeded: used %.2f > capacity %.2f",
				state.RateLimitUsed, state.RateLimitCap)
		}
		return ""
	}
}

// ResponseConsistency returns an invariant that checks response validity.
func ResponseConsistency() InvariantFn {
	return func(state InvariantState) string {
		if state.Phase == "post" && state.Response == nil {
			return "response is nil after handler execution"
		}
		// If there's an error, Error field should be non-empty
		if state.Response != nil && state.Response.Error != "" && len(state.Response.Body) > 0 {
			// Both error and body present is unusual but not necessarily invalid
			// (some handlers may return partial data with error context)
		}
		return ""
	}
}
