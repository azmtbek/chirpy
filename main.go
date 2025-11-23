package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/azmtbek/chirpy/internal/database"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type apiConfig struct {
	fileserverHits atomic.Int32
	db             *database.Queries
	platform       string
}

func main() {
	const filepathRoot = "."
	const port = "8080"

	godotenv.Load()
	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("DB_URL must be set")
	}
	dbConn, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("Error opening databse: %s", err)
	}
	dbQueries := database.New(dbConn)

	env := os.Getenv("PLATFORM")
	if env == "" {
		log.Printf("PLATFORM is empty")
	}
	apiConf := apiConfig{
		fileserverHits: atomic.Int32{},
		db:             dbQueries,
		platform:       env,
	}

	mux := http.NewServeMux()
	fsHandler := apiConf.middlewareMetricsInc(
		http.StripPrefix(
			"/app",
			http.FileServer(http.Dir(filepathRoot)),
		),
	)
	mux.Handle("/app/", fsHandler)

	mux.HandleFunc("GET /api/healthz", handlerReadiness)
	mux.HandleFunc("POST /api/users", apiConf.handlerUsers)
	mux.HandleFunc("POST /api/login", apiConf.handlerLogin)
	mux.HandleFunc("POST /api/chirps", apiConf.handlerChirpsCreate)
	mux.HandleFunc("GET /api/chirps", apiConf.handlerChirpsGetAll)
	mux.HandleFunc("GET /api/chirps/{chirpID}", apiConf.handlerChirpsGetOne)

	mux.HandleFunc("GET /admin/metics", apiConf.handlerMetrics)
	mux.HandleFunc("POST /admin/reset", apiConf.handlerReset)

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
	log.Fatal(srv.ListenAndServe())
}
