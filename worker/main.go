package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
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

type Task struct {
	TaskName string `json:"task_name"`
	TaskValue int `json:"task_value"`
}

func handlePostTask(w http.ResponseWriter, r*http.Request) {
	var task Task
	json.NewDecoder(r.Body).Decode(&task)

	log.Println("Received", task)
	log.Printf("%s %s %d\n", r.Method, r.URL.Path, http.StatusOK)
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "Task Processed")
}

func main() {
	mux := http.NewServeMux()
	
	mux.HandleFunc("GET /", handleRoot)
	mux.HandleFunc("GET /task", handleTask)
	mux.HandleFunc("POST /task", handlePostTask)

	log.Println("Worker Server running on port 8001...")
	http.ListenAndServe(":8001", mux)
}