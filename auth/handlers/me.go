package auth

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/FooxyS/auth-service/auth/models"
	"github.com/FooxyS/auth-service/pkg/consts"
)

func MeHandler(w http.ResponseWriter, r *http.Request) {

	id, ok := r.Context().Value(consts.USER_ID_KEY).(string)
	if !ok || id == "" {
		log.Println("wrong type or empty context")
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	jsonresp := models.UserJsonID{
		UserID: id,
	}
	errParseJson := json.NewEncoder(w).Encode(jsonresp)
	if errParseJson != nil {
		log.Printf("error with parsing json response: %v\n", errParseJson)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}
