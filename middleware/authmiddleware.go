package middleware

import (
	"context"
	"log"
	"net/http"

	"github.com/FooxyS/auth-service/auth/apperrors"
	"github.com/FooxyS/auth-service/auth/services"
	"github.com/FooxyS/auth-service/pkg/consts"
)

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		pgpool := services.GetDBPoolFromContext(r)

		bearerToken := services.GetValueFromHeader(r, "Authorization")

		errIsEmpty := services.ValidateStringNotEmpty(bearerToken)
		if errIsEmpty != nil {
			services.UnauthorizedResponseWithReason(w, errIsEmpty)
			return
		}

		access, errParseBearer := services.ParseBearerToToken(bearerToken)
		if errParseBearer != nil {
			services.UnauthorizedResponseWithReason(w, errParseBearer)
			return
		}

		accessClaims, token, errWithParse := services.ParseAccesstoken(access)
		if errWithParse != nil {
			services.UnauthorizedResponseWithReason(w, errWithParse)
			return
		}

		isValid := services.ValidateAccessToken(token)
		if !isValid {
			services.UnauthorizedResponseWithReason(w, apperrors.ErrInvalidToken)
			return
		}

		userSession, errGetSession := services.GetSessionByUserID(r.Context(), pgpool, accessClaims.UserID)
		if errGetSession != nil {
			log.Printf("error with Scan: %v\n", errGetSession)
			services.UnauthorizedResponseWithReason(w, errGetSession)
			return
		}

		host, errSplitHost := services.GetClientIPFromRemoteAddr(r)
		if errSplitHost != nil {
			log.Printf("error with remoteAddr: %v\n", errSplitHost)
			services.UnauthorizedResponseWithReason(w, errSplitHost)
			return
		}

		errIP := services.CompareIP(host, userSession.IP)
		if errIP != nil {
			errDelete := services.DeleteUserByID(pgpool, r, accessClaims.UserID)
			if errDelete != nil {
				log.Printf("error with delete session (user_id=%s): %v\n", accessClaims.UserID, errDelete)
			}
			services.UnauthorizedResponseWithReason(w, errDelete)
			return
		}

		agent := services.GetValueFromHeader(r, "User-Agent")
		errAgent := services.CompareUserAgent(agent, userSession.UserAgent)
		if errAgent != nil {
			errDelete := services.DeleteUserByID(pgpool, r, accessClaims.UserID)
			if errDelete != nil {
				log.Printf("error with delete session (user_id=%s): %v\n", accessClaims.UserID, errDelete)
			}
			services.UnauthorizedResponseWithReason(w, errAgent)
			return
		}

		//кладём в контекст userID, отправляем в обработчик
		ctx := context.WithValue(r.Context(), consts.USER_ID_KEY, accessClaims.UserID)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}
