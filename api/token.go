package api

import (
	"net/http"
	"context"
	"github.com/ataboo/gowssrv/storage"
	"github.com/satori/go.uuid"
	"github.com/op/go-logging"
	"github.com/ataboo/gowssrv/models"
)

func tokenMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, err := userFromToken(r.Header.Get("Authorization"))
		if err != nil {
			jsonErrorResponse(w, 401, "not authorized")
			return
		}

		// Add user to request for future handlers.
		r = r.WithContext(context.WithValue(r.Context(), "user", user))

		next(w, r)
	})
}

func userFromToken(token string) (models.User, error) {
	parsed, err := uuid.FromString(token)
	if err != nil || parsed.Version() != uuid.V4 {
		logging.MustGetLogger("gowssrv").Debug("Token not valid UUID4")

	}

	user, err := storage.Users.BySession(token)
	if err != nil {
		logging.MustGetLogger("gowssrv").Debug("No user found with session id")
		return user, err
	}

	return user, nil
}

func makeToken() string {
	return uuid.Must(uuid.NewV4()).String()
}