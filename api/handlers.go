package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"encoding/json"
	"github.com/rajsekharde/file-processor/shared"
	// "github.com/google/uuid"
)

var WorkerURL = "http://worker:8001"

func HandleRoot(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "API Server Running")
}

func HandleGetTask(w http.ResponseWriter, r *http.Request) {

	// Get task status from worker via an HTTP request
	resp, err := http.Get(WorkerURL + "/task")

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

func HandlePostTask(w http.ResponseWriter, r *http.Request) {
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Empty or invalid body", http.StatusBadRequest)
		return
	}

	// Send task to worker via an HTTP request
	url := WorkerURL + "/task"
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

	// Set header for json payload
	w.Header().Set("Content-Type", "application/json")

	// Copy status code
	w.WriteHeader(resp.StatusCode)

	// Copy json response
	w.Write(respData)
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

// Redis test func
func HandleEnqueue(w http.ResponseWriter, r *http.Request) {
	var task shared.Task
	json.NewDecoder(r.Body).Decode(&task)

	// convert task struct to JSON
	bytes, _ := json.Marshal(task)
	jobJSON := string(bytes)

	// push task to redis
	err := rdb.LPush(ctx, "email_jobs", jobJSON).Err()
	if err != nil {
		http.Error(w, "Failed to enqueue job", http.StatusInternalServerError)
		return
	}
}