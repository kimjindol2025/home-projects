package gateway

import (
	"encoding/binary"
	"fmt"
	"sync"

	"github.com/kimjindol2025/freelang-go/pkg/rpc"
)

// Gateway wraps an RPC server with a full security middleware stack.
type Gateway struct {
	server       *rpc.Server
	keyStore     *KeyStore
	limiter      map[string]*TokenBucket      // Per-key rate limiters
	signer       *Signer
	policy       *Policy
	verifier     *Verifier
	auditLog     *AuditLog
	methodPerms  map[string]Permission        // Method name -> required permission
	requestCount int64                        // Total requests processed
	mu           sync.RWMutex
	errors       []string
}

// NewGateway creates a new security gateway wrapping an RPC server.
func NewGateway(srv *rpc.Server, ks *KeyStore, policy *Policy, monitor *InvariantMonitor) *Gateway {
	signer := NewSigner()
	verifier := NewVerifier(ks, signer, policy, monitor)

	g := &Gateway{
		server:      srv,
		keyStore:    ks,
		limiter:     make(map[string]*TokenBucket),
		signer:      signer,
		policy:      policy,
		verifier:    verifier,
		auditLog:    NewAuditLog(),
		methodPerms: make(map[string]Permission),
		errors:      make([]string, 0),
	}

	return g
}

// RegisterMethod registers an RPC method with its required permission.
// This method should be called before any requests are processed.
func (g *Gateway) RegisterMethod(method string, perm Permission, handler rpc.HandlerFunc) error {
	g.mu.Lock()
	g.methodPerms[method] = perm
	g.mu.Unlock()

	// Wrap the handler with security middleware
	wrappedHandler := g.wrapHandler(method, perm, handler)

	// Register wrapped handler on the underlying RPC server
	return g.server.Register(method, wrappedHandler)
}

// wrapHandler creates a security middleware around an RPC handler.
// Middleware chain: authenticate → verify_sig → rate_limit → check_policy →
//                  pre_conditions → invariant_pre → handler →
//                  post_conditions → invariant_post → audit_log
func (g *Gateway) wrapHandler(method string, perm Permission, handler rpc.HandlerFunc) rpc.HandlerFunc {
	return func(args []byte) ([]byte, error) {
		// Parse request from wire format
		req, err := g.parseRequest(method, args)
		if err != nil {
			g.auditLog.Append("unknown", method, "authentication_failed")
			return nil, fmt.Errorf("parse request: %w", err)
		}

		// Verify request (pre-conditions)
		if err := g.verifier.VerifyRequest(req); err != nil {
			g.auditLog.Append(req.KeyID, method, fmt.Sprintf("denied: %v", err))
			return nil, err
		}

		// Get or create rate limiter for this key
		rateLimiter := g.getRateLimiter(req.KeyID)

		// Check rate limit
		if !rateLimiter.Allow() {
			g.auditLog.Append(req.KeyID, method, "rate_limit_exceeded")
			return nil, ErrRateLimitExceeded
		}

		// Increment request counter
		g.mu.Lock()
		g.requestCount++
		requestCount := g.requestCount
		g.mu.Unlock()

		// Call handler with original args (the body from the request)
		respBody, handlerErr := handler(req.Body)

		// Verify response (post-conditions)
		resp := &Response{
			Body:  respBody,
			Error: "",
		}
		if handlerErr != nil {
			resp.Error = handlerErr.Error()
		}

		rateLimitUsed := 1.0 // One token was consumed
		rateLimitCap := float64(g.getKeyRateLimit(req.KeyID))

		if err := g.verifier.VerifyResponse(req, resp, float64(requestCount), rateLimitUsed, rateLimitCap); err != nil {
			g.auditLog.Append(req.KeyID, method, fmt.Sprintf("post_condition_failed: %v", err))
			return nil, fmt.Errorf("post-condition: %w", err)
		}

		// Audit log with outcome
		outcome := "success"
		if handlerErr != nil {
			outcome = "handler_error"
		}
		g.auditLog.Append(req.KeyID, method, outcome)

		// Encode response
		encoded, err := g.encodeResponse(resp)
		if err != nil {
			return nil, fmt.Errorf("encode response: %w", err)
		}

		return encoded, handlerErr
	}
}

