package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
)

func HandleRoot(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "API Server Running")
}

func HandleGetTask(w http.ResponseWriter, r *http.Request) {
	resp, err := http.Get("http://localhost:8001/task")

	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Could not fetch API", http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	responseData, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Could not read data", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, string(responseData))
}

type Task struct {
	TaskName string `json:"task_name"`
	TaskValue int `json:"task_value"`
}

func HandlePostTask(w http.ResponseWriter, r *http.Request) {
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Empty or invalid body", http.StatusBadRequest)
		return
	}

	url := "http://localhost:8001/task"
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(bodyBytes))
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Failed to send task to worker", http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	respData, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Failed to read request", http.StatusInternalServerError)
		return
	}
	log.Println("Worker:", string(respData))

	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "Task Accepted")
}