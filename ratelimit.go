package main

import (
	"sync"
	"time"
)

// RateLimiter tracks requests per IP and per account with exponential backoff
type RateLimiter struct {
	ipLimits      map[string]*ipLimit
	accountLimits map[string]*accountLimit
	mu            sync.Mutex
}

type ipLimit struct {
	count   int
	lastAt  time.Time
	backoff int
}

type accountLimit struct {
	failures int
	lastAt   time.Time
	lockedAt time.Time
}

var limiter = &RateLimiter{
	ipLimits:      make(map[string]*ipLimit),
	accountLimits: make(map[string]*accountLimit),
}

// CheckLoginRate returns true if request should be allowed, false if rate limited
func (rl *RateLimiter) CheckLoginRate(ip string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()
	now := time.Now()
	limit, ok := rl.ipLimits[ip]
	if !ok {
		rl.ipLimits[ip] = &ipLimit{count: 1, lastAt: now, backoff: 0}
		return true
	}
	if now.Sub(limit.lastAt) > 15*time.Minute {
		rl.ipLimits[ip] = &ipLimit{count: 1, lastAt: now, backoff: 0}
		return true
	}
	if limit.count > 5 {
		return false
	}
	limit.count++
	limit.lastAt = now
	return true
}

// CheckAccountLockout returns true if account is currently locked (call before auth attempt)
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

// RecordAuthFailure increments failure counter for account (call after failed auth)
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

// RecordSuccess clears rate limit on successful login
func (rl *RateLimiter) RecordSuccess(ip, account string) {
	rl.mu.Lock()
	defer rl.mu.Unlock()
	delete(rl.ipLimits, ip)
	delete(rl.accountLimits, account)
}

// CheckLobbyPasswordRate returns true if request should be allowed
func (rl *RateLimiter) CheckLobbyPasswordRate(ip string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()
	now := time.Now()
	limit, ok := rl.ipLimits[ip]
	if !ok {
		rl.ipLimits[ip] = &ipLimit{count: 1, lastAt: now, backoff: 0}
		return true
	}
	if now.Sub(limit.lastAt) > 10*time.Minute {
		rl.ipLimits[ip] = &ipLimit{count: 1, lastAt: now, backoff: 0}
		return true
	}
	if limit.count > 10 {
		return false
	}
	limit.count++
	return true
}
