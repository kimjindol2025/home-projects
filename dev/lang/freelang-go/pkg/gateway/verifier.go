package gateway

import (
	"fmt"
	"time"
)

// Verifier performs formal verification of gateway operations.
type Verifier struct {
	keyStore   *KeyStore
	signer     *Signer
	policy     *Policy
	monitor    *InvariantMonitor
}

// NewVerifier creates a new verifier.
func NewVerifier(ks *KeyStore, s *Signer, p *Policy, m *InvariantMonitor) *Verifier {
	return &Verifier{
		keyStore: ks,
		signer:   s,
		policy:   p,
		monitor:  m,
	}
}

// VerifyRequest performs all pre-condition checks before handler execution.
// Returns nil if request passes all checks, or an error with the reason for failure.
func (v *Verifier) VerifyRequest(req *Request) error {
	// Pre-condition 1: Key must exist and be active
	if req.KeyID == "" {
		return ErrUnauthorized
	}

	key, err := v.keyStore.Get(req.KeyID)
	if err != nil {
		return ErrUnauthorized
	}

	if !key.IsActive() {
		return ErrKeyRevoked
	}

	// Pre-condition 2: Signature must be valid
	if !v.signer.Verify(req.Method, req.Timestamp, req.Body, key.Secret, req.Signature) {
		return ErrInvalidSignature
	}

	// Pre-condition 3: Timestamp must be within replay window
	if err := v.signer.CheckReplay(req.Timestamp, time.Now()); err != nil {
		return err
	}

	// Pre-condition 4: Policy must allow the request
	allowed, _ := v.policy.Evaluate(req)
	if !allowed {
		return ErrPermissionDenied
	}

	// Pre-condition 5: Invariants at pre-phase
	preState := InvariantState{
		Method: req.Method,
		KeyID:  req.KeyID,
		Phase:  "pre",
	}
	if err := v.monitor.Check(preState); err != nil {
		return fmt.Errorf("pre-condition failed: %w", err)
	}

	return nil
}

// VerifyResponse performs post-condition checks after handler execution.
// Returns nil if response is valid, or an error with the reason for failure.
func (v *Verifier) VerifyResponse(req *Request, resp *Response, requestCount float64, rateLimitUsed float64, rateLimitCap float64) error {
	// Post-condition 1: Response must not be nil
	if resp == nil {
		return fmt.Errorf("response is nil after handler execution")
	}

	// Post-condition 2: Invariants at post-phase
	postState := InvariantState{
		Method:        req.Method,
		KeyID:         req.KeyID,
		Phase:         "post",
		RequestCount:  int64(requestCount),
		RateLimitUsed: rateLimitUsed,
		RateLimitCap:  rateLimitCap,
		Response:      resp,
	}
	if err := v.monitor.Check(postState); err != nil {
		return fmt.Errorf("post-condition failed: %w", err)
	}

	return nil
}

// ValidatePolicy checks the policy for contradictions at startup.
func (v *Verifier) ValidatePolicy() error {
	return ValidatePolicy(v.policy)
}
