package routers

import (
	"example/cmd/server/core"
	"example/internal/auth"
	"example/internal/book"
	"example/internal/refresh"
	"example/internal/user"
)

type Modules struct {
	Auth    *auth.AuthModule
	User    *user.UserModule
	Refresh *refresh.RefreshModule
	Book    *book.BookModule
}

func InitModules(core *core.CoreServices) *Modules {

	// Khởi tạo các module KHÔNG phụ thuộc lẫn nhau
	userModule := user.NewUserModule(core)
	refreshModule := refresh.NewRefreshModule(core, userModule.Usecase)
	authModule := auth.NewAuthModule(core, userModule.Usecase, refreshModule.Usecase)
	refreshModule.SetAuthUseCase(authModule.Usecase)
	bookModule := book.NewBookModule(core)

	return &Modules{
		Auth:    authModule,
		User:    userModule,
		Refresh: refreshModule, Book: bookModule,
	}
}
