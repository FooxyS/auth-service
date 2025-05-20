package services

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/FooxyS/auth-service/auth/apperrors"
	"github.com/FooxyS/auth-service/auth/models"
	"github.com/jackc/pgx/v5"
)

func CompareRequestMethod(r *http.Request, method string) error {
	if r.Method != method {
		return apperrors.ErrMethodNotSupport
	}
	return nil
}

// координирует процесс регистрации пользователя
func RegisterProccess(ctx context.Context, ud models.NewUser, db models.DBClient) error {
	//валидируем данные пользователя
	errEmpty := ud.Validate()
	if errEmpty != nil {
		return errEmpty
	}

	//проверяем есть ли он в базе
	_, errFind := db.FindUser(ctx, ud.GetEmail())

	if errFind == nil {
		return apperrors.ErrUserExist
	}

	if !errors.Is(errFind, pgx.ErrNoRows) {
		return fmt.Errorf("error with FindUser: %v", errFind)
	}

	//записываем в базу
	errWithSave := db.SaveUser(ctx, ud.Save())
	if errWithSave != nil {
		return fmt.Errorf("error with SaveUser: %v", errWithSave)
	}

	return nil
}
