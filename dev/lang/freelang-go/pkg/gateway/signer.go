package gateway

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"
)

// ReplayWindow defines the maximum age of a timestamp before considering it a replay attack.
const ReplayWindow = 30 * time.Second

// Signer handles HMAC-SHA256 signing and verification of requests.
type Signer struct{}

// NewSigner creates a new Signer.
func NewSigner() *Signer {
	return &Signer{}
}

// Sign computes HMAC-SHA256 signature for a request.
// Canonical format: method|timestamp|hex(body)
func (s *Signer) Sign(method string, timestamp int64, body []byte, secret []byte) []byte {
	canonical := fmt.Sprintf("%s|%d|%s", method, timestamp, hex.EncodeToString(body))
	h := hmac.New(sha256.New, secret)
	h.Write([]byte(canonical))
	return h.Sum(nil)
}

// Verify checks if the provided signature matches the expected signature.
// Uses constant-time comparison to prevent timing attacks.
func (s *Signer) Verify(method string, timestamp int64, body []byte, secret []byte, signature []byte) bool {
	expected := s.Sign(method, timestamp, body, secret)
	return hmac.Equal(expected, signature)
}

// CheckReplay returns an error if the timestamp is outside the replay window.
func (s *Signer) CheckReplay(timestamp int64, now time.Time) error {
	msgTime := time.Unix(0, timestamp)
	age := now.Sub(msgTime).Abs()
	if age > ReplayWindow {
		return ErrReplayAttack
	}
	return nil
}
