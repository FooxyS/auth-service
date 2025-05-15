package auth_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/FooxyS/auth-service/handlers/auth"
)

func TestMeHandler(t *testing.T) {
	resp := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/auth/me", nil)

	userID := "db0e693f-c5e9-4e9e-bdb2-cd79a63a0fb8"
	ctx := context.WithValue(req.Context(), "UserIDKey", userID)
	reqWithCtx := req.WithContext(ctx)

	auth.MeHandler(resp, reqWithCtx)

	respjson := new(auth.UserJsonID)

	errDecodeJson := json.NewDecoder(resp.Body).Decode(respjson)
	if errDecodeJson != nil {
		t.Errorf("error with decoding json: %v\n", errDecodeJson)
		return
	}

	respHeader := resp.Header().Get("Content-Type")
	expectedHeader := "application/json"
	if respHeader != expectedHeader {
		t.Errorf("wrong header from response: got %s, want %s\n", respHeader, expectedHeader)
	}

	respStatusCode := resp.Code
	expectedStatusCode := http.StatusOK
	if respStatusCode != expectedStatusCode {
		t.Errorf("wrong status code from response: got %d, want %d\n", respStatusCode, expectedStatusCode)
		return
	}

	if respjson.UserID != userID {
		t.Errorf("wrong id from handler: got %s, want %s\n", respjson.UserID, userID)
		return
	}
}
