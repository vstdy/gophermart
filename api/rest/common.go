package rest

import (
	"context"
	"fmt"
	"net/http"

	"github.com/go-chi/jwtauth/v5"
	"github.com/google/uuid"

	"github.com/vstdy/gophermart/api/rest/model"
	canonical "github.com/vstdy/gophermart/model"
)

// addJWTCookie adds a jwt cookie to the response.
func (h Handler) addJWTCookie(w http.ResponseWriter, obj canonical.User) error {
	_, token, err := h.jwtAuth.Encode(model.NewJWTClaims(obj))
	if err != nil {
		return fmt.Errorf("auth cookie: %v", err)
	}

	cookie := http.Cookie{
		Name:  "jwt",
		Value: token,
		Path:  "/",
	}
	http.SetCookie(w, &cookie)

	if _, err = w.Write([]byte(token)); err != nil {
		return err
	}

	return nil
}

// getUserID retrieves the user ID from the context.
func (h Handler) getUserID(ctx context.Context) (uuid.UUID, error) {
	_, claims, err := jwtauth.FromContext(ctx)
	if err != nil {
		return uuid.Nil, err
	}
	userID, err := uuid.Parse(claims["id"].(string))
	if err != nil {
		return uuid.Nil, err
	}

	return userID, nil
}
