package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

// version is injected at build time via -ldflags "-X main.version=...".
// It defaults to "dev" for local runs.
var version = "dev"

// startTime is used to report uptime on the health endpoint.
var startTime = time.Now()

// healthResponse is the JSON payload returned by /health.
type healthResponse struct {
	Status string `json:"status"`
	Uptime string `json:"uptime"`
}

// versionResponse is the JSON payload returned by /version.
type versionResponse struct {
	Version  string `json:"version"`
	Hostname string `json:"hostname"`
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", rootHandler)
	mux.HandleFunc("/health", healthHandler)
	mux.HandleFunc("/version", versionHandler)

	srv := &http.Server{
		Addr:              ":" + port,
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
	}

	log.Printf("rumbles starting: version=%s port=%s", version, port)
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("server failed: %v", err)
	}
}

// rootHandler returns a simple greeting so a browser hit shows something useful.
func rootHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	fmt.Fprintf(w, "rumbles %s\n", version)
}

// healthHandler is the liveness/readiness probe target and the signal source
// for Argo Rollouts analysis gates.
func healthHandler(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, healthResponse{
		Status: "ok",
		Uptime: time.Since(startTime).Round(time.Second).String(),
	})
}

// versionHandler exposes the running version. During a canary rollout this lets
// you observe old vs new versions serving traffic side by side.
func versionHandler(w http.ResponseWriter, r *http.Request) {
	hostname, _ := os.Hostname()
	writeJSON(w, http.StatusOK, versionResponse{
		Version:  version,
		Hostname: hostname,
	})
}

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(payload); err != nil {
		log.Printf("failed to encode response: %v", err)
	}
}
