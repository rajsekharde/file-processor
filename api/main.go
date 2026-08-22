package main

import (
	"log"
	"net/http"
	"context"
	"github.com/redis/go-redis/v9"
	"github.com/rajsekharde/file-processor/shared"
)

var ctx = context.Background()
var rdb = redis.NewClient(&redis.Options{
	Addr: "localhost:6379",
	Password: "",
	DB: 0,
})

func main() {
	if err := rdb.Ping(ctx).Err(); err != nil {
        log.Fatalf("Failed to connect to Redis: %v", err)
    }
    log.Println("Successfully connected to Redis")

	mux := http.NewServeMux()

	root := http.FileServer(http.Dir("/frontend"))
	task := http.HandlerFunc(HandleGetTask)
	postTask := http.HandlerFunc(HandlePostTask)
	uploadFile := http.HandlerFunc(HandleUpload)
	downloadFile := http.HandlerFunc(HandleDownload)
	enqueue := http.HandlerFunc(HandleEnqueue)

	mux.Handle("GET /", shared.LoggerMiddleware(root))
	mux.Handle("GET /task", shared.LoggerMiddleware(task))
	mux.Handle("POST /task", shared.LoggerMiddleware(postTask))
	mux.Handle("POST /upload", shared.LoggerMiddleware(uploadFile))
	mux.Handle("GET /download/{filename}", shared.LoggerMiddleware(downloadFile))
	mux.Handle("POST /enqueue", shared.LoggerMiddleware(enqueue))

	log.Printf("API Server running on port 8000...\n\n")
	http.ListenAndServe(":8000", mux)
}
