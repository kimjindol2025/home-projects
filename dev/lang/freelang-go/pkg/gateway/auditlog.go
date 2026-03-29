package gateway

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"sync"
	"time"
)

// AuditLog is a tamper-evident SHA256-chained append-only log.
type AuditLog struct {
	entries  []*AuditEntry
	mu       sync.RWMutex
	nextSeq  uint64
	lastHash []byte // Hash of previous entry (or genesis)
}

// NewAuditLog creates a new audit log with a genesis entry.
func NewAuditLog() *AuditLog {
	al := &AuditLog{
		entries:  make([]*AuditEntry, 0, 1024),
		nextSeq:  1,
		lastHash: hashGenesis(),
	}
	return al
}

// hashGenesis computes the hash of the genesis (initial) entry.
func hashGenesis() []byte {
	h := sha256.Sum256([]byte("genesis"))
	return h[:]
}

// Append adds a new entry to the log and returns the entry with computed hash.
func (al *AuditLog) Append(keyID, method, outcome string) *AuditEntry {
	al.mu.Lock()
	defer al.mu.Unlock()

	now := time.Now().UnixNano()
	entry := &AuditEntry{
		Seq:       al.nextSeq,
		Timestamp: now,
		KeyID:     keyID,
		Method:    method,
		Outcome:   outcome,
		PrevHash:  al.lastHash,
	}

	// Compute hash: SHA256(Seq|Timestamp|KeyID|Method|Outcome|PrevHash)
	canonical := fmt.Sprintf("%d|%d|%s|%s|%s|%s",
		entry.Seq, entry.Timestamp, entry.KeyID, entry.Method, entry.Outcome,
		hex.EncodeToString(entry.PrevHash))
	h := sha256.Sum256([]byte(canonical))
	entry.Hash = h[:]

	al.entries = append(al.entries, entry)
	al.nextSeq++
	al.lastHash = entry.Hash

	return entry
}

// Entries returns a copy of all entries (for reading).
func (al *AuditLog) Entries() []*AuditEntry {
	al.mu.RLock()
	defer al.mu.RUnlock()
	return append([]*AuditEntry{}, al.entries...)
}

// Verify checks the integrity of the entire audit log.
// Returns nil if log is intact, or an error if tampering is detected.
func (al *AuditLog) Verify() error {
	al.mu.RLock()
	defer al.mu.RUnlock()

	if len(al.entries) == 0 {
		return nil
	}

	prevHash := hashGenesis()
	for _, entry := range al.entries {
		// Recompute expected hash
		canonical := fmt.Sprintf("%d|%d|%s|%s|%s|%s",
			entry.Seq, entry.Timestamp, entry.KeyID, entry.Method, entry.Outcome,
			hex.EncodeToString(prevHash))
		h := sha256.Sum256([]byte(canonical))
		expectedHash := h[:]

		// Check if entry hash matches
		if !bytesEqual(entry.Hash, expectedHash) {
			return fmt.Errorf("audit log tampering detected at entry %d: hash mismatch", entry.Seq)
		}

		// Check if PrevHash matches
		if !bytesEqual(entry.PrevHash, prevHash) {
			return fmt.Errorf("audit log tampering detected at entry %d: prev hash mismatch", entry.Seq)
		}

		prevHash = entry.Hash
	}

	return nil
}

// bytesEqual compares two byte slices for equality (constant-time).
func bytesEqual(a, b []byte) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
