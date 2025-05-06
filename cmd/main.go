package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/didsqq/user_api/internal/handler"
	"github.com/didsqq/user_api/internal/repository"
	"github.com/didsqq/user_api/internal/service"
	_ "github.com/lib/pq"
)

//	@title			User Service
//	@version		1.0
//	@description	Simple API server for user

//	@host		localhost:8080
//	@BasePath	/

func main() {
	// if err := godotenv.Load(); err != nil {
	// 	log.Fatal("Error loading .env file")
	// }

	dbPass := getEnvVar("DB_PASSWORD")
	dbUser := getEnvVar("DB_USER")
	dbName := getEnvVar("DB_NAME")
	dbHost := getEnvVar("DB_HOST")
	dbPort := getEnvVar("DB_PORT")
	srvPort := getEnvVar("SRV_PORT")

	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", dbUser, dbPass, dbHost, dbPort, dbName)

	db, err := repository.NewPostgresDB(dsn)
	if err != nil {
		log.Fatalf("Error open db: %s", err)
	}
	defer db.Close()

	repo := repository.NewRepository(db)
	service := service.NewService(repo)
	handler := handler.NewHandler(service)

	r := handler.InitRoutes()

	srv := &http.Server{
		Addr:    ":" + srvPort,
		Handler: r,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	log.Print("server is running")

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	<-stop
	log.Print("stopping server")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("failed to stop server: %v", err)

		return
	}

	log.Print("server stopped")
}

func getEnvVar(key string) string {
	value := os.Getenv(key)
	if value == "" {
		log.Fatalf("Ошибка: переменная окружения %s не найдена", key)
	}
	return value
}
