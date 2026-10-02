package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/macemacemace/seraphim/backend/internal/auth"

	"github.com/macemacemace/seraphim/backend/internal/db"
)

func main() {

	if err := godotenv.Load(); err != nil {
		log.Println("no .env file found, using environment variables")
	}

	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		log.Fatal("DATABASE_URL is not set")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	pool, err := pgxpool.New(ctx, databaseURL)
	if err != nil {
		log.Fatalf("create db pool: %v", err)
	}
	defer pool.Close()

	if err := pool.Ping(ctx); err != nil {
		log.Fatalf("connect to db: %v", err)
	}
	log.Println("connected to database")

	router := gin.Default()

	api := router.Group("/api")

	queries := db.New(pool)

	authService := auth.NewService(queries)

	authHandler := auth.NewHandler(authService)

	api.POST("/auth/register", authHandler.Register)

	api.GET("/health", func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c.Request.Context(), 2*time.Second)
		defer cancel()

		if err := pool.Ping(ctx); err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{"status": "unavailable", "database": "down"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": "ok", "database": "up"})
	})

	log.Printf("listening on :%s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("server stopped: %v", err)
	}
}
