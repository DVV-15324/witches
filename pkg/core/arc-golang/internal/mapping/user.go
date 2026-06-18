package mapping

import (
	dtoUser "arc-golang/internal/dto/user/response"

	entityUser "arc-golang/internal/entity/user"
	uid "arc-golang/internal/utils"
)

func ToDtoUser(entityUser *entityUser.User) *dtoUser.User {
	return &dtoUser.User{
		Id:        uid.NewUID(uint32(entityUser.Id), objectUser).ToBase58(),
		Name:      entityUser.Name,
		Role:      entityUser.Role,
		Banned:    entityUser.Banned,
		CreatedAt: entityUser.CreatedAt,
		UpdatedAt: entityUser.UpdatedAt,
	}
}
