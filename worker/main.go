package main

import (
	"context"
	// "fmt"
	"log"
	"time"
	"os"
	// "github.com/rajsekharde/file-processor/shared"
	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()
var rdb *redis.Client

func main() {
	var redisAddr = os.Getenv("REDIS_HOST")
	if redisAddr == "" {
		redisAddr = "localhost:6379"
	}

	rdb = redis.NewClient(&redis.Options{
		Addr: redisAddr,
		Password: "",
		DB: 0,
	})

	// Ping redis server to check if it's active
	if err := rdb.Ping(ctx).Err(); err != nil {
        log.Fatalf("Failed to connect to Redis: %v", err)
    }
    log.Printf("Successfully connected to Redis\n")

	log.Printf("Worker started. Waiting for jobs...\n\n")
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
			"status": "processing",
		})
		log.Printf("Processing Task: %s\n", taskID)

		// get task metadata
		metadata, err := rdb.HGetAll(ctx, hashKey).Result()
		if err != nil {
			rdb.HSet(ctx, hashKey, map[string]interface{}{
				"status": "failed",
			})
			log.Printf("Failed to fetch metadata for Task: %s\n\n", taskID)
			continue
		}

		// call image format conversion function
		fileName := metadata["file_name"]
		outputFormat := metadata["output_format"]
		res, message := ConvertFormat(fileName, outputFormat)
		if res != 0 {
			rdb.HSet(ctx, hashKey, map[string]interface{}{
				"status": "failed",
				"error": message,
			})
			log.Printf("Error: %s\n\n", message)
			continue
		}

		rdb.HSet(ctx, hashKey, map[string]interface{}{
			"status": "completed",
			"output_file": message,
		})
		log.Printf("Task Completed\n\n")
	}
}