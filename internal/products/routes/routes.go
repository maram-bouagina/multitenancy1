package routes

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"

	"multitenancypfe/internal/middleware"
	"multitenancypfe/internal/products/handlers"
	"multitenancypfe/internal/products/repo"
	"multitenancypfe/internal/products/services"
)

func RegisterProductRoutes(app *fiber.App, db *gorm.DB) {
	store := app.Group("/api/stores/:storeId",
		middleware.RequireAuth(),
		middleware.TenantDB(),
	)

	registerProducts(store, db)
	registerCategories(store, db)
	registerCollections(store, db)
}

func registerProducts(store fiber.Router, db *gorm.DB) {
	h := handlers.NewProductHandler(services.NewProductService(repo.NewProductRepository(db)))

	g := store.Group("/products")
	g.Post("/", h.Create)
	g.Get("/", h.GetAll)
	g.Get("/:id", h.GetByID)
	g.Put("/:id", h.Update)
	g.Delete("/:id", h.Delete)
}

func registerCategories(store fiber.Router, db *gorm.DB) {
	h := handlers.NewCategoryHandler(services.NewCategoryService(repo.NewCategoryRepository(db)))

	g := store.Group("/categories")
	g.Post("/", h.Create)
	g.Get("/", h.GetTree)
	g.Get("/:id", h.GetByID)
	g.Put("/:id", h.Update)
	g.Delete("/:id", h.Delete)
}

func registerCollections(store fiber.Router, db *gorm.DB) {
	h := handlers.NewCollectionHandler(services.NewCollectionService(repo.NewCollectionRepository(db)))

	g := store.Group("/collections")
	g.Post("/", h.Create)
	g.Get("/", h.GetAll)
	g.Get("/:id", h.GetByID)
	g.Put("/:id", h.Update)
	g.Delete("/:id", h.Delete)
	g.Get("/:id/products", h.GetProducts)
	g.Post("/:id/products/:productId", h.AddProduct)
	g.Delete("/:id/products/:productId", h.RemoveProduct)
}
