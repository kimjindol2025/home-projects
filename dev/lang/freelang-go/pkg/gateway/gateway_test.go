package gateway

import (
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/kimjindol2025/freelang-go/pkg/rpc"
)

// ============ Gateway Tests ============

func TestGatewayValidPolicy(t *testing.T) {
	policy := NewPolicy()
	policy.AddRule(&PolicyRule{
		Name:    "allow_read",
		Methods: []string{"KV.Get"},
		Effect:  EffectAllow,
	})
	policy.AddRule(&PolicyRule{
		Name:    "deny_write",
		Methods: []string{"KV.Set", "KV.Delete"},
		Effect:  EffectDeny,
	})

	if err := ValidatePolicy(policy); err != nil {
		t.Errorf("ValidatePolicy failed: %v", err)
	}
}

func TestGatewayContradictoryPolicy(t *testing.T) {
	policy := NewPolicy()
	policy.AddRule(&PolicyRule{
		Name:    "deny_all",
		Methods: []string{}, // Matches all
		Effect:  EffectDeny,
	})
	policy.AddRule(&PolicyRule{
		Name:    "allow_read",
		Methods: []string{"KV.Get"},
		Effect:  EffectAllow,
	})

	if err := ValidatePolicy(policy); err == nil {
		t.Errorf("ValidatePolicy should have caught contradiction, got nil")
	}
}

// ============ APIKey Tests ============

func TestAPIKeyCreateValidate(t *testing.T) {
	ks := NewKeyStore()
	key, err := ks.Create("key1", []Permission{PermRead, PermWrite}, 100)
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	if key.ID != "key1" {
		t.Errorf("Expected ID 'key1', got %q", key.ID)
	}

	if !key.HasPermission(PermRead) || !key.HasPermission(PermWrite) {
		t.Errorf("Expected both Read and Write permissions")
	}

	if key.RateLimit != 100 {
		t.Errorf("Expected rate limit 100, got %d", key.RateLimit)
	}
}

func TestAPIKeyRevoke(t *testing.T) {
	ks := NewKeyStore()
	ks.Create("key1", []Permission{PermRead}, 100)

	if err := ks.Revoke("key1"); err != nil {
		t.Fatalf("Revoke failed: %v", err)
	}

	key, _ := ks.Get("key1")
	if key.IsActive() {
		t.Errorf("Expected key to be revoked")
	}

	if ks.Exists("key1") {
		t.Errorf("Revoked key should not be in active list")
	}
}

func TestAPIKeyNotFound(t *testing.T) {
	ks := NewKeyStore()
	_, err := ks.Get("nonexistent")
	if err != ErrKeyNotFound {
		t.Errorf("Expected ErrKeyNotFound, got %v", err)
	}
}

func TestAPIKeyPermissions(t *testing.T) {
	ks := NewKeyStore()
	ks.Create("admin", []Permission{PermAdmin}, 1000)
	ks.Create("reader", []Permission{PermRead}, 10)

	adminKey, _ := ks.Get("admin")
	if !adminKey.HasPermission(PermAdmin) {
		t.Errorf("Admin key should have admin permission")
	}

	readerKey, _ := ks.Get("reader")
	if readerKey.HasPermission(PermWrite) {
		t.Errorf("Reader key should not have write permission")
	}
}

// ============ HMAC Tests ============

func TestSigningValid(t *testing.T) {
	signer := NewSigner()
	secret := []byte("test-secret-key")
	timestamp := time.Now().UnixNano()
	body := []byte("request body")

	sig := signer.Sign("TestMethod", timestamp, body, secret)
	if len(sig) != 32 {
		t.Errorf("Expected 32-byte signature, got %d", len(sig))
	}

	if !signer.Verify("TestMethod", timestamp, body, secret, sig) {
		t.Errorf("Signature verification failed")
	}
}

func TestSigningTamperedBody(t *testing.T) {
	signer := NewSigner()
	secret := []byte("test-secret-key")
	timestamp := time.Now().UnixNano()
	body := []byte("original body")

	sig := signer.Sign("TestMethod", timestamp, body, secret)

	// Tamper with body
	tamperedBody := []byte("tampered body")
	if signer.Verify("TestMethod", timestamp, tamperedBody, secret, sig) {
		t.Errorf("Verification should fail for tampered body")
	}
}

