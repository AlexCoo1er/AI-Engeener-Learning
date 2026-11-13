package sdd

import (
	"reflect"
	"sync"
	"testing"
	"time"
)

func code(err error) int {
	if err == nil {
		return 200
	}
	return err.(*Error).Status
}
func TestIdentityAndTransitions(t *testing.T) {
	s := NewService(false, time.Now)
	if _, e := s.Create("mallory", "handbook", 0); code(e) != 401 {
		t.Fatalf("identity: %v", e)
	}
	if _, e := s.Decide("bob", "missing", "approved"); code(e) != 404 {
		t.Fatalf("missing decision: %v", e)
	}
	r, e := s.Create("alice", "handbook", 0)
	if e != nil {
		t.Fatal(e)
	}
	if _, e = s.Get("bob", r.ID); code(e) != 403 {
		t.Fatalf("foreign read: %v", e)
	}
	if _, e = s.Decide("alice", r.ID, "approved"); code(e) != 403 {
		t.Fatalf("employee decision: %v", e)
	}
	if _, e := s.Create("alice", "unknown", 0); code(e) != 400 {
		t.Fatalf("resource: %v", e)
	}
	if _, e := s.Get("alice", "missing"); code(e) != 404 {
		t.Fatalf("missing: %v", e)
	}
	if view, e := s.Get("carol", r.ID); e != nil || view.ID != r.ID {
		t.Fatalf("approver read: %v", e)
	}
	own, _ := s.Create("carol", "runbook", 0)
	if _, e = s.Decide("carol", own.ID, "approved"); code(e) != 403 {
		t.Fatalf("self approval: %v", e)
	}
	if _, e = s.Decide("carol", own.ID, "rejected"); code(e) != 403 {
		t.Fatalf("self reject: %v", e)
	}
	first, e := s.Decide("carol", r.ID, "approved")
	if e != nil {
		t.Fatal(e)
	}
	repeated, e := s.Decide("carol", r.ID, "approved")
	if e != nil || !reflect.DeepEqual(repeated, first) {
		t.Fatalf("repeat changed request: %v", e)
	}
	if _, e = s.Decide("carol", r.ID, "rejected"); code(e) != 409 {
		t.Fatalf("opposite decision: %v", e)
	}
	if allowed, e := s.Access("alice", "alice", "handbook"); !allowed || e != nil {
		t.Fatalf("grant missing: %v", e)
	}
	if _, e = s.Access("bob", "alice", "handbook"); code(e) != 403 {
		t.Fatalf("foreign access: %v", e)
	}
	if _, e = s.Revoke("carol", r.ID, "done"); code(e) != 404 {
		t.Fatalf("baseline revoke: %v", e)
	}
	rejected, _ := s.Create("bob", "runbook", 0)
	rejectedFirst, _ := s.Decide("carol", rejected.ID, "rejected")
	if repeated, e := s.Decide("carol", rejected.ID, "rejected"); e != nil || !reflect.DeepEqual(repeated, rejectedFirst) {
		t.Fatalf("repeat rejected: %v", e)
	}
	if _, e = s.Decide("carol", rejected.ID, "approved"); code(e) != 409 {
		t.Fatalf("approve rejected: %v", e)
	}
	if allowed, _ := s.Access("bob", "bob", "runbook"); allowed {
		t.Fatal("rejection granted access")
	}
}
func TestExpiryLegacyAndReplay(t *testing.T) {
	now := time.Date(2026, 11, 5, 10, 0, 0, 0, time.UTC)
	s := NewService(false, func() time.Time { return now })
	legacy, _ := s.Create("bob", "handbook", 0)
	s.Decide("carol", legacy.ID, "approved")
	pending, _ := s.Create("alice", "runbook", 0)
	s.Upgrade()
	for _, hours := range []int{-1, 0, 25} {
		if _, e := s.Create("alice", "handbook", hours); code(e) != 400 {
			t.Fatalf("hours %d: %v", hours, e)
		}
	}
	for _, hours := range []int{1, 24} {
		if _, e := s.Create("alice", "handbook", hours); e != nil {
			t.Fatal(e)
		}
	}
	// Only already-approved legacy grants remain indefinite.
	if _, e := s.Decide("carol", pending.ID, "approved"); code(e) != 400 {
		t.Fatalf("pending legacy: %v", e)
	}
	if _, e := s.Decide("carol", pending.ID, "rejected"); e != nil {
		t.Fatal(e)
	}
	r, _ := s.Create("alice", "handbook", 1)
	approved, _ := s.Decide("carol", r.ID, "approved")
	want := now.Add(time.Hour)
	if approved.ExpiresAt == nil || !approved.ExpiresAt.Equal(want) {
		t.Fatalf("expiry: %+v", approved)
	}
	now = want.Add(-time.Nanosecond)
	if allowed, _ := s.Access("alice", "alice", "handbook"); !allowed {
		t.Fatal("expired early")
	}
	now = want
	if allowed, _ := s.Access("alice", "alice", "handbook"); allowed {
		t.Fatal("access at expiry")
	}
	replay, _ := s.Decide("carol", r.ID, "approved")
	if !replay.ExpiresAt.Equal(want) {
		t.Fatal("approval extended expiry")
	}
	now = now.Add(365 * 24 * time.Hour)
	if allowed, _ := s.Access("bob", "bob", "handbook"); !allowed {
		t.Fatal("legacy grant expired")
	}
}
func TestRevocationAndIndependentGrants(t *testing.T) {
	now := time.Date(2026, 11, 5, 10, 0, 0, 0, time.UTC)
	s := NewService(true, func() time.Time { return now })
	if _, e := s.Revoke("bob", "missing", "done"); code(e) != 404 {
		t.Fatalf("missing revoke: %v", e)
	}
	own, _ := s.Create("carol", "runbook", 1)
	if _, e := s.Revoke("carol", own.ID, "done"); code(e) != 403 {
		t.Fatalf("self revoke: %v", e)
	}
	denied, _ := s.Create("bob", "runbook", 1)
	s.Decide("carol", denied.ID, "rejected")
	if _, e := s.Revoke("carol", denied.ID, "done"); code(e) != 409 {
		t.Fatalf("rejected revoke: %v", e)
	}
	r, _ := s.Create("alice", "handbook", 1)
	if _, e := s.Revoke("carol", r.ID, "done"); code(e) != 409 {
		t.Fatalf("pending revoke: %v", e)
	}
	s.Decide("carol", r.ID, "approved")
	if _, e := s.Revoke("alice", r.ID, "done"); code(e) != 403 {
		t.Fatalf("employee revoke: %v", e)
	}
	if _, e := s.Revoke("carol", r.ID, " \t"); code(e) != 400 {
		t.Fatalf("empty reason: %v", e)
	}
	revoked, e := s.Revoke("carol", r.ID, "done")
	if e != nil || revoked.Status != "approved" {
		t.Fatalf("revoke: %+v %v", revoked, e)
	}
	now = now.Add(time.Minute)
	repeat, _ := s.Revoke("carol", r.ID, "  done\t")
	if !reflect.DeepEqual(repeat, revoked) {
		t.Fatal("repeat revoke changed record")
	}
	if _, e := s.Revoke("carol", r.ID, "other"); code(e) != 409 {
		t.Fatalf("changed reason: %v", e)
	}
	s.Decide("carol", r.ID, "approved")
	if allowed, _ := s.Access("alice", "alice", "handbook"); allowed {
		t.Fatal("replay revived revoked grant")
	}
	r2, _ := s.Create("alice", "handbook", 1)
	s.Decide("carol", r2.ID, "approved")
	if allowed, _ := s.Access("alice", "alice", "handbook"); !allowed {
		t.Fatal("independent grant lost")
	}
}
func TestConcurrentConflictingDecisions(t *testing.T) {
	s := NewService(true, time.Now)
	r, _ := s.Create("alice", "handbook", 1)
	var wg sync.WaitGroup
	results := make(chan int, 2)
	for _, decision := range []string{"approved", "rejected"} {
		wg.Add(1)
		go func(d string) { defer wg.Done(); _, e := s.Decide("carol", r.ID, d); results <- code(e) }(decision)
	}
	wg.Wait()
	close(results)
	sum := 0
	for status := range results {
		sum += status
	}
	if sum != 609 {
		t.Fatalf("want one success and one conflict; sum=%d", sum)
	}
	// Concurrent readers and independent creates exercise the synchronized store under -race.
	for i := 0; i < 20; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			s.Get("alice", r.ID)
			s.Access("alice", "alice", "handbook")
			s.Create("bob", "runbook", 1)
		}()
	}
	wg.Wait()
}

