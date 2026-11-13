package course.m17sdd;

import java.time.Clock;
import java.time.Instant;
import java.time.temporal.ChronoUnit;
import java.util.HashMap;
import java.util.Map;
import java.util.Set;

/** In-memory lab service. X-User identities are fixtures, not authentication. */
public final class Access {
    public record Request(String id, String user, String resource, int hours, String status,
                          Instant approvedAt, Instant expiresAt, Instant revokedAt, String revokeReason) {}

    public static final class Failure extends RuntimeException {
        public final int status;
        Failure(int status, String message) { super(message); this.status = status; }
    }
    static Failure fail(int status, String message) { return new Failure(status, message); }
    static void identity(String user) {
        if (user == null || !Set.of("alice", "bob", "carol").contains(user))
            throw fail(401, "unknown or missing demo identity");
    }
    static boolean resourceValid(String resource) {
        return "handbook".equals(resource) || "runbook".equals(resource);
    }

    // ponytail: one monitor and a linear grant scan; use a database beyond the lab.
    private final Map<String, Request> requests = new HashMap<>();
    private final Clock clock;
    private boolean reference;
    private int next;

    public Access(boolean reference, Clock clock) { this.reference = reference; this.clock = clock; }
    /** Retains already-approved legacy grants; pending legacy requests need replacement. */
    public synchronized void upgrade() { reference = true; }

    public synchronized Request create(String user, String resource, int hours) {
        identity(user);
        if (!resourceValid(resource)) throw fail(400, "unknown resource");
        if (reference ? hours < 1 || hours > 24 : hours != 0) throw fail(400, "invalid hours for mode");
        var r = new Request(Integer.toString(++next), user, resource, hours, "pending", null, null, null, null);
        requests.put(r.id(), r);
        return r;
    }
    public synchronized Request get(String user, String id) {
        identity(user); var r = find(id);
        if (!user.equals(r.user()) && !user.equals("carol")) throw fail(403, "foreign request");
        return r;
    }
    public synchronized Request decide(String user, String id, String decision) {
        identity(user); var r = find(id); approver(user);
        if (user.equals(r.user())) throw fail(403, "self decision forbidden");
        if (!"approved".equals(decision) && !"rejected".equals(decision)) throw fail(400, "invalid decision");
        if (r.status().equals(decision)) return r;
        if (!r.status().equals("pending")) throw fail(409, "decision already recorded");
        Instant approved = null, expires = null;
        if (decision.equals("approved")) {
            if (reference && r.hours() == 0) throw fail(400, "legacy pending request needs a new timed request");
            approved = clock.instant();
            if (r.hours() > 0) expires = approved.plus(r.hours(), ChronoUnit.HOURS);
        }
        var decided = new Request(r.id(), r.user(), r.resource(), r.hours(), decision, approved, expires, null, null);
        requests.put(id, decided);
        return decided;
    }
    public synchronized Request revoke(String user, String id, String reason) {
        identity(user);
        if (!reference) throw fail(404, "endpoint unavailable in baseline");
        var r = find(id); approver(user);
        if (user.equals(r.user())) throw fail(403, "self revoke forbidden");
        reason = reason == null ? "" : reason.strip();
        if (reason.isEmpty()) throw fail(400, "reason required");
        if (!r.status().equals("approved")) throw fail(409, "approved request required");
        if (r.revokedAt() != null) {
            if (!reason.equals(r.revokeReason())) throw fail(409, "revocation reason already recorded");
            return r;
        }
        var revoked = new Request(r.id(), r.user(), r.resource(), r.hours(), r.status(),
            r.approvedAt(), r.expiresAt(), clock.instant(), reason);
        requests.put(id, revoked);
        return revoked;
    }
    public synchronized boolean access(String actor, String user, String resource) {
        identity(actor);
        if (user == null || !Set.of("alice", "bob", "carol").contains(user) || !resourceValid(resource))
            throw fail(400, "unknown user or resource");
        if (!actor.equals(user) && !actor.equals("carol")) throw fail(403, "foreign access check");
        Instant now = clock.instant();
        return requests.values().stream().anyMatch(r -> r.user().equals(user) && r.resource().equals(resource)
            && r.status().equals("approved") && r.revokedAt() == null
            && (r.expiresAt() == null || now.isBefore(r.expiresAt())));
    }
    private Request find(String id) {
        var r = requests.get(id); if (r == null) throw fail(404, "request not found"); return r;
    }
    private void approver(String user) { if (!user.equals("carol")) throw fail(403, "approver required"); }
}
