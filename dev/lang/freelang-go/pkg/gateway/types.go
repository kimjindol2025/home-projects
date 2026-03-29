package gateway

import "time"

// Permission represents an access level.
type Permission string

const (
	PermRead  Permission = "read"
	PermWrite Permission = "write"
	PermAdmin Permission = "admin"
)

// APIKey represents a registered API key with metadata.
type APIKey struct {
	ID          string        // Unique key identifier
	Secret      []byte        // 32-byte HMAC secret
	Permissions []Permission  // List of granted permissions
	RateLimit   int           // Tokens per second
	CreatedAt   time.Time
	RevokedAt   *time.Time // nil = active
}

// IsActive returns true if the key is not revoked.
func (k *APIKey) IsActive() bool {
	return k.RevokedAt == nil
}

// HasPermission checks if the key has the given permission.
func (k *APIKey) HasPermission(p Permission) bool {
	for _, perm := range k.Permissions {
		if perm == p {
			return true
		}
	}
	return false
}

// Request represents a normalized gateway request.
type Request struct {
	Method    string // RPC method name
	Body      []byte // Request body (encoded args)
	KeyID     string // API key ID
	Signature []byte // HMAC-SHA256 (32 bytes)
	Timestamp int64  // Unix nanoseconds (for replay detection)
}

// Response represents a gateway response.
type Response struct {
	Body  []byte // Encoded result or error
	Error string // Non-empty if error occurred
}

// AuditEntry represents a single audit log entry in the SHA256 chain.
type AuditEntry struct {
	Seq       uint64    // Sequential entry number
	Timestamp int64     // Unix nanoseconds
	KeyID     string    // API key used
	Method    string    // RPC method called
	Outcome   string    // "success", "denied", "error", etc.
	PrevHash  []byte    // Hash of previous entry (32 bytes)
	Hash      []byte    // SHA256 of this entry (32 bytes)
}

// PolicyEffect represents allow/deny decision.
type PolicyEffect int

const (
	EffectAllow PolicyEffect = iota
	EffectDeny
)
