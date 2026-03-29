package gateway

import (
	"fmt"
)

// Predicate is a function that evaluates a condition on a request.
type Predicate func(req *Request) bool

// PolicyRule represents a single access control rule.
type PolicyRule struct {
	Name      string
	Methods   []string        // Which methods this rule applies to
	Effect    PolicyEffect    // Allow or Deny
	Condition Predicate       // Condition function (optional)
}

// Policy is a collection of policy rules that are evaluated in order.
type Policy struct {
	rules []*PolicyRule
}

// NewPolicy creates an empty policy.
func NewPolicy() *Policy {
	return &Policy{
		rules: make([]*PolicyRule, 0),
	}
}

// AddRule appends a new rule to the policy.
func (p *Policy) AddRule(rule *PolicyRule) {
	p.rules = append(p.rules, rule)
}

// Evaluate checks if a request is allowed by the policy.
// Rules are checked in order; first match wins.
// Returns (allowed, ruleNameThatMatched).
func (p *Policy) Evaluate(req *Request) (bool, string) {
	for _, rule := range p.rules {
		// Check if rule applies to this method
		if !methodMatches(rule.Methods, req.Method) {
			continue
		}

		// Check condition if present
		if rule.Condition != nil && !rule.Condition(req) {
			continue
		}

		// Rule matches; return its effect
		if rule.Effect == EffectAllow {
			return true, rule.Name
		}
		return false, rule.Name
	}

	// No rule matched; default deny
	return false, ""
}

// methodMatches checks if the method matches any of the rule methods.
// Empty methods list matches all methods.
func methodMatches(ruleMethods []string, method string) bool {
	if len(ruleMethods) == 0 {
		return true // Empty methods = match all
	}
	for _, m := range ruleMethods {
		if m == method || m == "*" {
			return true
		}
	}
	return false
}

// ValidatePolicy checks for contradictions (e.g., Deny-all before Allow-specific).
// Returns an error if contradictions are found.
func ValidatePolicy(p *Policy) error {
	if len(p.rules) == 0 {
		return nil
	}

	// Check if a Deny-all rule appears before an Allow-specific rule
	for i, rule := range p.rules {
		if rule.Effect == EffectDeny && isDenyAll(rule) {
			// Check if any later rule is Allow
			for j := i + 1; j < len(p.rules); j++ {
				laterRule := p.rules[j]
				if laterRule.Effect == EffectAllow && hasOverlap(rule, laterRule) {
					return fmt.Errorf("%w: rule %q (Deny-all) at position %d blocks rule %q (Allow) at position %d",
						ErrInvalidPolicy, rule.Name, i, laterRule.Name, j)
				}
			}
		}
	}

	return nil
}

// isDenyAll checks if a rule is a deny-all (no methods, no condition).
func isDenyAll(rule *PolicyRule) bool {
	return rule.Effect == EffectDeny && len(rule.Methods) == 0 && rule.Condition == nil
}

// hasOverlap checks if two rules have overlapping methods.
func hasOverlap(r1, r2 *PolicyRule) bool {
	if len(r1.Methods) == 0 || len(r2.Methods) == 0 {
		return true // One or both match all methods
	}

	for _, m1 := range r1.Methods {
		for _, m2 := range r2.Methods {
			if m1 == m2 || m1 == "*" || m2 == "*" {
				return true
			}
		}
	}
	return false
}

// Built-in predicates for common checks

// RequirePermission returns a predicate that checks for a specific permission.
// This requires the key to have the permission.
func RequirePermission(keyStore *KeyStore, perm Permission) Predicate {
	return func(req *Request) bool {
		key, err := keyStore.Get(req.KeyID)
		if err != nil {
			return false
		}
		return key.HasPermission(perm)
	}
}

// AllowIfActive returns a predicate that checks if the key is active (not revoked).
func AllowIfActive(keyStore *KeyStore) Predicate {
	return func(req *Request) bool {
		key, err := keyStore.Get(req.KeyID)
		if err != nil {
			return false
		}
		return key.IsActive()
	}
}
