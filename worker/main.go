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
	res := ConvertFormat(inputPath, task.OutputFormat)
	status := http.StatusOK
	message := "File converted successfully"
	if(res != 0) {
		status = http.StatusBadRequest
		message = "File conversion failed"
	}

	w.WriteHeader(status)
	fmt.Fprintln(w, message)
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