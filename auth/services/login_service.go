package services

import (
	"context"

	"github.com/FooxyS/auth-service/auth/models"
)

// координирует процесс авторизации пользователя
func LoginProccess(ctx context.Context, ud models.NewUser, db models.DBClient, tokens models.GenerateTokens) error {
	//валидирует
	errValid := ud.Validate()
	if errValid != nil {
		return errValid
	}

	//сверяет с базой
	id, errFindUser := db.FindUser(ctx, ud.GetEmail())
	if errFindUser != nil {
		return errFindUser
	}
	user, errGetUser := db.GetUserByID(ctx, id)
	if errGetUser != nil {
		return errGetUser
	}
	errCompare := ud.ComparePassword(*user)
	if errCompare != nil {
		return errCompare
	}
	ud.WriteID(id)

	//создаёт токены
	tokens.CreateTokens(id)

	//отправляет токены
	tokens.Sendtokens()

	return nil
}
