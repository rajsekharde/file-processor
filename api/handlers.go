package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/rajsekharde/file-processor/shared"
)

var WorkerURL = "http://worker:8001"

func HandleRoot(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "API Server Running")
}

// New POST /task handler using Redis as Task Queue
func HandlePostTask(w http.ResponseWriter, r *http.Request) {
	var task shared.TaskRequest
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
        http.Error(w, "Invalid request body", http.StatusBadRequest)
        return
    }

	taskID := uuid.New().String()
	taskData := map[string]interface{}{
		"id": taskID,
		"file_name": task.FileName,
		"output_format": task.OutputFormat,
		"status": "queued",
		"created_at": time.Now(),
	}
	hashKey := "task:" + taskID

	// Create hash using hash key and details
	rdb.HSet(ctx, hashKey, taskData)

	// Push the task id to list
	rdb.LPush(ctx, "tasks", taskID)

	log.Printf("Job ID: %s: Hash created and pushed to queue\n", taskID)

	response := map[string]string {
		"task_id": taskID,
		"status": "queued",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

// Get task details
func HandleGetTask(w http.ResponseWriter, r *http.Request) {
	taskID := r.PathValue("task_id")
	if taskID == "" {
		http.Error(w, "Task ID is required", http.StatusBadRequest)
		return
	}

	hashKey := "task:" + taskID

	// returns map[string]string of {field-name:value} pairs
	metadata, err := rdb.HGetAll(ctx, hashKey).Result()
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	if len(metadata) == 0 {
		http.Error(w, "Task not found", http.StatusNotFound)
		return
	}

	response := map[string]string{
		"task_id": taskID,
		"status":  metadata["status"],
		"created_at": metadata["created_at"],
	}

	if(metadata["status"] == "completed") {
		response["output_file"] = metadata["output_file"]
	}
	if(metadata["status"] == "failed") {
		response["error"] = metadata["error"]
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func HandleUpload(w http.ResponseWriter, r *http.Request) {
	file, handler, err := r.FormFile("newFile")
	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer file.Close()

	// create new local file at "../file-storage/uploads/" with same name
	dst, err := os.Create("../file-storage/uploads/" +  handler.Filename)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer dst.Close()

	// copy contents to new file
	io.Copy(dst, file)

	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "File Uploaded")
}

func HandleDownload(w http.ResponseWriter, r *http.Request) {
	fileName := r.PathValue("filename")

	if fileName == "" {
		http.Error(w, "Missing filename in path", http.StatusBadRequest)
		return
	}
	filePath := "../file-storage/completed/" + fileName

	// Check if file exists and if it's a regular file and not a directory
	info, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		http.Error(w, "File not found", http.StatusNotFound)
		return
	} else if err != nil || info.IsDir() {
		http.Error(w, "Invalid file request", http.StatusBadRequest)
		return
	}

	// Force browser to download file to disk instead of opening it in a new tab
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", info.Name()))

	http.ServeFile(w, r, filePath)
}