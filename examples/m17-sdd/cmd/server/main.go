package main

import (
	"flag"
	"log"
	"net"
	"net/http"
	"time"

	sdd "aiec/examples/m17-sdd"
)

func main() {
	mode := flag.String("mode", "baseline", "baseline or reference")
	addr := flag.String("addr", "127.0.0.1:8080", "loopback listen address")
	flag.Parse()
	if *mode != "baseline" && *mode != "reference" {
		log.Fatal("mode must be baseline or reference")
	}
	host, _, err := net.SplitHostPort(*addr)
	if err != nil {
		log.Fatal(err)
	}
	ip := net.ParseIP(host)
	if ip == nil || !ip.IsLoopback() {
		log.Fatal("addr must use a loopback IP")
	}
	server := &http.Server{Addr: *addr, Handler: sdd.NewHandler(sdd.NewService(*mode == "reference", time.Now)), ReadHeaderTimeout: 5 * time.Second, ReadTimeout: 10 * time.Second, WriteTimeout: 10 * time.Second, IdleTimeout: 30 * time.Second}
	log.Printf("demo %s server at http://%s (in-memory; X-User: alice/bob/carol)", *mode, *addr)
	log.Fatal(server.ListenAndServe())
}
