package auth

import (
	"net/http"

	"github.com/FooxyS/auth-service/auth/models"
	"github.com/FooxyS/auth-service/auth/services"
)

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	//проверка метода
	errMethod := services.CompareRequestMethod(r, http.MethodPost)
	if errMethod != nil {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	//получение клинта БД
	db := new(models.Postgres)
	postgres := db.FromContext(r.Context())

	//основная логика
	LoginProccess()

	//отправка ответа

}
