package auth

import "net/http"

//эндпоинт, который должен регистрировать пользователя: принимать POST-запрос с данными пользователя, регистрировать, если ещё нет в базе
func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

}
