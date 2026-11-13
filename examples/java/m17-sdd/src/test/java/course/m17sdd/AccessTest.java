package course.m17sdd;

import org.junit.jupiter.api.Test;
import java.time.*;
import java.util.concurrent.*;
import static org.junit.jupiter.api.Assertions.*;

class AccessTest {
    static class TestClock extends Clock {
        Instant now = Instant.parse("2026-11-05T10:00:00Z");
        public ZoneId getZone() { return ZoneOffset.UTC; }
        public Clock withZone(ZoneId zone) { return this; }
        public Instant instant() { return now; }
    }
    void status(int expected, Runnable action) {
        assertEquals(expected, assertThrows(Access.Failure.class, action::run).status);
    }
    @Test void identityAndTransitions() {
        var s = new Access(false, Clock.systemUTC());
        status(401, () -> s.create("mallory", "handbook", 0));
        status(404, () -> s.decide("bob", "missing", "approved"));
        var r = s.create("alice", "handbook", 0);
        status(403, () -> s.get("bob", r.id()));
        status(403, () -> s.decide("alice", r.id(), "approved"));
        status(400, () -> s.create("alice", "unknown", 0));
        status(404, () -> s.get("alice", "missing"));
        assertEquals(r, s.get("carol", r.id()));
        var own = s.create("carol", "runbook", 0);
        status(403, () -> s.decide("carol", own.id(), "approved"));
        status(403, () -> s.decide("carol", own.id(), "rejected"));
        var first = s.decide("carol", r.id(), "approved");
        assertEquals(first, s.decide("carol", r.id(), "approved"));
        status(409, () -> s.decide("carol", r.id(), "rejected"));
        assertTrue(s.access("alice", "alice", "handbook"));
        status(403, () -> s.access("bob", "alice", "handbook"));
        status(404, () -> s.revoke("carol", r.id(), "done"));
        var rejected = s.create("bob", "runbook", 0);
        var firstRejected = s.decide("carol", rejected.id(), "rejected");
        assertEquals(firstRejected, s.decide("carol", rejected.id(), "rejected"));
        status(409, () -> s.decide("carol", rejected.id(), "approved"));
        assertFalse(s.access("bob", "bob", "runbook"));
    }
    @Test void expiryLegacyAndReplay() {
        var clock = new TestClock(); var s = new Access(false, clock);
        var legacy = s.create("bob", "handbook", 0); s.decide("carol", legacy.id(), "approved");
        var pending = s.create("alice", "runbook", 0); s.upgrade();
        for (int hours : new int[]{-1,0,25}) status(400, () -> s.create("alice", "handbook", hours));
        for (int hours : new int[]{1,24}) assertNotNull(s.create("alice", "handbook", hours));
        status(400, () -> s.decide("carol", pending.id(), "approved"));
        assertEquals("rejected", s.decide("carol", pending.id(), "rejected").status());
        var r = s.create("alice", "handbook", 1);
        var approved = s.decide("carol", r.id(), "approved");
        Instant expiry = Instant.parse("2026-11-05T11:00:00Z");
        assertEquals(expiry, approved.expiresAt());
        clock.now = expiry.minusNanos(1); assertTrue(s.access("alice", "alice", "handbook"));
        clock.now = expiry; assertFalse(s.access("alice", "alice", "handbook"));
        assertEquals(approved, s.decide("carol", r.id(), "approved"));
        clock.now = expiry.plus(Duration.ofDays(365)); assertTrue(s.access("bob", "bob", "handbook"));
    }
    @Test void revokeAndIndependentGrants() {
        var clock = new TestClock(); var s = new Access(true, clock);
        status(404, () -> s.revoke("bob", "missing", "done"));
        var own = s.create("carol", "runbook", 1);
        status(403, () -> s.revoke("carol", own.id(), "done"));
        var rejected = s.create("bob", "runbook", 1); s.decide("carol", rejected.id(), "rejected");
        status(409, () -> s.revoke("carol", rejected.id(), "done"));
        var r = s.create("alice", "handbook", 1);
        status(409, () -> s.revoke("carol", r.id(), "done"));
        s.decide("carol", r.id(), "approved");
        status(403, () -> s.revoke("alice", r.id(), "done"));
        status(400, () -> s.revoke("carol", r.id(), " \t"));
        var revoked = s.revoke("carol", r.id(), "done"); assertEquals("approved", revoked.status());
        clock.now = clock.now.plusSeconds(60); assertEquals(revoked, s.revoke("carol", r.id(), "  done\t"));
        status(409, () -> s.revoke("carol", r.id(), "other"));
        assertEquals(revoked, s.decide("carol", r.id(), "approved"));
        assertFalse(s.access("alice", "alice", "handbook"));
        var other = s.create("alice", "handbook", 1); s.decide("carol", other.id(), "approved");
        assertTrue(s.access("alice", "alice", "handbook"));
    }
    @Test void approvalClockAndLegacyRevocation() {
        var clock = new TestClock(); var s = new Access(false, clock);
        var legacy = s.create("bob", "handbook", 0); s.decide("carol", legacy.id(), "approved"); s.upgrade();
        s.revoke("carol", legacy.id(), "retired"); assertFalse(s.access("bob", "bob", "handbook"));
        var timed = s.create("alice", "runbook", 24); clock.now = clock.now.plusSeconds(7200);
        var approved = s.decide("carol", timed.id(), "approved");
        var expiry = Instant.parse("2026-11-06T12:00:00Z"); assertEquals(expiry, approved.expiresAt());
        clock.now = expiry; assertEquals("approved", s.revoke("carol", timed.id(), "expired").status());
    }
    @Test void concurrentConflictingDecisions() throws Exception {
        var s = new Access(true, Clock.systemUTC()); var r = s.create("alice", "handbook", 1);
        var start = new CountDownLatch(1); var pool = Executors.newFixedThreadPool(2);
        try {
            var a = pool.submit(() -> decision(s, r.id(), "approved", start));
            var b = pool.submit(() -> decision(s, r.id(), "rejected", start));
            start.countDown(); assertEquals(609, a.get(5, TimeUnit.SECONDS) + b.get(5, TimeUnit.SECONDS));
            for (int i=0;i<20;i++) pool.submit(() -> {s.get("alice",r.id());s.access("alice","alice","handbook");s.create("bob","runbook",1);}).get();
        } finally { pool.shutdownNow(); }
    }
    int decision(Access s, String id, String decision, CountDownLatch start) throws InterruptedException {
        start.await(); try { s.decide("carol", id, decision); return 200; } catch (Access.Failure f) { return f.status; }
    }
}
