package main

import (
	"net/http"

	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerChirpsGetOne(w http.ResponseWriter, r *http.Request) {
	chirpIDString := r.PathValue("chirpID")
	chirpID, err := uuid.Parse(chirpIDString)

	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid chirp ID", err)
		return
	}

	dbChirp, err := cfg.db.GetOneChirp(r.Context(), chirpID)

	if err != nil {
		respondWithError(w, http.StatusNotFound, "Could not get a chirp", err)
		return
	}

	respondWithJSON(w, http.StatusOK,
		Chirp{
			ID:        dbChirp.ID,
			CreatedAt: dbChirp.CreatedAt,
			UpdatedAt: dbChirp.UpdatedAt,
			Body:      dbChirp.Body,
			UserID:    dbChirp.UserID,
		},
	)
}

func (cfg *apiConfig) handlerChirpsGetAll(w http.ResponseWriter, r *http.Request) {
	dbChirps, err := cfg.db.GetAllChirps(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not get chirps", err)
		return
	}
	
	chirps := []Chirp{}
	for _, chirp := range dbChirps {
		chirps = append(chirps, Chirp{
			ID:        chirp.ID,
			CreatedAt: chirp.CreatedAt,
			UpdatedAt: chirp.UpdatedAt,
			Body:      chirp.Body,
			UserID:    chirp.UserID,
		})
	}

	respondWithJSON(w, http.StatusOK, chirps)

}


