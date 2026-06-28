package mapping

import (
	dtoAuth "example/internal/dto/auth/request"
	entityAuth "example/internal/entity/auth"
)

func FromRegisterToEntityAuth(dtoAuth *dtoAuth.Register) *entityAuth.Auth {
	return &entityAuth.Auth{
		Email:    dtoAuth.Email,
		Password: dtoAuth.Password,
	}
}

func FromLoginToEntityAuth(dtoAuth *dtoAuth.Login) *entityAuth.Auth {
	return &entityAuth.Auth{
		Email:    dtoAuth.Email,
		Password: dtoAuth.Password,
	}
}