// parseRequest decodes a signed request from wire format.
// Wire format: [2]KeyID_len + [N]KeyID + [2]Sig_len(32) + [32]Signature + [8]Timestamp + [4]Body_len + [M]Body
func (g *Gateway) parseRequest(method string, args []byte) (*Request, error) {
	if len(args) < 2 {
		return nil, fmt.Errorf("args too short for KeyID length")
	}

	offset := 0

	// Read KeyID length
	keyIDLen := int(binary.BigEndian.Uint16(args[offset : offset+2]))
	offset += 2

	// Read KeyID
	if offset+keyIDLen > len(args) {
		return nil, fmt.Errorf("args truncated while reading KeyID")
	}
	keyID := string(args[offset : offset+keyIDLen])
	offset += keyIDLen

	// Read Signature length
	if offset+2 > len(args) {
		return nil, fmt.Errorf("args too short for Signature length")
	}
	sigLen := int(binary.BigEndian.Uint16(args[offset : offset+2]))
	offset += 2

	// Read Signature
	if offset+sigLen > len(args) {
		return nil, fmt.Errorf("args truncated while reading Signature")
	}
	signature := make([]byte, sigLen)
	copy(signature, args[offset:offset+sigLen])
	offset += sigLen

	// Read Timestamp
	if offset+8 > len(args) {
		return nil, fmt.Errorf("args too short for Timestamp")
	}
	timestamp := int64(binary.BigEndian.Uint64(args[offset : offset+8]))
	offset += 8

	// Read Body length
	if offset+4 > len(args) {
		return nil, fmt.Errorf("args too short for Body length")
	}
	bodyLen := int(binary.BigEndian.Uint32(args[offset : offset+4]))
	offset += 4

	// Read Body
	if offset+bodyLen > len(args) {
		return nil, fmt.Errorf("args truncated while reading Body")
	}
	body := args[offset : offset+bodyLen]

	return &Request{
		Method:    method,
		Body:      body,
		KeyID:     keyID,
		Signature: signature,
		Timestamp: timestamp,
	}, nil
}

// encodeResponse encodes a response into wire format.
// For now, returns just the body (actual RPC framework encodes further).
func (g *Gateway) encodeResponse(resp *Response) ([]byte, error) {
	// Simple format: [4]Body_len + [M]Body + [4]Error_len + [N]Error
	dst := make([]byte, 0, 256)

	// Encode Body
	dst = binary.BigEndian.AppendUint32(dst, uint32(len(resp.Body)))
	dst = append(dst, resp.Body...)

	// Encode Error
	dst = binary.BigEndian.AppendUint32(dst, uint32(len(resp.Error)))
	dst = append(dst, []byte(resp.Error)...)

	return dst, nil
}

// getRateLimiter gets or creates a rate limiter for a key.
func (g *Gateway) getRateLimiter(keyID string) *TokenBucket {
	g.mu.Lock()
	defer g.mu.Unlock()

	if limiter, exists := g.limiter[keyID]; exists {
		return limiter
	}

	// Create new limiter with key's rate limit
	rateLimit := g.getKeyRateLimit(keyID)
	limiter := NewTokenBucket(float64(rateLimit), float64(rateLimit))
	g.limiter[keyID] = limiter
	return limiter
}

// getKeyRateLimit returns the rate limit for a key (tokens/sec).
func (g *Gateway) getKeyRateLimit(keyID string) int {
	key, err := g.keyStore.Get(keyID)
	if err != nil {
		return 10 // Default
	}
	return key.RateLimit
}

// AuditLog returns the audit log for inspection.
func (g *Gateway) AuditLog() *AuditLog {
	return g.auditLog
}

// Errors returns accumulated errors.
func (g *Gateway) Errors() []string {
	g.mu.Lock()
	defer g.mu.Unlock()
	return append([]string{}, g.errors...)
}