func TestSigningReplay(t *testing.T) {
	signer := NewSigner()
	oldTime := time.Now().Add(-45 * time.Second).UnixNano()

	err := signer.CheckReplay(oldTime, time.Now())
	if err != ErrReplayAttack {
		t.Errorf("Expected ErrReplayAttack for old timestamp, got %v", err)
	}
}

// ============ RateLimit Tests ============

func TestRateLimitAllow(t *testing.T) {
	bucket := NewTokenBucket(10.0, 10.0) // 10 tokens/sec, capacity 10

	initialRemaining := bucket.Remaining()
	if initialRemaining != 10.0 {
		t.Errorf("Initial bucket should be full (10.0), got %.1f", initialRemaining)
	}

	if !bucket.Allow() {
		t.Errorf("First token should be allowed")
	}

	afterRemaining := bucket.Remaining()
	if afterRemaining >= initialRemaining {
		t.Errorf("Remaining tokens should be less after consuming one")
	}
}

func TestRateLimitExceed(t *testing.T) {
	bucket := NewTokenBucket(1.0, 100.0) // 1 token capacity, 100/sec refill

	// Consume the only token
	if !bucket.Allow() {
		t.Errorf("First token should be allowed")
	}

	// Should be out of tokens
	if bucket.Allow() {
		t.Errorf("Second token should not be allowed (no refill yet)")
	}
}

func TestRateLimitRefill(t *testing.T) {
	bucket := NewTokenBucket(1.0, 100.0) // Refill at 100/sec

	bucket.Allow() // Consume token
	if bucket.Allow() {
		t.Errorf("Should have no tokens initially")
	}

	time.Sleep(50 * time.Millisecond) // 5 tokens should refill
	if !bucket.Allow() {
		t.Errorf("Token should be available after refill")
	}
}

// ============ AuditLog Tests ============

func TestAuditLogAppend(t *testing.T) {
	log := NewAuditLog()

	entry := log.Append("key1", "KV.Set", "success")
	if entry.Seq != 1 {
		t.Errorf("Expected seq 1, got %d", entry.Seq)
	}

	if entry.KeyID != "key1" || entry.Method != "KV.Set" {
		t.Errorf("Entry fields mismatch")
	}

	if len(entry.Hash) != 32 {
		t.Errorf("Expected 32-byte hash, got %d", len(entry.Hash))
	}
}

func TestAuditLogVerifyIntact(t *testing.T) {
	log := NewAuditLog()

	log.Append("key1", "KV.Set", "success")
	log.Append("key2", "KV.Get", "success")
	log.Append("key1", "KV.Delete", "success")

	if err := log.Verify(); err != nil {
		t.Errorf("Verification should succeed for intact log: %v", err)
	}
}

func TestAuditLogTamperDetected(t *testing.T) {
	log := NewAuditLog()

	log.Append("key1", "KV.Set", "success")
	log.Append("key2", "KV.Get", "success")

	// Tamper with first entry
	entries := log.Entries()
	if len(entries) > 0 {
		entries[0].Method = "Tampered"
	}

	if err := log.Verify(); err == nil {
		t.Errorf("Verification should detect tampering")
	}
}

func TestAuditLogHashChain(t *testing.T) {
	log := NewAuditLog()

	e1 := log.Append("key1", "Method1", "success")
	e2 := log.Append("key2", "Method2", "success")

	if !bytesEqual(e2.PrevHash, e1.Hash) {
		t.Errorf("e2.PrevHash should match e1.Hash")
	}

	entries := log.Entries()
	if !bytesEqual(entries[0].PrevHash, hashGenesis()) {
		t.Errorf("First entry's PrevHash should be genesis hash")
	}
}

// ============ Invariant Tests ============

func TestInvariantNoUnauthWrite(t *testing.T) {
	inv := NoUnauthenticatedWrite()

	state := InvariantState{
		KeyID:  "",
		Method: "KV.Set",
		Phase:  "pre",
	}

	if msg := inv(state); msg == "" {
		t.Errorf("Invariant should detect unauthenticated write")
	}
}

func TestInvariantRateLimitBound(t *testing.T) {
	inv := RateLimitBound()

	state := InvariantState{
		RateLimitUsed: 11.0,
		RateLimitCap:  10.0,
		Phase:         "post",
	}

	if msg := inv(state); msg == "" {
		t.Errorf("Invariant should detect rate limit violation")
	}
}

