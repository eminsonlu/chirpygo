package main

import (
	"fmt"
	"github.com/eminsonlu/chirpygo/internal/database"
	"log"
	"net/http"
)

type apiConfig struct {
	fileserverHits int
	DB             *database.DB
}

func main() {
	mux := http.NewServeMux()

	db, err := database.NewDB("chirps.json")
	if err != nil {
		log.Fatalf("Failed to create database: %v", err)
	}
	cfg := &apiConfig{
		fileserverHits: 0,
		DB:             db,
	}

	mux.Handle("/app/*", http.StripPrefix("/app", cfg.middlewareMetricsInc(http.FileServer(http.Dir(".")))))
	mux.HandleFunc("GET /api/healthz", handlerReadiness)
	mux.HandleFunc("GET /admin/metrics", cfg.handlerMetrics)
	mux.HandleFunc("/api/reset", cfg.handlerReset)
	mux.HandleFunc("POST /api/validate_chirp", validateChirp)

	mux.HandleFunc("POST /api/chirps", cfg.handlerChirpsCreate)
	mux.HandleFunc("GET /api/chirps", cfg.handlerChirpsRetrieve)

	server := &http.Server{
		Handler: mux,
		Addr:    "localhost:8080",
	}

	log.Printf("Server starting on %s", server.Addr)
	err = server.ListenAndServe()
	if err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits++
		next.ServeHTTP(w, r)
	})
}

func (cfg *apiConfig) handlerMetrics(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf(`
		<html>
		<body>
		    <h1>Welcome, Chirpy Admin</h1>
    		<p>Chirpy has been visited %d times!</p>
		</body>
		</html>
	`, cfg.fileserverHits)))
}

func (cfg *apiConfig) handlerReset(w http.ResponseWriter, r *http.Request) {
	cfg.fileserverHits = 0
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hits reset to 0"))
}
