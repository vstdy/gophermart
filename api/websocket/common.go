package websocket

import (
	"strings"

	"github.com/go-chi/jwtauth/v5"
	"github.com/lestrrat-go/jwx/jwt"

	"github.com/vstdy/gophermart/pkg"
)

// verifyToken verifies jwt.
func verifyToken(jwtAuth *jwtauth.JWTAuth, bearer string) (jwt.Token, error) {
	if len(bearer) < 8 || strings.ToUpper(bearer[0:6]) != "BEARER" {
		return nil, nil
	}

	return jwtauth.VerifyToken(jwtAuth, bearer[7:])
}

// getUserID retrieves the user ID from the token.
func getUserID(token jwt.Token) (string, error) {
	claims := token.PrivateClaims()
	userID, ok := claims["id"].(string)
	if !ok {
		return "", pkg.ErrWrongCredentials
	}

	return userID, nil
}
