package auth

import (
	"errors"
	"github.com/FooxyS/auth-service/pkg/apperrors"
	"log"
	"net/http"

	"github.com/FooxyS/auth-service/auth/models"
	"github.com/FooxyS/auth-service/auth/services"
)

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	//проверка метода запроса
	errMethod := services.CompareRequestMethod(r, http.MethodPost)
	if errMethod != nil {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	// достаём клиент БД
	postgres := new(models.Postgres).FromContext(r.Context())

	//парсим данные пользователя
	us := new(models.UserData)
	errParse := us.ParseJson(r.Body)
	if errParse != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	//вызывается контроллер, который парсит, валидирует, проверяет и сохраняет информацию о пользователе
	errProccess := services.RegisterProccess(r.Context(), us, postgres)

	if errors.Is(errProccess, apperrors.ErrEmptyField) {
		http.Error(w, "Empty Fields", http.StatusBadRequest)
		return
	}

	if errors.Is(errProccess, apperrors.ErrUserExist) {
		http.Error(w, "User Already Exists", http.StatusConflict)
		return
	}

	if errProccess != nil {
		log.Printf("error with RegisterProccess: %v\n", errProccess)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusCreated)
	_, errWrite := w.Write([]byte("User successfully created"))
	if errWrite != nil {
		log.Printf("error writing response: %v\n", errWrite)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}
