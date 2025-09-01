package main

import (
	"net/http"

	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerChirpsGetOne(w http.ResponseWriter, r *http.Request) {
	id := uuid.MustParse(r.PathValue("chirpID"))

	/*
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "ChirpId is wrong", err)
		return
	}
	*/

	chirp, err := cfg.db.GetOneChirp(
		r.Context(),
		id,
	)

	if err != nil {
		respondWithError(w, http.StatusNotFound, "Could not get a chirp", err)
		return
	}

	respondWithJSON(w, http.StatusOK,
		Chirp{
			ID:        chirp.ID,
			CreatedAt: chirp.CreatedAt,
			UpdatedAt: chirp.UpdatedAt,
			Body:      chirp.Body,
			UserID:    chirp.UserID,
		},
	)
}
