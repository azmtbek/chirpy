package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync/atomic"
)

type apiConfig struct {
	fileserverHits atomic.Int32
}

func main() {
	const filepathRoot = "."
	const port = "8080"

	apiConf := &apiConfig{}
	apiConf.fileserverHits.Store(0)
	mux := http.NewServeMux()

	mux.Handle("/app/", apiConf.middlewareMetricsInc(
		http.StripPrefix(
			"/app",
			http.FileServer(http.Dir(filepathRoot)),
		),
	))

	mux.HandleFunc("GET /api/healthz", handlerReadiness)
	mux.HandleFunc("POST /api/validate_chirp", handlerChirpsValidate)

	mux.HandleFunc("GET /admin/metics", apiConf.handlerMetrics)
	mux.HandleFunc("POST /admin/reset", apiConf.handlerReset)

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	log.Printf("Server is running on port: %s \n", port)

	log.Fatal(srv.ListenAndServe())
}
