package domain

type User struct {
	UserID       string
	Email        string
	PasswordHash string
}

func (u User) CheckPassword(password string) bool {
	if password == u.PasswordHash {
		return true
	}
	return false
}
