package course.m17sdd;

import com.sun.net.httpserver.HttpExchange;
import com.sun.net.httpserver.HttpServer;
import java.io.IOException;
import java.net.InetSocketAddress;
import java.net.URLDecoder;
import java.nio.charset.StandardCharsets;
import java.time.Clock;
import java.util.LinkedHashMap;
import java.util.Map;
import java.util.Set;

/** Local demo server: baseline by default, --mode reference for the completed lab. */
public final class Server {
    public static HttpServer start(Access service, int port) throws IOException {
        var server = HttpServer.create(new InetSocketAddress("127.0.0.1", port), 0);
        server.createContext("/", exchange -> handle(service, exchange));
        server.start(); return server;
    }
    private static void handle(Access service, HttpExchange x) throws IOException {
        try {
            String user = x.getRequestHeaders().getFirst("X-User"); Access.identity(user);
            String method = x.getRequestMethod(), path = x.getRequestURI().getPath();
            String[] parts = path.split("/", -1);
            Object result; int status = 200;
            if (method.equals("POST") && path.equals("/requests")) {
                var body = body(x, Set.of("resource", "hours")); int hours = 0;
                if (body.containsKey("hours")) {
                    Object h = body.get("hours");
                    if (!(h instanceof Long n) || n < Integer.MIN_VALUE || n > Integer.MAX_VALUE)
                        throw Access.fail(400, "hours must be an integer");
                    hours = ((Long)h).intValue();
                }
                result = service.create(user, string(body, "resource"), hours); status = 201;
            } else if (method.equals("GET") && parts.length == 3 && parts[1].equals("requests")) {
                result = service.get(user, parts[2]);
            } else if (method.equals("POST") && parts.length == 4 && parts[1].equals("requests") && parts[3].equals("decision")) {
                result = service.decide(user, parts[2], string(body(x, Set.of("decision")), "decision"));
            } else if (method.equals("POST") && parts.length == 4 && parts[1].equals("requests") && parts[3].equals("revoke")) {
                result = service.revoke(user, parts[2], string(body(x, Set.of("reason")), "reason"));
            } else if (method.equals("GET") && path.equals("/access")) {
                var query = query(x.getRequestURI().getRawQuery());
                String target = query.get("user"), resource = query.get("resource");
                boolean allowed = service.access(user, target, resource);
                result = Map.of("user", target, "resource", resource, "allowed", allowed);
            } else throw Access.fail(404, "endpoint not found");
            send(x, status, result instanceof Access.Request r ? fields(r) : cast(result));
        } catch (Access.Failure f) {
            send(x, f.status, Map.of("error", f.getMessage()));
        } finally { x.close(); }
    }
    @SuppressWarnings("unchecked")
    private static Map<String, Object> cast(Object value) { return (Map<String, Object>)value; }
    private static Map<String, Object> body(HttpExchange x, Set<String> keys) throws IOException {
        byte[] bytes = x.getRequestBody().readNBytes(8193);
        if (bytes.length > 8192) throw Access.fail(400, "JSON object too large");
        var body = Json.parse(new String(bytes, StandardCharsets.UTF_8));
        if (!keys.containsAll(body.keySet())) throw Access.fail(400, "unknown JSON field");
        return body;
    }
    private static String string(Map<String, Object> body, String key) {
        Object value = body.get(key);
        if (value == null) return "";
        if (!(value instanceof String)) throw Access.fail(400, key + " must be a string");
        return (String)value;
    }
    private static Map<String, String> query(String raw) {
        var query = new LinkedHashMap<String, String>();
        if (raw == null) return query;
        try {
            for (String item : raw.split("&")) {
                String[] pair = item.split("=", 2);
                String key = URLDecoder.decode(pair[0], StandardCharsets.UTF_8);
                String value = pair.length == 2 ? URLDecoder.decode(pair[1], StandardCharsets.UTF_8) : "";
                query.putIfAbsent(key, value);
            }
        } catch (IllegalArgumentException e) { throw Access.fail(400, "invalid query"); }
        return query;
    }
    private static Map<String, Object> fields(Access.Request r) {
        var out = new LinkedHashMap<String, Object>();
        out.put("id",r.id()); out.put("user",r.user()); out.put("resource",r.resource());
        out.put("hours",r.hours()); out.put("status",r.status());
        if (r.approvedAt()!=null) out.put("approved_at",r.approvedAt());
        if (r.expiresAt()!=null) out.put("expires_at",r.expiresAt());
        if (r.revokedAt()!=null) out.put("revoked_at",r.revokedAt());
        if (r.revokeReason()!=null) out.put("revoke_reason",r.revokeReason());
        return out;
    }
    private static void send(HttpExchange x, int status, Map<String, Object> value) throws IOException {
        byte[] bytes = Json.encode(value).getBytes(StandardCharsets.UTF_8);
        x.getResponseHeaders().set("Content-Type", "application/json; charset=utf-8");
        x.sendResponseHeaders(status, bytes.length); x.getResponseBody().write(bytes);
    }
    public static void main(String[] args) throws IOException {
        String mode = "baseline"; int port = 8080;
        for (int i=0;i<args.length;i+=2) {
            if (i+1 == args.length) throw new IllegalArgumentException("flags require values");
            switch(args[i]) { case "--mode" -> mode=args[i+1]; case "--port" -> port=Integer.parseInt(args[i+1]);
                default -> throw new IllegalArgumentException("unknown flag " + args[i]); }
        }
        if (!mode.equals("baseline") && !mode.equals("reference")) throw new IllegalArgumentException("mode must be baseline or reference");
        var server = start(new Access(mode.equals("reference"), Clock.systemUTC()), port);
        Runtime.getRuntime().addShutdownHook(new Thread(() -> server.stop(0)));
        System.out.printf("demo %s server at http://127.0.0.1:%d (in-memory; X-User: alice/bob/carol)%n",mode,server.getAddress().getPort());
    }
}
