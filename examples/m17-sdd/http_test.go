package sdd

import (
	"encoding/json"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestHTTPContract(t *testing.T) {
	h := NewHandler(NewService(true, time.Now))
	call := func(method, path, user, body string, want int) map[string]any {
		t.Helper()
		r := httptest.NewRequest(method, path, strings.NewReader(body))
		r.Header.Set("X-User", user)
		w := httptest.NewRecorder()
		h.ServeHTTP(w, r)
		if w.Code != want {
			t.Fatalf("%s %s: got %d want %d: %s", method, path, w.Code, want, w.Body.String())
		}
		var out map[string]any
		if e := json.Unmarshal(w.Body.Bytes(), &out); e != nil {
			t.Fatal(e)
		}
		return out
	}
	call("POST", "/requests", "", `{"resource":"handbook","hours":1}`, 401)
	for _, body := range []string{`{}`, `{"resource":"other","hours":1}`, `{"resource":"handbook","hours":1.5}`, `{"resource":"handbook","hours":"1"}`, `{"resource":"handbook","hours":null}`, `{"resource":"handbook","hours":false}`, `{"resource":"handbook"}`, `{"resource":"handbook","hours":1,"role":"approver"}`, `{"resource":"handbook","hours":1} {}`, `[]`} {
		call("POST", "/requests", "alice", body, 400)
	}
	r := call("POST", "/requests", "alice", `{"resource":"handbook","hours":1}`, 201)
	id := r["id"].(string)
	call("GET", "/requests/"+id, "bob", "", 403)
	call("GET", "/requests/missing", "alice", "", 404)
	call("POST", "/requests/"+id+"/decision", "bob", `{"decision":"approved"}`, 403)
	call("POST", "/requests/"+id+"/decision", "carol", `{"decision":"maybe"}`, 400)
	call("POST", "/requests/"+id+"/decision", "carol", `{"decision":"approved"}`, 200)
	if out := call("GET", "/access?user=alice&resource=handbook", "alice", "", 200); out["allowed"] != true {
		t.Fatal(out)
	}
	call("GET", "/access?user=alice&resource=handbook", "bob", "", 403)
	call("POST", "/requests/"+id+"/revoke", "carol", `{"reason":"done"}`, 200)
	if out := call("GET", "/access?user=alice&resource=handbook", "alice", "", 200); out["allowed"] != false {
		t.Fatal(out)
	}
	call("GET", "/requests/"+id+"/extra", "alice", "", 404)
}

func TestBaselineHTTP(t *testing.T) {
	h := NewHandler(NewService(false, time.Now))
	for _, body := range []string{`{"resource":"handbook"}`, `{"resource":"runbook","hours":0}`} {
		r := httptest.NewRequest("POST", "/requests", strings.NewReader(body))
		r.Header.Set("X-User", "alice")
		w := httptest.NewRecorder()
		h.ServeHTTP(w, r)
		if w.Code != 201 {
			t.Fatalf("baseline create: %d %s", w.Code, w.Body.String())
		}
		var created Request
		if e := json.Unmarshal(w.Body.Bytes(), &created); e != nil {
			t.Fatal(e)
		}
		r = httptest.NewRequest("POST", "/requests/"+created.ID+"/revoke", strings.NewReader(`{"reason":"done"}`))
		r.Header.Set("X-User", "carol")
		w = httptest.NewRecorder()
		h.ServeHTTP(w, r)
		if w.Code != 404 {
			t.Fatalf("baseline revoke: %d", w.Code)
		}
	}
	r := httptest.NewRequest("POST", "/requests", strings.NewReader(`{"resource":"handbook","hours":1}`))
	r.Header.Set("X-User", "alice")
	w := httptest.NewRecorder()
	h.ServeHTTP(w, r)
	if w.Code != 400 {
		t.Fatalf("baseline accepted timed request: %d", w.Code)
	}
}
