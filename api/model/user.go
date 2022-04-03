package model

import (
	"github.com/vstdy0/go-diploma/model"
)

type RegisterBody struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

// ToCanonical converts a API model to canonical model.
func (b RegisterBody) ToCanonical() model.User {
	obj := model.User{
		Login:    b.Login,
		Password: b.Password,
	}

	return obj
}

func NewJWTClaims(obj model.User) map[string]interface{} {
	return map[string]interface{}{
		"id": obj.ID,
	}
}
