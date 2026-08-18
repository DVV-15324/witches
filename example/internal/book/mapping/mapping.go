package mapping

import (
	domainBook "example/internal/shared/domain"
	"example/internal/shared/utils"
	dtoBook "example/internal/book/dto/response"
	modelBook "example/internal/book/model"
	w_utils "github.com/DVV-15324/witches/pkg/core/utils"
)

// DTO ↔ Domain

func FromDtoToDomainBook(dto *dtoBook.Book) *domainBook.Book {
	if dto == nil {
		return nil
	}
	return &domainBook.Book{
		ID:        int(utils.DecodeFromBase58(dto.ID).LocalID),
		CreatedAt: dto.CreatedAt,
		UpdatedAt: dto.UpdatedAt,
	}
}

// List DTO → List domain (dùng generic)
func FromDtoToDomainBookList(dtos []*dtoBook.Book) []*domainBook.Book {
	return w_utils.MapPtrSlice(dtos, FromDtoToDomainBook)
}

//  Domain ↔ Model

func FromDomainToModelBook(domain *domainBook.Book) *modelBook.Book {
	if domain == nil {
		return nil
	}
	return &modelBook.Book{
		ID:        domain.ID,
		CreatedAt: domain.CreatedAt,
		UpdatedAt: domain.UpdatedAt,
	}
}

// List Domain → List Model (dùng generic)
func FromDomainToModelBookList(domains []*domainBook.Book) []*modelBook.Book {
	return w_utils.MapPtrSlice(domains, FromDomainToModelBook)
}

func FromModelToDomainBook(model *modelBook.Book) *domainBook.Book {
	if model == nil {
		return nil
	}
	return &domainBook.Book{
		ID:        model.ID,
		CreatedAt: model.CreatedAt,
		UpdatedAt: model.UpdatedAt,
	}
}

// List Model → List Domain (dùng generic)
func FromModelToDomainBookList(entities []*modelBook.Book) []*domainBook.Book {
	return w_utils.MapPtrSlice(entities, FromModelToDomainBook)
}

// List Domain → List DTO
func FromDomainToDtoBook(domain *domainBook.Book) *dtoBook.Book {
	if domain == nil {
		return nil
	}
	return &dtoBook.Book{
		ID:        utils.NewUID(uint32(domain.ID), utils.ObjectBook).ToBase58(),
		CreatedAt: domain.CreatedAt,
		UpdatedAt: domain.UpdatedAt,
	}
}

// List Domain → List DTO (dùng generic)
func FromDomainToDtoBookList(domains []*domainBook.Book) []*dtoBook.Book {
	return w_utils.MapPtrSlice(domains, FromDomainToDtoBook)
}
