package mapping

import (
	dtoUser "example/internal/dto/user/response"

	entityUser "example/internal/entity/user"
	uid "github.com/DVV-15324/witches/pkg/core/utils"
)

func FromEntityToDtoUser(entityUser *entityUser.User) *dtoUser.User {
	return &dtoUser.User{
		Id:        uid.NewUID(uint32(entityUser.Id), objectUser).ToBase58(),
		Name:      entityUser.Name,
		Role:      entityUser.Role,
		Banned:    entityUser.Banned,
		CreatedAt: entityUser.CreatedAt,
		UpdatedAt: entityUser.UpdatedAt,
	}
}
