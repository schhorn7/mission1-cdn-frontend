package main

import (
	_ "embed"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

//go:embed static/index.html
var indexHTML []byte

func main() {
	_ = godotenv.Load()

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	centralAPIURL := os.Getenv("CENTRAL_API_URL")
	if centralAPIURL == "" {
		centralAPIURL = "http://localhost:8080"
	}

	mux := http.NewServeMux()

	// 1. Dynamic JavaScript configuration for browser
	mux.HandleFunc("/config.js", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/javascript")
		_, _ = fmt.Fprintf(w, "window.CENTRAL_API_URL = %q;\n", centralAPIURL)
	})

	// 2. Serve Web-App Dashboard
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(indexHTML)
	})

	log.Printf("[Server #1: Standalone Frontend] Running on http://0.0.0.0:%s", port)
	log.Printf("[Server #1: Standalone Frontend] Configured Central API: %s", centralAPIURL)
	if err := http.ListenAndServe(":"+port, mux); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
