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

	result, message := ConvertFormat(task.FileName, task.OutputFormat)
	var status int
	var body string
	if(result != 0) {
		status = http.StatusBadRequest
		body = fmt.Sprintf(`{"status": "failed", "error": %q}`, message)
	} else {
		status = http.StatusOK
		body = fmt.Sprintf(`{"status": "success", "output_file": %q}`, message)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write([]byte(body))
}