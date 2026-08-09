// internal/user-service/mapping/user_mapping.go
package mapping

import (
	modelUser "example/internal/shared/model"
	"example/internal/shared/utils"
	dtoUser "example/internal/user-service/dto/response"
	entityUser "example/internal/user-service/entity"
)

// 1. DTO ↔ Model

func FromDtoToModelUser(dto *dtoUser.User) *modelUser.User {
	if dto == nil {
		return nil
	}
	return &modelUser.User{
		Id:        int(utils.DecodeFromBase58(dto.Id).LocalID),
		Name:      dto.Name,
		Email:     dto.Email,
		Role:      dto.Role,
		CreatedAt: dto.CreatedAt,
		UpdatedAt: dto.UpdatedAt,
	}
}

// 1.1 List DTO → List Model (dùng generic)
func FromDtoToModelUserList(dtos []*dtoUser.User) []*modelUser.User {
	return utils.MapPtrSlice(dtos, FromDtoToModelUser)
}

// 2. Model ↔ Entity

func FromModelToEntityUser(model *modelUser.User) *entityUser.User {
	if model == nil {
		return nil
	}
	return &entityUser.User{
		Id:        model.Id,
		Name:      model.Name,
		Email:     model.Email,
		Role:      model.Role,
		CreatedAt: model.CreatedAt,
		UpdatedAt: model.UpdatedAt,
	}
}

// List Model → List Entity (dùng generic)
func FromModelToEntityUserList(models []*modelUser.User) []*entityUser.User {
	return utils.MapPtrSlice(models, FromModelToEntityUser)
}

func FromEntityToModelUser(entity *entityUser.User) *modelUser.User {
	if entity == nil {
		return nil
	}
	return &modelUser.User{
		Id:        entity.Id,
		Name:      entity.Name,
		Email:     entity.Email,
		Role:      entity.Role,
		CreatedAt: entity.CreatedAt,
		UpdatedAt: entity.UpdatedAt,
	}
}

// List Entity → List Model (dùng generic)
func FromEntityToModelUserList(entities []*entityUser.User) []*modelUser.User {
	return utils.MapPtrSlice(entities, FromEntityToModelUser)
}

// List Model → List DTO
func FromModelToDtoUser(model *modelUser.User) *dtoUser.User {
	if model == nil {
		return nil
	}
	return &dtoUser.User{
		Id:        utils.NewUID(uint32(model.Id), utils.ObjectUser).ToBase58(),
		Name:      model.Name,
		Email:     model.Email,
		Role:      model.Role,
		CreatedAt: model.CreatedAt,
		UpdatedAt: model.UpdatedAt,
	}
}

// List Model → List DTO (dùng generic)
func FromModelToDtoUserList(models []*modelUser.User) []*dtoUser.User {
	return utils.MapPtrSlice(models, FromModelToDtoUser)
}
