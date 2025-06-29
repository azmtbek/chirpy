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
	// const filepathAssets = "./assets/"
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
	// mux.Handle("assets/", http.FileServer(http.Dir(filepathAssets)))

	mux.HandleFunc("GET /api/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)

		w.Write([]byte("OK"))
	})

	mux.HandleFunc("GET /admin/metics", apiConf.handlerMetrics)
	mux.HandleFunc("POST /admin/reset", apiConf.handlerReset)

	mux.HandleFunc("POST /api/validate_chirp", func(w http.ResponseWriter, r *http.Request) {
		type parameters struct {
			Body string `json:"body"`
		}
		type errorMsg struct {
			Error string `json:"error"`
		}
		returnBody := errorMsg{
			Error: "Something went wrong",
		}
		errorData, err := json.Marshal(returnBody)
		if err != nil {
			log.Printf("Error while marshaling error msg")
			w.WriteHeader(500)
			return
		}

		decoder := json.NewDecoder(r.Body)
		params := parameters{}
		err = decoder.Decode(&params)
		if err != nil {
			log.Printf("Error decoding parametes: $s", err)
			w.WriteHeader(500)
			w.Write(errorData)
			return
		}

		if len(params.Body) >= 140 {
			log.Printf("Chirp is too long")
			returnBody = errorMsg{
				Error: "Chirp is too long",
			}
			errorData, err = json.Marshal(returnBody)
			if err != nil {
				log.Printf("Error while marshaling error msg")
				w.WriteHeader(500)
				return
			}

			w.WriteHeader(400)
			w.Write(errorData)
			return
		}

		type valid struct {
			Valid bool `json:"valid"`
		}
		returnValue := valid{
			Valid: true,
		}
		valueData, err := json.Marshal(returnValue)
		if err != nil {
			log.Printf("Error while marshaling valueData")
			w.WriteHeader(500)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(valueData)
	})

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	log.Printf("Server is running on port: %s \n", port)

	log.Fatal(srv.ListenAndServe())
}
