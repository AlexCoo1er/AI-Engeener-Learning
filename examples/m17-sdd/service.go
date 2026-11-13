// Package sdd models the lab's demo identities and in-memory access requests.
package sdd

import (
	"fmt"
	"strings"
	"sync"
	"time"
)

type Error struct {
	Status  int
	Message string
}

func (e *Error) Error() string                 { return e.Message }
func failure(status int, message string) error { return &Error{status, message} }

type Request struct {
	ID           string     `json:"id"`
	User         string     `json:"user"`
	Resource     string     `json:"resource"`
	Hours        int        `json:"hours"`
	Status       string     `json:"status"`
	ApprovedAt   *time.Time `json:"approved_at,omitempty"`
	ExpiresAt    *time.Time `json:"expires_at,omitempty"`
	RevokedAt    *time.Time `json:"revoked_at,omitempty"`
	RevokeReason string     `json:"revoke_reason,omitempty"`
}

// Service serializes transitions so one conflicting decision wins.
// ponytail: one mutex and a linear grant scan; use a database if this stops being a lab.
type Service struct {
	mu        sync.Mutex
	reference bool
	clock     func() time.Time
	next      int
	requests  map[string]Request
}

func NewService(reference bool, clock func() time.Time) *Service {
	return &Service{reference: reference, clock: clock, requests: make(map[string]Request)}
}

// Upgrade changes policy in memory while retaining already-approved legacy grants.
func (s *Service) Upgrade() { s.mu.Lock(); defer s.mu.Unlock(); s.reference = true }
func Identity(user string) error {
	switch user {
	case "alice", "bob", "carol":
		return nil
	}
	return failure(401, "unknown or missing demo identity")
}
func resourceValid(resource string) bool { return resource == "handbook" || resource == "runbook" }
func copyRequest(r Request) Request {
	if r.ApprovedAt != nil {
		v := *r.ApprovedAt
		r.ApprovedAt = &v
	}
	if r.ExpiresAt != nil {
		v := *r.ExpiresAt
		r.ExpiresAt = &v
	}
	if r.RevokedAt != nil {
		v := *r.RevokedAt
		r.RevokedAt = &v
	}
	return r
}
func (s *Service) Create(user, resource string, hours int) (Request, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if e := Identity(user); e != nil {
		return Request{}, e
	}
	if !resourceValid(resource) {
		return Request{}, failure(400, "unknown resource")
	}
	if (s.reference && (hours < 1 || hours > 24)) || (!s.reference && hours != 0) {
		return Request{}, failure(400, "invalid hours for mode")
	}
	s.next++
	r := Request{ID: fmt.Sprint(s.next), User: user, Resource: resource, Hours: hours, Status: "pending"}
	s.requests[r.ID] = r
	return r, nil
}
func (s *Service) Get(user, id string) (Request, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if e := Identity(user); e != nil {
		return Request{}, e
	}
	r, ok := s.requests[id]
	if !ok {
		return Request{}, failure(404, "request not found")
	}
	if user != r.User && user != "carol" {
		return Request{}, failure(403, "foreign request")
	}
	return copyRequest(r), nil
}
func (s *Service) Decide(user, id, decision string) (Request, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if e := Identity(user); e != nil {
		return Request{}, e
	}
	r, ok := s.requests[id]
	if !ok {
		return Request{}, failure(404, "request not found")
	}
	if user != "carol" {
		return Request{}, failure(403, "approver required")
	}
	if user == r.User {
		return Request{}, failure(403, "self decision forbidden")
	}
	if decision != "approved" && decision != "rejected" {
		return Request{}, failure(400, "invalid decision")
	}
	if r.Status == decision {
		return copyRequest(r), nil
	}
	if r.Status != "pending" {
		return Request{}, failure(409, "decision already recorded")
	}
	if decision == "approved" {
		if s.reference && r.Hours == 0 {
			return Request{}, failure(400, "legacy pending request needs a new timed request")
		}
		now := s.clock().UTC()
		r.ApprovedAt = &now
		if r.Hours > 0 {
			expiry := now.Add(time.Duration(r.Hours) * time.Hour)
			r.ExpiresAt = &expiry
		}
	}
	r.Status = decision
	s.requests[id] = r
	return copyRequest(r), nil
}
func (s *Service) Revoke(user, id, reason string) (Request, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if e := Identity(user); e != nil {
		return Request{}, e
	}
	if !s.reference {
		return Request{}, failure(404, "endpoint unavailable in baseline")
	}
	r, ok := s.requests[id]
	if !ok {
		return Request{}, failure(404, "request not found")
	}
	if user != "carol" {
		return Request{}, failure(403, "approver required")
	}
	if user == r.User {
		return Request{}, failure(403, "self revoke forbidden")
	}
	reason = strings.TrimSpace(reason)
	if reason == "" {
		return Request{}, failure(400, "reason required")
	}
	if r.Status != "approved" {
		return Request{}, failure(409, "approved request required")
	}
	if r.RevokedAt != nil {
		if reason != r.RevokeReason {
			return Request{}, failure(409, "revocation reason already recorded")
		}
		return copyRequest(r), nil
	}
	now := s.clock().UTC()
	r.RevokedAt = &now
	r.RevokeReason = reason
	s.requests[id] = r
	return copyRequest(r), nil
}
func (s *Service) Access(actor, user, resource string) (bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if e := Identity(actor); e != nil {
		return false, e
	}
	if Identity(user) != nil || !resourceValid(resource) {
		return false, failure(400, "unknown user or resource")
	}
	if actor != user && actor != "carol" {
		return false, failure(403, "foreign access check")
	}
	now := s.clock()
	for _, r := range s.requests {
		if r.User == user && r.Resource == resource && r.Status == "approved" && r.RevokedAt == nil && (r.ExpiresAt == nil || now.Before(*r.ExpiresAt)) {
			return true, nil
		}
	}
	return false, nil
}
