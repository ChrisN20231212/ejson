package main

import (
	"flag"
	"log"
	"net/http"
	"os"

	"ejson/internal/server"
)

func main() {
	listenAddr := flag.String("listen", "0.0.0.0:8080", "HTTP listen address")
	flag.Parse()

	addr := resolveListenAddr(*listenAddr)
	log.Printf("ejson listening on http://%s", addr)
	if err := http.ListenAndServe(addr, server.New()); err != nil {
		log.Fatal(err)
	}
}

func resolveListenAddr(flagValue string) string {
	if port := os.Getenv("PORT"); port != "" {
		return "0.0.0.0:" + port
	}
	return flagValue
}
