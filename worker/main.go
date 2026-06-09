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
	log.Printf("%s %s %d\n", r.Method, r.URL.Path, http.StatusOK)
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "Task Accepted")
}

func main() {
	mux := http.NewServeMux()
	
	mux.HandleFunc("GET /", handleRoot)
	mux.HandleFunc("GET /task", handleTask)
	mux.HandleFunc("POST /task", handlePostTask)

	log.Println("Worker Server running on port 8001...")
	http.ListenAndServe(":8001", mux)
}