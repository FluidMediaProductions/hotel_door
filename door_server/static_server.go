package main

import (
	"net/http"
	"log"
)

func serveStatic(addr string, dir string) {
	fs := http.FileServer(http.Dir(dir))
	http.Handle("/", fs)

	log.Printf("Listening on %s\n", addr)
	http.ListenAndServe(addr, nil)
}
