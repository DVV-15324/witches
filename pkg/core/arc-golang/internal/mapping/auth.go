package mapping

import (
	dtoAuth "arc-golang/internal/dto/auth/request"
	entityAuth "arc-golang/internal/entity/auth"
)

// Register_1
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
