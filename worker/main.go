package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/rajsekharde/file-processor/shared"
)

func handleRoot(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "Worker Server Running")
}

func handleTask(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "Task Processed")
}

func handlePostTask(w http.ResponseWriter, r*http.Request) {
	var task shared.Task
	json.NewDecoder(r.Body).Decode(&task)

	log.Println("Received", task)

	inputPath := "../file-storage/uploads/" + task.FileName
	result, outputPath := ConvertFormat(inputPath, task.OutputFormat)
	var status int
	var message string
	if(result != 0) {
		status = http.StatusBadRequest
		message = fmt.Sprintf(`{"status": "failed", "error": %q}`, outputPath)
	} else {
		status = http.StatusOK
		message = fmt.Sprintf(`{"status": "success", "output_path": %q}`, outputPath)
	}

	w.WriteHeader(status)
	w.Write([]byte(message))
}

func main() {
	mux := http.NewServeMux()

	root := http.HandlerFunc(handleRoot)
	getTask := http.HandlerFunc(handleTask)
	postTask := http.HandlerFunc(handlePostTask)
	
	mux.Handle("GET /", shared.LoggerMiddleware(root))
	mux.Handle("GET /task", shared.LoggerMiddleware(getTask))
	mux.Handle("POST /task", shared.LoggerMiddleware(postTask))

	log.Printf("Worker Server running on port 8001...\n\n")
	http.ListenAndServe(":8001", mux)
}