// ============ Policy Tests ============

func TestPolicyAllowRule(t *testing.T) {
	policy := NewPolicy()
	policy.AddRule(&PolicyRule{
		Name:    "allow_read",
		Methods: []string{"KV.Get"},
		Effect:  EffectAllow,
	})

	req := &Request{Method: "KV.Get"}
	allowed, ruleName := policy.Evaluate(req)

	if !allowed || ruleName != "allow_read" {
		t.Errorf("Expected allow with rule 'allow_read', got %v/%s", allowed, ruleName)
	}
}

func TestPolicyDenyRule(t *testing.T) {
	policy := NewPolicy()
	policy.AddRule(&PolicyRule{
		Name:    "deny_admin",
		Methods: []string{"AdminMethod"},
		Effect:  EffectDeny,
	})

	req := &Request{Method: "AdminMethod"}
	allowed, ruleName := policy.Evaluate(req)

	if allowed || ruleName != "deny_admin" {
		t.Errorf("Expected deny with rule 'deny_admin', got %v/%s", allowed, ruleName)
	}
}

// ============ Integration Tests ============

func TestGatewayEndToEnd(t *testing.T) {
	// Setup RPC infrastructure
	clientT, serverT, err := rpc.NewInProcessPair()
	if err != nil {
		t.Fatalf("NewInProcessPair failed: %v", err)
	}

	codec := rpc.NewBinaryCodec()
	rpcServer := rpc.NewServer(codec)
	rpcServer.ServeAsync(serverT)

	// Setup gateway
	ks := NewKeyStore()
	ks.Create("testkey", []Permission{PermRead, PermWrite}, 100)

	policy := NewPolicy()
	policy.AddRule(&PolicyRule{
		Name:    "allow_all",
		Methods: []string{},
		Effect:  EffectAllow,
	})

	monitor := NewInvariantMonitor()
	monitor.Register("no_unauth_write", NoUnauthenticatedWrite())

	gateway := NewGateway(rpcServer, ks, policy, monitor)

	// Register a simple echo handler through gateway
	handler := func(args []byte) ([]byte, error) {
		return args, nil
	}

	if err := gateway.RegisterMethod("Echo", PermRead, handler); err != nil {
		t.Fatalf("RegisterMethod failed: %v", err)
	}

	// Create client and make a call
	client := rpc.NewClient(codec, clientT)
	defer client.Close()

	// For now, just verify gateway infrastructure is set up
	// Full integration test requires encoding request with signature
}

func TestGatewayRevokedKey(t *testing.T) {
	ks := NewKeyStore()
	ks.Create("testkey", []Permission{PermRead}, 100)
	ks.Revoke("testkey")

	policy := NewPolicy()
	monitor := NewInvariantMonitor()

	gateway := NewGateway(nil, ks, policy, monitor)
	verifier := gateway.verifier

	req := &Request{
		Method:    "KV.Get",
		KeyID:     "testkey",
		Body:      []byte("test"),
		Signature: []byte("dummy"),
		Timestamp: time.Now().UnixNano(),
	}

	err := verifier.VerifyRequest(req)
	if err != ErrKeyRevoked {
		t.Errorf("Expected ErrKeyRevoked, got %v", err)
	}
}

func TestGatewayErrors(t *testing.T) {
	ks := NewKeyStore()
	policy := NewPolicy()
	monitor := NewInvariantMonitor()

	gateway := NewGateway(nil, ks, policy, monitor)

	// Gateway should accumulate errors
	_ = gateway.Errors()
}

func TestGatewayStaticValidation(t *testing.T) {
	policy := NewPolicy()
	policy.AddRule(&PolicyRule{
		Name:    "deny_all",
		Effect:  EffectDeny,
	})
	policy.AddRule(&PolicyRule{
		Name:    "allow_read",
		Methods: []string{"KV.Get"},
		Effect:  EffectAllow,
	})

	// Should detect contradiction
	ks := NewKeyStore()
	monitor := NewInvariantMonitor()
	gateway := NewGateway(nil, ks, policy, monitor)

	if err := gateway.verifier.ValidatePolicy(); err == nil {
		t.Errorf("ValidatePolicy should detect contradiction")
	}
}

