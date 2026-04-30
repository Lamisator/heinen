package main

import (
	"sync"
	"time"
)

// RateLimiter tracks requests per IP and per account with exponential backoff.
// Login and lobby-password attempts use separate maps to prevent cross-interference.
type RateLimiter struct {
	ipLoginLimits  map[string]*ipLimit
	ipLobbyLimits  map[string]*ipLimit
	accountLimits  map[string]*accountLimit
	mu             sync.Mutex
}

type ipLimit struct {
	count  int
	lastAt time.Time
}

type accountLimit struct {
	failures int
	lastAt   time.Time
	lockedAt time.Time
}

var limiter = &RateLimiter{
	ipLoginLimits: make(map[string]*ipLimit),
	ipLobbyLimits: make(map[string]*ipLimit),
	accountLimits: make(map[string]*accountLimit),
}

// CheckLoginRate returns true if the login attempt should be allowed.
func (rl *RateLimiter) CheckLoginRate(ip string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()
	now := time.Now()
	limit, ok := rl.ipLoginLimits[ip]
	if !ok {
		rl.ipLoginLimits[ip] = &ipLimit{count: 1, lastAt: now}
		return true
	}
	if now.Sub(limit.lastAt) > 15*time.Minute {
		rl.ipLoginLimits[ip] = &ipLimit{count: 1, lastAt: now}
		return true
	}
	if limit.count >= 5 {
		return false
	}
	limit.count++
	limit.lastAt = now
	return true
}

// CheckAccountLockout returns true if the account is currently locked.
func (rl *RateLimiter) CheckAccountLockout(account string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()
	now := time.Now()
	acct, ok := rl.accountLimits[account]
	if !ok {
		rl.accountLimits[account] = &accountLimit{failures: 0, lastAt: now}
		return false
	}
	if !acct.lockedAt.IsZero() {
		if now.Sub(acct.lockedAt) > 15*time.Minute {
			rl.accountLimits[account] = &accountLimit{failures: 0, lastAt: now}
			return false
		}
		return true
	}
	return false
}

// RecordAuthFailure increments the failure counter for an account.
func (rl *RateLimiter) RecordAuthFailure(account string) {
	rl.mu.Lock()
	defer rl.mu.Unlock()
	now := time.Now()
	acct, ok := rl.accountLimits[account]
	if !ok {
		rl.accountLimits[account] = &accountLimit{failures: 1, lastAt: now}
		return
	}
	if now.Sub(acct.lastAt) > 15*time.Minute {
		rl.accountLimits[account] = &accountLimit{failures: 1, lastAt: now}
		return
	}
	acct.failures++
	acct.lastAt = now
	if acct.failures >= 5 {
		acct.lockedAt = now
	}
}

// RecordSuccess clears rate-limit state on successful login.
func (rl *RateLimiter) RecordSuccess(ip, account string) {
	rl.mu.Lock()
	defer rl.mu.Unlock()
	delete(rl.ipLoginLimits, ip)
	delete(rl.accountLimits, account)
}

// CheckLobbyPasswordRate returns true if the lobby-password attempt should be allowed.
func (rl *RateLimiter) CheckLobbyPasswordRate(ip string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()
	now := time.Now()
	limit, ok := rl.ipLobbyLimits[ip]
	if !ok {
		rl.ipLobbyLimits[ip] = &ipLimit{count: 1, lastAt: now}
		return true
	}
	if now.Sub(limit.lastAt) > 10*time.Minute {
		rl.ipLobbyLimits[ip] = &ipLimit{count: 1, lastAt: now}
		return true
	}
	if limit.count >= 10 {
		return false
	}
	limit.count++
	limit.lastAt = now
	return true
}

// CleanExpired removes stale entries older than the given threshold to bound memory usage.
func (rl *RateLimiter) CleanExpired() {
	rl.mu.Lock()
	defer rl.mu.Unlock()
	threshold := 30 * time.Minute
	now := time.Now()
	for ip, l := range rl.ipLoginLimits {
		if now.Sub(l.lastAt) > threshold {
			delete(rl.ipLoginLimits, ip)
		}
	}
	for ip, l := range rl.ipLobbyLimits {
		if now.Sub(l.lastAt) > threshold {
			delete(rl.ipLobbyLimits, ip)
		}
	}
	for acct, a := range rl.accountLimits {
		cutoff := a.lastAt
		if !a.lockedAt.IsZero() && a.lockedAt.After(cutoff) {
			cutoff = a.lockedAt
		}
		if now.Sub(cutoff) > threshold {
			delete(rl.accountLimits, acct)
		}
	}
}
