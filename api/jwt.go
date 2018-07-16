package api

import (
	"net/http"
	"github.com/dgrijalva/jwt-go"
	"github.com/ataboo/gowssrv/config"
	"context"
	"github.com/ataboo/gowssrv/models"
	"github.com/ataboo/gowssrv/storage"
	"fmt"
)

func jwtRequired(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, err := parseToken(r.Header.Get("Authorization"))

		if err != nil || !tokenValid(token) {
			jsonErrorResponse(w, 401, "not authorized")
			return
		}

		claims := token.Claims.(jwt.MapClaims)
		user, err := storage.Users.Find(claims["user_id"].(string))

		if err != nil || user.SessionId != token.Signature {
			jsonErrorResponse(w, 401, "not authorized on later one")

			fmt.Printf("\nSig: %s\nSession: %s, userID: %s, err: %s", token.Signature, user.SessionId, claims["user_id"], err)
			return
		}

		// Add user_token to request for next handlers.
		r = r.WithContext(context.WithValue(r.Context(), "user_token", token))

		next(w, r)
	})
}

func parseToken(rawToken string) (*jwt.Token, error) {
	token, err := jwt.Parse(rawToken, jwt.Keyfunc(func(token *jwt.Token) (interface{}, error) {
		return []byte(config.Config.Api.JwtSecret), nil
	}))

	return token, err
}

func tokenValid(token *jwt.Token) bool {
	//TODO may include checks against user id/account

	return token.Valid
}

func createToken(user models.User) (string, error) {
	claims := jwt.MapClaims{
		"user_id": user.ID,
		"username": user.Username,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(config.Config.Api.JwtSecret))
}