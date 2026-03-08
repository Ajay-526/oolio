package routes

import (
	"net/http"
	"oolio/internal/config"
	"oolio/internal/interfaces/input/api/rest/handler"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
)

func InitRoutes(
	productHandler *handler.ProductHandler,
	orderHandler *handler.OrderHandler,
) http.Handler {
	router := chi.NewRouter()
	env := config.Env()

	origins := []string{}
	if env.APP_ENV == "dev" {
		origins = []string{"http://localhost:3000"}
	}

	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   origins,
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-User-Id", "X-Api-Key", "X-Timestamp", "X-Signature"},
		AllowCredentials: true,
	}))

	// product
	router.Group(func(r chi.Router) {
		r.Get("/products", productHandler.GetAllProducts)
		r.Get("/products/{id}", productHandler.GetProductByID)
		r.Post("/products", productHandler.CreateProduct)
	})
	// order
	router.Group(func(r chi.Router) {
		r.Post("/order", orderHandler.PlaceOrder)
	})

	return router
}
