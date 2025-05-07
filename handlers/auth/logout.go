package auth

import "net/http"

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	/*
		полное удаление сессии из БД (refresh-token(hash), user-agent...)
	*/
}
