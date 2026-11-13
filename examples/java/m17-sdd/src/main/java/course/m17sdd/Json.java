package course.m17sdd;

import java.util.LinkedHashMap;
import java.util.Map;
import java.util.regex.Pattern;

/** Deliberately small JSON boundary: flat objects with scalar values only. */
final class Json {
    private static final Pattern NUMBER = Pattern.compile("-?(?:0|[1-9][0-9]*)(?:\\.[0-9]+)?(?:[eE][+-]?[0-9]+)?");
    private final String text;
    private int pos;
    private Json(String text) { this.text = text; }

    static Map<String, Object> parse(String text) {
        try { return new Json(text).object(); }
        catch (IllegalArgumentException e) { throw Access.fail(400, "invalid JSON object"); }
    }
    private Map<String, Object> object() {
        var out = new LinkedHashMap<String, Object>();
        expect('{'); skip();
        if (!take('}')) {
            do {
                skip(); String key = string(); expect(':'); skip();
                Object value = scalar();
                if (out.containsKey(key)) throw new IllegalArgumentException();
                out.put(key, value); skip();
                if (take('}')) break;
                expect(',');
            } while (true);
        }
        skip(); if (pos != text.length()) throw new IllegalArgumentException();
        return out;
    }
    private Object scalar() {
        if (pos < text.length() && text.charAt(pos) == '"') return string();
        int start = pos;
        while (pos < text.length() && ",} \t\r\n".indexOf(text.charAt(pos)) < 0) pos++;
        String token = text.substring(start, pos);
        return switch (token) {
            case "true" -> true;
            case "false" -> false;
            case "null" -> null;
            default -> {
                if (!NUMBER.matcher(token).matches()) throw new IllegalArgumentException();
                if (token.indexOf('.') < 0 && token.indexOf('e') < 0 && token.indexOf('E') < 0)
                    yield Long.parseLong(token);
                yield new java.math.BigDecimal(token);
            }
        };
    }
    private String string() {
        if (!take('"')) throw new IllegalArgumentException();
        var out = new StringBuilder();
        while (pos < text.length()) {
            char c = text.charAt(pos++);
            if (c == '"') return out.toString();
            if (c < 32) throw new IllegalArgumentException();
            if (c != '\\') { out.append(c); continue; }
            if (pos == text.length()) throw new IllegalArgumentException();
            char escaped = text.charAt(pos++);
            switch (escaped) {
                case '"', '\\', '/' -> out.append(escaped);
                case 'b' -> out.append('\b'); case 'f' -> out.append('\f');
                case 'n' -> out.append('\n'); case 'r' -> out.append('\r'); case 't' -> out.append('\t');
                case 'u' -> {
                    if (pos + 4 > text.length()) throw new IllegalArgumentException();
                    out.append((char) Integer.parseInt(text.substring(pos, pos + 4), 16)); pos += 4;
                }
                default -> throw new IllegalArgumentException();
            }
        }
        throw new IllegalArgumentException();
    }
    private void skip() { while (pos < text.length() && " \t\r\n".indexOf(text.charAt(pos)) >= 0) pos++; }
    private boolean take(char c) { if (pos < text.length() && text.charAt(pos) == c) { pos++; return true; } return false; }
    private void expect(char c) { skip(); if (!take(c)) throw new IllegalArgumentException(); }
    static String encode(Map<String, Object> values) {
        var out = new StringBuilder("{");
        for (var entry : values.entrySet()) {
            if (out.length() > 1) out.append(',');
            out.append(quote(entry.getKey())).append(':');
            Object value = entry.getValue();
            out.append(value instanceof Boolean || value instanceof Number ? value.toString() : quote(value.toString()));
        }
        return out.append('}').toString();
    }
    private static String quote(String value) {
        var out = new StringBuilder("\"");
        for (char c : value.toCharArray()) {
            switch (c) {
                case '\\' -> out.append("\\\\"); case '"' -> out.append("\\\"");
                default -> { if (c < 32) out.append(String.format("\\u%04x", (int)c)); else out.append(c); }
            }
        }
        return out.append('"').toString();
    }
}
