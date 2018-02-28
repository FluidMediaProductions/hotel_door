package main

import (
	"net/http"
	"log"

	"github.com/graphql-go/handler"
	"github.com/rs/cors"
)

func serveStatic(addr string, dir string) {
	initGraphql()
	h := handler.New(&handler.Config{
		Schema: &schema,
		Pretty: true,
		GraphiQL: true,
	})

	corsH := cors.Default().Handler(h)

	http.Handle("/graphql", corsH)

	fs := http.FileServer(http.Dir(dir))
	http.Handle("/", fs)

	log.Printf("Listening on %s\n", addr)
	http.ListenAndServe(addr, nil)
}
