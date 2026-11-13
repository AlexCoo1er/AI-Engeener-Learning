package sdd

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"
)

func NewHandler(service *Service) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		write := func(status int, value any) { w.WriteHeader(status); _ = json.NewEncoder(w).Encode(value) }
		fail := func(e error) { f := e.(*Error); write(f.Status, map[string]string{"error": f.Message}) }
		actor := r.Header.Get("X-User")
		if e := Identity(actor); e != nil {
			fail(e)
			return
		}
		var result any
		var err error
		status := 200
		path := strings.Split(strings.TrimPrefix(r.URL.Path, "/"), "/")
		switch {
		case r.Method == "POST" && r.URL.Path == "/requests":
			var body struct {
				Resource string          `json:"resource"`
				Hours    json.RawMessage `json:"hours"`
			}
			if e := decode(r, w, &body); e != nil {
				fail(e)
				return
			}
			hours := 0
			if body.Hours != nil {
				if string(body.Hours) == "null" {
					fail(failure(400, "hours must be an integer"))
					return
				}
				if e := json.Unmarshal(body.Hours, &hours); e != nil {
					fail(failure(400, "hours must be an integer"))
					return
				}
			}
			result, err = service.Create(actor, body.Resource, hours)
			status = 201
		case r.Method == "GET" && len(path) == 2 && path[0] == "requests":
			result, err = service.Get(actor, path[1])
		case r.Method == "POST" && len(path) == 3 && path[0] == "requests" && path[2] == "decision":
			var body struct {
				Decision string `json:"decision"`
			}
			if e := decode(r, w, &body); e != nil {
				fail(e)
				return
			}
			result, err = service.Decide(actor, path[1], body.Decision)
		case r.Method == "POST" && len(path) == 3 && path[0] == "requests" && path[2] == "revoke":
			var body struct {
				Reason string `json:"reason"`
			}
			if e := decode(r, w, &body); e != nil {
				fail(e)
				return
			}
			result, err = service.Revoke(actor, path[1], body.Reason)
		case r.Method == "GET" && r.URL.Path == "/access":
			user, resource := r.URL.Query().Get("user"), r.URL.Query().Get("resource")
			var allowed bool
			allowed, err = service.Access(actor, user, resource)
			result = map[string]any{"user": user, "resource": resource, "allowed": allowed}
		default:
			err = failure(404, "endpoint not found")
		}
		if err != nil {
			fail(err)
			return
		}
		write(status, result)
	})
}
func decode(r *http.Request, w http.ResponseWriter, value any) error {
	d := json.NewDecoder(http.MaxBytesReader(w, r.Body, 8192))
	d.DisallowUnknownFields()
	if e := d.Decode(value); e != nil {
		return failure(400, "invalid JSON object")
	}
	var extra any
	if e := d.Decode(&extra); e != io.EOF {
		return failure(400, "one JSON object required")
	}
	return nil
}
