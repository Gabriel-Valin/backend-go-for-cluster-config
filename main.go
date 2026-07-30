package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
)

// version é injetada no build via ldflags (o SHA do commit)
var version = "dev"

func main() {
	mux := http.NewServeMux()

	// healthcheck — usado pelas probes do Kubernetes
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	})

	// endpoint principal — mostra a versão em execução
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		resp := map[string]string{
			"message": "hello from meu-backend",
			"version": version,
			"host":    hostname(),
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("iniciando na porta %s (version=%s)", port, version)
	if err := http.ListenAndServe(":"+port, mux); err != nil {
		log.Fatal(err)
	}
}

func hostname() string {
	h, _ := os.Hostname()
	return h
}