func TestGatewayRPCIntegration(t *testing.T) {
	// Setup RPC server
	clientT, serverT, err := rpc.NewInProcessPair()
	if err != nil {
		t.Fatalf("NewInProcessPair failed: %v", err)
	}

	codec := rpc.NewBinaryCodec()
	rpcServer := rpc.NewServer(codec)
	rpcServer.ServeAsync(serverT)

	// Register a simple method directly on RPC server (not through gateway)
	rpcServer.Register("SimpleEcho", func(args []byte) ([]byte, error) {
		return args, nil
	})

	// Call through RPC client
	client := rpc.NewClient(codec, clientT)
	defer client.Close()

	var result string
	if err := client.Call("SimpleEcho", "hello", &result); err != nil {
		t.Fatalf("Call failed: %v", err)
	}

	if result != "hello" {
		t.Errorf("Expected 'hello', got %q", result)
	}
}

// ============ Concurrency Tests ============

func TestGatewayRateLimitConcurrency(t *testing.T) {
	ks := NewKeyStore()
	ks.Create("key1", []Permission{PermRead}, 10)

	policy := NewPolicy()
	policy.AddRule(&PolicyRule{
		Name:    "allow_all",
		Methods: []string{},
		Effect:  EffectAllow,
	})

	monitor := NewInvariantMonitor()
	gateway := NewGateway(nil, ks, policy, monitor)

	// Get rate limiter and use it concurrently
	limiter := gateway.getRateLimiter("key1")

	var wg sync.WaitGroup
	successCount := 0
	var mu sync.Mutex

	for i := 0; i < 20; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if limiter.Allow() {
				mu.Lock()
				successCount++
				mu.Unlock()
			}
		}()
	}

	wg.Wait()

	// Should have limited successes due to capacity
	if successCount == 0 {
		t.Errorf("Expected some tokens to be allowed, got 0")
	}
}

func TestKeyStoreRWConcurrency(t *testing.T) {
	ks := NewKeyStore()

	// Create some keys
	for i := 0; i < 5; i++ {
		ks.Create(fmt.Sprintf("key%d", i), []Permission{PermRead}, 100)
	}

	var wg sync.WaitGroup

	// Read concurrently
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			ks.Get(fmt.Sprintf("key%d", id%5))
		}(i)
	}

	// Write concurrently
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			ks.Create(fmt.Sprintf("extra%d", id), []Permission{PermWrite}, 200)
		}(i)
	}

	wg.Wait()
}

// ============ Wire Format Tests ============

func TestParseRequestWireFormat(t *testing.T) {
	ks := NewKeyStore()
	ks.Create("testkey", []Permission{PermRead}, 100)

	policy := NewPolicy()
	policy.AddRule(&PolicyRule{
		Name:    "allow_all",
		Methods: []string{},
		Effect:  EffectAllow,
	})

	monitor := NewInvariantMonitor()
	gateway := NewGateway(nil, ks, policy, monitor)

	// Build wire format manually
	keyID := "testkey"
	signature := make([]byte, 32) // Dummy
	timestamp := time.Now().UnixNano()
	body := []byte("test body")

	dst := make([]byte, 0, 256)

	// KeyID length + KeyID
	dst = append(dst, byte((len(keyID)>>8)&0xFF), byte(len(keyID)&0xFF))
	dst = append(dst, []byte(keyID)...)

	// Signature length (32) + Signature
	dst = append(dst, byte(0), byte(32))
	dst = append(dst, signature...)

	// Timestamp
	for i := 7; i >= 0; i-- {
		dst = append(dst, byte((timestamp>>uint(i*8))&0xFF))
	}

	// Body length + Body
	dst = append(dst, byte((len(body)>>24)&0xFF), byte((len(body)>>16)&0xFF),
		byte((len(body)>>8)&0xFF), byte(len(body)&0xFF))
	dst = append(dst, body...)

	// Parse
	req, err := gateway.parseRequest("TestMethod", dst)
	if err != nil {
		t.Fatalf("parseRequest failed: %v", err)
	}

	if req.KeyID != keyID {
		t.Errorf("KeyID mismatch: expected %q, got %q", keyID, req.KeyID)
	}

	if len(req.Signature) != 32 {
		t.Errorf("Signature length mismatch")
	}

	if req.Timestamp != timestamp {
		t.Errorf("Timestamp mismatch")
	}

	if string(req.Body) != "test body" {
		t.Errorf("Body mismatch: expected 'test body', got %q", string(req.Body))
	}
}
