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

	fmt.Println(cfg)

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

	pwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	fmt.Println(pwd)

	promocodeService, err := usecase.NewPromoService([]string{
		pwd + "/static/mine/file1.txt",
		pwd + "/static/mine/file2.txt",
		pwd + "/static/mine/file3.txt",
	}, redisService)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to create promocode service", "error", err.Error())
		os.Exit(1)
	}
	orderService := usecase.NewOrderService(promocodeService, redisService, &productService, &productRepo)

	productHandler := handler.NewProductHandler(&productService)
	orderHandler := handler.NewOrderHandler(orderService)

	handler := routes.InitRoutes(productHandler, orderHandler)

	// start server
	port := fmt.Sprintf(":%s", cfg.APP_PORT)
	if err := http.ListenAndServe(port, handler); err != nil {
		slog.ErrorContext(ctx, "Failed to start server", "error", err.Error())
		os.Exit(1)
	}
}
