package main

import (
	"context"
	// "fmt"
	"log"
	"time"
	"os"
	"encoding/json"
	"github.com/rajsekharde/file-processor/shared"
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


	// infinite loop blocking and listening on redis list
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
		taskType := metadata["type"]

		var res int
		var message string

		// call function based on operation type
		switch taskType {

		// resize image
		case "resize_image":
			var payload shared.ResizeImagePayload
			if err := json.Unmarshal([] byte(metadata["payload"]), &payload); err != nil {
				res, message = 1, "Invalid payload for image resizing"
				break
			}
			fileName := payload.FileName
			targetWidth := payload.TargetWidth
			targetHeight := payload.TargetHeight
			res, message = resizeImage(taskID, fileName, targetWidth, targetHeight)

		// convert image format
		case "convert_format":
			var payload shared.ConvertFormatPayload
			if err := json.Unmarshal([] byte(metadata["payload"]), &payload); err != nil {
				res, message = 1, "Invalid payload for format conversion"
				break
			}
			fileName := payload.FileName
			outputFormat := payload.OutputFormat
			res, message = convertFormat(taskID, fileName, outputFormat)

		// default response & message
		default:
			res = 1
			message = "Invalid task type"
		}
		
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