package main

import (
	"flag"
	"fmt"
	"github.com/eminsonlu/chirpygo/internal/database"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
)

type apiConfig struct {
	fileserverHits int
	DB             *database.DB
	JWTSecret      string
	PolkaKey       string
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	jwtSecret := os.Getenv("JWT_SECRET")
	polkaKey := os.Getenv("POLKA_KEY")

	mux := http.NewServeMux()

	dbg := flag.Bool("debug", false, "Enable debug mode")
	flag.Parse()

	if *dbg {
		log.Printf("Debug mode enabled")
		err := os.Remove("database.json")
		if err != nil {
		}
	}

	db, err := database.NewDB("database.json")

	if err != nil {
		log.Fatalf("Failed to create database: %v", err)
	}

	cfg := &apiConfig{
		fileserverHits: 0,
		DB:             db,
		JWTSecret:      jwtSecret,
		PolkaKey:       polkaKey,
	}

	mux.Handle("/app/*", http.StripPrefix("/app", cfg.middlewareMetricsInc(http.FileServer(http.Dir(".")))))
	mux.HandleFunc("GET /api/healthz", handlerReadiness)
	mux.HandleFunc("GET /admin/metrics", cfg.handlerMetrics)
	mux.HandleFunc("/api/reset", cfg.handlerReset)
	mux.HandleFunc("POST /api/validate_chirp", ValidateChirp)

	mux.HandleFunc("POST /api/chirps", cfg.handlerChirpsCreate)
	mux.HandleFunc("GET /api/chirps", cfg.handlerChirpsRetrieve)

	mux.HandleFunc("POST /api/users", cfg.handlerUsersCreate)
	mux.HandleFunc("GET /api/users", cfg.handlerUsersRetrieve)

	mux.HandleFunc("POST /api/login", cfg.handlerUsersLogin)
	mux.HandleFunc("PUT /api/users", cfg.handlerUsersUpdate)
	mux.HandleFunc("POST /api/refresh", cfg.handlerUsersRefresh)
	mux.HandleFunc("POST /api/revoke", cfg.handlerUsersRevoke)
	mux.HandleFunc("GET /api/chirps/{chirpID}", cfg.handlerChirpsGet)
	mux.HandleFunc("DELETE /api/chirps/{chirpID}", cfg.handlerChirpsDelete)

	mux.HandleFunc("POST /api/polka/webhooks", cfg.handlerPolka)

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
