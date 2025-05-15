package auth

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/FooxyS/auth-service/auth/models"
	"github.com/FooxyS/auth-service/pkg/consts"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

// эндпоинт, который должен регистрировать пользователя: принимать POST-запрос с данными пользователя, регистрировать, если ещё нет в базе
func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	pgpool := r.Context().Value(consts.CTX_KEY_DB).(*pgxpool.Pool)

	registrationData := new(models.UserData)

	errWithDecode := json.NewDecoder(r.Body).Decode(registrationData)
	if errWithDecode != nil {
		log.Printf("error with decoding json data: %v\n", errWithDecode)
		http.Error(w, "error with json", http.StatusBadRequest)
		return
	}

	if registrationData.Email == "" || registrationData.Password == "" {
		http.Error(w, "json fields are empty", http.StatusBadRequest)
		return
	}

	UserForCheck := new(models.UserData)
	errQueryRow := pgpool.QueryRow(r.Context(), "select user_id from users where email=$1", registrationData.Email).Scan(&UserForCheck.UserID)
	if errQueryRow == nil {
		w.WriteHeader(http.StatusOK)
		message := fmt.Sprintf("Пользователь %v уже зарегистрирован!\n", UserForCheck.UserID)
		w.Write([]byte(message))
		return
	}

	newIDForUser := uuid.New()

	hashedPassword, errWithHashing := bcrypt.GenerateFromPassword([]byte(registrationData.Password), bcrypt.DefaultCost)
	if errWithHashing != nil {
		log.Printf("error with hashing: %v\n", errWithHashing)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	creationDate := time.Now()

	_, errExec := pgpool.Exec(r.Context(), "insert into users (user_id, email, password, creation_date) VALUES ($1, $2, $3, $4)", newIDForUser.String(), registrationData.Email, hashedPassword, creationDate)
	if errExec != nil {
		log.Printf("error with Exec: %v\n", errExec)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	message := "Пользователь успешно зарегистрирован"

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(message))
}
