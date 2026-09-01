package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/rajsekharde/file-processor/shared"
	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()
var rdb *redis.Client

func main() {
	var redisAddr = os.Getenv("REDIS_HOST")
	if redisAddr == "" {
		redisAddr = "localhost:6380"
	}

	rdb = redis.NewClient(&redis.Options{
		Addr: redisAddr,
		Password: "",
		DB: 0,
	})
	
	// Ping the redis server to check if it's active
	if err := rdb.Ping(ctx).Err(); err != nil {
        log.Fatalf("Failed to connect to Redis: %v", err)
    }
    log.Println("Successfully connected to Redis")

	frontendPath := os.Getenv("FRONTEND_PATH")
	if frontendPath == "" {
		frontendPath = "../frontend"
	}

	mux := http.NewServeMux()

	root := http.FileServer(http.Dir(frontendPath))
	postTask := http.HandlerFunc(HandlePostTask)
	getTask := http.HandlerFunc(HandleGetTask)
	uploadFile := http.HandlerFunc(HandleUpload)
	downloadFile := http.HandlerFunc(HandleDownload)

	mux.Handle("GET /", shared.LoggerMiddleware(root))
	mux.Handle("POST /task", shared.LoggerMiddleware(postTask))
	mux.Handle("GET /task/{task_id}", shared.LoggerMiddleware(getTask))
	mux.Handle("POST /upload", shared.LoggerMiddleware(uploadFile))
	mux.Handle("GET /download/{filename}", shared.LoggerMiddleware(downloadFile))

	server := &http.Server{
		Addr: ":8000",
		Handler: mux,
	}

	go func() {
		log.Printf("API Server running on port 8000...\n\n")
		err := server.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Server shutting down...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10 * time.Second)
	defer cancel()

	err := server.Shutdown(shutdownCtx)
	if err != nil {
		log.Fatalf("Error while shutting down server: %v\n", err.Error())
	}

	if err := rdb.Close(); err != nil {
		log.Printf("Error closing Redis connection: %v", err)
	} else {
		log.Println("Redis connection closed successfully.")
	}

	log.Println("Server exited cleanly.")
}
