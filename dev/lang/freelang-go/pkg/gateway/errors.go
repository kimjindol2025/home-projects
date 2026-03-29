package gateway

import "errors"

// Sentinel errors for gateway operations.
var (
	ErrUnauthorized     = errors.New("unauthorized: missing or invalid API key")
	ErrInvalidSignature = errors.New("invalid signature: HMAC verification failed")
	ErrReplayAttack     = errors.New("replay attack: timestamp outside valid window")
	ErrRateLimitExceeded = errors.New("rate limit exceeded: quota exhausted for this key")
	ErrPermissionDenied = errors.New("permission denied: insufficient privileges for this method")
	ErrKeyRevoked       = errors.New("key revoked: this API key has been revoked")
	ErrKeyNotFound      = errors.New("key not found: API key ID does not exist")
	ErrDuplicateKey     = errors.New("duplicate key: API key ID already registered")
	ErrInvalidPolicy    = errors.New("invalid policy: contradictory rules detected")
	ErrInvariantViolation = errors.New("invariant violated: security constraint broken")
)
