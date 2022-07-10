package pgbun

import (
	"context"
	"database/sql"
	"errors"

	"github.com/uptrace/bun/driver/pgdriver"
	"github.com/vstdy/gophermart/model"
	"github.com/vstdy/gophermart/pkg"
	"github.com/vstdy/gophermart/storage/pgbun/schema"
)

const userTableName = "user"

// CreateUser adds given url objects to storage
func (st *Storage) CreateUser(ctx context.Context, rawObj model.User) (model.User, error) {
	logger := st.Logger(withTable(userTableName), withOperation("insert"))

	dbObj := schema.NewUserFromCanonical(rawObj)

	if err := dbObj.EncryptPassword(); err != nil {
		return model.User{}, err
	}

	_, err := st.db.NewInsert().
		Model(&dbObj).
		Returning("*").
		Exec(ctx)
	if err != nil {
		pgErr := &pgdriver.Error{}
		if errors.As(err, pgErr) {
			if pgErr.IntegrityViolation() {
				return model.User{}, pkg.ErrAlreadyExists
			}
		}
		return model.User{}, err
	}

	obj, err := dbObj.ToCanonical()
	if err != nil {
		return model.User{}, err
	}

	logger.Info().Msgf("User added %+v", obj)

	return obj, nil
}

// AuthenticateUser verifies the identity of credentials.
func (st *Storage) AuthenticateUser(ctx context.Context, rawObj model.User) (model.User, error) {
	logger := st.Logger(withTable(userTableName), withOperation("login"))

	dbObj := schema.NewUserFromCanonical(rawObj)

	err := st.db.NewSelect().
		Model(&dbObj).
		WherePK("login").
		Scan(ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return model.User{}, pkg.ErrWrongCredentials
		}
		return model.User{}, err
	}

	if err = dbObj.ComparePasswords(rawObj.Password); err != nil {
		return model.User{}, err
	}

	obj, err := dbObj.ToCanonical()
	if err != nil {
		return model.User{}, err
	}

	logger.Info().Msgf("User logged in %+v", obj)

	return obj, nil
}
