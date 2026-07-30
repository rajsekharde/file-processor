package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/rajsekharde/file-processor/shared"
	// "github.com/rajsekharde/file-processor/shared"
)

func handleRoot(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "Worker Server Running")
	log.Printf("%s %s %d\n", r.Method, r.URL.Path, http.StatusOK)
}

func handleTask(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "Task Processed")
	log.Printf("%s %s %d\n", r.Method, r.URL.Path, http.StatusOK)
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
	log.Printf("%s %s %d\n\n", r.Method, r.URL.Path, status)
}

func main() {
	mux := http.NewServeMux()
	
	mux.HandleFunc("GET /", handleRoot)
	mux.HandleFunc("GET /task", handleTask)
	mux.HandleFunc("POST /task", handlePostTask)

	log.Printf("Worker Server running on port 8001...\n\n")
	http.ListenAndServe(":8001", mux)
}