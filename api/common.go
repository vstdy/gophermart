package api

import (
	"context"
	"fmt"
	"net/http"

	"github.com/go-chi/jwtauth/v5"
	"github.com/google/uuid"

	"github.com/vstdy/gophermart/api/model"
	canonical "github.com/vstdy/gophermart/model"
)

func (h Handler) setAuthCookie(w http.ResponseWriter, obj canonical.User) error {
	_, token, err := h.tokenAuth.Encode(model.NewJWTClaims(obj))
	if err != nil {
		return fmt.Errorf("auth cookie: %v", err)
	}

	cookie := http.Cookie{
		Name:  "jwt",
		Value: token,
		Path:  "/",
	}
	http.SetCookie(w, &cookie)

	return nil
}

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

func (h Handler) addOrder(ctx context.Context, userID uuid.UUID, orderID string) (canonical.Order, error) {
	order := model.Order{
		UserID: userID,
		Number: orderID,
	}

	obj, err := order.ToCanonical()
	if err != nil {
		return canonical.Order{}, err
	}

	dbObj, err := h.service.AddOrder(ctx, obj)
	if err != nil {
		return dbObj, err
	}

	return dbObj, nil
}
