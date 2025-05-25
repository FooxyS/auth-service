package handler

type RegisterRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type MessageResponse struct {
	Message string `json:"message"`
}

type TokenResponse struct {
	Access string `json:"access"`
}
