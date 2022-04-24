package gophermart

import (
	"context"
	"fmt"

	"github.com/vstdy/gophermart/model"
	"github.com/vstdy/gophermart/pkg"
	"github.com/vstdy/gophermart/service/gophermart/v1/validator"
)

// CreateUser creates a new model.User.
func (svc *Service) CreateUser(ctx context.Context, rawObj model.User) (model.User, error) {
	if err := validator.ValidateLogin(rawObj.Login); err != nil {
		return model.User{}, fmt.Errorf("%w: login: %v", pkg.ErrInvalidInput, err)
	}
	if err := validator.ValidatePassword(rawObj.Password); err != nil {
		return model.User{}, fmt.Errorf("%w: password: %v", pkg.ErrInvalidInput, err)
	}

	obj, err := svc.storage.CreateUser(ctx, rawObj)
	if err != nil {
		return model.User{}, fmt.Errorf("creating user: %w", err)
	}

	return obj, nil
}

// AuthenticateUser verifies the identity of credentials.
func (svc *Service) AuthenticateUser(ctx context.Context, rawObj model.User) (model.User, error) {
	if err := validator.ValidateLogin(rawObj.Login); err != nil {
		return model.User{}, fmt.Errorf("%w: login: %v", pkg.ErrInvalidInput, err)
	}
	if err := validator.ValidatePassword(rawObj.Password); err != nil {
		return model.User{}, fmt.Errorf("%w: password: %v", pkg.ErrInvalidInput, err)
	}

	obj, err := svc.storage.AuthenticateUser(ctx, rawObj)
	if err != nil {
		return model.User{}, fmt.Errorf("authenticating user: %w", err)
	}

	return obj, nil
}
