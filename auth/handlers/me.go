package auth

import (
	"encoding/json"
	"log"
	"net/http"
)

type UserJsonID struct {
	UserID string `json:"userid"`
}

func MeHandler(w http.ResponseWriter, r *http.Request) {

	id, ok := r.Context().Value("UserIDKey").(string)
	if !ok || id == "" {
		log.Println("wrong type or empty context")
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	jsonresp := UserJsonID{
		UserID: id,
	}
	errParseJson := json.NewEncoder(w).Encode(jsonresp)
	if errParseJson != nil {
		log.Printf("error with parsing json response: %v\n", errParseJson)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}