func TestApprovalClockAndLegacyRevocation(t *testing.T) {
	now := time.Date(2026, 11, 5, 10, 0, 0, 0, time.UTC)
	s := NewService(false, func() time.Time { return now })
	legacy, _ := s.Create("bob", "handbook", 0)
	s.Decide("carol", legacy.ID, "approved")
	s.Upgrade()
	if _, e := s.Revoke("carol", legacy.ID, "retired"); e != nil {
		t.Fatal(e)
	}
	if allowed, _ := s.Access("bob", "bob", "handbook"); allowed {
		t.Fatal("legacy revoked grant active")
	}
	timed, _ := s.Create("alice", "runbook", 24)
	now = now.Add(2 * time.Hour)
	approved, _ := s.Decide("carol", timed.ID, "approved")
	want := time.Date(2026, 11, 6, 12, 0, 0, 0, time.UTC)
	if !approved.ExpiresAt.Equal(want) {
		t.Fatalf("approval should start 24h lifetime: %+v", approved)
	}
	// Snapshots cannot mutate the stored approval timestamp or expiration.
	*approved.ExpiresAt = now.Add(100 * time.Hour)
	approved.ID = "forged"
	approved.Status = "rejected"
	stored, e := s.Get("alice", timed.ID)
	if e != nil || !stored.ExpiresAt.Equal(want) || stored.Status != "approved" {
		t.Fatalf("mutable snapshot: %+v %v", stored, e)
	}
	now = want
	if _, e = s.Revoke("carol", timed.ID, "expired"); e != nil {
		t.Fatal(e)
	}
}
