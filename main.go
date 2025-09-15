package main

import (
	"database/sql"
	"http-chirpy/internal/database"
	"log"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type apiConfig struct {
	fileserverHits atomic.Int32
	db             *database.Queries
	platform       string
	jwtSecret      string
}

func main() {
	const filepathRoot = "."
	const port = "8080"

	godotenv.Load()

	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatalf("DB_URL must be set")
	}

	platform := os.Getenv("PLATFORM")
	if platform == "" {
		log.Fatalf("PLATFORM must be set")
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Fatalf("JWT_SECRET must be set")
	}

	dbConn, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}
	dbQueries := database.New(dbConn)

	apiCfg := apiConfig{
		fileserverHits: atomic.Int32{},
		db:             dbQueries,
		platform:       platform,
		jwtSecret:      jwtSecret,
	}

	mux := http.NewServeMux()
	mux.Handle("/app/", apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot)))))

	mux.HandleFunc("GET /api/healthz", handlerReadiness)

	mux.HandleFunc("POST /api/login", apiCfg.handlerLogin)
	mux.HandleFunc("POST /api/refresh", apiCfg.handlerRefreshToken)
	mux.HandleFunc("POST /api/revoke", apiCfg.handlerRevokeRefreshToken)

	mux.HandleFunc("POST /api/users", apiCfg.handlerCreateUser)
	mux.HandleFunc("PUT /api/users", apiCfg.handlerUpdateUser)

	mux.HandleFunc("POST /api/chirps", apiCfg.handlerCreateChirp)
	mux.HandleFunc("GET /api/chirps", apiCfg.handlerGetChirps)
	mux.HandleFunc("GET /api/chirps/{chirpID}", apiCfg.handlerGetChirpByID)

	mux.HandleFunc("GET /admin/metrics", apiCfg.handlerMetrics)
	mux.HandleFunc("POST /admin/reset", apiCfg.handlerReset)

	server := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	log.Printf("Serving on port %s:\n", port)
	log.Fatal(server.ListenAndServe())
}
