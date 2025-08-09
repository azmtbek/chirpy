package main

import (
	"net/http"
)

func (cfg *apiConfig) handlerReset(w http.ResponseWriter, r *http.Request) {
	if cfg.platform != "dev" {
		respondWithError(w, 403, "Method not Allowed in prod", nil)
		return
	}

	cfg.fileserverHits.Store(0)
	err := cfg.db.Reset(r.Context())
	if err != nil {
		respondWithError(w, 400, "Error deleting users ", err)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hit reset to 0 and database is reset to initial state."))
}
