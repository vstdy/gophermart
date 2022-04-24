package schema

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
	"golang.org/x/crypto/bcrypt"

	"github.com/vstdy/gophermart/model"
	"github.com/vstdy/gophermart/pkg"
)

// User keeps user data.
type User struct {
	bun.BaseModel `bun:"users,alias:u"`
	ID            uuid.UUID `bun:"id,pk"`
	Login         string    `bun:"login,unique,notnull"`
	Password      string    `bun:"password,notnull"`
	CreatedAt     time.Time `bun:"created_at,nullzero,notnull,default:current_timestamp"`
	UpdatedAt     time.Time `bun:"updated_at,nullzero,notnull,default:current_timestamp"`
	DeletedAt     time.Time `bun:"deleted_at,nullzero,soft_delete"`
}

func (u *User) EncryptPassword() error {
	hash, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("encrypting password: %w", err)
	}
	u.Password = string(hash)

	return nil
}

func (u *User) ComparePasswords(password string) error {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	if err != nil {
		return pkg.ErrWrongCredentials
	}

	return nil
}

// NewUserFromCanonical creates a new User DB object from canonical model.
func NewUserFromCanonical(obj model.User) User {
	return User{
		ID:        obj.ID,
		Login:     obj.Login,
		Password:  obj.Password,
		CreatedAt: obj.CreatedAt,
		UpdatedAt: obj.UpdatedAt,
		DeletedAt: obj.DeletedAt,
	}
}

// ToCanonical converts a DB object to canonical model.
func (u User) ToCanonical() (model.User, error) {
	return model.User{
		ID:        u.ID,
		Login:     u.Login,
		Password:  u.Password,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
		DeletedAt: u.DeletedAt,
	}, nil
}
