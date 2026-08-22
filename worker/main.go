package main

import (
	"context"
	// "encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/rajsekharde/file-processor/shared"
	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()
var rdb = redis.NewClient(&redis.Options{
	Addr: "localhost:6379",
	Password: "",
	DB: 0,
})

func main() {
	mux := http.NewServeMux()

	root := http.HandlerFunc(handleRoot)
	getTask := http.HandlerFunc(handleTask)
	postTask := http.HandlerFunc(handlePostTask)
	
	mux.Handle("GET /", shared.LoggerMiddleware(root))
	mux.Handle("GET /task", shared.LoggerMiddleware(getTask))
	mux.Handle("POST /task", shared.LoggerMiddleware(postTask))

	// log.Printf("Worker Server running on port 8001...\n\n")
	// http.ListenAndServe(":8001", mux)

	if err := rdb.Ping(ctx).Err(); err != nil {
        log.Fatalf("Failed to connect to Redis: %v", err)
    }
    log.Printf("Successfully connected to Redis\n")

	log.Printf("Worker started. Waiting for jobs...\n")
	for {
		result, err := rdb.BRPop(ctx, 0*time.Second, "tasks").Result()
		if err != nil {
			log.Printf("Error popping job: %v", err)
			time.Sleep(1 * time.Second) // Prevents hard loop on connection errors
			continue
		}

		taskID := result[1]
		hashKey := "task:" + taskID

		rdb.HSet(ctx, hashKey, map[string]interface{}{
			"status":     string("processing"),
		})

		fmt.Printf("[WORKER] Task Accepted. ID: %v\n", taskID)
	}
}