package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"oolio/internal/adaptors/persistance/psql"
	"oolio/internal/adaptors/persistance/redis"
	"oolio/internal/config"
	"oolio/internal/interfaces/input/api/rest/handler"
	"oolio/internal/interfaces/input/api/rest/routes"
	"oolio/internal/usecase"
	"oolio/pkg/prettylog"
	"os"
)

func main() {
	ctx := context.Background()
	cfg := config.Env()

	slog.SetDefault(slog.New(prettylog.NewHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
		Level:     slog.LevelDebug,
	})))

	// database
	database, err := psql.NewDatabase()
	if err != nil {
		slog.ErrorContext(ctx, "Failed to connect to database", "error", err.Error())
		os.Exit(1)
	}

	// redis
	redisService, err := redis.NewRedisService()
	if err != nil {
		slog.ErrorContext(ctx, "Failed to create redis service", "error", err.Error())
		os.Exit(1)
	}
	if err := redisService.Connect(ctx); err != nil {
		slog.ErrorContext(ctx, "Failed to connect to Redis", "error", err.Error())
		os.Exit(1)
	}
	defer redisService.Close()

	productRepo := psql.NewProductRepo(database)
	productService := usecase.NewProductService(productRepo)

	productHandler := handler.NewProductHandler(&productService)

	handler := routes.InitRoutes(productHandler)

	// start server
	if err := http.ListenAndServe(cfg.APP_PORT, handler); err != nil {
		slog.ErrorContext(ctx, "Failed to start server", "error", err.Error())
		os.Exit(1)
	}
	fmt.Println("Serving listening on port", cfg.APP_PORT)
}
