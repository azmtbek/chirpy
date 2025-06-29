package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func handlerChirpsValidate(w http.ResponseWriter, r *http.Request) {
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
}
