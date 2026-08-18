package book

import (
	core "example/cmd/server/core"

	handler "example/internal/book/handler"
	"example/internal/book/repository"
	usecase_book "example/internal/book/usecase"

	w_handl "github.com/DVV-15324/witches/pkg/core/handler"
	w_resp "github.com/DVV-15324/witches/pkg/core/response"
	"github.com/gin-gonic/gin"
	"github.com/ulule/limiter/v3"
)

type BookModule struct {
	Handler *handler.BookHandler
	Usecase *usecase_book.BookUsecase
	Repo    *repository.BookRepository
}

func NewBookModule(
	core *core.CoreServices,
) *BookModule {
	repo := repository.NewBookRepository(core.DB, core.TxManager, core.Redis.GetClient())
	usecase := usecase_book.NewBookUsecase(repo)
	handler := handler.NewBookHandler(usecase, core.Logger, core.Config)
	return &BookModule{
		Handler: handler,
		Usecase: usecase,
		Repo:    repo,
	}
}

func (m *BookModule) RegisterPublicRoutes(gen *w_handl.SwaggerGenerator, rateLimit *limiter.Rate) {
	// Public routes (nếu có)
	// Ví dụ: gen.POST("/v1/book/public", ...)
}

func (m *BookModule) RegisterProtectedRoutes(gen *w_handl.SwaggerGenerator, rateLimit *limiter.Rate, authMiddleware gin.HandlerFunc) {
	// --- Create ---
	gen.POST("/v1/book/create").
		Summary("Create Book").
		Description("Create a new Book").
		Tags("book").
		HeaderParam("User-Agent", "User-Agent", false).
		HeaderParam("Accept-Language", "Preferred language (e.g., vi-VN, en-US)", false).
		HeaderParam("X-Timezone", "Timezone (e.g., Asia/Ho_Chi_Minh)", false).
		HeaderParam("Authorization", "Bearer token for authentication", false).
		RateLimit(*rateLimit).
		Body(map[string]interface{}{}, "Request body").
		Response(200, w_resp.AppResponse{}, "Created successfully").
		Response(400, w_resp.AppResponse{}, "Bad request").
		Use(authMiddleware).
		Handler(m.Handler.Create()).
		Build()

	// --- Get by ID ---
	gen.POST("/v1/book/get").
		Summary("Get Book by ID").
		Description("Get Book details by ID").
		Tags("book").
		HeaderParam("User-Agent", "User-Agent", false).
		HeaderParam("Accept-Language", "Preferred language (e.g., vi-VN, en-US)", false).
		HeaderParam("X-Timezone", "Timezone (e.g., Asia/Ho_Chi_Minh)", false).
		HeaderParam("Authorization", "Bearer token for authentication", false).
		RateLimit(*rateLimit).
		Body(map[string]int{"id": 0}, "Book ID").
		Response(200, w_resp.AppResponse{}, "Book found").
		Response(404, w_resp.AppResponse{}, "Book not found").
		Use(authMiddleware).
		Handler(m.Handler.GetByID()).
		Build()

	// --- Get All ---
	gen.POST("/v1/book/get_all").
		Summary("Get all Books").
		Description("Get list of all Books with pagination").
		Tags("book").
		HeaderParam("User-Agent", "User-Agent", false).
		HeaderParam("Accept-Language", "Preferred language (e.g., vi-VN, en-US)", false).
		HeaderParam("X-Timezone", "Timezone (e.g., Asia/Ho_Chi_Minh)", false).
		HeaderParam("Authorization", "Bearer token for authentication", false).
		RateLimit(*rateLimit).
		QueryParam("page", "Page number", false).
		QueryParam("limit", "Items per page", false).
		QueryParam("search", "Search keyword", false).
		QueryParam("sort", "Sort field", false).
		QueryParam("order", "Sort order (asc/desc)", false).
		Response(200, w_resp.AppResponse{}, "List of Books").
		Use(authMiddleware).
		Handler(m.Handler.GetAll()).
		Build()

	// --- Update ---
	gen.PUT("/v1/book/update/:id").
		Summary("Update Book").
		Description("Update an existing Book").
		Tags("book").
		HeaderParam("User-Agent", "User-Agent", false).
		HeaderParam("Accept-Language", "Preferred language (e.g., vi-VN, en-US)", false).
		HeaderParam("X-Timezone", "Timezone (e.g., Asia/Ho_Chi_Minh)", false).
		HeaderParam("Authorization", "Bearer token for authentication", false).
		RateLimit(*rateLimit).
		PathParam("id", "Book ID", false).
		Body(map[string]interface{}{}, "Update data").
		Response(200, w_resp.AppResponse{}, "Updated successfully").
		Response(400, w_resp.AppResponse{}, "Bad request").
		Response(404, w_resp.AppResponse{}, "Book not found").
		Use(authMiddleware).
		Handler(m.Handler.Update()).
		Build()

	// --- Delete ---
	gen.DELETE("/v1/book/delete/:id").
		Summary("Delete Book").
		Description("Delete a Book by ID").
		Tags("book").
		HeaderParam("User-Agent", "User-Agent", false).
		HeaderParam("Accept-Language", "Preferred language (e.g., vi-VN, en-US)", false).
		HeaderParam("X-Timezone", "Timezone (e.g., Asia/Ho_Chi_Minh)", false).
		HeaderParam("Authorization", "Bearer token for authentication", false).
		RateLimit(*rateLimit).
		PathParam("id", "Book ID", false).
		Response(200, w_resp.AppResponse{}, "Deleted successfully").
		Response(404, w_resp.AppResponse{}, "Book not found").
		Use(authMiddleware).
		Handler(m.Handler.Delete()).
		Build()
}